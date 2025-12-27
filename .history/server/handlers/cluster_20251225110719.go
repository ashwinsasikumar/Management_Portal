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

	query := `
		SELECT c.id, c.name, c.academic_year, c.regulation_no
		FROM curriculum c
		LEFT JOIN department_overview do ON c.id = do.regulation_id
		LEFT JOIN cluster_departments cd ON do.id = cd.department_id
		WHERE cd.id IS NULL
		ORDER BY c.name
	`

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
		var name, academicYear, regulationNo string
		if err := rows.Scan(&id, &name, &academicYear, &regulationNo); err == nil {
			departments = append(departments, map[string]interface{}{
				"id":            id,
				"name":          name,
				"academic_year": academicYear,
				"regulation_no": regulationNo,
			})
		}
	}

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
		       do.regulation_id, c.name as regulation_name
		FROM cluster_departments cd
		LEFT JOIN department_overview do ON cd.department_id = do.id
		LEFT JOIN curriculum c ON do.regulation_id = c.id
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
		var regulationID sql.NullInt64
		var regulationName sql.NullString
		if err := rows.Scan(&cd.ID, &cd.ClusterID, &cd.DepartmentID, &cd.CreatedAt, &regulationID, &regulationName); err == nil {
			dept := map[string]interface{}{
				"id":            cd.ID,
				"cluster_id":    cd.ClusterID,
				"department_id": cd.DepartmentID,
				"created_at":    cd.CreatedAt,
			}
			if regulationID.Valid {
				dept["regulation_id"] = regulationID.Int64
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

	// Get or create department_overview entry for this regulation
	var deptOverviewID int
	err = db.DB.QueryRow("SELECT id FROM department_overview WHERE regulation_id = ?", req.DepartmentID).Scan(&deptOverviewID)
	if err == sql.ErrNoRows {
		// Create department_overview entry if it doesn't exist
		result, err := db.DB.Exec("INSERT INTO department_overview (regulation_id, vision) VALUES (?, '')", req.DepartmentID)
		if err != nil {
			log.Println("Error creating department_overview:", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create department overview"})
			return
		}
		id, _ := result.LastInsertId()
		deptOverviewID = int(id)
	} else if err != nil {
		log.Println("Error checking department_overview:", err)
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

	// Add department to cluster using department_overview id
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

	query := "DELETE FROM cluster_departments WHERE cluster_id = ? AND department_id = ?"
	result, err := db.DB.Exec(query, clusterID, deptID)
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
