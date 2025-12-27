package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"server/db"
	"server/models"
)

// GetRegulations retrieves all regulations from the database
func GetRegulations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	query := "SELECT id, name, academic_year, created_at FROM curriculum ORDER BY created_at DESC"
	rows, err := db.DB.Query(query)
	if err != nil {
		log.Println("Error querying curriculum:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch curriculum"})
		return
	}
	defer rows.Close()

	var regulations []models.Regulation = make([]models.Regulation, 0)
	for rows.Next() {
		var reg models.Regulation
		err := rows.Scan(&reg.ID, &reg.Name, &reg.AcademicYear, &reg.CreatedAt)
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

// CreateRegulation creates a new regulation in the database
func CreateRegulation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	var reg models.Regulation
	err := json.NewDecoder(r.Body).Decode(&reg)
	if err != nil {
		log.Println("Error decoding request body:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	query := "INSERT INTO regulations (name, academic_year) VALUES (?, ?)"
	result, err := db.DB.Exec(query, reg.Name, reg.AcademicYear)
	if err != nil {
		log.Println("Error inserting regulation:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create regulation"})
		return
	}

	id, _ := result.LastInsertId()
	reg.ID = int(id)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(reg)
}

// DeleteRegulation deletes a regulation from the database
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

	query := "DELETE FROM regulations WHERE id = ?"
	result, err := db.DB.Exec(query, id)
	if err != nil {
		log.Println("Error deleting regulation:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete regulation"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Regulation not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Regulation deleted successfully"})
}
