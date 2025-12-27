package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"server/db"

	"github.com/gorilla/mux"
)

// GetDepartmentSharingInfo gets all sharable content for a department with visibility status
func GetDepartmentSharingInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	regulationID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid regulation ID"})
		return
	}

	// Get or create department_overview entry
	var deptOverviewID int
	err = db.DB.QueryRow("SELECT id FROM department_overview WHERE regulation_id = ?", regulationID).Scan(&deptOverviewID)
	if err == sql.ErrNoRows {
		// Create if doesn't exist
		result, err := db.DB.Exec("INSERT INTO department_overview (regulation_id, vision) VALUES (?, '')", regulationID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create department overview"})
			return
		}
		id, _ := result.LastInsertId()
		deptOverviewID = int(id)
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
		return
	}

	// Check if department is in a cluster
	var clusterID sql.NullInt64
	var clusterName sql.NullString
	db.DB.QueryRow(`
		SELECT c.id, c.name 
		FROM cluster_departments cd 
		JOIN clusters c ON cd.cluster_id = c.id 
		WHERE cd.department_id = ?
	`, deptOverviewID).Scan(&clusterID, &clusterName)

	// Get all items with their visibility
	response := map[string]interface{}{
		"department_id": deptOverviewID,
		"regulation_id": regulationID,
		"in_cluster":    clusterID.Valid,
	}

	if clusterID.Valid {
		response["cluster_id"] = clusterID.Int64
		response["cluster_name"] = clusterName.String
	}

	// Fetch mission items
	response["mission"] = fetchItemsWithVisibility(deptOverviewID, "department_mission", "mission_text")
	// Fetch PEOs
	response["peos"] = fetchItemsWithVisibility(deptOverviewID, "department_peos", "peo_text")
	// Fetch POs
	response["pos"] = fetchItemsWithVisibility(deptOverviewID, "department_pos", "po_text")
	// Fetch PSOs
	response["psos"] = fetchItemsWithVisibility(deptOverviewID, "department_psos", "pso_text")

	json.NewEncoder(w).Encode(response)
}

// Helper to fetch items with visibility
func fetchItemsWithVisibility(deptID int, tableName, columnName string) []map[string]interface{} {
	query := fmt.Sprintf("SELECT id, %s, visibility, position FROM %s WHERE department_id = ? ORDER BY position", columnName, tableName)
	rows, err := db.DB.Query(query, deptID)
	if err != nil {
		return []map[string]interface{}{}
	}
	defer rows.Close()

	items := []map[string]interface{}{}
	for rows.Next() {
		var id, position int
		var text, visibility string
		if err := rows.Scan(&id, &text, &visibility, &position); err == nil {
			items = append(items, map[string]interface{}{
				"id":         id,
				"text":       text,
				"visibility": visibility,
				"position":   position,
			})
		}
	}
	return items
}

// UpdateItemVisibility updates the visibility of a specific item
func UpdateItemVisibility(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	var reqData struct {
		ItemType   string `json:"item_type"` // "mission", "peos", "pos", "psos"
		ItemID     int    `json:"item_id"`
		Visibility string `json:"visibility"` // "UNIQUE" or "CLUSTER"
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Validate visibility
	if reqData.Visibility != "UNIQUE" && reqData.Visibility != "CLUSTER" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid visibility value"})
		return
	}

	// Map item type to table name
	tableMap := map[string]string{
		"mission": "department_mission",
		"peos":    "department_peos",
		"pos":     "department_pos",
		"psos":    "department_psos",
	}

	tableName, ok := tableMap[reqData.ItemType]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid item type"})
		return
	}

	// Update visibility
	query := fmt.Sprintf("UPDATE %s SET visibility = ? WHERE id = ?", tableName)
	result, err := db.DB.Exec(query, reqData.Visibility, reqData.ItemID)
	if err != nil {
		log.Println("Error updating visibility:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update visibility"})
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Item not found"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message":    "Visibility updated successfully",
		"visibility": reqData.Visibility,
	})
}

// GetClusterSharedContent gets all shared content visible to a cluster
func GetClusterSharedContent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	clusterID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid cluster ID"})
		return
	}

	// Get all departments in cluster
	deptQuery := `
		SELECT cd.department_id, do.regulation_id, c.name
		FROM cluster_departments cd
		JOIN department_overview do ON cd.department_id = do.id
		JOIN curriculum c ON do.regulation_id = c.id
		WHERE cd.cluster_id = ?
	`
	rows, err := db.DB.Query(deptQuery, clusterID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch cluster departments"})
		return
	}
	defer rows.Close()

	departments := []map[string]interface{}{}
	for rows.Next() {
		var deptID, regID int
		var name string
		if err := rows.Scan(&deptID, &regID, &name); err == nil {
			dept := map[string]interface{}{
				"department_id": deptID,
				"regulation_id": regID,
				"name":          name,
				"mission":       fetchSharedItems(deptID, "department_mission", "mission_text"),
				"peos":          fetchSharedItems(deptID, "department_peos", "peo_text"),
				"pos":           fetchSharedItems(deptID, "department_pos", "po_text"),
				"psos":          fetchSharedItems(deptID, "department_psos", "pso_text"),
			}
			departments = append(departments, dept)
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"cluster_id":  clusterID,
		"departments": departments,
	})
}

// Helper to fetch only shared (CLUSTER visibility) items
func fetchSharedItems(deptID int, tableName, columnName string) []map[string]interface{} {
	query := fmt.Sprintf("SELECT id, %s FROM %s WHERE department_id = ? AND visibility = 'CLUSTER' ORDER BY position", columnName, tableName)
	rows, err := db.DB.Query(query, deptID)
	if err != nil {
		return []map[string]interface{}{}
	}
	defer rows.Close()

	items := []map[string]interface{}{}
	for rows.Next() {
		var id int
		var text string
		if err := rows.Scan(&id, &text); err == nil {
			items = append(items, map[string]interface{}{
				"id":   id,
				"text": text,
			})
		}
	}
	return items
}
