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
		
		// Get all departments in this cluster for selective sharing
		response["cluster_departments"] = fetchClusterDepartments(int(clusterID.Int64), deptOverviewID)
	}

	// Fetch mission items
	response["mission"] = fetchItemsWithVisibility(deptOverviewID, "department_mission", "mission_text")
	// Fetch PEOs
	response["peos"] = fetchItemsWithVisibility(deptOverviewID, "department_peos", "peo_text")
	// Fetch POs
	response["pos"] = fetchItemsWithVisibility(deptOverviewID, "department_pos", "po_text")
	// Fetch PSOs
	response["psos"] = fetchItemsWithVisibility(deptOverviewID, "department_psos", "pso_text")

	// Fetch semesters with visibility
	response["semesters"] = fetchSemestersWithVisibility(regulationID)

	json.NewEncoder(w).Encode(response)
}

// fetchSemestersWithVisibility fetches semesters with their visibility status and courses
func fetchSemestersWithVisibility(regulationID int) []map[string]interface{} {
	query := "SELECT id, semester_number, COALESCE(visibility, 'UNIQUE') as visibility, source_department_id FROM semesters WHERE regulation_id = ? ORDER BY semester_number"
	rows, err := db.DB.Query(query, regulationID)
	if err != nil {
		log.Println("Error fetching semesters:", err)
		return []map[string]interface{}{}
	}
	defer rows.Close()

	// Get department ID for this regulation
	var deptID int
	db.DB.QueryRow("SELECT id FROM department_overview WHERE regulation_id = ?", regulationID).Scan(&deptID)

	semesters := []map[string]interface{}{}
	for rows.Next() {
		var id, semNum int
		var visibility string
		var sourceDeptID sql.NullInt64
		if err := rows.Scan(&id, &semNum, &visibility, &sourceDeptID); err == nil {
			semester := map[string]interface{}{
				"id":              id,
				"semester_number": semNum,
				"visibility":      visibility,
				"is_owner":        !sourceDeptID.Valid || sourceDeptID.Int64 == int64(deptID),
			}
			if sourceDeptID.Valid && sourceDeptID.Int64 != int64(deptID) {
				semester["source_department_id"] = sourceDeptID.Int64
			}
			// Fetch courses for this semester
			semester["courses"] = fetchCoursesForSemester(regulationID, id)
			semesters = append(semesters, semester)
		}
	}
	return semesters
}

// fetchCoursesForSemester fetches courses for a specific semester with visibility
func fetchCoursesForSemester(regulationID, semesterID int) []map[string]interface{} {
	// Note: Courses don't have source_department_id, they're shared globally
	// We'll determine ownership by checking if the semester is owned
	query := `
		SELECT c.course_id, c.course_code, c.course_name, COALESCE(c.visibility, 'UNIQUE') as visibility
		FROM courses c
		JOIN curriculum_courses cc ON c.course_id = cc.course_id
		WHERE cc.regulation_id = ? AND cc.semester_id = ?
		ORDER BY c.course_code
	`
	rows, err := db.DB.Query(query, regulationID, semesterID)
	if err != nil {
		log.Println("Error fetching courses:", err)
		return []map[string]interface{}{}
	}
	defer rows.Close()

	// Check if the semester is owned by this department
	var semesterSourceDeptID sql.NullInt64
	var deptID int
	db.DB.QueryRow("SELECT id FROM department_overview WHERE regulation_id = ?", regulationID).Scan(&deptID)
	db.DB.QueryRow("SELECT source_department_id FROM semesters WHERE id = ?", semesterID).Scan(&semesterSourceDeptID)

	semesterIsOwned := !semesterSourceDeptID.Valid || semesterSourceDeptID.Int64 == int64(deptID)

	courses := []map[string]interface{}{}
	for rows.Next() {
		var courseID int
		var courseCode, courseName, visibility string
		if err := rows.Scan(&courseID, &courseCode, &courseName, &visibility); err == nil {
			courses = append(courses, map[string]interface{}{
				"id":          courseID,
				"course_code": courseCode,
				"course_name": courseName,
				"visibility":  visibility,
				"is_owner":    semesterIsOwned, // Course ownership follows semester ownership
			})
		}
	}
	return courses
}

// Helper to fetch items with visibility
func fetchItemsWithVisibility(deptID int, tableName, columnName string) []map[string]interface{} {
	query := fmt.Sprintf("SELECT id, %s, visibility, position, source_department_id FROM %s WHERE department_id = ? ORDER BY position", columnName, tableName)
	rows, err := db.DB.Query(query, deptID)
	if err != nil {
		return []map[string]interface{}{}
	}
	defer rows.Close()

	items := []map[string]interface{}{}
	for rows.Next() {
		var id, position int
		var text, visibility string
		var sourceDeptID sql.NullInt64
		if err := rows.Scan(&id, &text, &visibility, &position, &sourceDeptID); err == nil {
			item := map[string]interface{}{
				"id":         id,
				"text":       text,
				"visibility": visibility,
				"position":   position,
				"is_owner":   !sourceDeptID.Valid || sourceDeptID.Int64 == int64(deptID), // Owner if no source or source is self
			}
			if sourceDeptID.Valid && sourceDeptID.Int64 != int64(deptID) {
				item["source_department_id"] = sourceDeptID.Int64
			}
			items = append(items, item)
		}
	}
	return items
}

// fetchClusterDepartments gets all departments in a cluster except the current one
func fetchClusterDepartments(clusterID, currentDeptID int) []map[string]interface{} {
	query := `
		SELECT cd.department_id, do.regulation_id, c.name
		FROM cluster_departments cd
		JOIN department_overview do ON cd.department_id = do.id
		JOIN curriculum c ON do.regulation_id = c.id
		WHERE cd.cluster_id = ? AND cd.department_id != ?
	`
	rows, err := db.DB.Query(query, clusterID, currentDeptID)
	if err != nil {
		return []map[string]interface{}{}
	}
	defer rows.Close()

	depts := []map[string]interface{}{}
	for rows.Next() {
		var deptID, regID int
		var name string
		if err := rows.Scan(&deptID, &regID, &name); err == nil {
			depts = append(depts, map[string]interface{}{
				"department_id": deptID,
				"regulation_id": regID,
				"name":          name,
			})
		}
	}
	return depts
}

// UpdateItemVisibility updates the visibility of a specific item
func UpdateItemVisibility(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}       string `json:"item_type"`   // "mission", "peos", "pos", "psos", "semester", "course"
		ItemID            int    `json:"item_id"`
		Visibility        string `json:"visibility"`  // "UNIQUE" or "CLUSTER"
		TargetDepartments []int  `json:"target_departments,omitempty"` // Optional: specific departments to share with
		ItemType   string `json:"item_type"` // "mission", "peos", "pos", "psos", "semester", "course"
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

	// Handle semester and course separately
	if reqData.ItemType == "semester" {, reqData.TargetDepartments); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update semester visibility"})
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Semester visibility updated successfully"})
		return
	}

	if reqData.ItemType == "course" {
		if err := updateCourseVisibility(reqData.ItemID, reqData.Visibility, reqData.TargetDepartments
		if err := updateCourseVisibility(reqData.ItemID, reqData.Visibility); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update course visibility"})
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Course visibility updated successfully"})
		return
	}

	// Map item type to table name and column name
	tableMap := map[string]struct {
		Table  string
		Column string
	}{
		"mission": {"department_mission", "mission_text"},
		"peos":    {"department_peos", "peo_text"},
		"pos":     {"department_pos", "po_text"},
		"psos":    {"department_psos", "pso_text"},
	}

	tableInfo, ok := tableMap[reqData.ItemType]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid item type"})
		return
	}

	// Get item details and department info
	var deptID int
	var currentVisibility string
	var sourceDeptID sql.NullInt64
	query := fmt.Sprintf("SELECT department_id, visibility, source_department_id FROM %s WHERE id = ?", tableInfo.Table)
	err := db.DB.QueryRow(query, reqData.ItemID).Scan(&deptID, &currentVisibility, &sourceDeptID)
	if err != nil {
		log.Println("Error fetching item:", err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Item not found"})
		return
	}

	// Check if this item is owned by this department
	isOwned := !sourceDeptID.Valid || sourceDeptID.Int64 == int64(deptID)
	if !isOwned {
		// This is a received item from another department - cannot toggle visibility
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "Cannot change visibility of received items. Only the source department can modify sharing."})
		return
	}

	// If changing to CLUSTER, copy to other departments
	if reqData.Visibility == "CLUSTER" {
		if err := shareItemToCluster(deptID, reqData.ItemID, tableInfo.Table, tableInfo.Column, reqData.TargetDepartments); err != nil {
			log.Println("Error sharing item:", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to share item"})
			return
		}
	} else {
		// If changing to UNIQUE, remove from other departments
		if err := unshareItemFromCluster(deptID, reqData.ItemID, tableInfo.Table); err != nil {
			log.Println("Error unsharing item:", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to unshare item"})
			return
		}
	}

	// Update visibility in source item
	updateQuery := fmt.Sprintf("UPDATE %s SET visibility = ? WHERE id = ?", tableInfo.Table)
	_, err = db.DB.Exec(updateQuery, reqData.Visibility, reqData.ItemID)
	if err != nil {
		log.Println("Error updating visibility:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update visibility"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message":    "Visibility updated successfully",
		"visibility": reqData.Visibility,
	})
}

// shareItemToCluster copies an item to selected or all other departments in the same cluster
func shareItemToCluster(sourceDeptID, itemID int, tableName, columnName string, targetDepartments []int) error {
	// Get cluster ID for this department
	var clusterID sql.NullInt64
	err := db.DB.QueryRow(`
		SELECT cluster_id FROM cluster_departments WHERE department_id = ?
	`, sourceDeptID).Scan(&clusterID)
	if err != nil || !clusterID.Valid {
		return fmt.Errorf("department not in cluster")
	}

	// Get departments to share with
	var targetDeptQuery string
	var queryArgs []interface{}
	
	if len(targetDepartments) > 0 {
		// Selective sharing - only to specified departments
		placeholders := ""
		for i := range targetDepartments {
			if i > 0 {
				placeholders += ","
			}
			placeholders += "?"
			queryArgs = append(queryArgs, targetDepartments[i])
		}
		targetDeptQuery = fmt.Sprintf(`
			SELECT department_id FROM cluster_departments 
			WHERE cluster_id = ? AND department_id != ? AND department_id IN (%s)
		`, placeholders)
		queryArgs = append([]interface{}{clusterID.Int64, sourceDeptID}, queryArgs...)
	} else {
		// Share with all departments in cluster
		targetDeptQuery = `
			SELECT department_id FROM cluster_departments 
			WHERE cluster_id = ? AND department_id != ?
		`
		queryArgs = []interface{}{clusterID.Int64, sourceDeptID}
	}

	rows, err := db.DB.Query(targetDeptQuery, queryArgs...)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Get the item data
	var text string
	var position int
	getQuery := fmt.Sprintf("SELECT %s, position FROM %s WHERE id = ?", columnName, tableName)
	err = db.DB.QueryRow(getQuery, itemID).Scan(&text, &position)
	if err != nil {
		return err
	}

	// Determine item type based on table name
	itemType := ""
	switch tableName {
	case "department_mission":
		itemType = "mission"
	case "department_peos":
		itemType = "peos"
	case "department_pos":
		itemType = "pos"
	case "department_psos":
		itemType = "psos"
	}

	// Copy to each department
	insertQuery := fmt.Sprintf(`
		INSERT INTO %s (department_id, %s, visibility, position, source_department_id)
		VALUES (?, ?, 'CLUSTER', ?, ?)
	`, tableName, columnName)

	for rows.Next() {
		var targetDeptID int
		if err := rows.Scan(&targetDeptID); err != nil {
			continue
		}

		// Check if already exists (avoid duplicates)
		var existsID int
		checkQuery := fmt.Sprintf(`
			SELECT id FROM %s 
			WHERE department_id = ? AND source_department_id = ? AND %s = ?
		`, tableName, columnName)
		err := db.DB.QueryRow(checkQuery, targetDeptID, sourceDeptID, text).Scan(&existsID)
		if err == nil {
			// Already exists, update tracking table
			_, _ = db.DB.Exec(`
				INSERT INTO sharing_tracking (source_dept_id, target_dept_id, item_type, source_item_id, copied_item_id)
				VALUES (?, ?, ?, ?, ?)
				ON DUPLICATE KEY UPDATE copied_item_id = VALUES(copied_item_id)
			`, sourceDeptID, targetDeptID, itemType, itemID, existsID)
			continue
		}

		// Insert the shared item
		result, err := db.DB.Exec(insertQuery, targetDeptID, text, position, sourceDeptID)
		if err != nil {
			log.Printf("Error copying item to dept %d: %v\n", targetDeptID, err)
			continue
		}

		// Get the new item ID
		copiedItemID, _ := result.LastInsertId()

		// Record in tracking table
		_, err = db.DB.Exec(`
			INSERT INTO sharing_tracking (source_dept_id, target_dept_id, item_type, source_item_id, copied_item_id)
			VALUES (?, ?, ?, ?, ?)
		`, sourceDeptID, targetDeptID, itemType, itemID, copiedItemID)
		if err != nil {
			log.Printf("Error recording sharing tracking: %v\n", err)
		}
	}

	return nil
}

// unshareItemFromCluster removes shared copies from other departments
func unshareItemFromCluster(sourceDeptID, itemID int, tableName string) error {
	// Determine item type based on table name
	itemType := ""
	switch tableName {
	case "department_mission":
		itemType = "mission"
	case "department_peos":
		itemType = "peos"
	case "department_pos":
		itemType = "pos"
	case "department_psos":
		itemType = "psos"
	}

	// Get all copied items from the tracking table
	rows, err := db.DB.Query(`
		SELECT target_dept_id, copied_item_id 
		FROM sharing_tracking 
		WHERE source_dept_id = ? AND item_type = ? AND source_item_id = ?
	`, sourceDeptID, itemType, itemID)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Delete each copied item
	deleteQuery := fmt.Sprintf(`DELETE FROM %s WHERE id = ? AND department_id = ?`, tableName)

	var deletedCount int
	for rows.Next() {
		var targetDeptID, copiedItemID int
		if err := rows.Scan(&targetDeptID, &copiedItemID); err != nil {
			continue
		}

		// Delete the copied item
		_, err = db.DB.Exec(deleteQuery, copiedItemID, targetDeptID)
		if err != nil {
			log.Printf("Error removing shared item %d from dept %d: %v\n", copiedItemID, targetDeptID, err)
		} else {
			deletedCount++
		}
	}

	// Remove tracking records
	_, err = db.DB.Exec(`
		DELETE FROM sharing_tracking 
		WHERE source_dept_id = ? AND item_type = ? AND source_item_id = ?
	`, sourceDeptID, itemType, itemID)
	if err != nil {
		log.Printf("Error removing tracking records: %v\n", err)
	}

	log.Printf("Unshared item %d from %d departments\n", itemID, deletedCount)
	return nil
}

// updateSharedItemInCluster updates the text of a shared item in all cluster departments
func updateSharedItemInCluster(sourceDeptID, itemID int, tableName, columnName, newText string) error {
	// Determine item type based on table name
	itemType := ""
	switch tableName {
	case "department_mission":
		itemType = "mission"
	case "department_peos":
		itemType = "peos"
	case "department_pos":
		itemType = "pos"
	case "department_psos":
		itemType = "psos"
	}

	// Get all copied items from the tracking table
	rows, err := db.DB.Query(`
		SELECT target_dept_id, copied_item_id 
		FROM sharing_tracking 
		WHERE source_dept_id = ? AND item_type = ? AND source_item_id = ?
	`, sourceDeptID, itemType, itemID)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Update each copied item with the new text
	updateQuery := fmt.Sprintf(`UPDATE %s SET %s = ? WHERE id = ?`, tableName, columnName)

	var updatedCount int
	for rows.Next() {
		var targetDeptID, copiedItemID int
		if err := rows.Scan(&targetDeptID, &copiedItemID); err != nil {
			continue
		}

		// Update the copied item with the new text
		_, err = db.DB.Exec(updateQuery, newText, copiedItemID)
		if err != nil {
			log.Printf("Error updating shared item %d in dept %d: %v\n", copiedItemID, targetDeptID, err)
		} else {
			updatedCount++
		}
	}

	log.Printf("Updated item %d text in %d shared departments\n", itemID, updatedCount)
	return nil
}

// updateSemesterVisibility updates the visibility of a semester and replicates/removes data
func updateSemesterVisibility(semesterID int, visibility string, targetDepartments []int) error {
	// Get semester and regulation info
	var regulationID, semesterNum int
	var sourceDeptID sql.NullInt64
	err := db.DB.QueryRow("SELECT regulation_id, semester_number, source_department_id FROM semesters WHERE id = ?", semesterID).Scan(&regulationID, &semesterNum, &sourceDeptID)
	if err != nil {
		return err
	}

	// Get department_id from regulation
	var deptID int
	err = db.DB.QueryRow("SELECT id FROM department_overview WHERE regulation_id = ?", regulationID).Scan(&deptID)
	if err != nil {
		return err
	}

	// Check if this semester is owned by this department
	isOwned := !sourceDeptID.Valid || sourceDeptID.Int64 == int64(deptID)
	if !isOwned {
		return fmt.Errorf("cannot change visibility of received semester")
	}

	if visibility == "CLUSTER" {
		// Share semester - copy to selected or all cluster departments
		if err := shareSemesterToCluster(deptID, regulationID, semesterID, semesterNum, targetDepartments); err != nil {
			log.Printf("Error sharing semester: %v\n", err)
			return err
		}
	} else {
		// Unshare semester - remove copies
		if err := unshareSemesterFromCluster(deptID, semesterID); err != nil {
			log.Printf("Error unsharing semester: %v\n", err)
			return err
		}
	}

	// Update semester visibility
	_, err = db.DB.Exec("UPDATE semesters SET visibility = ? WHERE id = ?", visibility, semesterID)
	return err
}

// updateCourseVisibility updates the visibility of a course and replicates/removes data
func updateCourseVisibility(courseID int, visibility string, targetDepartments []int) error {
	// Get course info and check its semester ownership
	var courseCode, courseName string
	err := db.DB.QueryRow("SELECT course_code, course_name FROM courses WHERE course_id = ?", courseID).Scan(&courseCode, &courseName)
	if err != nil {
		return err
	}

	// Get regulation_id and semester_id from curriculum_courses
	var regulationID, semesterID int
	err = db.DB.QueryRow("SELECT regulation_id, semester_id FROM curriculum_courses WHERE course_id = ? LIMIT 1", courseID).Scan(&regulationID, &semesterID)
	if err != nil {
		return err
	}

	// Get department_id
	var deptID int
	err = db.DB.QueryRow("SELECT id FROM department_overview WHERE regulation_id = ?", regulationID).Scan(&deptID)
	if err != nil {
		return err
	}

	// Check if the semester (and thus course) is owned by this department
	var semesterSourceDeptID sql.NullInt64
	err = db.DB.QueryRow("SELECT source_department_id FROM semesters WHERE id = ?", semesterID).Scan(&semesterSourceDeptID)
	if err != nil {
		return err
	}

	isOwned := !semesterSourceDeptID.Valid || semesterSourceDeptID.Int64 == int64(deptID)
	if !isOwned {
		return fmt.Errorf("cannot change visibility of received course")
	}

	if visibility == "CLUSTER" {
		// Share course - copy to selected or all cluster departments
		if err := shareCourseToCluster(deptID, regulationID, courseID, targetDepartments); err != nil {
			log.Printf("Error sharing course: %v\n", err)
			return err
		}
	} else {
		// Unshare course - remove copies
		if err := unshareCourseFromCluster(deptID, courseID); err != nil {
			log.Printf("Error unsharing course: %v\n", err)
			return err
		}
	}

	// Update course visibility
	_, err = db.DB.Exec("UPDATE courses SET visibility = ? WHERE course_id = ?", visibility, courseID)
	return err
}

// shareSemesterToCluster copies a semester to selected or all cluster departments
func shareSemesterToCluster(sourceDeptID, sourceRegulationID, semesterID, semesterNum int, targetDepartments []int) error {
	// Get cluster ID for this department
	var clusterID sql.NullInt64
	err := db.DB.QueryRow(`
		SELECT cluster_id FROM cluster_departments WHERE department_id = ?
	`, sourceDeptID).Scan(&clusterID)
	if err != nil || !clusterID.Valid {
		return fmt.Errorf("department not in cluster")
	}

	// Get departments to share with
	var targetDeptQuery string
	var queryArgs []interface{}
	
	if len(targetDepartments) > 0 {
		// Selective sharing
		placeholders := ""
		for i := range targetDepartments {
			if i > 0 {
				placeholders += ","
			}
			placeholders += "?"
			queryArgs = append(queryArgs, targetDepartments[i])
		}
		targetDeptQuery = fmt.Sprintf(`
			SELECT cd.department_id, do.regulation_id 
			FROM cluster_departments cd
			JOIN department_overview do ON cd.department_id = do.id
			WHERE cd.cluster_id = ? AND cd.department_id != ? AND cd.department_id IN (%s)
		`, placeholders)
		queryArgs = append([]interface{}{clusterID.Int64, sourceDeptID}, queryArgs...)
	} else {
		// Share with all departments
		targetDeptQuery = `
			SELECT cd.department_id, do.regulation_id 
			FROM cluster_departments cd
			JOIN department_overview do ON cd.department_id = do.id
			WHERE cd.cluster_id = ? AND cd.department_id != ?
		`
		queryArgs = []interface{}{clusterID.Int64, sourceDeptID}
	}

	rows, err := db.DB.Query(targetDeptQuery, queryArgs...)
	if err != nil {
		return err
	}
	defer rows.Close()

	// For each target department
	for rows.Next() {
		var targetDeptID, targetRegulationID int
		if err := rows.Scan(&targetDeptID, &targetRegulationID); err != nil {
			continue
		}

		// Check if semester already exists in target regulation
		var existingSemID int
		err := db.DB.QueryRow(`
			SELECT id FROM semesters 
			WHERE regulation_id = ? AND semester_number = ? AND source_department_id = ?
		`, targetRegulationID, semesterNum, sourceDeptID).Scan(&existingSemID)

		if err == sql.ErrNoRows {
			// Create the semester copy
			result, err := db.DB.Exec(`
				INSERT INTO semesters (regulation_id, semester_number, visibility, source_department_id)
				VALUES (?, ?, 'CLUSTER', ?)
			`, targetRegulationID, semesterNum, sourceDeptID)
			if err != nil {
				log.Printf("Error creating semester copy for dept %d: %v\n", targetDeptID, err)
				continue
			}

			copiedSemID, _ := result.LastInsertId()

			// Record in tracking table
			_, _ = db.DB.Exec(`
				INSERT INTO sharing_tracking (source_dept_id, target_dept_id, item_type, source_item_id, copied_item_id)
				VALUES (?, ?, 'semester', ?, ?)
			`, sourceDeptID, targetDeptID, semesterID, copiedSemID)

			// Copy all courses from the source semester to the copied semester
			copyCoursesBetweenSemesters(sourceRegulationID, semesterID, targetRegulationID, int(copiedSemID))
		} else if err == nil {
			// Already exists, update tracking
			_, _ = db.DB.Exec(`
				INSERT INTO sharing_tracking (source_dept_id, target_dept_id, item_type, source_item_id, copied_item_id)
				VALUES (?, ?, 'semester', ?, ?)
				ON DUPLICATE KEY UPDATE copied_item_id = VALUES(copied_item_id)
			`, sourceDeptID, targetDeptID, semesterID, existingSemID)
		}
	}

	return nil
}

// unshareSemesterFromCluster removes semester copies from other departments
func unshareSemesterFromCluster(sourceDeptID, semesterID int) error {
	// Get all copied semesters from tracking table
	rows, err := db.DB.Query(`
		SELECT target_dept_id, copied_item_id 
		FROM sharing_tracking 
		WHERE source_dept_id = ? AND item_type = 'semester' AND source_item_id = ?
	`, sourceDeptID, semesterID)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Delete each copied semester and its courses
	for rows.Next() {
		var targetDeptID, copiedSemID int
		if err := rows.Scan(&targetDeptID, &copiedSemID); err != nil {
			continue
		}

		// Get regulation ID for the copied semester
		var targetRegulationID int
		db.DB.QueryRow("SELECT regulation_id FROM semesters WHERE id = ?", copiedSemID).Scan(&targetRegulationID)

		// Delete courses associated with this copied semester
		db.DB.Exec(`
			DELETE FROM curriculum_courses 
			WHERE regulation_id = ? AND semester_id = ?
		`, targetRegulationID, copiedSemID)

		// Delete the semester itself
		_, err = db.DB.Exec("DELETE FROM semesters WHERE id = ? AND source_department_id = ?", copiedSemID, sourceDeptID)
		if err != nil {
			log.Printf("Error removing shared semester %d: %v\n", copiedSemID, err)
		}
	}

	// Remove tracking records
	_, err = db.DB.Exec(`
		DELETE FROM sharing_tracking 
		WHERE source_dept_id = ? AND item_type = 'semester' AND source_item_id = ?
	`, sourceDeptID, semesterID)

	return err
}

// shareCourseToCluster copies a course to selected or all cluster departments
func shareCourseToCluster(sourceDeptID, sourceRegulationID, courseID int, targetDepartments []int) error {
	// Get cluster ID for this department
	var clusterID sql.NullInt64
	err := db.DB.QueryRow(`
		SELECT cluster_id FROM cluster_departments WHERE department_id = ?
	`, sourceDeptID).Scan(&clusterID)
	if err != nil || !clusterID.Valid {
		return fmt.Errorf("department not in cluster")
	}

	// Get course details
	var courseCode, courseName string
	var credits sql.NullFloat64
	var courseType sql.NullString
	err = db.DB.QueryRow(`
		SELECT course_code, course_name, credits, course_type 
		FROM courses WHERE course_id = ?
	`, courseID).Scan(&courseCode, &courseName, &credits, &courseType)
	if err != nil {
		return err
	}

	// Get semester information for this course
	var sourceSemesterID int
	err = db.DB.QueryRow(`
		SELECT semester_id FROM curriculum_courses 
		WHERE course_id = ? AND regulation_id = ?
	`, courseID, sourceRegulationID).Scan(&sourceSemesterID)
	if err != nil {
		return err
	}

	// Get semester number
	var semesterNum int
	err = db.DB.QueryRow("SELECT semester_number FROM semesters WHERE id = ?", sourceSemesterID).Scan(&semesterNum)
	if err != nil {
		return err
	}

	// Get departments to share with
	var targetDeptQuery string
	var queryArgs []interface{}
	
	if len(targetDepartments) > 0 {
		// Selective sharing
		placeholders := ""
		for i := range targetDepartments {
			if i > 0 {
				placeholders += ","
			}
			placeholders += "?"
			queryArgs = append(queryArgs, targetDepartments[i])
		}
		targetDeptQuery = fmt.Sprintf(`
			SELECT cd.department_id, do.regulation_id 
			FROM cluster_departments cd
			JOIN department_overview do ON cd.department_id = do.id
			WHERE cd.cluster_id = ? AND cd.department_id != ? AND cd.department_id IN (%s)
		`, placeholders)
		queryArgs = append([]interface{}{clusterID.Int64, sourceDeptID}, queryArgs...)
	} else {
		// Share with all departments
		targetDeptQuery = `
			SELECT cd.department_id, do.regulation_id 
			FROM cluster_departments cd
			JOIN department_overview do ON cd.department_id = do.id
			WHERE cd.cluster_id = ? AND cd.department_id != ?
		`
		queryArgs = []interface{}{clusterID.Int64, sourceDeptID}
	}

	rows, err := db.DB.Query(targetDeptQuery, queryArgs...)
	if err != nil {
		return err
	}
	defer rows.Close()

	// For each target department
	for rows.Next() {
		var targetDeptID, targetRegulationID int
		if err := rows.Scan(&targetDeptID, &targetRegulationID); err != nil {
			continue
		}

		// Find or create the corresponding semester in target regulation
		var targetSemesterID int
		err := db.DB.QueryRow(`
			SELECT id FROM semesters 
			WHERE regulation_id = ? AND semester_number = ?
		`, targetRegulationID, semesterNum).Scan(&targetSemesterID)

		if err == sql.ErrNoRows {
			// Create semester if it doesn't exist
			result, err := db.DB.Exec(`
				INSERT INTO semesters (regulation_id, semester_number, visibility)
				VALUES (?, ?, 'UNIQUE')
			`, targetRegulationID, semesterNum)
			if err != nil {
				log.Printf("Error creating semester for dept %d: %v\n", targetDeptID, err)
				continue
			}
			semID, _ := result.LastInsertId()
			targetSemesterID = int(semID)
		}

		// Check if course already exists
		var existingCourseID int
		err = db.DB.QueryRow(`
			SELECT course_id FROM courses 
			WHERE course_code = ? AND course_name = ?
		`, courseCode, courseName).Scan(&existingCourseID)

		var targetCourseID int
		if err == sql.ErrNoRows {
			// Create new course
			result, err := db.DB.Exec(`
				INSERT INTO courses (course_code, course_name, credits, course_type, visibility)
				VALUES (?, ?, ?, ?, 'CLUSTER')
			`, courseCode, courseName, credits, courseType)
			if err != nil {
				log.Printf("Error creating course copy: %v\n", err)
				continue
			}
			cID, _ := result.LastInsertId()
			targetCourseID = int(cID)
		} else {
			targetCourseID = existingCourseID
			// Update visibility to CLUSTER
			db.DB.Exec("UPDATE courses SET visibility = 'CLUSTER' WHERE course_id = ?", targetCourseID)
		}

		// Link course to target regulation and semester
		_, err = db.DB.Exec(`
			INSERT IGNORE INTO curriculum_courses (regulation_id, semester_id, course_id)
			VALUES (?, ?, ?)
		`, targetRegulationID, targetSemesterID, targetCourseID)

		// Record in tracking table
		_, _ = db.DB.Exec(`
			INSERT INTO sharing_tracking (source_dept_id, target_dept_id, item_type, source_item_id, copied_item_id)
			VALUES (?, ?, 'course', ?, ?)
			ON DUPLICATE KEY UPDATE copied_item_id = VALUES(copied_item_id)
		`, sourceDeptID, targetDeptID, courseID, targetCourseID)
	}

	return nil
}

// unshareCourseFromCluster removes course copies from other departments
func unshareCourseFromCluster(sourceDeptID, courseID int) error {
	// Get all copied courses from tracking table
	rows, err := db.DB.Query(`
		SELECT target_dept_id, copied_item_id 
		FROM sharing_tracking 
		WHERE source_dept_id = ? AND item_type = 'course' AND source_item_id = ?
	`, sourceDeptID, courseID)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Remove each copied course
	for rows.Next() {
		var targetDeptID, copiedCourseID int
		if err := rows.Scan(&targetDeptID, &copiedCourseID); err != nil {
			continue
		}

		// Get target regulation ID
		var targetRegulationID int
		err := db.DB.QueryRow(`
			SELECT do.regulation_id 
			FROM department_overview do
			WHERE do.id = ?
		`, targetDeptID).Scan(&targetRegulationID)
		if err != nil {
			continue
		}

		// Remove course from curriculum_courses
		_, err = db.DB.Exec(`
			DELETE FROM curriculum_courses 
			WHERE regulation_id = ? AND course_id = ?
		`, targetRegulationID, copiedCourseID)
		if err != nil {
			log.Printf("Error removing course %d from curriculum: %v\n", copiedCourseID, err)
		}

		// Optionally update course visibility back to UNIQUE if not shared elsewhere
		// (Only if this was the only sharing instance)
		var shareCount int
		db.DB.QueryRow(`
			SELECT COUNT(*) FROM sharing_tracking 
			WHERE copied_item_id = ? AND item_type = 'course'
		`, copiedCourseID).Scan(&shareCount)

		if shareCount == 1 {
			db.DB.Exec("UPDATE courses SET visibility = 'UNIQUE' WHERE course_id = ?", copiedCourseID)
		}
	}

	// Remove tracking records
	_, err = db.DB.Exec(`
		DELETE FROM sharing_tracking 
		WHERE source_dept_id = ? AND item_type = 'course' AND source_item_id = ?
	`, sourceDeptID, courseID)

	return err
}

// copyCoursesBetweenSemesters copies all courses from source semester to target semester
func copyCoursesBetweenSemesters(sourceRegID, sourceSemID, targetRegID, targetSemID int) error {
	// Get all courses from source semester
	rows, err := db.DB.Query(`
		SELECT c.course_id, c.course_code, c.course_name, c.credits, c.course_type
		FROM courses c
		JOIN curriculum_courses cc ON c.course_id = cc.course_id
		WHERE cc.regulation_id = ? AND cc.semester_id = ?
	`, sourceRegID, sourceSemID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var courseID int
		var courseCode, courseName string
		var credits sql.NullFloat64
		var courseType sql.NullString

		if err := rows.Scan(&courseID, &courseCode, &courseName, &credits, &courseType); err != nil {
			continue
		}

		// Check if course already exists in system
		var existingCourseID int
		err := db.DB.QueryRow(`
			SELECT course_id FROM courses 
			WHERE course_code = ? AND course_name = ?
		`, courseCode, courseName).Scan(&existingCourseID)

		var targetCourseID int
		if err == sql.ErrNoRows {
			// Create new course
			result, err := db.DB.Exec(`
				INSERT INTO courses (course_code, course_name, credits, course_type, visibility)
				VALUES (?, ?, ?, ?, 'CLUSTER')
			`, courseCode, courseName, credits, courseType)
			if err != nil {
				log.Printf("Error creating course: %v\n", err)
				continue
			}
			cID, _ := result.LastInsertId()
			targetCourseID = int(cID)
		} else {
			targetCourseID = existingCourseID
		}

		// Link to target semester
		_, _ = db.DB.Exec(`
			INSERT IGNORE INTO curriculum_courses (regulation_id, semester_id, course_id)
			VALUES (?, ?, ?)
		`, targetRegID, targetSemID, targetCourseID)
	}

	return nil
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
				"semesters":     fetchSharedSemesters(regID),
			}
			departments = append(departments, dept)
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"cluster_id":  clusterID,
		"departments": departments,
	})
}

// fetchSharedSemesters fetches semesters with CLUSTER visibility and their shared courses
func fetchSharedSemesters(regulationID int) []map[string]interface{} {
	query := "SELECT id, semester_number FROM semesters WHERE regulation_id = ? AND visibility = 'CLUSTER' ORDER BY semester_number"
	rows, err := db.DB.Query(query, regulationID)
	if err != nil {
		return []map[string]interface{}{}
	}
	defer rows.Close()

	semesters := []map[string]interface{}{}
	for rows.Next() {
		var id, semNum int
		if err := rows.Scan(&id, &semNum); err == nil {
			semester := map[string]interface{}{
				"id":              id,
				"semester_number": semNum,
				"courses":         fetchSharedCourses(regulationID, id),
			}
			semesters = append(semesters, semester)
		}
	}
	return semesters
}

// fetchSharedCourses fetches courses with CLUSTER visibility for a semester
func fetchSharedCourses(regulationID, semesterID int) []map[string]interface{} {
	query := `
		SELECT c.course_id, c.course_code, c.course_name
		FROM courses c
		JOIN curriculum_courses cc ON c.course_id = cc.course_id
		WHERE cc.regulation_id = ? AND cc.semester_id = ? AND c.visibility = 'CLUSTER'
		ORDER BY c.course_code
	`
	rows, err := db.DB.Query(query, regulationID, semesterID)
	if err != nil {
		return []map[string]interface{}{}
	}
	defer rows.Close()

	courses := []map[string]interface{}{}
	for rows.Next() {
		var courseID int
		var courseCode, courseName string
		if err := rows.Scan(&courseID, &courseCode, &courseName); err == nil {
			courses = append(courses, map[string]interface{}{
				"id":          courseID,
				"course_code": courseCode,
				"course_name": courseName,
			})
		}
	}
	return courses
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
