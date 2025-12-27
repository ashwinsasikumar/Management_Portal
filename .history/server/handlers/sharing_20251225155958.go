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

	// Fetch semesters with visibility
	response["semesters"] = fetchSemestersWithVisibility(regulationID)

	json.NewEncoder(w).Encode(response)
}

// fetchSemestersWithVisibility fetches semesters with their visibility status and courses
func fetchSemestersWithVisibility(regulationID int) []map[string]interface{} {
	query := "SELECT id, semester_number, COALESCE(visibility, 'UNIQUE') as visibility FROM semesters WHERE regulation_id = ? ORDER BY semester_number"
	rows, err := db.DB.Query(query, regulationID)
	if err != nil {
		log.Println("Error fetching semesters:", err)
		return []map[string]interface{}{}
	}
	defer rows.Close()

	semesters := []map[string]interface{}{}
	for rows.Next() {
		var id, semNum int
		var visibility string
		if err := rows.Scan(&id, &semNum, &visibility); err == nil {
			semester := map[string]interface{}{
				"id":              id,
				"semester_number": semNum,
				"visibility":      visibility,
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
			})
		}
	}
	return courses
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
	if reqData.ItemType == "semester" {
		if err := updateSemesterVisibility(reqData.ItemID, reqData.Visibility); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update semester visibility"})
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Semester visibility updated successfully"})
		return
	}

	if reqData.ItemType == "course" {
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
	query := fmt.Sprintf("SELECT department_id, visibility FROM %s WHERE id = ?", tableInfo.Table)
	err := db.DB.QueryRow(query, reqData.ItemID).Scan(&deptID, &currentVisibility)
	if err != nil {
		log.Println("Error fetching item:", err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Item not found"})
		return
	}

	// If changing to CLUSTER, copy to other departments
	if reqData.Visibility == "CLUSTER" {
		if err := shareItemToCluster(deptID, reqData.ItemID, tableInfo.Table, tableInfo.Column); err != nil {
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

// shareItemToCluster copies an item to all other departments in the same cluster
func shareItemToCluster(sourceDeptID, itemID int, tableName, columnName string) error {
	// Get cluster ID for this department
	var clusterID sql.NullInt64
	err := db.DB.QueryRow(`
		SELECT cluster_id FROM cluster_departments WHERE department_id = ?
	`, sourceDeptID).Scan(&clusterID)
	if err != nil || !clusterID.Valid {
		return fmt.Errorf("department not in cluster")
	}

	// Get all other departments in the cluster
	rows, err := db.DB.Query(`
		SELECT department_id FROM cluster_departments 
		WHERE cluster_id = ? AND department_id != ?
	`, clusterID.Int64, sourceDeptID)
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
			// Already exists, skip
			continue
		}
		
		// Insert the shared item
		_, err = db.DB.Exec(insertQuery, targetDeptID, text, position, sourceDeptID)
		if err != nil {
			log.Printf("Error copying item to dept %d: %v\n", targetDeptID, err)
		}
	}

	return nil
}

// unshareItemFromCluster removes shared copies from other departments
func unshareItemFromCluster(sourceDeptID, itemID int, tableName string) error {
	// Delete all items in other departments that have this as source
	deleteQuery := fmt.Sprintf(`
		DELETE FROM %s 
		WHERE source_department_id = ? 
		AND department_id IN (
			SELECT cd2.department_id 
			FROM cluster_departments cd1 
			JOIN cluster_departments cd2 ON cd1.cluster_id = cd2.cluster_id 
			WHERE cd1.department_id = ? AND cd2.department_id != ?
		)
	`, tableName)
	
	_, err := db.DB.Exec(deleteQuery, sourceDeptID, sourceDeptID, sourceDeptID)
	return err
}

// updateSemesterVisibility updates the visibility of a semester and optionally all its courses
func updateSemesterVisibility(semesterID int, visibility string) error {
	// Update semester visibility
	_, err := db.DB.Exec("UPDATE semesters SET visibility = ? WHERE id = ?", visibility, semesterID)
	if err != nil {
		return err
	}
	
	// If sharing semester, also share all courses in that semester
	if visibility == "CLUSTER" {
		// Get regulation_id for this semester
		var regulationID int
		err = db.DB.QueryRow("SELECT regulation_id FROM semesters WHERE id = ?", semesterID).Scan(&regulationID)
		if err != nil {
			return err
		}
		
		// Update all courses in this semester to CLUSTER visibility
		_, err = db.DB.Exec(`
			UPDATE courses c
			JOIN curriculum_courses cc ON c.course_id = cc.course_id
			SET c.visibility = ?
			WHERE cc.semester_id = ? AND cc.regulation_id = ?
		`, visibility, semesterID, regulationID)
	}
	
	return err
}

// updateCourseVisibility updates the visibility of a course and all its related data
func updateCourseVisibility(courseID int, visibility string) error {
	_, err := db.DB.Exec("UPDATE courses SET visibility = ? WHERE course_id = ?", visibility, courseID)
	return err
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
