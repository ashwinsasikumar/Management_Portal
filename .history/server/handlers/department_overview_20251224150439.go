package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"server/db"
	"server/models"

	"github.com/gorilla/mux"
)

// GetDepartmentOverview retrieves department overview data for a regulation
func GetDepartmentOverview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	regulationID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid regulation ID"})
		return
	}

	query := `SELECT id, regulation_id, vision FROM department_overview WHERE regulation_id = ?`

	var overview models.DepartmentOverview

	err = db.DB.QueryRow(query, regulationID).Scan(
		&overview.ID,
		&overview.RegulationID,
		&overview.Vision,
	)

	if err == sql.ErrNoRows {
		// Return empty structure with default values
		overview = models.DepartmentOverview{
			RegulationID: regulationID,
			Vision:       "",
			Mission:      []string{},
			PEOs:         []string{},
			POs:          []string{},
			PSOs:         []string{},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(overview)
		return
	} else if err != nil {
		log.Println("Error querying department overview:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch department overview"})
		return
	}

	departmentID := overview.ID

	// Fetch mission items ordered by position
	overview.Mission = fetchDepartmentList(departmentID, "department_mission", "mission_text")

	// Fetch PEOs ordered by position
	overview.PEOs = fetchDepartmentList(departmentID, "department_peos", "peo_text")

	// Fetch POs ordered by position
	overview.POs = fetchDepartmentList(departmentID, "department_pos", "po_text")

	// Fetch PSOs ordered by position
	overview.PSOs = fetchDepartmentList(departmentID, "department_psos", "pso_text")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(overview)
}

// Helper function to fetch list items from normalized tables
func fetchDepartmentList(departmentID int, tableName, columnName string) []string {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE department_id = ? ORDER BY position", columnName, tableName)
	rows, err := db.DB.Query(query, departmentID)
	if err != nil {
		return []string{}
	}
	defer rows.Close()

	items := []string{}
	for rows.Next() {
		var text string
		if err := rows.Scan(&text); err == nil {
			items = append(items, text)
		}
	}
	return items
}

// SaveDepartmentOverview creates or updates department overview data
func SaveDepartmentOverview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	vars := mux.Vars(r)
	regulationID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid regulation ID"})
		return
	}

	var overview models.DepartmentOverview
	err = json.NewDecoder(r.Body).Decode(&overview)
	if err != nil {
		log.Println("Error decoding request body:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	overview.RegulationID = regulationID

	// Convert arrays to JSON strings
	missionJSON, _ := json.Marshal(overview.Mission)
	peosJSON, _ := json.Marshal(overview.PEOs)
	posJSON, _ := json.Marshal(overview.POs)
	psosJSON, _ := json.Marshal(overview.PSOs)

	// Fetch existing data for diff
	var oldOverview models.DepartmentOverview
	var oldMissionJSON, oldPeosJSON, oldPosJSON, oldPsosJSON string
	fetchQuery := "SELECT id, vision, mission, peos, pos, psos FROM department_overview WHERE regulation_id = ?"
	fetchErr := db.DB.QueryRow(fetchQuery, regulationID).Scan(&oldOverview.ID, &oldOverview.Vision,
		&oldMissionJSON, &oldPeosJSON, &oldPosJSON, &oldPsosJSON)

	hasExisting := fetchErr != sql.ErrNoRows
	if hasExisting {
		json.Unmarshal([]byte(oldMissionJSON), &oldOverview.Mission)
		json.Unmarshal([]byte(oldPeosJSON), &oldOverview.PEOs)
		json.Unmarshal([]byte(oldPosJSON), &oldOverview.POs)
		json.Unmarshal([]byte(oldPsosJSON), &oldOverview.PSOs)
	}

	// Check if record exists
	var existingID int
	checkQuery := "SELECT id FROM department_overview WHERE regulation_id = ?"
	err = db.DB.QueryRow(checkQuery, regulationID).Scan(&existingID)

	if err == sql.ErrNoRows {
		// INSERT new record
		insertQuery := `INSERT INTO department_overview (regulation_id, vision, mission, peos, pos, psos) 
		                VALUES (?, ?, ?, ?, ?, ?)`
		result, err := db.DB.Exec(insertQuery, regulationID, overview.Vision, missionJSON, peosJSON, posJSON, psosJSON)
		if err != nil {
			log.Println("Error inserting department overview:", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to save department overview"})
			return
		}
		id, _ := result.LastInsertId()
		overview.ID = int(id)
	} else if err != nil {
		log.Println("Error checking existing record:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to save department overview"})
		return
	} else {
		// UPDATE existing record
		updateQuery := `UPDATE department_overview 
		                SET vision = ?, mission = ?, peos = ?, pos = ?, psos = ? 
		                WHERE regulation_id = ?`
		_, err := db.DB.Exec(updateQuery, overview.Vision, missionJSON, peosJSON, posJSON, psosJSON, regulationID)
		if err != nil {
			log.Println("Error updating department overview:", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update department overview"})
			return
		}
		overview.ID = existingID
	}

	// Generate granular diff and log the activity
	if hasExisting {
		// Vision change
		if oldOverview.Vision != overview.Vision {
			diff := map[string]map[string]interface{}{
				"vision": {"old": oldOverview.Vision, "new": overview.Vision},
			}
			LogCurriculumActivityWithDiff(regulationID, "Vision Updated",
				"Updated department vision", "System", diff)
		}

		// Mission changes (per index)
		detectArrayChanges(regulationID, "Mission", oldOverview.Mission, overview.Mission)

		// PEO changes (per index)
		detectArrayChanges(regulationID, "PEO", oldOverview.PEOs, overview.PEOs)

		// PO changes (per index)
		detectArrayChanges(regulationID, "PO", oldOverview.POs, overview.POs)

		// PSO changes (per index)
		detectArrayChanges(regulationID, "PSO", oldOverview.PSOs, overview.PSOs)
	} else {
		LogCurriculumActivity(regulationID, "Department Overview Created",
			"Created department vision, mission, PEOs, POs, and PSOs", "System")
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(overview)
}

// Helper function to compare string arrays
func stringArraysEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// detectArrayChanges compares two string arrays and logs individual item changes
func detectArrayChanges(regulationID int, label string, oldArray, newArray []string) {
	maxLen := len(oldArray)
	if len(newArray) > maxLen {
		maxLen = len(newArray)
	}

	for i := 0; i < maxLen; i++ {
		var oldVal, newVal string
		if i < len(oldArray) {
			oldVal = oldArray[i]
		}
		if i < len(newArray) {
			newVal = newArray[i]
		}

		// Detect changes
		if oldVal != newVal {
			diff := map[string]map[string]interface{}{
				fmt.Sprintf("%s[%d]", label, i): {"old": oldVal, "new": newVal},
			}

			var action, description string
			if oldVal == "" {
				action = fmt.Sprintf("%s[%d] Added", label, i)
				description = fmt.Sprintf("Added %s item at index %d", label, i)
			} else if newVal == "" {
				action = fmt.Sprintf("%s[%d] Deleted", label, i)
				description = fmt.Sprintf("Deleted %s item at index %d", label, i)
			} else {
				action = fmt.Sprintf("%s[%d] Updated", label, i)
				description = fmt.Sprintf("Updated %s item at index %d", label, i)
			}

			LogCurriculumActivityWithDiff(regulationID, action, description, "System", diff)
		}
	}
}
