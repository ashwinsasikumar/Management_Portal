package curriculum

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/db"
	"server/models"
	"strconv"

	"github.com/gorilla/mux"
)

// GetPEOPOMapping handles GET /curriculum/:id/peo-po-mapping
func GetPEOPOMapping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	curriculumID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid curriculum ID", http.StatusBadRequest)
		return
	}

	// Fetch existing PEO-PO mappings
	matrix := make(map[string]int)
	rows, err := db.DB.Query("SELECT peo_index, po_index, mapping_value FROM peo_po_mapping WHERE curriculum_id = ?", curriculumID)
	if err != nil {
		log.Println("Error fetching PEO-PO mappings:", err)
		http.Error(w, "Failed to fetch mappings", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var peoIndex, poIndex, value int
		if err := rows.Scan(&peoIndex, &poIndex, &value); err == nil {
			// Convert from 1-based (database) to 0-based (frontend)
			key := fmt.Sprintf("%d-%d", peoIndex-1, poIndex-1)
			matrix[key] = value
		}
	}

	response := models.PEOPOMappingResponse{
		Matrix: matrix,
	}

	json.NewEncoder(w).Encode(response)
}

// SavePEOPOMapping handles POST /curriculum/:id/peo-po-mapping
func SavePEOPOMapping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	curriculumID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid curriculum ID", http.StatusBadRequest)
		return
	}

	var request models.PEOPOMappingRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Start transaction
	tx, err := db.DB.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		http.Error(w, "Failed to save mappings", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Delete existing mappings for this curriculum
	_, err = tx.Exec("DELETE FROM peo_po_mapping WHERE curriculum_id = ?", curriculumID)
	if err != nil {
		log.Println("Error deleting existing mappings:", err)
		http.Error(w, "Failed to save mappings", http.StatusInternalServerError)
		return
	}

	// Insert new mappings
	for _, mapping := range request.Mappings {
		_, err = tx.Exec(`
			INSERT INTO peo_po_mapping (curriculum_id, peo_index, po_index, mapping_value)
			VALUES (?, ?, ?, ?)
		`, curriculumID, mapping.PEOIndex, mapping.POIndex, mapping.MappingValue)
		if err != nil {
			log.Println("Error inserting PEO-PO mapping:", err)
			http.Error(w, "Failed to save mappings", http.StatusInternalServerError)
			return
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		log.Println("Error committing transaction:", err)
		http.Error(w, "Failed to save mappings", http.StatusInternalServerError)
		return
	}

	// Log the activity
	LogCurriculumActivity(curriculumID, "PEO-PO Mapping Saved",
		"Updated PEO-PO mappings for the curriculum", "System")

	json.NewEncoder(w).Encode(map[string]string{"message": "PEO-PO mappings saved successfully"})
}
