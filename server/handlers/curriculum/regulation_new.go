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

// GetRegulationsNew retrieves all regulations (new system)
func GetRegulationsNew(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.DB.Query(`
		SELECT id, code, name, status, created_at, updated_at 
		FROM regulations 
		ORDER BY created_at DESC
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	regulations := []models.Regulation{}
	for rows.Next() {
		var reg models.Regulation
		if err := rows.Scan(&reg.ID, &reg.Code, &reg.Name, &reg.Status, &reg.CreatedAt, &reg.UpdatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		regulations = append(regulations, reg)
	}

	json.NewEncoder(w).Encode(regulations)
}

// GetRegulationByID retrieves a single regulation by ID
func GetRegulationByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	var reg models.Regulation
	err := db.DB.QueryRow(`
		SELECT id, code, name, status, created_at, updated_at 
		FROM regulations 
		WHERE id = ?
	`, id).Scan(&reg.ID, &reg.Code, &reg.Name, &reg.Status, &reg.CreatedAt, &reg.UpdatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "Regulation not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(reg)
}

// CreateRegulationNew creates a new regulation
func CreateRegulationNew(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var reg models.Regulation
	if err := json.NewDecoder(r.Body).Decode(&reg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if reg.Code == "" || reg.Name == "" {
		http.Error(w, "Code and Name are required", http.StatusBadRequest)
		return
	}

	// Set default status if not provided
	if reg.Status == "" {
		reg.Status = "DRAFT"
	}

	result, err := db.DB.Exec(`
		INSERT INTO regulations (code, name, status, created_at, updated_at) 
		VALUES (?, ?, ?, NOW(), NOW())
	`, reg.Code, reg.Name, reg.Status)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reg.ID = int(id)
	json.NewEncoder(w).Encode(reg)
}

// UpdateRegulationNew updates an existing regulation
func UpdateRegulationNew(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	var reg models.Regulation
	if err := json.NewDecoder(r.Body).Decode(&reg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.DB.Exec(`
		UPDATE regulations 
		SET code = ?, name = ?, status = ?, updated_at = NOW() 
		WHERE id = ?
	`, reg.Code, reg.Name, reg.Status, id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reg.ID, _ = strconv.Atoi(id)
	json.NewEncoder(w).Encode(reg)
}

// DeleteRegulationNew deletes a regulation
func DeleteRegulationNew(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := db.DB.Exec("DELETE FROM regulations WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Regulation deleted successfully"})
}

// GetRegulationClauses retrieves all clauses for a regulation
func GetRegulationClauses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	regulationID := vars["id"]

	rows, err := db.DB.Query(`
		SELECT id, regulation_id, section_no, clause_no, title, content, created_at, updated_at 
		FROM regulation_clauses 
		WHERE regulation_id = ? 
		ORDER BY section_no, clause_no
	`, regulationID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	clauses := []models.RegulationClause{}
	for rows.Next() {
		var clause models.RegulationClause
		if err := rows.Scan(&clause.ID, &clause.RegulationID, &clause.SectionNo, &clause.ClauseNo, &clause.Title, &clause.Content, &clause.CreatedAt, &clause.UpdatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		clauses = append(clauses, clause)
	}

	json.NewEncoder(w).Encode(clauses)
}

// CreateRegulationClause creates a new clause
func CreateRegulationClause(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	regulationID := vars["id"]

	var clause models.RegulationClause
	if err := json.NewDecoder(r.Body).Decode(&clause); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	regID, _ := strconv.Atoi(regulationID)
	clause.RegulationID = regID

	result, err := db.DB.Exec(`
		INSERT INTO regulation_clauses (regulation_id, section_no, clause_no, title, content, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, NOW(), NOW())
	`, clause.RegulationID, clause.SectionNo, clause.ClauseNo, clause.Title, clause.Content)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	clause.ID = int(id)

	json.NewEncoder(w).Encode(clause)
}

// UpdateRegulationClause updates a clause and logs the change
func UpdateRegulationClause(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	clauseID := vars["clauseId"]

	var clause models.RegulationClause
	if err := json.NewDecoder(r.Body).Decode(&clause); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get old content for history
	var oldContent string
	err := db.DB.QueryRow("SELECT content FROM regulation_clauses WHERE id = ?", clauseID).Scan(&oldContent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update clause
	_, err = db.DB.Exec(`
		UPDATE regulation_clauses 
		SET section_no = ?, clause_no = ?, title = ?, content = ?, updated_at = NOW() 
		WHERE id = ?
	`, clause.SectionNo, clause.ClauseNo, clause.Title, clause.Content, clauseID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Log to history
	_, _ = db.DB.Exec(`
		INSERT INTO regulation_clause_history (clause_id, old_content, new_content, changed_by, changed_at, change_reason) 
		VALUES (?, ?, ?, ?, NOW(), ?)
	`, clauseID, oldContent, clause.Content, "system", "Updated via API")

	clause.ID, _ = strconv.Atoi(clauseID)
	json.NewEncoder(w).Encode(clause)
}

// DeleteRegulationClause deletes a clause
func DeleteRegulationClause(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	clauseID := vars["clauseId"]

	_, err := db.DB.Exec("DELETE FROM regulation_clauses WHERE id = ?", clauseID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Clause deleted successfully"})
}

// package handlers

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"
// 	"strconv"

// 	"server/db"
// 	"server/models"
// )

// // GetRegulations retrieves all regulations from the database
// func GetRegulations(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Header().Set("Access-Control-Allow-Origin", "*")

// 	query := "SELECT id, name, academic_year, max_credits, created_at FROM curriculum ORDER BY created_at DESC"
// 	rows, err := db.DB.Query(query)
// 	if err != nil {
// 		log.Println("Error querying curriculum:", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch curriculum"})
// 		return
// 	}
// 	defer rows.Close()

// 	var regulations []models.Regulation = make([]models.Regulation, 0)
// 	for rows.Next() {
// 		var reg models.Regulation
// 		err := rows.Scan(&reg.ID, &reg.Name, &reg.AcademicYear, &reg.MaxCredits, &reg.CreatedAt)
// 		if err != nil {
// 			log.Println("Error scanning curriculum:", err)
// 			continue
// 		}
// 		regulations = append(regulations, reg)
// 	}

// 	if err := rows.Err(); err != nil {
// 		log.Println("Rows iteration error:", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch curriculum"})
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(regulations)
// }

// // CreateRegulation creates a new regulation in the database
// func CreateRegulation(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Header().Set("Access-Control-Allow-Origin", "*")

// 	if r.Method != http.MethodPost {
// 		w.WriteHeader(http.StatusMethodNotAllowed)
// 		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
// 		return
// 	}

// 	var reg models.Regulation
// 	err := json.NewDecoder(r.Body).Decode(&reg)
// 	if err != nil {
// 		log.Println("Error decoding request body:", err)
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
// 		return
// 	}

// 	query := "INSERT INTO curriculum (name, academic_year, max_credits) VALUES (?, ?, ?)"
// 	result, err := db.DB.Exec(query, reg.Name, reg.AcademicYear, reg.MaxCredits)
// 	if err != nil {
// 		log.Println("Error inserting curriculum:", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create curriculum"})
// 		return
// 	}

// 	id, _ := result.LastInsertId()
// 	reg.ID = int(id)

// 	// Log the activity
// 	LogCurriculumActivity(int(id), "Curriculum Created",
// 		"Created new curriculum: "+reg.Name+" ("+reg.AcademicYear+")", "System")

// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(reg)
// }

// // DeleteRegulation deletes a regulation from the database
// func DeleteRegulation(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Header().Set("Access-Control-Allow-Origin", "*")

// 	if r.Method != http.MethodDelete {
// 		w.WriteHeader(http.StatusMethodNotAllowed)
// 		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
// 		return
// 	}

// 	idStr := r.URL.Query().Get("id")
// 	if idStr == "" {
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode(map[string]string{"error": "ID parameter is required"})
// 		return
// 	}

// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid ID"})
// 		return
// 	}

// 	query := "DELETE FROM curriculum WHERE id = ?"
// 	result, err := db.DB.Exec(query, id)
// 	if err != nil {
// 		log.Println("Error deleting curriculum:", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete curriculum"})
// 		return
// 	}

// 	rowsAffected, _ := result.RowsAffected()
// 	if rowsAffected == 0 {
// 		w.WriteHeader(http.StatusNotFound)
// 		json.NewEncoder(w).Encode(map[string]string{"error": "Curriculum not found"})
// 		return
// 	}

// 	// Note: Log will be automatically deleted due to CASCADE on foreign key

// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(map[string]string{"message": "Curriculum deleted successfully"})
// }
