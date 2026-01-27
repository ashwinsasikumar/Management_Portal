package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"server/db"
	"server/models"

	"github.com/gorilla/mux"
)

// GetCourseAllocations retrieves courses for a specific semester and academic year with their faculty assignments
func GetCourseAllocations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	semesterID := r.URL.Query().Get("semester_id")
	academicYear := r.URL.Query().Get("academic_year")

	if semesterID == "" || academicYear == "" {
		http.Error(w, "semester_id and academic_year are required", http.StatusBadRequest)
		return
	}

	// 1. Fetch all courses linked to this semester
	courseQuery := `
		SELECT c.course_id, c.course_code, c.course_name, c.course_type, c.credit
		FROM courses c
		JOIN curriculum_courses cc ON c.course_id = cc.course_id
		WHERE cc.semester_id = ? AND cc.status = 1 AND c.status = 1
	`
	rows, err := db.DB.Query(courseQuery, semesterID)
	if err != nil {
		log.Printf("Error fetching courses for allocation: %v", err)
		http.Error(w, "Failed to fetch courses", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var courses []models.CourseWithAllocations
	for rows.Next() {
		var c models.CourseWithAllocations
		if err := rows.Scan(&c.CourseID, &c.CourseCode, &c.CourseName, &c.CourseType, &c.Credit); err != nil {
			log.Printf("Error scanning course row: %v", err)
			continue
		}
		c.Allocations = []models.CourseAllocation{}
		courses = append(courses, c)
	}

	// 2. Fetch all allocations for these courses in this academic year
	allocationQuery := `
		SELECT ca.id, ca.course_id, ca.teacher_id, t.name, ca.academic_year, ca.semester, ca.section, ca.role
		FROM course_allocations ca
		JOIN teachers t ON ca.teacher_id = t.id
		WHERE ca.academic_year = ? AND ca.status = 1
	`
	aRows, err := db.DB.Query(allocationQuery, academicYear)
	if err != nil {
		log.Printf("Error fetching allocations: %v", err)
		// We still return courses even if allocations fetch fails
	} else {
		defer aRows.Close()
		allocMap := make(map[int][]models.CourseAllocation)
		for aRows.Next() {
			var a models.CourseAllocation
			if err := aRows.Scan(&a.ID, &a.CourseID, &a.TeacherID, &a.TeacherName, &a.AcademicYear, &a.Semester, &a.Section, &a.Role); err != nil {
				continue
			}
			allocMap[a.CourseID] = append(allocMap[a.CourseID], a)
		}

		// 3. Merge allocations into courses
		for i := range courses {
			if allocs, ok := allocMap[courses[i].CourseID]; ok {
				courses[i].Allocations = allocs
			}
		}
	}

	json.NewEncoder(w).Encode(courses)
}

// CreateAllocation assigns a teacher to a course
func CreateAllocation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var alloc models.CourseAllocation
	if err := json.NewDecoder(r.Body).Decode(&alloc); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if alloc.CourseID == 0 || alloc.TeacherID == 0 || alloc.AcademicYear == "" {
		http.Error(w, "CourseID, TeacherID, and AcademicYear are required", http.StatusBadRequest)
		return
	}

	if alloc.Section == "" {
		alloc.Section = "A"
	}
	if alloc.Role == "" {
		alloc.Role = "Primary"
	}

	query := `
		INSERT INTO course_allocations (course_id, teacher_id, academic_year, semester, section, role, status)
		VALUES (?, ?, ?, ?, ?, ?, 1)
		ON DUPLICATE KEY UPDATE status = 1, role = VALUES(role), section = VALUES(section)
	`
	_, err := db.DB.Exec(query, alloc.CourseID, alloc.TeacherID, alloc.AcademicYear, alloc.Semester, alloc.Section, alloc.Role)
	if err != nil {
		log.Printf("Error creating allocation: %v", err)
		http.Error(w, "Failed to create allocation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Allocation successful"})
}

// DeleteAllocation performs a soft delete of an allocation
func DeleteAllocation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	id := vars["id"]

	query := `UPDATE course_allocations SET status = 0 WHERE id = ?`
	_, err := db.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting allocation: %v", err)
		http.Error(w, "Failed to delete allocation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Allocation removed successfully"})
}

// UpdateAllocation updates an existing allocation
func UpdateAllocation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	id := vars["id"]

	var alloc models.CourseAllocation
	if err := json.NewDecoder(r.Body).Decode(&alloc); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	query := `
		UPDATE course_allocations 
		SET teacher_id = ?, role = ?, section = ?, semester = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND status = 1
	`
	_, err := db.DB.Exec(query, alloc.TeacherID, alloc.Role, alloc.Section, alloc.Semester, id)
	if err != nil {
		log.Printf("Error updating allocation: %v", err)
		http.Error(w, "Failed to update allocation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Allocation updated successfully"})
}

// GetTeacherCourses retrieves all courses assigned to a specific teacher
func GetTeacherCourses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	teacherID := vars["id"]

	academicYear := r.URL.Query().Get("academic_year")
	semester := r.URL.Query().Get("semester")

	query := `
		SELECT 
			ca.id, ca.course_id, c.course_code, c.course_name, c.course_type, 
			c.credit, ca.academic_year, ca.semester, ca.section, ca.role
		FROM course_allocations ca
		JOIN courses c ON ca.course_id = c.course_id
		WHERE ca.teacher_id = ? AND ca.status = 1
	`

	args := []interface{}{teacherID}
	if academicYear != "" {
		query += " AND ca.academic_year = ?"
		args = append(args, academicYear)
	}
	if semester != "" {
		query += " AND ca.semester = ?"
		args = append(args, semester)
	}
	query += " ORDER BY ca.semester, c.course_code"

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		log.Printf("Error fetching teacher courses: %v", err)
		http.Error(w, "Failed to fetch courses", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type TeacherCourse struct {
		ID           int    `json:"id"`
		CourseID     int    `json:"course_id"`
		CourseCode   string `json:"course_code"`
		CourseName   string `json:"course_name"`
		CourseType   string `json:"course_type"`
		Credit       int    `json:"credit"`
		AcademicYear string `json:"academic_year"`
		Semester     int    `json:"semester"`
		Section      string `json:"section"`
		Role         string `json:"role"`
	}

	var courses []TeacherCourse
	for rows.Next() {
		var course TeacherCourse
		err := rows.Scan(&course.ID, &course.CourseID, &course.CourseCode, &course.CourseName,
			&course.CourseType, &course.Credit, &course.AcademicYear, &course.Semester,
			&course.Section, &course.Role)
		if err != nil {
			log.Printf("Error scanning course row: %v", err)
			continue
		}
		courses = append(courses, course)
	}

	if courses == nil {
		courses = []TeacherCourse{}
	}

	json.NewEncoder(w).Encode(courses)
}

// GetCourseTeachers retrieves all teachers assigned to a specific course
func GetCourseTeachers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	courseID := vars["id"]

	academicYear := r.URL.Query().Get("academic_year")

	query := `
		SELECT 
			ca.id, ca.teacher_id, t.name, t.email, t.dept, d.department_name,
			ca.academic_year, ca.semester, ca.section, ca.role
		FROM course_allocations ca
		JOIN teachers t ON ca.teacher_id = t.id
		LEFT JOIN departments d ON t.dept = d.id
		WHERE ca.course_id = ? AND ca.status = 1
	`

	args := []interface{}{courseID}
	if academicYear != "" {
		query += " AND ca.academic_year = ?"
		args = append(args, academicYear)
	}
	query += " ORDER BY ca.role DESC, t.name"

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		log.Printf("Error fetching course teachers: %v", err)
		http.Error(w, "Failed to fetch teachers", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type CourseTeacher struct {
		ID             int     `json:"id"`
		TeacherID      int     `json:"teacher_id"`
		TeacherName    string  `json:"teacher_name"`
		Email          string  `json:"email"`
		DeptID         *int    `json:"dept_id"`
		DepartmentName *string `json:"department_name"`
		AcademicYear   string  `json:"academic_year"`
		Semester       int     `json:"semester"`
		Section        string  `json:"section"`
		Role           string  `json:"role"`
	}

	var teachers []CourseTeacher
	for rows.Next() {
		var teacher CourseTeacher
		err := rows.Scan(&teacher.ID, &teacher.TeacherID, &teacher.TeacherName, &teacher.Email,
			&teacher.DeptID, &teacher.DepartmentName, &teacher.AcademicYear, &teacher.Semester,
			&teacher.Section, &teacher.Role)
		if err != nil {
			log.Printf("Error scanning teacher row: %v", err)
			continue
		}
		teachers = append(teachers, teacher)
	}

	if teachers == nil {
		teachers = []CourseTeacher{}
	}

	json.NewEncoder(w).Encode(teachers)
}

// GetUnassignedCourses retrieves courses without teacher assignments
func GetUnassignedCourses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	semesterID := r.URL.Query().Get("semester_id")
	academicYear := r.URL.Query().Get("academic_year")

	if semesterID == "" || academicYear == "" {
		http.Error(w, "semester_id and academic_year are required", http.StatusBadRequest)
		return
	}

	query := `
		SELECT c.course_id, c.course_code, c.course_name, c.course_type, c.credit
		FROM courses c
		JOIN curriculum_courses cc ON c.course_id = cc.course_id
		WHERE cc.semester_id = ? AND c.status = 1 AND cc.status = 1
		AND NOT EXISTS (
			SELECT 1 FROM course_allocations ca
			WHERE ca.course_id = c.course_id 
			AND ca.academic_year = ?
			AND ca.status = 1
			AND ca.role = 'Primary'
		)
		ORDER BY c.course_code
	`

	rows, err := db.DB.Query(query, semesterID, academicYear)
	if err != nil {
		log.Printf("Error fetching unassigned courses: %v", err)
		http.Error(w, "Failed to fetch unassigned courses", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var courses []models.CourseWithAllocations
	for rows.Next() {
		var c models.CourseWithAllocations
		err := rows.Scan(&c.CourseID, &c.CourseCode, &c.CourseName, &c.CourseType, &c.Credit)
		if err != nil {
			log.Printf("Error scanning course row: %v", err)
			continue
		}
		c.Allocations = []models.CourseAllocation{}
		courses = append(courses, c)
	}

	if courses == nil {
		courses = []models.CourseWithAllocations{}
	}

	json.NewEncoder(w).Encode(courses)
}

// GetAllocationSummary retrieves allocation summary statistics
func GetAllocationSummary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	semesterID := r.URL.Query().Get("semester_id")
	academicYear := r.URL.Query().Get("academic_year")

	if semesterID == "" || academicYear == "" {
		http.Error(w, "semester_id and academic_year are required", http.StatusBadRequest)
		return
	}

	type Summary struct {
		TotalCourses      int `json:"total_courses"`
		AssignedCourses   int `json:"assigned_courses"`
		UnassignedCourses int `json:"unassigned_courses"`
		TotalTeachers     int `json:"total_teachers"`
		ActiveTeachers    int `json:"active_teachers"`
	}

	var summary Summary

	// Total courses
	err := db.DB.QueryRow(`
		SELECT COUNT(DISTINCT c.course_id)
		FROM courses c
		JOIN curriculum_courses cc ON c.course_id = cc.course_id
		WHERE cc.semester_id = ? AND c.status = 1 AND cc.status = 1
	`, semesterID).Scan(&summary.TotalCourses)
	if err != nil {
		log.Printf("Error counting total courses: %v", err)
	}

	// Assigned courses (with primary teacher)
	err = db.DB.QueryRow(`
		SELECT COUNT(DISTINCT ca.course_id)
		FROM course_allocations ca
		JOIN curriculum_courses cc ON ca.course_id = cc.course_id
		WHERE cc.semester_id = ? AND ca.academic_year = ? 
		AND ca.status = 1 AND ca.role = 'Primary'
	`, semesterID, academicYear).Scan(&summary.AssignedCourses)
	if err != nil {
		log.Printf("Error counting assigned courses: %v", err)
	}

	summary.UnassignedCourses = summary.TotalCourses - summary.AssignedCourses

	// Total teachers
	err = db.DB.QueryRow(`SELECT COUNT(*) FROM teachers WHERE status = 1`).Scan(&summary.TotalTeachers)
	if err != nil {
		log.Printf("Error counting total teachers: %v", err)
	}

	// Active teachers (assigned to at least one course in this semester)
	err = db.DB.QueryRow(`
		SELECT COUNT(DISTINCT ca.teacher_id)
		FROM course_allocations ca
		JOIN curriculum_courses cc ON ca.course_id = cc.course_id
		WHERE cc.semester_id = ? AND ca.academic_year = ? AND ca.status = 1
	`, semesterID, academicYear).Scan(&summary.ActiveTeachers)
	if err != nil {
		log.Printf("Error counting active teachers: %v", err)
	}

	json.NewEncoder(w).Encode(summary)
}
