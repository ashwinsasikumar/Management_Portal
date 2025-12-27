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

	rows, err := db.DB.Query("SELECT id, name, description, created_at FROM clusters ORDER BY name")
	if err != nil {
		log.Println("Error fetching clusters:", err)
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

	result, err := db.DB.Exec("INSERT INTO clusters (name, description) VALUES (?, ?)",
		cluster.Name, cluster.Description)
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

	rows, err := db.DB.Query(`
		SELECT cd.id, cd.cluster_id, cd.department_id, d.regulation_id 
		FROM cluster_departments cd
		JOIN department_overview d ON cd.department_id = d.id
		WHERE cd.cluster_id = ?`, clusterID)
	if err != nil {
		log.Println("Error fetching cluster departments:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch departments"})
		return
	}
	defer rows.Close()

	departments := []models.ClusterDepartment{}
	for rows.Next() {
		var cd models.ClusterDepartment
		var regulationID int
		if err := rows.Scan(&cd.ID, &cd.ClusterID, &cd.DepartmentID, &regulationID); err == nil {
			departments = append(departments, cd)
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

	result, err := db.DB.Exec("INSERT INTO cluster_departments (cluster_id, department_id) VALUES (?, ?)",
		clusterID, req.DepartmentID)
	if err != nil {
		log.Println("Error adding department to cluster:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to add department"})
		return
	}

	id, _ := result.LastInsertId()
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id, "cluster_id": clusterID, "department_id": req.DepartmentID})
}

// RemoveDepartmentFromCluster removes a department from a cluster
func RemoveDepartmentFromCluster(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	clusterID, err1 := strconv.Atoi(vars["id"])
	departmentID, err2 := strconv.Atoi(vars["deptId"])
	if err1 != nil || err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid IDs"})
		return
	}

	_, err := db.DB.Exec("DELETE FROM cluster_departments WHERE cluster_id = ? AND department_id = ?",
		clusterID, departmentID)
	if err != nil {
		log.Println("Error removing department from cluster:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to remove department"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// GetClusterObjects retrieves cluster-level shared objects
func GetClusterObjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	clusterID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid cluster ID"})
		return
	}

	objects := models.ClusterObjects{
		ClusterID: clusterID,
		Mission:   fetchClusterListItems(clusterID, "cluster_mission", "mission_text"),
		PEOs:      fetchClusterListItems(clusterID, "cluster_peos", "peo_text"),
		POs:       fetchClusterListItems(clusterID, "cluster_pos", "po_text"),
		PSOs:      fetchClusterListItems(clusterID, "cluster_psos", "pso_text"),
	}

	json.NewEncoder(w).Encode(objects)
}

// SaveClusterObjects saves cluster-level shared objects
func SaveClusterObjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	clusterID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid cluster ID"})
		return
	}

	var objects models.ClusterObjects
	if err := json.NewDecoder(r.Body).Decode(&objects); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Save cluster objects
	saveClusterListItems(clusterID, "cluster_mission", "mission_text", objects.Mission)
	saveClusterListItems(clusterID, "cluster_peos", "peo_text", objects.PEOs)
	saveClusterListItems(clusterID, "cluster_pos", "po_text", objects.POs)
	saveClusterListItems(clusterID, "cluster_psos", "pso_text", objects.PSOs)

	objects.ClusterID = clusterID
	json.NewEncoder(w).Encode(objects)
}

// Helper function to fetch cluster-level list items
func fetchClusterListItems(clusterID int, tableName, columnName string) []string {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE cluster_id = ? ORDER BY position", columnName, tableName)
	rows, err := db.DB.Query(query, clusterID)
	if err != nil {
		return []string{}
	}
	defer rows.Close()

	items := []string{}
	for rows.Next() {
		var text string
		if err := rows.Scan(&text); err == nil {
			items = append(items, text)
		}
	}
	return items
}

// Helper function to save cluster-level list items
func saveClusterListItems(clusterID int, tableName, columnName string, items []string) error {
	// Delete existing items
	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE cluster_id = ?", tableName)
	_, err := db.DB.Exec(deleteQuery, clusterID)
	if err != nil {
		return err
	}

	// Insert new items with position
	insertQuery := fmt.Sprintf("INSERT INTO %s (cluster_id, %s, position) VALUES (?, ?, ?)", tableName, columnName)
	for i, text := range items {
		if text == "" {
			continue
		}
		_, err := db.DB.Exec(insertQuery, clusterID, text, i)
		if err != nil {
			return err
		}
	}
	return nil
}
