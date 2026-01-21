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

// GetClusters retrieves all clusters
func GetClusters(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	query := "SELECT id, name, description, created_at FROM clusters ORDER BY name"
	rows, err := db.DB.Query(query)
	if err != nil {
		log.Println("Error querying clusters:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch clusters"})
		return
	}
	defer rows.Close()

	clusters := []models.Cluster{}
	for rows.Next() {
		var cluster models.Cluster
		if err := rows.Scan(&cluster.ID, &cluster.Name, &cluster.Description, &cluster.CreatedAt); err == nil {
			clusters = append(clusters, cluster)
		}
	}

	json.NewEncoder(w).Encode(clusters)
}

// GetAvailableDepartments retrieves departments that are not in any cluster
func GetAvailableDepartments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// First, ensure all curriculums have a curriculum_vision entry
	ensureQuery := `
		INSERT INTO curriculum_vision (curriculum_id, vision)
		SELECT c.id, ''
		FROM curriculum c
		LEFT JOIN curriculum_vision do ON c.id = do.curriculum_id
		WHERE do.id IS NULL
	`
	result, err := db.DB.Exec(ensureQuery)
	if err != nil {
		log.Printf("Error ensuring curriculum_vision entries: %v\n", err)
	} else {
		affected, _ := result.RowsAffected()
		log.Printf("Created %d new curriculum_vision entries\n", affected)
	}

	// Debug: Check all curriculums
	var totalCurriculums int
	db.DB.QueryRow("SELECT COUNT(*) FROM curriculum").Scan(&totalCurriculums)
	log.Printf("Total curriculums in database: %d\n", totalCurriculums)

	// Debug: Check all curriculum_visions
	var totalDeptOverviews int
	db.DB.QueryRow("SELECT COUNT(*) FROM curriculum_vision").Scan(&totalDeptOverviews)
	log.Printf("Total curriculum_visions: %d\n", totalDeptOverviews)

	// Debug: Check curriculums without curriculum_vision
	var curriculumsWithoutDO int
	db.DB.QueryRow(`
		SELECT COUNT(*) FROM curriculum c 
		LEFT JOIN curriculum_vision do ON c.id = do.curriculum_id 
		WHERE do.id IS NULL
	`).Scan(&curriculumsWithoutDO)
	log.Printf("Curriculums without curriculum_vision: %d\n", curriculumsWithoutDO)

	// Now fetch departments not in any cluster
	query := `
		SELECT c.id, c.name, c.academic_year, do.id as dept_overview_id
		FROM curriculum c
		INNER JOIN curriculum_vision do ON c.id = do.curriculum_id
		LEFT JOIN cluster_departments cd ON c.id = cd.department_id
		WHERE cd.id IS NULL
		ORDER BY c.name
	`

	log.Printf("Querying available departments...")
	rows, err := db.DB.Query(query)
	if err != nil {
		log.Println("Error querying available departments:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch available departments"})
		return
	}
	defer rows.Close()

	departments := []map[string]interface{}{}
	for rows.Next() {
		var id, deptOverviewID int
		var name, academicYear string
		if err := rows.Scan(&id, &name, &academicYear, &deptOverviewID); err == nil {
			departments = append(departments, map[string]interface{}{
				"id":                   id,
				"name":                 name,
				"academic_year":        academicYear,
				"curriculum_vision_id": deptOverviewID,
			})
		}
	}

	log.Printf("Found %d available departments\n", len(departments))
	json.NewEncoder(w).Encode(departments)
}

// CreateCluster creates a new cluster
func CreateCluster(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var cluster models.Cluster
	if err := json.NewDecoder(r.Body).Decode(&cluster); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	query := "INSERT INTO clusters (name, description) VALUES (?, ?)"
	result, err := db.DB.Exec(query, cluster.Name, cluster.Description)
	if err != nil {
		log.Println("Error creating cluster:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create cluster"})
		return
	}

	id, _ := result.LastInsertId()
	cluster.ID = int(id)

	json.NewEncoder(w).Encode(cluster)
}

// GetClusterDepartments retrieves all departments in a cluster
func GetClusterDepartments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	clusterID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid cluster ID"})
		return
	}

	query := `
		SELECT cd.id, cd.cluster_id, cd.department_id, cd.created_at,
		       c.id as regulation_id, c.name as regulation_name
		FROM cluster_departments cd
		LEFT JOIN curriculum c ON cd.department_id = c.id
		WHERE cd.cluster_id = ?
		ORDER BY cd.created_at
	`

	rows, err := db.DB.Query(query, clusterID)
	if err != nil {
		log.Println("Error querying cluster departments:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch cluster departments"})
		return
	}
	defer rows.Close()

	departments := []map[string]interface{}{}
	for rows.Next() {
		var cd models.ClusterDepartment
		var curriculumID sql.NullInt64
		var regulationName sql.NullString
		if err := rows.Scan(&cd.ID, &cd.ClusterID, &cd.DepartmentID, &cd.CreatedAt, &curriculumID, &regulationName); err == nil {
			dept := map[string]interface{}{
				"id":              cd.DepartmentID,
				"cluster_id":      cd.ClusterID,
				"department_id":   cd.DepartmentID,
				"created_at":      cd.CreatedAt,
			}
			if curriculumID.Valid {
				dept["curriculum_id"] = curriculumID.Int64
			}
			if regulationName.Valid {
				dept["name"] = regulationName.String
			}
			departments = append(departments, dept)
		}
	}

	json.NewEncoder(w).Encode(departments)
}

// AddDepartmentToCluster adds a department to a cluster
func AddDepartmentToCluster(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	clusterID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid cluster ID"})
		return
	}

	var req struct {
		DepartmentID int `json:"department_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Verify regulation exists in curriculum table
	var regulationExists bool
	err = db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM curriculum WHERE id = ?)", req.DepartmentID).Scan(&regulationExists)
	if err != nil || !regulationExists {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid department ID - regulation not found"})
		return
	}

	// Get or create curriculum_vision entry for this regulation
	var deptOverviewID int
	err = db.DB.QueryRow("SELECT id FROM curriculum_vision WHERE curriculum_id = ?", req.DepartmentID).Scan(&deptOverviewID)
	if err == sql.ErrNoRows {
		// Create curriculum_vision entry if it doesn't exist
		result, err := db.DB.Exec("INSERT INTO curriculum_vision (curriculum_id, vision) VALUES (?, '')", req.DepartmentID)
		if err != nil {
			log.Println("Error creating curriculum_vision:", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create department overview"})
			return
		}
		id, _ := result.LastInsertId()
		deptOverviewID = int(id)
	} else if err != nil {
		log.Println("Error checking curriculum_vision:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
		return
	}

	// Check if department is already in a cluster
	var existingClusterID sql.NullInt64
	var existingClusterName sql.NullString
	err = db.DB.QueryRow(`
		SELECT cd.cluster_id, c.name 
		FROM cluster_departments cd 
		JOIN clusters c ON cd.cluster_id = c.id 
		WHERE cd.department_id = ?
	`, deptOverviewID).Scan(&existingClusterID, &existingClusterName)

	if err == nil && existingClusterID.Valid {
		// Department is already in a cluster
		clusterName := "unknown"
		if existingClusterName.Valid {
			clusterName = existingClusterName.String
		}
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Department is already in cluster '%s' (ID: %d). A department can only belong to one cluster at a time.", clusterName, existingClusterID.Int64),
		})
		return
	} else if err != nil && err != sql.ErrNoRows {
		// Actual database error (not just "no rows")
		log.Println("Error checking existing cluster:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to check cluster membership"})
		return
	}

	// Add department to cluster using curriculum_vision id
	query := "INSERT INTO cluster_departments (cluster_id, department_id) VALUES (?, ?)"
	result, err := db.DB.Exec(query, clusterID, deptOverviewID)
	if err != nil {
		log.Println("Error adding department to cluster:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to add department to cluster"})
		return
	}

	id, _ := result.LastInsertId()
	response := models.ClusterDepartment{
		ID:           int(id),
		ClusterID:    clusterID,
		DepartmentID: deptOverviewID,
	}

	json.NewEncoder(w).Encode(response)
}

// RemoveDepartmentFromCluster removes a department from a cluster
func RemoveDepartmentFromCluster(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	clusterID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid cluster ID"})
		return
	}

	deptID, err := strconv.Atoi(vars["deptId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid department ID"})
		return
	}

	// Start a transaction to ensure all operations succeed or fail together
	tx, err := db.DB.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	// Step 1: Convert all received data to owned data for this department
	// This makes all received semesters/courses/honour cards now belong to the department
	err = convertReceivedDataToOwned(tx, deptID)
	if err != nil {
		log.Println("Error converting received data:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to convert received data"})
		return
	}

	// Step 2: Unshare all data owned by this department from other departments in the cluster
	err = unshareAllDepartmentData(tx, clusterID, deptID)
	if err != nil {
		log.Println("Error unsharing department data:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to unshare department data"})
		return
	}

	// Step 3: Remove department from cluster
	query := "DELETE FROM cluster_departments WHERE cluster_id = ? AND department_id = ?"
	result, err := tx.Exec(query, clusterID, deptID)
	if err != nil {
		log.Println("Error removing department from cluster:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to remove department from cluster"})
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Department not found in cluster"})
		return
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Println("Error committing transaction:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to commit transaction"})
		return
	}

	log.Printf("Department %d removed from cluster %d successfully", deptID, clusterID)
	json.NewEncoder(w).Encode(map[string]string{"message": "Department removed from cluster successfully"})
}

// DeleteCluster deletes a cluster
func DeleteCluster(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	clusterID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid cluster ID"})
		return
	}

	query := "DELETE FROM clusters WHERE id = ?"
	result, err := db.DB.Exec(query, clusterID)
	if err != nil {
		log.Println("Error deleting cluster:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete cluster"})
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Cluster not found"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Cluster deleted successfully"})
}

// convertReceivedDataToOwned converts all received data to owned data for a department
// This transfers ownership of all shared items to the receiving department
func convertReceivedDataToOwned(tx *sql.Tx, deptID int) error {
	log.Printf("Converting received data to owned for department %d", deptID)

	// Get all curriculum for this department
	var curriculumIDs []int
	rows, err := tx.Query("SELECT curriculum_id FROM curriculum_vision WHERE id = ?", deptID)
	if err != nil {
		return fmt.Errorf("error fetching curriculum IDs: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var regID int
		if err := rows.Scan(&regID); err == nil {
			curriculumIDs = append(curriculumIDs, regID)
		}
	}

	// For each curriculum, convert received data
	for _, regID := range curriculumIDs {
		// Convert normal_cards (semesters)
		_, err = tx.Exec(`
			UPDATE normal_cards 
			SET source_department_id = NULL, visibility = 'UNIQUE' 
			WHERE regulation_id = ? AND source_department_id IS NOT NULL
		`, regID)
		if err != nil {
			return fmt.Errorf("error converting normal_cards: %v", err)
		}

		// Convert honour_cards
		_, err = tx.Exec(`
			UPDATE honour_cards 
			SET source_department_id = NULL, visibility = 'UNIQUE' 
			WHERE regulation_id = ? AND source_department_id IS NOT NULL
		`, regID)
		if err != nil {
			return fmt.Errorf("error converting honour_cards: %v", err)
		}

		// Convert courses
		_, err = tx.Exec(`
			UPDATE courses c
			INNER JOIN curriculum_courses cc ON c.course_id = cc.course_id
			SET c.source_curriculum_id = NULL, c.visibility = 'UNIQUE' 
			WHERE cc.curriculum_id = ? AND c.source_curriculum_id IS NOT NULL
		`, regID)
		if err != nil {
			return fmt.Errorf("error converting courses: %v", err)
		}

		// Convert department mission
		_, err = tx.Exec(`
			UPDATE curriculum_mission 
			SET source_department_id = NULL, visibility = 'UNIQUE' 
			WHERE department_id = ? AND source_department_id IS NOT NULL
		`, deptID)
		if err != nil {
			return fmt.Errorf("error converting curriculum_mission: %v", err)
		}

		// Convert department PEOs
		_, err = tx.Exec(`
			UPDATE curriculum_peos 
			SET source_department_id = NULL, visibility = 'UNIQUE' 
			WHERE department_id = ? AND source_department_id IS NOT NULL
		`, deptID)
		if err != nil {
			return fmt.Errorf("error converting curriculum_peos: %v", err)
		}

		// Convert department POs
		_, err = tx.Exec(`
			UPDATE curriculum_pos 
			SET source_department_id = NULL, visibility = 'UNIQUE' 
			WHERE department_id = ? AND source_department_id IS NOT NULL
		`, deptID)
		if err != nil {
			return fmt.Errorf("error converting curriculum_pos: %v", err)
		}

		// Convert department PSOs
		_, err = tx.Exec(`
			UPDATE curriculum_psos 
			SET source_department_id = NULL, visibility = 'UNIQUE' 
			WHERE department_id = ? AND source_department_id IS NOT NULL
		`, deptID)
		if err != nil {
			return fmt.Errorf("error converting curriculum_psos: %v", err)
		}
	}

	log.Printf("Successfully converted received data to owned for department %d", deptID)
	return nil
}

// unshareAllDepartmentData unshares all data owned by a department from other departments in the cluster
func unshareAllDepartmentData(tx *sql.Tx, clusterID, sourceDeptID int) error {
	log.Printf("Unsharing all data from department %d in cluster %d", sourceDeptID, clusterID)

	// Get all other departments in the cluster
	var targetDeptIDs []int
	rows, err := tx.Query(`
		SELECT department_id 
		FROM cluster_departments 
		WHERE cluster_id = ? AND department_id != ?
	`, clusterID, sourceDeptID)
	if err != nil {
		return fmt.Errorf("error fetching cluster departments: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var deptID int
		if err := rows.Scan(&deptID); err == nil {
			targetDeptIDs = append(targetDeptIDs, deptID)
		}
	}

	// Get regulation ID for source department
	var sourceRegID int
	err = tx.QueryRow("SELECT curriculum_id FROM curriculum_vision WHERE id = ?", sourceDeptID).Scan(&sourceRegID)
	if err != nil {
		return fmt.Errorf("error fetching source regulation ID: %v", err)
	}

	// For each target department, convert shared items from source department to owned
	for _, targetDeptID := range targetDeptIDs {
		// Get target department's regulation ID
		var targetRegID int
		err := tx.QueryRow("SELECT curriculum_id FROM curriculum_vision WHERE id = ?", targetDeptID).Scan(&targetRegID)
		if err != nil {
			log.Printf("Error fetching regulation for department %d: %v", targetDeptID, err)
			continue
		}

		// Convert normal_cards that were shared from source department to owned
		_, err = tx.Exec(`
			UPDATE normal_cards 
			SET source_department_id = NULL, visibility = 'UNIQUE'
			WHERE regulation_id = ? AND source_department_id = ?
		`, targetRegID, sourceDeptID)
		if err != nil {
			return fmt.Errorf("error converting shared normal_cards to owned: %v", err)
		}

		// Convert honour_cards that were shared from source department to owned
		_, err = tx.Exec(`
			UPDATE honour_cards 
			SET source_department_id = NULL, visibility = 'UNIQUE'
			WHERE regulation_id = ? AND source_department_id = ?
		`, targetRegID, sourceDeptID)
		if err != nil {
			return fmt.Errorf("error converting shared honour_cards to owned: %v", err)
		}

		// Convert courses that were shared from source department to owned
		_, err = tx.Exec(`
			UPDATE courses c
			JOIN curriculum_courses cc ON cc.course_id = c.course_id
			SET c.source_curriculum_id = NULL, c.visibility = 'UNIQUE'
			WHERE cc.curriculum_id = ? AND c.source_curriculum_id = ?
		`, targetRegID, sourceRegID)
		if err != nil {
			return fmt.Errorf("error converting shared courses to owned: %v", err)
		}
	}

	// Update visibility of source department's own items to UNIQUE
	_, err = tx.Exec(`
		UPDATE normal_cards 
		SET visibility = 'UNIQUE' 
		WHERE regulation_id = ? AND visibility = 'CLUSTER' AND source_department_id IS NULL
	`, sourceRegID)
	if err != nil {
		return fmt.Errorf("error updating normal_cards visibility: %v", err)
	}

	_, err = tx.Exec(`
		UPDATE honour_cards 
		SET visibility = 'UNIQUE' 
		WHERE regulation_id = ? AND visibility = 'CLUSTER' AND source_department_id IS NULL
	`, sourceRegID)
	if err != nil {
		return fmt.Errorf("error updating honour_cards visibility: %v", err)
	}

	_, err = tx.Exec(`
		UPDATE courses c
		JOIN curriculum_courses cc ON c.course_id = cc.course_id
		SET c.visibility = 'UNIQUE' 
		WHERE cc.curriculum_id = ? AND c.visibility = 'CLUSTER' AND c.source_curriculum_id IS NULL
	`, sourceRegID)
	if err != nil {
		return fmt.Errorf("error updating courses visibility: %v", err)
	}

	log.Printf("Successfully unshared all data from department %d", sourceDeptID)
	return nil
}
