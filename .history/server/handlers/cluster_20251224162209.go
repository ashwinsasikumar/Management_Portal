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
		       do.regulation_id
		FROM cluster_departments cd
		LEFT JOIN department_overview do ON cd.department_id = do.id
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
		if err := rows.Scan(&cd.ID, &cd.ClusterID, &cd.DepartmentID, &cd.CreatedAt, &regulationID); err == nil {
			dept := map[string]interface{}{
				"id":            cd.ID,
				"cluster_id":    cd.ClusterID,
				"department_id": cd.DepartmentID,
				"created_at":    cd.CreatedAt,
			}
			if regulationID.Valid {
				dept["regulation_id"] = regulationID.Int64
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

	// Verify department exists
	var exists bool
	err = db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM department_overview WHERE id = ?)", req.DepartmentID).Scan(&exists)
	if err != nil || !exists {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid department ID"})
		return
	}

	// Check if department is already in a cluster
	var existingClusterID sql.NullInt64
	db.DB.QueryRow("SELECT cluster_id FROM cluster_departments WHERE department_id = ?", req.DepartmentID).Scan(&existingClusterID)
	if existingClusterID.Valid {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Department is already in cluster %d", existingClusterID.Int64)})
		return
	}

	// Add department to cluster
	query := "INSERT INTO cluster_departments (cluster_id, department_id) VALUES (?, ?)"
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
		DepartmentID: req.DepartmentID,
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
