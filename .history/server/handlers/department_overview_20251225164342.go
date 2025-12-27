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

// Helper function to fetch list items from normalized tables
// Note: Shared items are now physically replicated in the database with source_department_id,
// so we only need to fetch items for this department (includes both owned and shared items)
func fetchDepartmentList(departmentID int, tableName, columnName string) []models.DepartmentListItem {
	query := fmt.Sprintf("SELECT id, %s, visibility, source_department_id FROM %s WHERE department_id = ? ORDER BY position", columnName, tableName)
	rows, err := db.DB.Query(query, departmentID)
	if err != nil {
		return []models.DepartmentListItem{}
	}
	defer rows.Close()

	items := []models.DepartmentListItem{}
	for rows.Next() {
		var item models.DepartmentListItem
		var sourceDeptID sql.NullInt64
		if err := rows.Scan(&item.ID, &item.Text, &item.Visibility, &sourceDeptID); err == nil {
			if sourceDeptID.Valid {
				item.SourceDepartmentID = int(sourceDeptID.Int64)
			}
			items = append(items, item)
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

// Helper function to save list items to normalized tables with proper update/insert/delete logic
func saveDepartmentList(departmentID int, tableName, columnName string, items []models.DepartmentListItem) error {
	// Fetch existing items from database
	existingItems := make(map[int]struct {
		text               string
		visibility         string
		sourceDepartmentID sql.NullInt64
	})
	
	query := fmt.Sprintf("SELECT id, %s, visibility, source_department_id FROM %s WHERE department_id = ?", columnName, tableName)
	rows, err := db.DB.Query(query, departmentID)
	if err != nil {
		log.Printf("Error fetching existing items: %v\n", err)
		return err
	}
	
	for rows.Next() {
		var id int
		var text, visibility string
		var sourceDeptID sql.NullInt64
		if err := rows.Scan(&id, &text, &visibility, &sourceDeptID); err == nil {
			existingItems[id] = struct {
				text               string
				visibility         string
				sourceDepartmentID sql.NullInt64
			}{text, visibility, sourceDeptID}
		}
	}
	rows.Close()

	// Track which IDs are still present
	presentIDs := make(map[int]bool)

	// Process each item in the input
	for i, item := range items {
		if item.Text == "" {
			continue
		}

		if item.ID > 0 {
			// Existing item - check if it needs updating
			presentIDs[item.ID] = true
			existing, exists := existingItems[item.ID]
			
			if !exists {
				log.Printf("Warning: Item ID %d not found in database\n", item.ID)
				continue
			}

			// IMPORTANT: Only block modification if item is from ANOTHER department
			// Items owned by THIS department (source_department_id IS NULL or equals departmentID) 
			// can ALWAYS be modified, even if they have CLUSTER visibility
			if existing.sourceDepartmentID.Valid && existing.sourceDepartmentID.Int64 != int64(departmentID) {
				// This is a shared item from another department - read-only, skip it
				log.Printf("Skipping modification of shared item ID %d from department %d (read-only)\n", 
					item.ID, existing.sourceDepartmentID.Int64)
				continue
			}

			// Item is owned by this department - allow modification
			// Check if text changed
			textChanged := existing.text != item.Text
			
			if textChanged && existing.visibility == "CLUSTER" {
				// Text was modified on a shared item - automatically unshare it
				log.Printf("Item %d was modified and is currently shared - unsharing from cluster and making UNIQUE\n", item.ID)
				
				// Unshare from cluster (removes copies from other departments)
				if err := unshareItemFromCluster(departmentID, item.ID, tableName); err != nil {
					log.Printf("Error unsharing modified item: %v\n", err)
				}
				
				// Change visibility to UNIQUE (no longer shared)
				item.Visibility = "UNIQUE"
			}

			// Update the item
			visibility := item.Visibility
			if visibility == "" {
				visibility = "UNIQUE"
			}
			
			updateQuery := fmt.Sprintf("UPDATE %s SET %s = ?, visibility = ?, position = ? WHERE id = ?", tableName, columnName)
			_, err := db.DB.Exec(updateQuery, item.Text, visibility, i, item.ID)
			if err != nil {
				log.Printf("Error updating item %d: %v\n", item.ID, err)
				return err
			}
		} else {
			// New item - insert it
			visibility := item.Visibility
			if visibility == "" {
				visibility = "UNIQUE"
			}
			
			insertQuery := fmt.Sprintf(`
				INSERT INTO %s (department_id, %s, visibility, position, source_department_id) 
				VALUES (?, ?, ?, ?, NULL)
			`, tableName, columnName)
			
			_, err := db.DB.Exec(insertQuery, departmentID, item.Text, visibility, i)
			if err != nil {
				log.Printf("Error inserting new item: %v\n", err)
				return err
			}
		}
	}

	// Delete items that are no longer present (and are owned by this department)
	for id, existing := range existingItems {
		if !presentIDs[id] {
			// Check if this item is owned by this department (not shared from elsewhere)
			if !existing.sourceDepartmentID.Valid || existing.sourceDepartmentID.Int64 == int64(departmentID) {
				// If it was shared (CLUSTER), unshare it first
				if existing.visibility == "CLUSTER" {
					log.Printf("Deleting item %d which was shared - unsharing first\n", id)
					if err := unshareItemFromCluster(departmentID, id, tableName); err != nil {
						log.Printf("Error unsharing deleted item: %v\n", err)
					}
				}
				
				// Delete the item
				deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE id = ?", tableName)
				_, err := db.DB.Exec(deleteQuery, id)
				if err != nil {
					log.Printf("Error deleting item %d: %v\n", id, err)
				}
			} else {
				// This is a shared item from another department - keep it
				log.Printf("Keeping shared item %d from department %d\n", id, existing.sourceDepartmentID.Int64)
			}
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
