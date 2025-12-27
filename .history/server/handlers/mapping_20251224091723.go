package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/db"
	"server/models"
	"strconv"

	"github.com/gorilla/mux"
)

// GetCourseMapping handles GET /course/:courseId/mapping
func GetCourseMapping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	courseID, err := strconv.Atoi(vars["courseId"])
	if err != nil {
		http.Error(w, "Invalid course ID", http.StatusBadRequest)
		return
	}

	// Fetch COs from syllabus
	var outcomesJSON []byte
	err = db.DB.QueryRow("SELECT outcomes FROM course_syllabus WHERE course_id = ?", courseID).Scan(&outcomesJSON)

	var cos []string
	if err == sql.ErrNoRows {
		cos = []string{}
	} else if err != nil {
		log.Println("Error fetching course outcomes:", err)
		http.Error(w, "Failed to fetch course outcomes", http.StatusInternalServerError)
		return
	} else {
		if err := json.Unmarshal(outcomesJSON, &cos); err != nil {
			cos = []string{}
		}
	}

	// Fetch existing CO-PO mappings
	coPoMatrix := make(map[string]int)
	rows, err := db.DB.Query("SELECT co_index, po_index, mapping_value FROM co_po_mapping WHERE course_id = ?", courseID)
	if err != nil {
		log.Println("Error fetching CO-PO mappings:", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var coIndex, poIndex, value int
			if err := rows.Scan(&coIndex, &poIndex, &value); err == nil {
				key := fmt.Sprintf("%d-%d", coIndex, poIndex)
				coPoMatrix[key] = value
			}
		}
	}

	// Fetch existing CO-PSO mappings
	coPsoMatrix := make(map[string]int)
	rows, err = db.DB.Query("SELECT co_index, pso_index, mapping_value FROM co_pso_mapping WHERE course_id = ?", courseID)
	if err != nil {
		log.Println("Error fetching CO-PSO mappings:", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var coIndex, psoIndex, value int
			if err := rows.Scan(&coIndex, &psoIndex, &value); err == nil {
				key := fmt.Sprintf("%d-%d", coIndex, psoIndex)
				coPsoMatrix[key] = value
			}
		}
	}

	response := models.MappingResponse{
		COs:         cos,
		COPOMatrix:  coPoMatrix,
		COPSOMatrix: coPsoMatrix,
	}

	json.NewEncoder(w).Encode(response)
}

// SaveCourseMapping handles POST /course/:courseId/mapping
func SaveCourseMapping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	courseID, err := strconv.Atoi(vars["courseId"])
	if err != nil {
		http.Error(w, "Invalid course ID", http.StatusBadRequest)
		return
	}

	var request models.MappingRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Fetch existing mappings for diff
	oldCOPO := make(map[string]int)
	oldCOPSO := make(map[string]int)
	oldRows, _ := db.DB.Query("SELECT co_index, po_index, mapping_value FROM co_po_mapping WHERE course_id = ?", courseID)
	if oldRows != nil {
		for oldRows.Next() {
			var coIndex, poIndex, value int
			if oldRows.Scan(&coIndex, &poIndex, &value) == nil {
				key := fmt.Sprintf("CO%d-PO%d", coIndex, poIndex)
				oldCOPO[key] = value
			}
		}
		oldRows.Close()
	}
	oldRows2, _ := db.DB.Query("SELECT co_index, pso_index, mapping_value FROM co_pso_mapping WHERE course_id = ?", courseID)
	if oldRows2 != nil {
		for oldRows2.Next() {
			var coIndex, psoIndex, value int
			if oldRows2.Scan(&coIndex, &psoIndex, &value) == nil {
				key := fmt.Sprintf("CO%d-PSO%d", coIndex, psoIndex)
				oldCOPSO[key] = value
			}
		}
		oldRows2.Close()
	}

	// Start transaction
	tx, err := db.DB.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		http.Error(w, "Failed to save mappings", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Delete existing CO-PO mappings
	_, err = tx.Exec("DELETE FROM co_po_mapping WHERE course_id = ?", courseID)
	if err != nil {
		log.Println("Error deleting existing CO-PO mappings:", err)
		http.Error(w, "Failed to save mappings", http.StatusInternalServerError)
		return
	}

	// Insert new CO-PO mappings
	for _, mapping := range request.COPOMatrix {
		_, err = tx.Exec(`
			INSERT INTO co_po_mapping (course_id, co_index, po_index, mapping_value)
			VALUES (?, ?, ?, ?)
		`, courseID, mapping.COIndex, mapping.POIndex, mapping.MappingValue)
		if err != nil {
			log.Println("Error inserting CO-PO mapping:", err)
			http.Error(w, "Failed to save CO-PO mappings", http.StatusInternalServerError)
			return
		}
	}

	// Delete existing CO-PSO mappings
	_, err = tx.Exec("DELETE FROM co_pso_mapping WHERE course_id = ?", courseID)
	if err != nil {
		log.Println("Error deleting existing CO-PSO mappings:", err)
		http.Error(w, "Failed to save mappings", http.StatusInternalServerError)
		return
	}

	// Insert new CO-PSO mappings
	for _, mapping := range request.COPSOMatrix {
		_, err = tx.Exec(`
			INSERT INTO co_pso_mapping (course_id, co_index, pso_index, mapping_value)
			VALUES (?, ?, ?, ?)
		`, courseID, mapping.COIndex, mapping.PSOIndex, mapping.MappingValue)
		if err != nil {
			log.Println("Error inserting CO-PSO mapping:", err)
			http.Error(w, "Failed to save CO-PSO mappings", http.StatusInternalServerError)
			return
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		log.Println("Error committing transaction:", err)
		http.Error(w, "Failed to save mappings", http.StatusInternalServerError)
		return
	}

	// Get curriculum ID and course name for logging
	var curriculumID int
	var courseName string
	db.DB.QueryRow(`
		SELECT rc.regulation_id, c.course_name 
		FROM curriculum_courses rc 
		JOIN courses c ON rc.course_id = c.course_id 
		WHERE rc.course_id = ? LIMIT 1
	`, courseID).Scan(&curriculumID, &courseName)

	if curriculumID > 0 {
		// Generate diff
		newCOPO := make(map[string]int)
		newCOPSO := make(map[string]int)
		for _, mapping := range request.COPOMatrix {
			key := fmt.Sprintf("CO%d-PO%d", mapping.COIndex, mapping.POIndex)
			newCOPO[key] = mapping.MappingValue
		}
		for _, mapping := range request.COPSOMatrix {
			key := fmt.Sprintf("CO%d-PSO%d", mapping.COIndex, mapping.PSOIndex)
			newCOPSO[key] = mapping.MappingValue
		}

		diff := make(map[string]interface{})
		diff["co_po_mappings"] = map[string]interface{}{"old": oldCOPO, "new": newCOPO}
		diff["co_pso_mappings"] = map[string]interface{}{"old": oldCOPSO, "new": newCOPSO}

		LogCurriculumActivityWithDiff(curriculumID, "CO-PO/PSO Mapping Saved",
			"Updated CO-PO and CO-PSO mappings for course: "+courseName, "System", diff)
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Mappings saved successfully"})
}
