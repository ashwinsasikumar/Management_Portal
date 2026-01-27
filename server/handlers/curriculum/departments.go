package curriculum

import (
	"encoding/json"
	"log"
	"net/http"
	"server/db"
)

// Department represents a department entity
type Department struct {
	ID             int    `json:"id"`
	DepartmentName string `json:"department_name"`
	Status         int    `json:"status"`
}

// GetDepartments retrieves all active departments
func GetDepartments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	query := `
		SELECT id, department_name, status
		FROM departments
		WHERE status = 1
		ORDER BY department_name
	`

	rows, err := db.DB.Query(query)
	if err != nil {
		log.Printf("Error querying departments: %v", err)
		http.Error(w, "Failed to fetch departments", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	departments := []Department{}
	for rows.Next() {
		var dept Department
		if err := rows.Scan(&dept.ID, &dept.DepartmentName, &dept.Status); err != nil {
			log.Printf("Error scanning department row: %v", err)
			continue
		}
		departments = append(departments, dept)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating departments: %v", err)
		http.Error(w, "Failed to process departments", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(departments)
}
