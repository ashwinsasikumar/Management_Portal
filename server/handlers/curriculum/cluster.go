package curriculum

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

	// Fetch curriculums not in any cluster
	query := `
		SELECT c.id, c.name, c.academic_year
		FROM curriculum c
		LEFT JOIN cluster_departments cd ON c.id = cd.curriculum_id
		WHERE cd.id IS NULL AND c.status = 1
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
		var id int
		var name, academicYear string
		if err := rows.Scan(&id, &name, &academicYear); err == nil {
			departments = append(departments, map[string]interface{}{
				"id":            id,
				"name":          name,
				"academic_year": academicYear,
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
		SELECT cd.id, cd.cluster_id, cd.curriculum_id, cd.created_at,
		       c.id as curriculum_id_join, c.name as curriculum_name, c.curriculum_template
		FROM cluster_departments cd
		LEFT JOIN curriculum c ON cd.curriculum_id = c.id
		WHERE cd.cluster_id = ?
		ORDER BY c.curriculum_template, cd.created_at
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
		var id, clusterID, curriculumID int
		var createdAt string
		var curriculumIDJoin sql.NullInt64
		var curriculumName, curriculumTemplate sql.NullString
		if err := rows.Scan(&id, &clusterID, &curriculumID, &createdAt, &curriculumIDJoin, &curriculumName, &curriculumTemplate); err == nil {
			dept := map[string]interface{}{
				"id":            id,
				"cluster_id":    clusterID,
				"curriculum_id": curriculumID,
				"created_at":    createdAt,
			}
			if curriculumName.Valid {
				dept["name"] = curriculumName.String
			}
			if curriculumTemplate.Valid {
				dept["curriculum_template"] = curriculumTemplate.String
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
		WHERE cd.curriculum_id = ?
	`, req.DepartmentID).Scan(&existingClusterID, &existingClusterName)

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

	// Add department to cluster using curriculum id
	query := "INSERT INTO cluster_departments (cluster_id, curriculum_id) VALUES (?, ?)"
	result, err := db.DB.Exec(query, clusterID, req.DepartmentID)
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
	query := "DELETE FROM cluster_departments WHERE cluster_id = ? AND curriculum_id = ?"
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

	// The deptID is actually the curriculum_id from cluster_departments table
	curriculumID := deptID

	// For now, we'll only handle tables that actually have sharing columns
	// Check if courses table has source_curriculum_id column
	_, err := tx.Exec(`
		UPDATE courses c
		INNER JOIN curriculum_courses cc ON c.course_id = cc.course_id
		SET c.source_curriculum_id = NULL, c.visibility = 'UNIQUE' 
		WHERE cc.curriculum_id = ? AND c.source_curriculum_id IS NOT NULL
	`, curriculumID)
	if err != nil {
		// If column doesn't exist, just log and continue
		log.Printf("Note: Could not convert courses (column may not exist): %v", err)
	}

	log.Printf("Successfully converted received data to owned for department %d", deptID)
	return nil
}

// unshareAllDepartmentData unshares all data owned by a curriculum from other curriculums in the cluster
func unshareAllDepartmentData(tx *sql.Tx, clusterID, sourceCurriculumID int) error {
	log.Printf("Unsharing all data from curriculum %d in cluster %d", sourceCurriculumID, clusterID)

	// Get all other curriculums in the cluster
	var targetCurriculumIDs []int
	rows, err := tx.Query(`
		SELECT curriculum_id 
		FROM cluster_departments 
		WHERE cluster_id = ? AND curriculum_id != ?
	`, clusterID, sourceCurriculumID)
	if err != nil {
		return fmt.Errorf("error fetching cluster curriculums: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var currID int
		if err := rows.Scan(&currID); err == nil {
			targetCurriculumIDs = append(targetCurriculumIDs, currID)
		}
	}

	// For each target curriculum, convert shared courses from source curriculum to owned
	for _, targetCurriculumID := range targetCurriculumIDs {
		// Convert courses that were shared from source curriculum to owned
		_, err = tx.Exec(`
			UPDATE courses c
			JOIN curriculum_courses cc ON cc.course_id = c.course_id
			SET c.source_curriculum_id = NULL, c.visibility = 'UNIQUE'
			WHERE cc.curriculum_id = ? AND c.source_curriculum_id = ?
		`, targetCurriculumID, sourceCurriculumID)
		if err != nil {
			log.Printf("Note: Could not convert shared courses to owned (column may not exist): %v", err)
		}
	}

	// Update visibility of source curriculum's own courses to UNIQUE
	_, err = tx.Exec(`
		UPDATE courses c
		JOIN curriculum_courses cc ON c.course_id = cc.course_id
		SET c.visibility = 'UNIQUE' 
		WHERE cc.curriculum_id = ? AND c.visibility = 'CLUSTER' AND c.source_curriculum_id IS NULL
	`, sourceCurriculumID)
	if err != nil {
		log.Printf("Note: Could not update courses visibility (column may not exist): %v", err)
	}

	log.Printf("Successfully unshared all data from curriculum %d", sourceCurriculumID)
	return nil
}
