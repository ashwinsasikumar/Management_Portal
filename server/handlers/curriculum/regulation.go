package curriculum

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"server/db"
	"server/models"
)

// GetRegulations retrieves all curriculum (legacy) from the database
func GetRegulations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	query := "SELECT id, name, academic_year, max_credits, curriculum_template, created_at FROM curriculum WHERE status = 1 ORDER BY created_at DESC"
	rows, err := db.DB.Query(query)
	if err != nil {
		log.Println("Error querying curriculum:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch curriculum"})
		return
	}
	defer rows.Close()

	var regulations []models.LegacyRegulation = make([]models.LegacyRegulation, 0)
	for rows.Next() {
		var reg models.LegacyRegulation
		err := rows.Scan(&reg.ID, &reg.Name, &reg.AcademicYear, &reg.MaxCredits, &reg.CurriculumTemplate, &reg.CreatedAt)
		if err != nil {
			log.Println("Error scanning curriculum:", err)
			continue
		}
		regulations = append(regulations, reg)
	}

	if err := rows.Err(); err != nil {
		log.Println("Rows iteration error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch curriculum"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(regulations)
}

// CreateRegulation creates a new curriculum (legacy) in the database
func CreateRegulation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	var reg models.LegacyRegulation
	err := json.NewDecoder(r.Body).Decode(&reg)
	if err != nil {
		log.Println("Error decoding request body:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	if reg.CurriculumTemplate == "" {
		reg.CurriculumTemplate = "2026"
	}

	query := "INSERT INTO curriculum (name, academic_year, max_credits, curriculum_template) VALUES (?, ?, ?, ?)"
	result, err := db.DB.Exec(query, reg.Name, reg.AcademicYear, reg.MaxCredits, reg.CurriculumTemplate)
	if err != nil {
		log.Println("Error inserting curriculum:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create curriculum"})
		return
	}

	id, _ := result.LastInsertId()
	reg.ID = int(id)

	// Log the activity
	LogCurriculumActivity(int(id), "Curriculum Created",
		"Created new curriculum: "+reg.Name+" ("+reg.AcademicYear+")", "System")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(reg)
}

// DeleteRegulation deletes a curriculum (legacy) from the database
func DeleteRegulation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ID parameter is required"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid ID"})
		return
	}

	query := "UPDATE curriculum SET status = 0 WHERE id = ? AND status = 1"
	result, err := db.DB.Exec(query, id)
	if err != nil {
		log.Println("Error deleting curriculum:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete curriculum"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Curriculum not found"})
		return
	}

	// Note: Log will be automatically deleted due to CASCADE on foreign key

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Curriculum deleted successfully"})
}
