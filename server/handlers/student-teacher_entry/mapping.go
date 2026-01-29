package studentteacher

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/db"
	"server/models"
	"strconv"
)

// GetMappingFilters retrieves departments and available years for filtering
func GetMappingFilters(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	type FilterResponse struct {
		Departments []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"departments"`
		Years []int `json:"years"`
	}

	response := FilterResponse{
		Departments: make([]struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}, 0),
		Years: []int{2024, 2025, 2026, 2027, 2028},
	}

	// Get departments
	deptQuery := `SELECT id, department_name FROM departments WHERE status = 1 ORDER BY department_name`
	rows, err := db.DB.Query(deptQuery)
	if err != nil {
		log.Printf("Error fetching departments: %v", err)
		http.Error(w, "Failed to fetch departments", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var dept struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}
		if err := rows.Scan(&dept.ID, &dept.Name); err != nil {
			log.Printf("Error scanning department: %v", err)
			continue
		}
		response.Departments = append(response.Departments, dept)
	}

	json.NewEncoder(w).Encode(response)
}

// GetMappingData retrieves teachers and students for a specific department and year
func GetMappingData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	departmentID := r.URL.Query().Get("department_id")
	year := r.URL.Query().Get("year")
	academicYear := r.URL.Query().Get("academic_year")

	log.Printf("[MAPPING DEBUG] Request params - department_id: %s, year: %s, academic_year: %s", departmentID, year, academicYear)

	if departmentID == "" || year == "" {
		http.Error(w, "department_id and year are required", http.StatusBadRequest)
		return
	}

	type MappingData struct {
		Teachers []models.TeacherWithStudentCount `json:"teachers"`
		Students []models.StudentWithMapping      `json:"students"`
	}

	data := MappingData{
		Teachers: make([]models.TeacherWithStudentCount, 0),
		Students: make([]models.StudentWithMapping, 0),
	}

	// Get teachers in this department
	teacherQuery := `
		SELECT 
			t.id, 
			t.name, 
			COALESCE(t.email, ''),
			COALESCE(t.profile_img, ''),
			COALESCE(t.desg, ''),
			COALESCE(COUNT(stm.id), 0) as student_count
		FROM teachers t
		INNER JOIN department_teachers dt ON t.id = dt.teacher_id
		LEFT JOIN student_teacher_mapping stm ON t.id = stm.teacher_id 
			AND stm.department_id = ? 
			AND stm.year = ?
			` + func() string {
		if academicYear != "" {
			return `AND stm.academic_year = ?`
		}
		return ""
	}() + `
		WHERE dt.department_id = ? AND t.status = 1 AND dt.status = 1
		GROUP BY t.id, t.name, t.email, t.profile_img, t.desg
		ORDER BY t.name
	`

	var teacherRows *sql.Rows
	var err error

	if academicYear != "" {
		teacherRows, err = db.DB.Query(teacherQuery, departmentID, year, academicYear, departmentID)
	} else {
		teacherRows, err = db.DB.Query(teacherQuery, departmentID, year, departmentID)
	}

	if err != nil {
		log.Printf("Error fetching teachers: %v", err)
		http.Error(w, "Failed to fetch teachers", http.StatusInternalServerError)
		return
	}
	defer teacherRows.Close()

	for teacherRows.Next() {
		var teacher models.TeacherWithStudentCount
		if err := teacherRows.Scan(
			&teacher.TeacherID,
			&teacher.TeacherName,
			&teacher.Email,
			&teacher.ProfileImg,
			&teacher.Designation,
			&teacher.StudentCount,
		); err != nil {
			log.Printf("Error scanning teacher: %v", err)
			continue
		}
		data.Teachers = append(data.Teachers, teacher)
	}

	log.Printf("[MAPPING DEBUG] Found %d teachers", len(data.Teachers))

	// Get students in this department and year
	// studentQuery := `
	// 	SELECT
	// 		s.student_id,
	// 		COALESCE(s.enrollment_no, ''),
	// 		s.student_name,
	// 		COALESCE(ad.department, ''),
	// 		COALESCE(ad.year, 0),
	// 		stm.teacher_id,
	// 		t.name as teacher_name
	// 	FROM students s
	// 	INNER JOIN academic_details ad ON s.student_id = ad.student_id
	// 	INNER JOIN departments d ON ad.department = d.department_name
	// 	LEFT JOIN student_teacher_mapping stm ON s.student_id = stm.student_id
	// 		AND stm.year = ?
	// 		` + func() string {
	// 	if academicYear != "" {
	// 		return `AND stm.academic_year = ?`
	// 	}
	// 	return ""
	// }() + `
	// 	LEFT JOIN teachers t ON stm.teacher_id = t.id
	// 	WHERE d.id = ? AND ad.year = ? AND s.status = 1
	// 	ORDER BY s.enrollment_no
	// `

	studentQuery := `
		SELECT 
			s.student_id,
			COALESCE(s.enrollment_no, ''),
			s.student_name,
			COALESCE(ad.department, ''),
			COALESCE(ad.year, 0),
			stm.teacher_id,
			t.name as teacher_name
		FROM students s
		INNER JOIN academic_details ad ON s.student_id = ad.student_id
		INNER JOIN departments d ON ad.department = d.department_name
		LEFT JOIN student_teacher_mapping stm ON s.student_id = stm.student_id 
			` + `
		LEFT JOIN teachers t ON stm.teacher_id = t.id
		WHERE d.id = ? AND s.status = 1
		ORDER BY s.enrollment_no
	`

	var studentRows *sql.Rows

	if academicYear != "" {
		studentRows, err = db.DB.Query(studentQuery, departmentID)
	} else {
		studentRows, err = db.DB.Query(studentQuery, departmentID)
	}

	if err != nil {
		log.Printf("[MAPPING ERROR] Error fetching students: %v", err)
		http.Error(w, "Failed to fetch students", http.StatusInternalServerError)
		return
	}
	defer studentRows.Close()

	log.Printf("[MAPPING DEBUG] Student query executed successfully")

	for studentRows.Next() {
		var student models.StudentWithMapping
		if err := studentRows.Scan(
			&student.StudentID,
			&student.EnrollmentNo,
			&student.StudentName,
			&student.Department,
			&student.Year,
			&student.TeacherID,
			&student.TeacherName,
		); err != nil {
			log.Printf("Error scanning student: %v", err)
			continue
		}
		data.Students = append(data.Students, student)
	}

	log.Printf("[MAPPING DEBUG] Found %d students", len(data.Students))
	log.Printf("[MAPPING DEBUG] Sending response: %d teachers, %d students", len(data.Teachers), len(data.Students))

	json.NewEncoder(w).Encode(data)
}

// AssignStudentsToTeachers automatically distributes students evenly among teachers
func AssignStudentsToTeachers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	log.Printf("[DEBUG][AUTO-ASSIGN] === START Auto-Assign Request ===")
	log.Printf("[DEBUG][AUTO-ASSIGN] Method: %s, URL: %s", r.Method, r.URL.String())

	if r.Method != http.MethodPost {
		log.Printf("[DEBUG][AUTO-ASSIGN] Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.StudentTeacherMappingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[DEBUG][AUTO-ASSIGN] Error decoding JSON: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("[DEBUG][AUTO-ASSIGN] Request Body: DepartmentID=%d, Year=%d, AcademicYear=%s",
		req.DepartmentID, req.Year, req.AcademicYear)

	if req.DepartmentID == 0 || req.Year == 0 || req.AcademicYear == "" {
		log.Printf("[DEBUG][AUTO-ASSIGN] Missing required fields")
		http.Error(w, "department_id, year, and academic_year are required", http.StatusBadRequest)
		return
	}

	// Start transaction
	log.Printf("[DEBUG][AUTO-ASSIGN] Starting database transaction")
	tx, err := db.DB.Begin()
	if err != nil {
		log.Printf("[DEBUG][AUTO-ASSIGN] Error starting transaction: %v", err)
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()
	log.Printf("[DEBUG][AUTO-ASSIGN] Transaction started successfully")

	// First, clear existing mappings for this department, year, and academic year
	deleteQuery := `DELETE FROM student_teacher_mapping WHERE department_id = ? AND year = ? AND academic_year = ?`
	log.Printf("[DEBUG][AUTO-ASSIGN] Clearing existing mappings with query: %s", deleteQuery)
	log.Printf("[DEBUG][AUTO-ASSIGN] Delete params: dept=%d, year=%d, academic_year=%s",
		req.DepartmentID, req.Year, req.AcademicYear)

	_, err = tx.Exec(deleteQuery, req.DepartmentID, req.Year, req.AcademicYear)
	if err != nil {
		log.Printf("[DEBUG][AUTO-ASSIGN] Error clearing existing mappings: %v", err)
		http.Error(w, "Failed to clear existing mappings", http.StatusInternalServerError)
		return
	}
	log.Printf("[DEBUG][AUTO-ASSIGN] Existing mappings cleared successfully")

	// Get all active teachers in this department
	teacherQuery := `
		SELECT t.id 
		FROM teachers t
		INNER JOIN department_teachers dt ON t.id = dt.teacher_id
		WHERE dt.department_id = ? AND t.status = 1 AND dt.status = 1
		ORDER BY t.id
	`
	log.Printf("[DEBUG][AUTO-ASSIGN] Fetching teachers with query: %s", teacherQuery)
	log.Printf("[DEBUG][AUTO-ASSIGN] Teacher query param: dept=%d", req.DepartmentID)

	teacherRows, err := tx.Query(teacherQuery, req.DepartmentID)
	if err != nil {
		log.Printf("[DEBUG][AUTO-ASSIGN] Error fetching teachers: %v", err)
		http.Error(w, "Failed to fetch teachers", http.StatusInternalServerError)
		return
	}
	defer teacherRows.Close()

	teachers := make([]int64, 0)
	teacherCount := 0
	for teacherRows.Next() {
		var teacherID int64
		if err := teacherRows.Scan(&teacherID); err != nil {
			log.Printf("[DEBUG][AUTO-ASSIGN] Error scanning teacher ID: %v", err)
			continue
		}
		teachers = append(teachers, teacherID)
		teacherCount++
		log.Printf("[DEBUG][AUTO-ASSIGN] Found teacher ID: %d", teacherID)
	}
	teacherRows.Close()

	log.Printf("[DEBUG][AUTO-ASSIGN] Total teachers found: %d", len(teachers))
	log.Printf("[DEBUG][AUTO-ASSIGN] Teacher IDs: %v", teachers)

	if len(teachers) == 0 {
		log.Printf("[DEBUG][AUTO-ASSIGN] No teachers found in department %d", req.DepartmentID)
		http.Error(w, "No teachers found in this department", http.StatusBadRequest)
		return
	}

	// Get all students in this department and year
	studentQuery := `
		SELECT s.student_id
		FROM students s
		INNER JOIN academic_details ad ON s.student_id = ad.student_id
		INNER JOIN departments d ON ad.department = d.department_name
		WHERE d.id = ? AND s.status = 1
		ORDER BY s.student_id
	`
	log.Printf("[DEBUG][AUTO-ASSIGN] Fetching students with query: %s", studentQuery)
	log.Printf("[DEBUG][AUTO-ASSIGN] Student query params: dept=%d, year=%d", req.DepartmentID, req.Year)

	studentRows, err := tx.Query(studentQuery, req.DepartmentID)
	if err != nil {
		log.Printf("[DEBUG][AUTO-ASSIGN] Error fetching students: %v", err)
		http.Error(w, "Failed to fetch students", http.StatusInternalServerError)
		return
	}
	defer studentRows.Close()

	students := make([]int, 0)
	studentCount := 0
	for studentRows.Next() {
		var studentID int
		if err := studentRows.Scan(&studentID); err != nil {
			log.Printf("[DEBUG][AUTO-ASSIGN] Error scanning student ID: %v", err)
			continue
		}
		students = append(students, studentID)
		studentCount++
		if studentCount <= 10 { // Log first 10 students to avoid too much logging
			log.Printf("[DEBUG][AUTO-ASSIGN] Found student ID: %d", studentID)
		}
	}
	studentRows.Close()

	log.Printf("[DEBUG][AUTO-ASSIGN] Total students found: %d", len(students))
	if len(students) > 10 {
		log.Printf("[DEBUG][AUTO-ASSIGN] First 10 student IDs: %v", students[:10])
	} else {
		log.Printf("[DEBUG][AUTO-ASSIGN] Student IDs: %v", students)
	}

	if len(students) == 0 {
		log.Printf("[DEBUG][AUTO-ASSIGN] No students found for dept=%d, year=%d", req.DepartmentID, req.Year)
		response := models.StudentTeacherMappingResponse{
			Success:         true,
			Message:         "No students found to assign",
			TotalStudents:   0,
			TotalTeachers:   len(teachers),
			MappingsCreated: 0,
		}
		json.NewEncoder(w).Encode(response)
		tx.Commit()
		log.Printf("[DEBUG][AUTO-ASSIGN] Transaction committed (no students case)")
		log.Printf("[DEBUG][AUTO-ASSIGN] === END Auto-Assign Request ===")
		return
	}

	// Calculate distribution
	totalStudents := len(students)
	totalTeachers := len(teachers)
	baseCount := totalStudents / totalTeachers
	remainder := totalStudents % totalTeachers

	log.Printf("[DEBUG][AUTO-ASSIGN] Distribution calculation:")
	log.Printf("[DEBUG][AUTO-ASSIGN]   Total Students: %d", totalStudents)
	log.Printf("[DEBUG][AUTO-ASSIGN]   Total Teachers: %d", totalTeachers)
	log.Printf("[DEBUG][AUTO-ASSIGN]   Base Count (students per teacher): %d", baseCount)
	log.Printf("[DEBUG][AUTO-ASSIGN]   Remainder (extra students): %d", remainder)
	log.Printf("[DEBUG][AUTO-ASSIGN]   First %d teachers will get %d students each",
		remainder, baseCount+1)
	log.Printf("[DEBUG][AUTO-ASSIGN]   Remaining %d teachers will get %d students each",
		totalTeachers-remainder, baseCount)

	// Distribute students to teachers
	insertQuery := `INSERT INTO student_teacher_mapping (student_id, teacher_id, department_id, year, academic_year) VALUES (?, ?, ?, ?, ?)`
	log.Printf("[DEBUG][AUTO-ASSIGN] Insert query: %s", insertQuery)

	studentIndex := 0
	mappingsCreated := 0

	log.Printf("[DEBUG][AUTO-ASSIGN] Starting distribution...")
	for teacherIdx, teacherID := range teachers {
		// First 'remainder' teachers get one extra student
		studentsForThisTeacher := baseCount
		if teacherIdx < remainder {
			studentsForThisTeacher++
		}

		log.Printf("[DEBUG][AUTO-ASSIGN] Teacher %d (index %d) will get %d students",
			teacherID, teacherIdx, studentsForThisTeacher)

		for i := 0; i < studentsForThisTeacher && studentIndex < totalStudents; i++ {
			studentID := students[studentIndex]

			log.Printf("[DEBUG][AUTO-ASSIGN]   Assigning student %d to teacher %d",
				studentID, teacherID)

			_, err := tx.Exec(insertQuery, studentID, teacherID, req.DepartmentID, req.Year, req.AcademicYear)
			if err != nil {
				log.Printf("[DEBUG][AUTO-ASSIGN] ERROR inserting mapping: student=%d, teacher=%d, error=%v",
					studentID, teacherID, err)
				http.Error(w, fmt.Sprintf("Failed to create mapping: %v", err), http.StatusInternalServerError)
				return
			}
			mappingsCreated++
			studentIndex++

			// Log every 10th assignment to track progress
			if mappingsCreated%10 == 0 {
				log.Printf("[DEBUG][AUTO-ASSIGN]   Progress: %d mappings created", mappingsCreated)
			}
		}
		log.Printf("[DEBUG][AUTO-ASSIGN] Teacher %d assigned %d students", teacherID, studentsForThisTeacher)
	}

	log.Printf("[DEBUG][AUTO-ASSIGN] Distribution complete: %d mappings created", mappingsCreated)
	log.Printf("[DEBUG][AUTO-ASSIGN] Student index after distribution: %d (should be %d)",
		studentIndex, totalStudents)

	// Commit transaction
	log.Printf("[DEBUG][AUTO-ASSIGN] Committing transaction...")
	if err := tx.Commit(); err != nil {
		log.Printf("[DEBUG][AUTO-ASSIGN] ERROR committing transaction: %v", err)
		http.Error(w, "Failed to commit mappings", http.StatusInternalServerError)
		return
	}
	log.Printf("[DEBUG][AUTO-ASSIGN] Transaction committed successfully")

	response := models.StudentTeacherMappingResponse{
		Success:         true,
		Message:         fmt.Sprintf("Successfully assigned %d students to %d teachers", totalStudents, totalTeachers),
		TotalStudents:   totalStudents,
		TotalTeachers:   totalTeachers,
		MappingsCreated: mappingsCreated,
	}

	log.Printf("[DEBUG][AUTO-ASSIGN] Sending response: %+v", response)
	log.Printf("[DEBUG][AUTO-ASSIGN] === END Auto-Assign Request ===")

	json.NewEncoder(w).Encode(response)
}

// ClearMappings removes all mappings for a specific department, year, and academic year
func ClearMappings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	departmentID := r.URL.Query().Get("department_id")
	year := r.URL.Query().Get("year")
	academicYear := r.URL.Query().Get("academic_year")

	if departmentID == "" || year == "" || academicYear == "" {
		http.Error(w, "department_id, year, and academic_year are required", http.StatusBadRequest)
		return
	}

	deptID, err := strconv.Atoi(departmentID)
	if err != nil {
		http.Error(w, "Invalid department_id", http.StatusBadRequest)
		return
	}

	yearInt, err := strconv.Atoi(year)
	if err != nil {
		http.Error(w, "Invalid year", http.StatusBadRequest)
		return
	}

	deleteQuery := `DELETE FROM student_teacher_mapping WHERE department_id = ? AND year = ? AND academic_year = ?`
	result, err := db.DB.Exec(deleteQuery, deptID, yearInt, academicYear)
	if err != nil {
		log.Printf("Error clearing mappings: %v", err)
		http.Error(w, "Failed to clear mappings", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()

	response := map[string]interface{}{
		"success":       true,
		"message":       fmt.Sprintf("Cleared %d mappings", rowsAffected),
		"rows_affected": rowsAffected,
	}

	json.NewEncoder(w).Encode(response)
}
