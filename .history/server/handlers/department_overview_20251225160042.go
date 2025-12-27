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
			Mission:      []models.DepartmentListItem{},
			PEOs:         []models.DepartmentListItem{},
			POs:          []models.DepartmentListItem{},
			PSOs:         []models.DepartmentListItem{},
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

// Helper function to fetch list items from normalized tables with cluster propagation
func fetchDepartmentList(departmentID int, tableName, columnName string) []models.DepartmentListItem {
	// Get cluster ID for this department (if any)
	var clusterID sql.NullInt64
	db.DB.QueryRow("SELECT cluster_id FROM cluster_departments WHERE department_id = ?", departmentID).Scan(&clusterID)

	// Fetch department-specific items
	query := fmt.Sprintf("SELECT %s, visibility FROM %s WHERE department_id = ? ORDER BY position", columnName, tableName)
	rows, err := db.DB.Query(query, departmentID)
	if err != nil {
		return []models.DepartmentListItem{}
	}
	defer rows.Close()

	items := []models.DepartmentListItem{}
	for rows.Next() {
		var item models.DepartmentListItem
		if err := rows.Scan(&item.Text, &item.Visibility); err == nil {
			items = append(items, item)
		}
	}

	// If department is in a cluster, fetch CLUSTER-visibility items from other departments in same cluster
	if clusterID.Valid {
		clusterQuery := fmt.Sprintf(`
			SELECT %s, visibility 
			FROM %s 
			WHERE department_id IN (
				SELECT department_id FROM cluster_departments 
				WHERE cluster_id = ? AND department_id != ?
			) AND visibility = 'CLUSTER'
			ORDER BY position
		`, columnName, tableName)
		
		clusterRows, err := db.DB.Query(clusterQuery, clusterID.Int64, departmentID)
		if err == nil {
			defer clusterRows.Close()
			for clusterRows.Next() {
				var item models.DepartmentListItem
				if err := clusterRows.Scan(&item.Text, &item.Visibility); err == nil {
					items = append(items, item)
				}
			}
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

	// Fetch existing data for diff
	var oldOverview models.DepartmentOverview
	var departmentID int
	fetchQuery := "SELECT id, vision FROM department_overview WHERE regulation_id = ?"
	fetchErr := db.DB.QueryRow(fetchQuery, regulationID).Scan(&departmentID, &oldOverview.Vision)

	hasExisting := fetchErr != sql.ErrNoRows
	if hasExisting {
		oldOverview.ID = departmentID
		oldOverview.Mission = fetchDepartmentList(departmentID, "department_mission", "mission_text")
		oldOverview.PEOs = fetchDepartmentList(departmentID, "department_peos", "peo_text")
		oldOverview.POs = fetchDepartmentList(departmentID, "department_pos", "po_text")
		oldOverview.PSOs = fetchDepartmentList(departmentID, "department_psos", "pso_text")
	}

	// Check if record exists
	var existingID int
	checkQuery := "SELECT id FROM department_overview WHERE regulation_id = ?"
	err = db.DB.QueryRow(checkQuery, regulationID).Scan(&existingID)

	if err == sql.ErrNoRows {
		// INSERT new record
		insertQuery := `INSERT INTO department_overview (regulation_id, vision) VALUES (?, ?)`
		result, err := db.DB.Exec(insertQuery, regulationID, overview.Vision)
		if err != nil {
			log.Println("Error inserting department overview:", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to save department overview"})
			return
		}
		id, _ := result.LastInsertId()
		overview.ID = int(id)
		departmentID = int(id)
	} else if err != nil {
		log.Println("Error checking existing record:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to save department overview"})
		return
	} else {
		// UPDATE existing record
		updateQuery := `UPDATE department_overview SET vision = ? WHERE regulation_id = ?`
		_, err := db.DB.Exec(updateQuery, overview.Vision, regulationID)
		if err != nil {
			log.Println("Error updating department overview:", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update department overview"})
			return
		}
		overview.ID = existingID
		departmentID = existingID
	}

	// Save list items to normalized tables
	saveDepartmentList(departmentID, "department_mission", "mission_text", overview.Mission)
	saveDepartmentList(departmentID, "department_peos", "peo_text", overview.PEOs)
	saveDepartmentList(departmentID, "department_pos", "po_text", overview.POs)
	saveDepartmentList(departmentID, "department_psos", "pso_text", overview.PSOs)

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

// Helper function to save list items to normalized tables with visibility
func saveDepartmentList(departmentID int, tableName, columnName string, items []models.DepartmentListItem) error {
	// Delete existing items for this department (only those created by this department, not shared ones)
	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE department_id = ? AND (source_department_id IS NULL OR source_department_id = ?)", tableName)
	_, err := db.DB.Exec(deleteQuery, departmentID, departmentID)
	if err != nil {
		log.Printf("Error deleting items from %s: %v\n", tableName, err)
		return err
	}

	// Insert new items
	insertQuery := fmt.Sprintf(`
		INSERT INTO %s (department_id, %s, visibility, position, source_department_id) 
		VALUES (?, ?, ?, ?, NULL)
	`, tableName, columnName)

	for i, item := range items {
		// If item text is modified and it was shared, it becomes UNIQUE
		visibility := item.Visibility
		if visibility == "" {
			visibility = "UNIQUE"
		}
		
		if item.Text == "" {
			continue
		}
		
		_, err := db.DB.Exec(insertQuery, departmentID, item.Text, visibility, i)
		if err != nil {
			log.Printf("Error inserting into %s: %v\n", tableName, err)
			return err
		}
	}

	return nil
}

// detectArrayChanges logs granular changes for array-type fields {
			continue
		}
		visibility := item.Visibility
		if visibility == "" {
			visibility = "UNIQUE" // Default to UNIQUE
		}
		_, err := db.DB.Exec(insertQuery, departmentID, item.Text, visibility, i)
		if err != nil {
			return err
		}
	}
	return nil
}

// Helper function to compare DepartmentListItem arrays
func listItemArraysEqual(a, b []models.DepartmentListItem) bool {
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

// detectArrayChanges compares two DepartmentListItem arrays and logs individual item changes
func detectArrayChanges(regulationID int, label string, oldArray, newArray []models.DepartmentListItem) {
	maxLen := len(oldArray)
	if len(newArray) > maxLen {
		maxLen = len(newArray)
	}

	for i := 0; i < maxLen; i++ {
		var oldVal, newVal string
		if i < len(oldArray) {
			oldVal = oldArray[i].Text
		}
		if i < len(newArray) {
			newVal = newArray[i].Text
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
