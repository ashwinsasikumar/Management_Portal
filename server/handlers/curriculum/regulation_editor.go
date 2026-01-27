package curriculum

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"server/db"
	"server/models"
	"strconv"

	"github.com/gorilla/mux"
)

// GetRegulationStructure retrieves the complete regulation with sections and clauses
func GetRegulationStructure(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	regulationID := vars["id"]

	// Get regulation details
	var reg models.Regulation
	err := db.DB.QueryRow(`
		SELECT id, code, name, status, created_at, updated_at 
		FROM regulations 
		WHERE id = ?
	`, regulationID).Scan(&reg.ID, &reg.Code, &reg.Name, &reg.Status, &reg.CreatedAt, &reg.UpdatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "Regulation not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get sections
	sectionRows, err := db.DB.Query(`
		SELECT id, regulation_id, section_no, title, display_order, created_at, updated_at 
		FROM regulation_sections 
		WHERE regulation_id = ? 
		ORDER BY display_order, section_no
	`, regulationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer sectionRows.Close()

	sections := []models.RegulationSectionWithClauses{}
	for sectionRows.Next() {
		var section models.RegulationSection
		if err := sectionRows.Scan(&section.ID, &section.RegulationID, &section.SectionNo,
			&section.Title, &section.DisplayOrder, &section.CreatedAt, &section.UpdatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Get clauses for this section
		clauseRows, err := db.DB.Query(`
			SELECT id, regulation_id, section_id, section_no, clause_no, title, content, display_order, created_at, updated_at 
			FROM regulation_clauses 
			WHERE section_id = ? 
			ORDER BY display_order, clause_no
		`, section.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		clauses := []models.RegulationClause{}
		for clauseRows.Next() {
			var clause models.RegulationClause
			if err := clauseRows.Scan(&clause.ID, &clause.RegulationID, &clause.SectionID, &clause.SectionNo,
				&clause.ClauseNo, &clause.Title, &clause.Content, &clause.DisplayOrder,
				&clause.CreatedAt, &clause.UpdatedAt); err != nil {
				clauseRows.Close()
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			clauses = append(clauses, clause)
		}
		clauseRows.Close()

		sectionWithClauses := models.RegulationSectionWithClauses{
			RegulationSection: section,
			Clauses:           clauses,
		}
		sections = append(sections, sectionWithClauses)
	}

	structure := models.RegulationStructure{
		Regulation: reg,
		Sections:   sections,
	}

	json.NewEncoder(w).Encode(structure)
}

// CreateSection creates a new section in a regulation
func CreateSection(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	regulationID := vars["id"]

	// Check if regulation is LOCKED
	var status string
	err := db.DB.QueryRow("SELECT status FROM regulations WHERE id = ?", regulationID).Scan(&status)
	if err != nil {
		http.Error(w, "Regulation not found", http.StatusNotFound)
		return
	}
	if status == "LOCKED" {
		http.Error(w, "Cannot edit LOCKED regulation", http.StatusForbidden)
		return
	}

	var section models.RegulationSection
	if err := json.NewDecoder(r.Body).Decode(&section); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	regID, _ := strconv.Atoi(regulationID)
	section.RegulationID = regID

	// Get max display order
	var maxOrder int
	db.DB.QueryRow("SELECT COALESCE(MAX(display_order), 0) FROM regulation_sections WHERE regulation_id = ?",
		regulationID).Scan(&maxOrder)
	section.DisplayOrder = maxOrder + 1

	result, err := db.DB.Exec(`
		INSERT INTO regulation_sections (regulation_id, section_no, title, display_order, created_at, updated_at) 
		VALUES (?, ?, ?, ?, NOW(), NOW())
	`, section.RegulationID, section.SectionNo, section.Title, section.DisplayOrder)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	section.ID = int(id)

	json.NewEncoder(w).Encode(section)
}

// UpdateSection updates a section
func UpdateSection(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	sectionID := vars["sectionId"]

	// Check if regulation is LOCKED
	var status string
	err := db.DB.QueryRow(`
		SELECT r.status FROM regulations r 
		JOIN regulation_sections s ON r.id = s.curriculum_id 
		WHERE s.id = ?
	`, sectionID).Scan(&status)
	if err != nil {
		http.Error(w, "Section not found", http.StatusNotFound)
		return
	}
	if status == "LOCKED" {
		http.Error(w, "Cannot edit LOCKED regulation", http.StatusForbidden)
		return
	}

	var section models.RegulationSection
	if err := json.NewDecoder(r.Body).Decode(&section); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec(`
		UPDATE regulation_sections 
		SET title = ?, display_order = ?, updated_at = NOW() 
		WHERE id = ?
	`, section.Title, section.DisplayOrder, sectionID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	section.ID, _ = strconv.Atoi(sectionID)
	json.NewEncoder(w).Encode(section)
}

// DeleteSection deletes a section (only if it has no clauses)
func DeleteSection(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	sectionID := vars["sectionId"]

	// Check if regulation is LOCKED
	var status string
	err := db.DB.QueryRow(`
		SELECT r.status FROM regulations r 
		JOIN regulation_sections s ON r.id = s.curriculum_id 
		WHERE s.id = ?
	`, sectionID).Scan(&status)
	if err != nil {
		http.Error(w, "Section not found", http.StatusNotFound)
		return
	}
	if status == "LOCKED" {
		http.Error(w, "Cannot edit LOCKED regulation", http.StatusForbidden)
		return
	}

	// Check if section has clauses
	var clauseCount int
	db.DB.QueryRow("SELECT COUNT(*) FROM regulation_clauses WHERE section_id = ?", sectionID).Scan(&clauseCount)
	if clauseCount > 0 {
		http.Error(w, "Cannot delete section with clauses", http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec("DELETE FROM regulation_sections WHERE id = ?", sectionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Section deleted successfully"})
}

// CreateClause creates a new clause in a section
func CreateClause(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	sectionID := vars["sectionId"]

	// Check if regulation is LOCKED
	var status string
	var regulationID int
	err := db.DB.QueryRow(`
		SELECT r.id, r.status FROM regulations r 
		JOIN regulation_sections s ON r.id = s.curriculum_id 
		WHERE s.id = ?
	`, sectionID).Scan(&regulationID, &status)
	if err != nil {
		http.Error(w, "Section not found", http.StatusNotFound)
		return
	}
	if status == "LOCKED" {
		http.Error(w, "Cannot edit LOCKED regulation", http.StatusForbidden)
		return
	}

	var clause models.RegulationClause
	if err := json.NewDecoder(r.Body).Decode(&clause); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	secID, _ := strconv.Atoi(sectionID)
	clause.SectionID = secID
	clause.RegulationID = regulationID

	// Get section number
	var sectionNo int
	db.DB.QueryRow("SELECT section_no FROM regulation_sections WHERE id = ?", sectionID).Scan(&sectionNo)
	clause.SectionNo = sectionNo

	// Get max display order
	var maxOrder int
	db.DB.QueryRow("SELECT COALESCE(MAX(display_order), 0) FROM regulation_clauses WHERE section_id = ?",
		sectionID).Scan(&maxOrder)
	clause.DisplayOrder = maxOrder + 1

	result, err := db.DB.Exec(`
		INSERT INTO regulation_clauses 
		(regulation_id, section_id, section_no, clause_no, title, content, display_order, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`, clause.RegulationID, clause.SectionID, clause.SectionNo, clause.ClauseNo,
		clause.Title, clause.Content, clause.DisplayOrder)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	clause.ID = int(id)

	json.NewEncoder(w).Encode(clause)
}

// UpdateClause updates a clause and creates history
func UpdateClause(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	clauseID := vars["clauseId"]

	// Check if regulation is LOCKED
	var status string
	err := db.DB.QueryRow(`
		SELECT r.status FROM regulations r 
		JOIN regulation_clauses c ON r.id = c.regulation_id 
		WHERE c.id = ?
	`, clauseID).Scan(&status)
	if err != nil {
		http.Error(w, "Clause not found", http.StatusNotFound)
		return
	}
	if status == "LOCKED" {
		http.Error(w, "Cannot edit LOCKED regulation", http.StatusForbidden)
		return
	}

	var clause models.RegulationClause
	if err := json.NewDecoder(r.Body).Decode(&clause); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get old content for history
	var oldContent string
	err = db.DB.QueryRow("SELECT content FROM regulation_clauses WHERE id = ?", clauseID).Scan(&oldContent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update clause
	_, err = db.DB.Exec(`
		UPDATE regulation_clauses 
		SET title = ?, content = ?, display_order = ?, updated_at = NOW() 
		WHERE id = ?
	`, clause.Title, clause.Content, clause.DisplayOrder, clauseID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Log to history only if content changed
	if oldContent != clause.Content {
		changedBy := "system" // TODO: Get from authentication
		if userHeader := r.Header.Get("X-User-Email"); userHeader != "" {
			changedBy = userHeader
		}

		_, _ = db.DB.Exec(`
			INSERT INTO regulation_clause_history 
			(clause_id, old_content, new_content, changed_by, changed_at, change_reason) 
			VALUES (?, ?, ?, ?, NOW(), ?)
		`, clauseID, oldContent, clause.Content, changedBy, "Updated via editor")
	}

	clause.ID, _ = strconv.Atoi(clauseID)
	json.NewEncoder(w).Encode(clause)
}

// DeleteClause deletes a clause
func DeleteClause(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	clauseID := vars["clauseId"]

	// Check if regulation is LOCKED
	var status string
	err := db.DB.QueryRow(`
		SELECT r.status FROM regulations r 
		JOIN regulation_clauses c ON r.id = c.regulation_id 
		WHERE c.id = ?
	`, clauseID).Scan(&status)
	if err != nil {
		http.Error(w, "Clause not found", http.StatusNotFound)
		return
	}
	if status == "LOCKED" {
		http.Error(w, "Cannot edit LOCKED regulation", http.StatusForbidden)
		return
	}

	_, err = db.DB.Exec("DELETE FROM regulation_clauses WHERE id = ?", clauseID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Clause deleted successfully"})
}

// GetClauseHistory retrieves the edit history of a clause
func GetClauseHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	clauseID := vars["clauseId"]

	rows, err := db.DB.Query(`
		SELECT id, clause_id, old_content, new_content, changed_by, changed_at, change_reason 
		FROM regulation_clause_history 
		WHERE clause_id = ? 
		ORDER BY changed_at DESC
	`, clauseID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	history := []models.RegulationClauseHistory{}
	for rows.Next() {
		var h models.RegulationClauseHistory
		if err := rows.Scan(&h.ID, &h.ClauseID, &h.OldContent, &h.NewContent, &h.ChangedBy, &h.ChangedAt, &h.ChangeReason); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		history = append(history, h)
	}

	json.NewEncoder(w).Encode(history)
}
