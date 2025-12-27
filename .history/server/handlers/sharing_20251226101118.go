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
	}

	var reqData struct {
		ItemType          string `json:"item_type"` // "mission", "peos", "pos", "psos", "semester", "course"
		ItemID            int    `json:"item_id"`
		Visibility        string `json:"visibility"`                   // "UNIQUE" or "CLUSTER"
		TargetDepartments []int  `json:"target_departments,omitempty"` // Optional: specific departments to share with
		SharingMode       string `json:"sharing_mode,omitempty"`       // "replace" (default), "add", or "remove"
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
	if reqData.ItemType == "semester" {
		// Default to "replace" mode if not specified
		mode := reqData.SharingMode
		if mode == "" {
			mode = "replace"
		}
		if err := updateSemesterVisibility(reqData.ItemID, reqData.Visibility, reqData.TargetDepartments, mode); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update semester visibility"})
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Semester visibility updated successfully"})
		return
	}

	if reqData.ItemType == "course" {
		// Default to "replace" mode if not specified
		mode := reqData.SharingMode
		if mode == "" {
			mode = "replace"
		}
		if err := updateCourseVisibility(reqData.ItemID, reqData.Visibility, reqData.TargetDepartments, mode); err != nil {
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
// mode can be: "replace" (default) - replace sharing list, "add" - add to existing list, "remove" - remove from list
func updateSemesterVisibility(semesterID int, visibility string, targetDepartments []int, mode string) error {
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
		// Handle different sharing modes
		switch mode {
		case "add":
			// Add new departments to existing sharing list
			if err := addDepartmentsToSemesterSharing(deptID, regulationID, semesterID, semesterNum, targetDepartments); err != nil {
				log.Printf("Error adding departments to semester sharing: %v\n", err)
				return err
			}
			// Don't update visibility - it should already be CLUSTER
			return nil
		case "remove":
			// Remove departments from sharing list
			if err := removeDepartmentsFromSemesterSharing(deptID, semesterID, targetDepartments); err != nil {
				log.Printf("Error removing departments from semester sharing: %v\n", err)
				return err
			}
			// Check if there are any remaining shared departments
			var shareCount int
			db.DB.QueryRow(`
				SELECT COUNT(*) FROM sharing_tracking 
				WHERE source_dept_id = ? AND source_item_id = ? AND item_type = 'semester'
			`, deptID, semesterID).Scan(&shareCount)

			// If no more shares, update to UNIQUE; otherwise keep as CLUSTER
			if shareCount == 0 {
				_, err = db.DB.Exec("UPDATE semesters SET visibility = 'UNIQUE' WHERE id = ?", semesterID)
				return err
			}
			return nil
		default: // "replace" or empty
			// Replace entire sharing list (original behavior)
			if err := shareSemesterToCluster(deptID, regulationID, semesterID, semesterNum, targetDepartments); err != nil {
				log.Printf("Error sharing semester: %v\n", err)
				return err
			}
		}
	} else {
		// Unshare semester - remove copies
		if err := unshareSemesterFromCluster(deptID, semesterID); err != nil {
			log.Printf("Error unsharing semester: %v\n", err)
			return err
		}
	}

	// Update semester visibility (only for replace mode or unshare)
	_, err = db.DB.Exec("UPDATE semesters SET visibility = ? WHERE id = ?", visibility, semesterID)
	return err
}

// updateCourseVisibility updates the visibility of a course and replicates/removes data
// mode can be: "replace" (default) - replace sharing list, "add" - add to existing list, "remove" - remove from list
func updateCourseVisibility(courseID int, visibility string, targetDepartments []int, mode string) error {
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
		// Handle different sharing modes
		switch mode {
		case "add":
			// Add new departments to existing sharing list
			if err := addDepartmentsToCourseSharing(deptID, regulationID, courseID, targetDepartments); err != nil {
				log.Printf("Error adding departments to course sharing: %v\n", err)
				return err
			}
			// Don't update visibility - it should already be CLUSTER
			return nil
		case "remove":
			// Remove departments from sharing list
			if err := removeDepartmentsFromCourseSharing(deptID, courseID, targetDepartments); err != nil {
				log.Printf("Error removing departments from course sharing: %v\n", err)
				return err
			}
			// Check if there are any remaining shared departments
			var shareCount int
			db.DB.QueryRow(`
				SELECT COUNT(*) FROM sharing_tracking 
				WHERE source_dept_id = ? AND source_item_id = ? AND item_type = 'course'
			`, deptID, courseID).Scan(&shareCount)

			// If no more shares, update to UNIQUE; otherwise keep as CLUSTER
			if shareCount == 0 {
				_, err = db.DB.Exec("UPDATE courses SET visibility = 'UNIQUE' WHERE course_id = ?", courseID)
				return err
			}
			return nil
		default: // "replace" or empty
			// Replace entire sharing list (original behavior)
			if err := shareCourseToCluster(deptID, regulationID, courseID, targetDepartments); err != nil {
				log.Printf("Error sharing course: %v\n", err)
				return err
			}
		}
	} else {
		// Unshare course - remove copies
		if err := unshareCourseFromCluster(deptID, courseID); err != nil {
			log.Printf("Error unsharing course: %v\n", err)
			return err
		}
	}

	// Update course visibility (only for replace mode or unshare)
	_, err = db.DB.Exec("UPDATE courses SET visibility = ? WHERE course_id = ?", visibility, courseID)
	return err
}

// shareSemesterToCluster copies a semester to selected or all cluster departments
func shareSemesterToCluster(sourceDeptID, sourceRegulationID, semesterID, semesterNum int, targetDepartments []int) error {
	log.Printf("=== Starting shareSemesterToCluster ===\n")
	log.Printf("Source Dept: %d, Source Reg: %d, Semester: %d (Sem #%d)\n",
		sourceDeptID, sourceRegulationID, semesterID, semesterNum)

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
			log.Printf("Creating new semester copy for dept %d, regulation %d\n", targetDeptID, targetRegulationID)
			result, err := db.DB.Exec(`
				INSERT INTO semesters (regulation_id, semester_number, visibility, source_department_id)
				VALUES (?, ?, 'CLUSTER', ?)
			`, targetRegulationID, semesterNum, sourceDeptID)
			if err != nil {
				log.Printf("Error creating semester copy for dept %d: %v\n", targetDeptID, err)
				continue
			}

			copiedSemID, _ := result.LastInsertId()
			log.Printf("Created semester copy with ID: %d\n", copiedSemID)

			// Record in tracking table
			_, _ = db.DB.Exec(`
				INSERT INTO sharing_tracking (source_dept_id, target_dept_id, item_type, source_item_id, copied_item_id)
				VALUES (?, ?, 'semester', ?, ?)
			`, sourceDeptID, targetDeptID, semesterID, copiedSemID)

			// Copy all courses from the source semester to the copied semester
			if err := copyCoursesBetweenSemesters(sourceRegulationID, semesterID, targetRegulationID, int(copiedSemID)); err != nil {
				log.Printf("Error copying courses from semester %d to %d: %v\n", semesterID, copiedSemID, err)
			} else {
				log.Printf("Successfully copied courses from semester %d to %d\n", semesterID, copiedSemID)
			}

			// Copy PEO-PO mappings from source to target regulation
			if err := copyPEOPOMappings(sourceRegulationID, targetRegulationID); err != nil {
				log.Printf("Warning: Failed to copy PEO-PO mappings for regulation %d: %v\n", targetRegulationID, err)
			} else {
				log.Printf("Successfully copied PEO-PO mappings from regulation %d to %d\n", sourceRegulationID, targetRegulationID)
			}
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
	var credit sql.NullInt64
	var courseType sql.NullString
	err = db.DB.QueryRow(`
		SELECT course_code, course_name, credit, course_type 
		FROM courses WHERE course_id = ?
	`, courseID).Scan(&courseCode, &courseName, &credit, &courseType)
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

	// Get semester number - this is critical for validation
	var semesterNum int
	err = db.DB.QueryRow("SELECT semester_number FROM semesters WHERE id = ?", sourceSemesterID).Scan(&semesterNum)
	if err != nil {
		return err
	}

	log.Printf("Sharing course %s from semester %d to cluster departments\n", courseCode, semesterNum)

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

		// CRITICAL VALIDATION: Check if the receiving department has the same semester
		var targetSemesterID int
		err := db.DB.QueryRow(`
			SELECT id FROM semesters 
			WHERE regulation_id = ? AND semester_number = ?
		`, targetRegulationID, semesterNum).Scan(&targetSemesterID)

		if err == sql.ErrNoRows {
			// Receiving department doesn't have this semester - cannot share course
			log.Printf("Cannot share course %s: Receiving dept %d (regulation %d) does not have semester %d\n",
				courseCode, targetDeptID, targetRegulationID, semesterNum)
			continue
		} else if err != nil {
			log.Printf("Error checking semester for dept %d: %v\n", targetDeptID, err)
			continue
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
				INSERT INTO courses (course_code, course_name, credit, course_type, visibility)
				VALUES (?, ?, ?, ?, 'CLUSTER')
			`, courseCode, courseName, credit, courseType)
			if err != nil {
				log.Printf("Error creating course copy: %v\n", err)
				continue
			}
			cID, _ := result.LastInsertId()
			targetCourseID = int(cID)

			// Copy syllabus data for the newly created course
			if err := copySyllabusData(courseID, targetCourseID); err != nil {
				log.Printf("Warning: Failed to copy syllabus for course %s: %v\n", courseCode, err)
			} else {
				log.Printf("Successfully copied syllabus data for course %s (ID: %d -> %d)\n", courseCode, courseID, targetCourseID)
			}
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
// This includes course details, syllabus data, and all related information
func copyCoursesBetweenSemesters(sourceRegID, sourceSemID, targetRegID, targetSemID int) error {
	// Get all courses from source semester with full details
	rows, err := db.DB.Query(`
		SELECT c.course_id, c.course_code, c.course_name, c.course_type,
		       c.category, c.credit, c.theory_hours, c.activity_hours, c.lecture_hours,
		       c.tutorial_hours, c.practical_hours, c.cia_marks, c.see_marks, 
		       c.total_marks, c.total_hours
		FROM courses c
		JOIN curriculum_courses cc ON c.course_id = cc.course_id
		WHERE cc.regulation_id = ? AND cc.semester_id = ?
	`, sourceRegID, sourceSemID)
	if err != nil {
		return err
	}
	defer rows.Close()

	log.Printf("Copying courses from semester %d (reg %d) to semester %d (reg %d)\n",
		sourceSemID, sourceRegID, targetSemID, targetRegID)

	courseCount := 0
	for rows.Next() {
		var courseID int
		var courseCode, courseName string
		var courseType, category sql.NullString
		var credit, theoryHours, activityHours, lectureHours sql.NullInt64
		var tutorialHours, practicalHours, ciaMarks, seeMarks sql.NullInt64
		var totalMarks, totalHours sql.NullInt64

		if err := rows.Scan(&courseID, &courseCode, &courseName, &courseType,
			&category, &credit, &theoryHours, &activityHours, &lectureHours,
			&tutorialHours, &practicalHours, &ciaMarks, &seeMarks,
			&totalMarks, &totalHours); err != nil {
			log.Printf("Error scanning course: %v\n", err)
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
			// Create new course with all details
			result, err := db.DB.Exec(`
				INSERT INTO courses (course_code, course_name, course_type, visibility,
					category, credit, theory_hours, activity_hours, lecture_hours,
					tutorial_hours, practical_hours, cia_marks, see_marks, total_marks, total_hours)
				VALUES (?, ?, ?, 'CLUSTER', ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`, courseCode, courseName, courseType, category, credit,
				theoryHours, activityHours, lectureHours, tutorialHours, practicalHours,
				ciaMarks, seeMarks, totalMarks, totalHours)
			if err != nil {
				log.Printf("Error creating course %s: %v\n", courseCode, err)
				continue
			}
			cID, _ := result.LastInsertId()
			targetCourseID = int(cID)

			// Copy syllabus data for the new course
			if err := copySyllabusData(courseID, targetCourseID); err != nil {
				log.Printf("Warning: Failed to copy syllabus for course %s: %v\n", courseCode, err)
			}
		} else {
			targetCourseID = existingCourseID
			// Update visibility to CLUSTER since it's now shared
			db.DB.Exec("UPDATE courses SET visibility = 'CLUSTER' WHERE course_id = ?", targetCourseID)
		}

		// Link to target semester
		_, err = db.DB.Exec(`
			INSERT IGNORE INTO curriculum_courses (regulation_id, semester_id, course_id)
			VALUES (?, ?, ?)
		`, targetRegID, targetSemID, targetCourseID)

		if err != nil {
			log.Printf("Error linking course %s to curriculum: %v\n", courseCode, err)
		} else {
			courseCount++
			log.Printf("Successfully linked course %s (ID: %d) to semester %d\n", courseCode, targetCourseID, targetSemID)
		}
	}

	log.Printf("Successfully copied %d courses from semester %d to semester %d\n", courseCount, sourceSemID, targetSemID)
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

// copySyllabusData copies all syllabus-related data from source course to target course
// This includes course_syllabus, syllabus_models, syllabus_titles, and syllabus_topics
func copySyllabusData(sourceCourseID, targetCourseID int) error {
	// First, check if source course has syllabus data
	var sourceSyllabusID int
	err := db.DB.QueryRow("SELECT id FROM course_syllabus WHERE course_id = ?", sourceCourseID).Scan(&sourceSyllabusID)

	if err == sql.ErrNoRows {
		// No syllabus data to copy
		log.Printf("No syllabus data found for source course %d\\n", sourceCourseID)
		return nil
	} else if err != nil {
		return fmt.Errorf("error checking source syllabus: %w", err)
	}

	// Get the syllabus header data
	var objectives, outcomes, referenceList, prerequisites sql.NullString
	err = db.DB.QueryRow(`
		SELECT objectives, outcomes, reference_list, prerequisites 
		FROM course_syllabus WHERE id = ?
	`, sourceSyllabusID).Scan(&objectives, &outcomes, &referenceList, &prerequisites)

	if err != nil {
		return fmt.Errorf("error fetching syllabus data: %w", err)
	}

	// Check if target course already has syllabus
	var targetSyllabusID int
	err = db.DB.QueryRow("SELECT id FROM course_syllabus WHERE course_id = ?", targetCourseID).Scan(&targetSyllabusID)

	if err == sql.ErrNoRows {
		// Create new syllabus for target course
		result, err := db.DB.Exec(`
			INSERT INTO course_syllabus (course_id, objectives, outcomes, reference_list, prerequisites)
			VALUES (?, ?, ?, ?, ?)
		`, targetCourseID, objectives, outcomes, referenceList, prerequisites)

		if err != nil {
			return fmt.Errorf("error creating target syllabus: %w", err)
		}

		syllabusID, _ := result.LastInsertId()
		targetSyllabusID = int(syllabusID)

		log.Printf("Created syllabus %d for target course %d\\n", targetSyllabusID, targetCourseID)
	} else if err != nil {
		return fmt.Errorf("error checking target syllabus: %w", err)
	} else {
		// Update existing syllabus
		_, err = db.DB.Exec(`
			UPDATE course_syllabus 
			SET objectives = ?, outcomes = ?, reference_list = ?, prerequisites = ?
			WHERE id = ?
		`, objectives, outcomes, referenceList, prerequisites, targetSyllabusID)

		if err != nil {
			return fmt.Errorf("error updating target syllabus: %w", err)
		}

		log.Printf("Updated existing syllabus %d for target course %d\\n", targetSyllabusID, targetCourseID)
	}

	// Copy syllabus models (modules)
	if err := copySyllabusModels(sourceSyllabusID, targetSyllabusID, sourceCourseID, targetCourseID); err != nil {
		return fmt.Errorf("error copying syllabus models: %w", err)
	}

	return nil
}

// copySyllabusModels copies syllabus models and their nested titles and topics
func copySyllabusModels(sourceSyllabusID, targetSyllabusID, sourceCourseID, targetCourseID int) error {
	// Get all models from source syllabus
	rows, err := db.DB.Query(`
		SELECT id, model_name, name, position
		FROM syllabus_models
		WHERE syllabus_id = ? AND course_id = ?
		ORDER BY position
	`, sourceSyllabusID, sourceCourseID)

	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var sourceModelID, position int
		var modelName, name sql.NullString

		if err := rows.Scan(&sourceModelID, &modelName, &name, &position); err != nil {
			log.Printf("Error scanning model: %v\\n", err)
			continue
		}

		// Create model in target syllabus
		result, err := db.DB.Exec(`
			INSERT INTO syllabus_models (syllabus_id, model_name, name, position, course_id)
			VALUES (?, ?, ?, ?, ?)
		`, targetSyllabusID, modelName, name, position, targetCourseID)

		if err != nil {
			log.Printf("Error creating model: %v\\n", err)
			continue
		}

		targetModelID, _ := result.LastInsertId()

		// Copy titles for this model
		if err := copySyllabusTitles(sourceModelID, int(targetModelID)); err != nil {
			log.Printf("Error copying titles for model %d: %v\\n", sourceModelID, err)
		}
	}

	return nil
}

// copySyllabusTitles copies syllabus titles and their nested topics
func copySyllabusTitles(sourceModelID, targetModelID int) error {
	// Get all titles from source model
	rows, err := db.DB.Query(`
		SELECT id, title_name, title, hours, position
		FROM syllabus_titles
		WHERE model_id = ?
		ORDER BY position
	`, sourceModelID)

	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var sourceTitleID, position int
		var titleName, title sql.NullString
		var hours sql.NullInt64

		if err := rows.Scan(&sourceTitleID, &titleName, &title, &hours, &position); err != nil {
			log.Printf("Error scanning title: %v\\n", err)
			continue
		}

		// Create title in target model
		result, err := db.DB.Exec(`
			INSERT INTO syllabus_titles (model_id, title_name, title, hours, position)
			VALUES (?, ?, ?, ?, ?)
		`, targetModelID, titleName, title, hours, position)

		if err != nil {
			log.Printf("Error creating title: %v\\n", err)
			continue
		}

		targetTitleID, _ := result.LastInsertId()

		// Copy topics for this title
		if err := copySyllabusTopics(sourceTitleID, int(targetTitleID)); err != nil {
			log.Printf("Error copying topics for title %d: %v\\n", sourceTitleID, err)
		}
	}

	return nil
}

// copySyllabusTopics copies syllabus topics from source to target title
func copySyllabusTopics(sourceTitleID, targetTitleID int) error {
	// Get all topics from source title
	rows, err := db.DB.Query(`
		SELECT topic, content, position
		FROM syllabus_topics
		WHERE title_id = ?
		ORDER BY position
	`, sourceTitleID)

	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var topic, content sql.NullString
		var position int

		if err := rows.Scan(&topic, &content, &position); err != nil {
			log.Printf("Error scanning topic: %v\\n", err)
			continue
		}

		// Create topic in target title
		_, err := db.DB.Exec(`
			INSERT INTO syllabus_topics (title_id, topic, content, position)
			VALUES (?, ?, ?, ?)
		`, targetTitleID, topic, content, position)

		if err != nil {
			log.Printf("Error creating topic: %v\\n", err)
		}
	}

	return nil
}

// copyPEOPOMappings copies PEO-PO mappings from source to target regulation
func copyPEOPOMappings(sourceRegulationID, targetRegulationID int) error {
	// Get all PEO-PO mappings from source regulation
	rows, err := db.DB.Query(`
		SELECT peo_index, po_index, mapping_value
		FROM peo_po_mapping
		WHERE regulation_id = ?
	`, sourceRegulationID)

	if err != nil {
		return err
	}
	defer rows.Close()

	mappings := []struct {
		peoIndex, poIndex, value int
	}{}

	for rows.Next() {
		var peoIndex, poIndex, value int
		if err := rows.Scan(&peoIndex, &poIndex, &value); err != nil {
			log.Printf("Error scanning PEO-PO mapping: %v\\n", err)
			continue
		}
		mappings = append(mappings, struct{ peoIndex, poIndex, value int }{peoIndex, poIndex, value})
	}

	// Only copy if there are mappings
	if len(mappings) == 0 {
		log.Printf("No PEO-PO mappings found for source regulation %d\\n", sourceRegulationID)
		return nil
	}

	// Delete existing mappings for target regulation
	_, err = db.DB.Exec("DELETE FROM peo_po_mapping WHERE regulation_id = ?", targetRegulationID)
	if err != nil {
		return fmt.Errorf("error deleting existing PEO-PO mappings: %w", err)
	}

	// Insert new mappings
	for _, m := range mappings {
		_, err := db.DB.Exec(`
			INSERT INTO peo_po_mapping (regulation_id, peo_index, po_index, mapping_value)
			VALUES (?, ?, ?, ?)
		`, targetRegulationID, m.peoIndex, m.poIndex, m.value)

		if err != nil {
			log.Printf("Error inserting PEO-PO mapping: %v\\n", err)
		}
	}

	log.Printf("Copied %d PEO-PO mappings from regulation %d to %d\\n", len(mappings), sourceRegulationID, targetRegulationID)
	return nil
}

// addDepartmentsToSemesterSharing adds new departments to an existing semester sharing
func addDepartmentsToSemesterSharing(sourceDeptID, sourceRegulationID, semesterID, semesterNum int, newDepartments []int) error {
	log.Printf("Adding departments %v to semester %d sharing\n", newDepartments, semesterID)

	// Simply call shareSemesterToCluster with the new departments
	// It will skip departments that already have the semester
	return shareSemesterToCluster(sourceDeptID, sourceRegulationID, semesterID, semesterNum, newDepartments)
}

// removeDepartmentsFromSemesterSharing removes specific departments from semester sharing
func removeDepartmentsFromSemesterSharing(sourceDeptID, semesterID int, departmentsToRemove []int) error {
	log.Printf("Removing departments %v from semester %d sharing\n", departmentsToRemove, semesterID)

	// Get all shared semesters for the specified target departments
	for _, targetDeptID := range departmentsToRemove {
		// Get the copied semester ID from tracking table
		var copiedSemID int
		err := db.DB.QueryRow(`
			SELECT copied_item_id 
			FROM sharing_tracking 
			WHERE source_dept_id = ? AND target_dept_id = ? 
			  AND item_type = 'semester' AND source_item_id = ?
		`, sourceDeptID, targetDeptID, semesterID).Scan(&copiedSemID)

		if err == sql.ErrNoRows {
			log.Printf("No shared semester found for dept %d\n", targetDeptID)
			continue
		} else if err != nil {
			log.Printf("Error finding shared semester for dept %d: %v\n", targetDeptID, err)
			continue
		}

		// Get regulation ID for the copied semester
		var targetRegulationID int
		db.DB.QueryRow("SELECT regulation_id FROM semesters WHERE id = ?", copiedSemID).Scan(&targetRegulationID)

		// Get all course IDs for this semester that need to be deleted (before deleting anything)
		rows, err := db.DB.Query(`
			SELECT course_id 
			FROM curriculum_courses 
			WHERE regulation_id = ? AND semester_id = ?
		`, targetRegulationID, copiedSemID)
		
		var courseIDs []int
		if err == nil {
			for rows.Next() {
				var courseID int
				if rows.Scan(&courseID) == nil {
					courseIDs = append(courseIDs, courseID)
				}
			}
			rows.Close()
		}

		// Delete syllabus data and courses for each course
		for _, courseID := range courseIDs {
			log.Printf("Deleting course %d and its syllabus data\n", courseID)
			
			// Get syllabus IDs first
			var syllabusIDs []int
			syllabusRows, _ := db.DB.Query("SELECT id FROM course_syllabus WHERE course_id = ?", courseID)
			for syllabusRows.Next() {
				var syllabusID int
				if syllabusRows.Scan(&syllabusID) == nil {
					syllabusIDs = append(syllabusIDs, syllabusID)
				}
			}
			syllabusRows.Close()

			// For each syllabus, delete its models and their content
			for _, syllabusID := range syllabusIDs {
				var modelIDs []int
				modelRows, _ := db.DB.Query("SELECT id FROM syllabus_models WHERE syllabus_id = ?", syllabusID)
				for modelRows.Next() {
					var modelID int
					if modelRows.Scan(&modelID) == nil {
						modelIDs = append(modelIDs, modelID)
					}
				}
				modelRows.Close()

				// For each model, delete its titles and topics
				for _, modelID := range modelIDs {
					var titleIDs []int
					titleRows, _ := db.DB.Query("SELECT id FROM syllabus_titles WHERE model_id = ?", modelID)
					for titleRows.Next() {
						var titleID int
						if titleRows.Scan(&titleID) == nil {
							titleIDs = append(titleIDs, titleID)
						}
					}
					titleRows.Close()

					// Delete topics for each title
					for _, titleID := range titleIDs {
						db.DB.Exec("DELETE FROM syllabus_topics WHERE title_id = ?", titleID)
					}

					// Delete titles
					db.DB.Exec("DELETE FROM syllabus_titles WHERE model_id = ?", modelID)
				}

				// Delete models
				db.DB.Exec("DELETE FROM syllabus_models WHERE syllabus_id = ?", syllabusID)
			}

			// Delete course syllabus
			db.DB.Exec("DELETE FROM course_syllabus WHERE course_id = ?", courseID)

			// Delete from curriculum_courses
			db.DB.Exec("DELETE FROM curriculum_courses WHERE course_id = ?", courseID)

			// Delete the course itself
			db.DB.Exec("DELETE FROM courses WHERE course_id = ?", courseID)
			
			log.Printf("Deleted course %d and its syllabus data\n", courseID)
		}

		// Delete the semester itself
		_, err = db.DB.Exec("DELETE FROM semesters WHERE id = ? AND source_department_id = ?", copiedSemID, sourceDeptID)
		if err != nil {
			log.Printf("Error removing shared semester %d: %v\n", copiedSemID, err)
		} else {
			log.Printf("Removed shared semester %d from dept %d\n", copiedSemID, targetDeptID)
		}

		// Remove tracking record for semester
		db.DB.Exec(`
			DELETE FROM sharing_tracking 
			WHERE source_dept_id = ? AND target_dept_id = ? 
			  AND item_type = 'semester' AND source_item_id = ?
		`, sourceDeptID, targetDeptID, semesterID)
		
		// Remove tracking records for all courses that were in this semester
		for _, courseID := range courseIDs {
			db.DB.Exec(`
				DELETE FROM sharing_tracking 
				WHERE source_dept_id = ? AND target_dept_id = ? 
				  AND item_type = 'course' AND copied_item_id = ?
			`, sourceDeptID, targetDeptID, courseID)
		}
	}

	return nil
}

// addDepartmentsToCourseSharing adds new departments to an existing course sharing
func addDepartmentsToCourseSharing(sourceDeptID, sourceRegulationID, courseID int, newDepartments []int) error {
	log.Printf("Adding departments %v to course %d sharing\n", newDepartments, courseID)

	// Simply call shareCourseToCluster with the new departments
	// It will handle departments that already have the course
	return shareCourseToCluster(sourceDeptID, sourceRegulationID, courseID, newDepartments)
}

// removeDepartmentsFromCourseSharing removes specific departments from course sharing
func removeDepartmentsFromCourseSharing(sourceDeptID, courseID int, departmentsToRemove []int) error {
	log.Printf("Removing departments %v from course %d sharing\n", departmentsToRemove, courseID)

	// Get all shared courses for the specified target departments
	for _, targetDeptID := range departmentsToRemove {
		// Get the copied course ID from tracking table
		var copiedCourseID int
		err := db.DB.QueryRow(`
			SELECT copied_item_id 
			FROM sharing_tracking 
			WHERE source_dept_id = ? AND target_dept_id = ? 
			  AND item_type = 'course' AND source_item_id = ?
		`, sourceDeptID, targetDeptID, courseID).Scan(&copiedCourseID)

		if err == sql.ErrNoRows {
			log.Printf("No shared course found for dept %d\n", targetDeptID)
			continue
		} else if err != nil {
			log.Printf("Error finding shared course for dept %d: %v\n", targetDeptID, err)
			continue
		}

		// Get target regulation ID
		var targetRegulationID int
		err = db.DB.QueryRow(`
			SELECT do.regulation_id 
			FROM department_overview do
			WHERE do.id = ?
		`, targetDeptID).Scan(&targetRegulationID)
		if err != nil {
			log.Printf("Error getting regulation for dept %d: %v\n", targetDeptID, err)
			continue
		}

		log.Printf("Deleting course %d and its syllabus data\n", copiedCourseID)
		
		// Get syllabus IDs first
		var syllabusIDs []int
		syllabusRows, _ := db.DB.Query("SELECT id FROM course_syllabus WHERE course_id = ?", copiedCourseID)
		for syllabusRows.Next() {
			var syllabusID int
			if syllabusRows.Scan(&syllabusID) == nil {
				syllabusIDs = append(syllabusIDs, syllabusID)
			}
		}
		syllabusRows.Close()

		// For each syllabus, delete its models and their content
		for _, syllabusID := range syllabusIDs {
			var modelIDs []int
			modelRows, _ := db.DB.Query("SELECT id FROM syllabus_models WHERE syllabus_id = ?", syllabusID)
			for modelRows.Next() {
				var modelID int
				if modelRows.Scan(&modelID) == nil {
					modelIDs = append(modelIDs, modelID)
				}
			}
			modelRows.Close()

			// For each model, delete its titles and topics
			for _, modelID := range modelIDs {
				var titleIDs []int
				titleRows, _ := db.DB.Query("SELECT id FROM syllabus_titles WHERE model_id = ?", modelID)
				for titleRows.Next() {
					var titleID int
					if titleRows.Scan(&titleID) == nil {
						titleIDs = append(titleIDs, titleID)
					}
				}
				titleRows.Close()

				// Delete topics for each title
				for _, titleID := range titleIDs {
					db.DB.Exec("DELETE FROM syllabus_topics WHERE title_id = ?", titleID)
				}

				// Delete titles
				db.DB.Exec("DELETE FROM syllabus_titles WHERE model_id = ?", modelID)
			}

			// Delete models
			db.DB.Exec("DELETE FROM syllabus_models WHERE syllabus_id = ?", syllabusID)
		}

		// Delete course syllabus
		db.DB.Exec("DELETE FROM course_syllabus WHERE course_id = ?", copiedCourseID)

		// Remove course from curriculum_courses
		_, err = db.DB.Exec(`
			DELETE FROM curriculum_courses 
			WHERE regulation_id = ? AND course_id = ?
		`, targetRegulationID, copiedCourseID)
		if err != nil {
			log.Printf("Error removing course %d from curriculum: %v\n", copiedCourseID, err)
		} else {
			log.Printf("Removed course %d from dept %d curriculum\n", copiedCourseID, targetDeptID)
		}

		// Delete the course itself
		db.DB.Exec("DELETE FROM courses WHERE course_id = ?", copiedCourseID)
		log.Printf("Deleted course %d and its syllabus data\n", copiedCourseID)

		// Check if the source course is still shared with other departments
		var shareCount int
		db.DB.QueryRow(`
			SELECT COUNT(*) FROM sharing_tracking 
			WHERE source_item_id = ? AND item_type = 'course' AND source_dept_id = ?
		`, courseID, sourceDeptID).Scan(&shareCount)

		// If this was the last sharing instance, update source course visibility back to UNIQUE
		if shareCount == 1 {
			db.DB.Exec("UPDATE courses SET visibility = 'UNIQUE' WHERE course_id = ?", courseID)
			log.Printf("Updated source course %d visibility to UNIQUE (no longer shared)\n", courseID)
		}

		// Remove tracking record
		db.DB.Exec(`
			DELETE FROM sharing_tracking 
			WHERE source_dept_id = ? AND target_dept_id = ? 
			  AND item_type = 'course' AND source_item_id = ?
		`, sourceDeptID, targetDeptID, courseID)
	}

	return nil
}

// GetItemSharedDepartments returns the list of departments that currently have a shared item
func GetItemSharedDepartments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	itemType := vars["item_type"]
	itemID, err := strconv.Atoi(vars["item_id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid item ID"})
		return
	}

	var sourceDeptID int

	// Get source department ID based on item type
	switch itemType {
	case "semester":
		err = db.DB.QueryRow(`
			SELECT do.id 
			FROM semesters s
			JOIN department_overview do ON s.regulation_id = do.regulation_id
			WHERE s.id = ? AND (s.source_department_id IS NULL OR s.source_department_id = do.id)
		`, itemID).Scan(&sourceDeptID)
	case "course":
		err = db.DB.QueryRow(`
			SELECT do.id 
			FROM courses c
			JOIN curriculum_courses cc ON c.course_id = cc.course_id
			JOIN department_overview do ON cc.regulation_id = do.regulation_id
			WHERE c.course_id = ?
			LIMIT 1
		`, itemID).Scan(&sourceDeptID)
	case "mission", "peos", "pos", "psos":
		tableName := map[string]string{
			"mission": "department_mission",
			"peos":    "department_peos",
			"pos":     "department_pos",
			"psos":    "department_psos",
		}[itemType]
		query := fmt.Sprintf("SELECT department_id FROM %s WHERE id = ?", tableName)
		err = db.DB.QueryRow(query, itemID).Scan(&sourceDeptID)
	default:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid item type"})
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Item not found"})
		return
	}

	// Get departments that have this item from sharing_tracking
	rows, err := db.DB.Query(`
		SELECT DISTINCT target_dept_id
		FROM sharing_tracking
		WHERE source_dept_id = ? AND source_item_id = ? AND item_type = ?
	`, sourceDeptID, itemID, itemType)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
		return
	}
	defer rows.Close()

	sharedDeptIDs := []int{}
	for rows.Next() {
		var deptID int
		if err := rows.Scan(&deptID); err == nil {
			sharedDeptIDs = append(sharedDeptIDs, deptID)
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"source_department_id": sourceDeptID,
		"shared_with":          sharedDeptIDs,
	})
}
