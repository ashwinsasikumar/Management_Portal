package studentteacher

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"server/db"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Teacher represents the teacher model
type Teacher struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Phone      *string   `json:"phone"`
	ProfileImg *string   `json:"profile_img"`
	Dept       *int      `json:"dept"`
	Department *string   `json:"department"` // For display purposes
	Desg       *string   `json:"designation"`
	LastLogin  time.Time `json:"last_login"`
	Status     int       `json:"status"` // 1 = active, 0 = deleted
}

// TeacherInput represents the input for creating/updating a teacher
type TeacherInput struct {
	Name       string  `json:"name"`
	Email      string  `json:"email"`
	Phone      *string `json:"phone"`
	ProfileImg *string `json:"profile_img"`
	Department string  `json:"department"` // Department name from frontend
	Desg       *string `json:"designation"`
}

// GetTeachers retrieves all teachers from the database
func GetTeachers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	query := `
		SELECT 
			t.id, t.name, t.email, t.phone, t.profile_img, 
			t.dept, d.department_name as department, t.desg, 
			t.last_login, t.status
		FROM teachers t
		LEFT JOIN departments d ON t.dept = d.id
		WHERE t.status = 1
		ORDER BY t.id DESC
	`

	rows, err := db.DB.Query(query)
	if err != nil {
		log.Printf("Error querying teachers: %v", err)
		http.Error(w, "Failed to fetch teachers", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var teachers []Teacher
	for rows.Next() {
		var teacher Teacher
		err := rows.Scan(
			&teacher.ID, &teacher.Name, &teacher.Email, &teacher.Phone,
			&teacher.ProfileImg, &teacher.Dept, &teacher.Department, &teacher.Desg,
			&teacher.LastLogin, &teacher.Status,
		)
		if err != nil {
			log.Printf("Error scanning teacher row: %v", err)
			continue
		}
		teachers = append(teachers, teacher)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating teacher rows: %v", err)
		http.Error(w, "Failed to process teachers", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(teachers)
}

// GetTeacherByID retrieves a single teacher by ID
func GetTeacherByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}

	query := `
		SELECT 
			t.id, t.name, t.email, t.phone, t.profile_img, 
			t.dept, d.department_name as department, t.desg, 
			t.last_login, t.status
		FROM teachers t
		LEFT JOIN departments d ON t.dept = d.id
		WHERE t.id = ? AND t.status = 1
	`

	var teacher Teacher
	err = db.DB.QueryRow(query, id).Scan(
		&teacher.ID, &teacher.Name, &teacher.Email, &teacher.Phone,
		&teacher.ProfileImg, &teacher.Dept, &teacher.Department, &teacher.Desg,
		&teacher.LastLogin, &teacher.Status,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Error querying teacher: %v", err)
		http.Error(w, "Failed to fetch teacher", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(teacher)
}

// CreateTeacher creates a new teacher
func CreateTeacher(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var input TeacherInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if input.Name == "" || input.Email == "" {
		http.Error(w, "Name and email are required", http.StatusBadRequest)
		return
	}

	// Get department ID from department name
	var deptID *int
	if input.Department != "" {
		var tempDeptID int
		err := db.DB.QueryRow("SELECT id FROM departments WHERE department_name = ?", input.Department).Scan(&tempDeptID)
		if err == nil {
			deptID = &tempDeptID
		} else {
			log.Printf("Warning: Department '%s' not found: %v", input.Department, err)
			// Try to insert the department if it doesn't exist
			result, insertErr := db.DB.Exec("INSERT INTO departments (department_name, status) VALUES (?, 1)", input.Department)
			if insertErr == nil {
				newID, _ := result.LastInsertId()
				tempDeptID = int(newID)
				deptID = &tempDeptID
				log.Printf("Created new department '%s' with ID %d", input.Department, tempDeptID)
			} else {
				log.Printf("Failed to create department '%s': %v", input.Department, insertErr)
			}
		}
	}

	// Insert teacher with status = 1 (active) by default
	query := `
		INSERT INTO teachers (name, email, phone, profile_img, dept, desg, status)
		VALUES (?, ?, ?, ?, ?, ?, 1)
	`

	result, err := db.DB.Exec(
		query,
		input.Name,
		input.Email,
		input.Phone,
		input.ProfileImg,
		deptID,
		input.Desg,
	)

	if err != nil {
		log.Printf("Error creating teacher: %v", err)
		if err.Error() == "Error 1062: Duplicate entry" {
			http.Error(w, "Teacher with this email already exists", http.StatusConflict)
		} else {
			http.Error(w, "Failed to create teacher", http.StatusInternalServerError)
		}
		return
	}

	teacherID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID: %v", err)
		http.Error(w, "Failed to get teacher ID", http.StatusInternalServerError)
		return
	}

	// Insert into department_teachers junction table if department is provided
	if deptID != nil {
		_, err = db.DB.Exec(
			"INSERT INTO department_teachers (teacher_id, department_id, status) VALUES (?, ?, 1)",
			teacherID,
			*deptID,
		)
		if err != nil {
			log.Printf("Warning: Failed to link teacher to department: %v", err)
		} else {
			log.Printf("Linked teacher ID %d to department ID %d", teacherID, *deptID)
		}
	}

	// Fetch and return the created teacher
	createdTeacher := Teacher{
		ID:         teacherID,
		Name:       input.Name,
		Email:      input.Email,
		Phone:      input.Phone,
		ProfileImg: input.ProfileImg,
		Dept:       deptID,
		Department: &input.Department,
		Desg:       input.Desg,
		Status:     1,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdTeacher)
}

// UpdateTeacher updates an existing teacher
func UpdateTeacher(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}

	var input TeacherInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if input.Name == "" || input.Email == "" {
		http.Error(w, "Name and email are required", http.StatusBadRequest)
		return
	}

	// Get department ID from department name
	var deptID *int
	if input.Department != "" {
		var tempDeptID int
		err := db.DB.QueryRow("SELECT id FROM departments WHERE department_name = ?", input.Department).Scan(&tempDeptID)
		if err == nil {
			deptID = &tempDeptID
		} else {
			log.Printf("Warning: Department '%s' not found: %v", input.Department, err)
			// Try to insert the department if it doesn't exist
			result, insertErr := db.DB.Exec("INSERT INTO departments (department_name, status) VALUES (?, 1)", input.Department)
			if insertErr == nil {
				newID, _ := result.LastInsertId()
				tempDeptID = int(newID)
				deptID = &tempDeptID
				log.Printf("Created new department '%s' with ID %d", input.Department, tempDeptID)
			} else {
				log.Printf("Failed to create department '%s': %v", input.Department, insertErr)
			}
		}
	}

	// Update teacher
	query := `
		UPDATE teachers 
		SET name = ?, email = ?, phone = ?, profile_img = ?, dept = ?, desg = ?
		WHERE id = ? AND status = 1
	`

	result, err := db.DB.Exec(
		query,
		input.Name,
		input.Email,
		input.Phone,
		input.ProfileImg,
		deptID,
		input.Desg,
		id,
	)

	if err != nil {
		log.Printf("Error updating teacher: %v", err)
		http.Error(w, "Failed to update teacher", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		http.Error(w, "Failed to update teacher", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	}

	// Update department_teachers junction table if department is provided
	if deptID != nil {
		// First, remove old department associations
		_, err = db.DB.Exec("DELETE FROM department_teachers WHERE teacher_id = ?", id)
		if err != nil {
			log.Printf("Warning: Failed to remove old department associations: %v", err)
		}

		// Then add the new department association
		_, err = db.DB.Exec(
			"INSERT INTO department_teachers (teacher_id, department_id, status) VALUES (?, ?, 1)",
			id,
			*deptID,
		)
		if err != nil {
			log.Printf("Warning: Failed to link teacher to department: %v", err)
		} else {
			log.Printf("Updated department link for teacher ID %d to department ID %d", id, *deptID)
		}
	}

	// Fetch and return the updated teacher
	updatedTeacher := Teacher{
		ID:         id,
		Name:       input.Name,
		Email:      input.Email,
		Phone:      input.Phone,
		ProfileImg: input.ProfileImg,
		Dept:       deptID,
		Department: &input.Department,
		Desg:       input.Desg,
		Status:     1,
	}

	json.NewEncoder(w).Encode(updatedTeacher)
}

// DeleteTeacher soft deletes a teacher by setting status to 0
func DeleteTeacher(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}

	// Soft delete: set status to 0
	query := "UPDATE teachers SET status = 0 WHERE id = ? AND status = 1"
	result, err := db.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting teacher: %v", err)
		http.Error(w, "Failed to delete teacher", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		http.Error(w, "Failed to delete teacher", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Teacher deleted successfully"})
}
