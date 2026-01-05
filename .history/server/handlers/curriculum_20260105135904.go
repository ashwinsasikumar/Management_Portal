package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"server/db"
	"server/models"

	"github.com/gorilla/mux"
)

// GetSemesters retrieves all semesters for a regulation
func GetSemesters(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	regulationID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid regulation ID"})
		return
	}

	query := "SELECT id, regulation_id, semester_number, COALESCE(card_type, 'semester') as card_type FROM normal_cards WHERE regulation_id = ? ORDER BY COALESCE(semester_number, 999), id"
	rows, err := db.DB.Query(query, regulationID)
	if err != nil {
		log.Println("Error querying semesters:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch semesters"})
		return
	}
	defer rows.Close()

	var semesters []models.Semester = make([]models.Semester, 0)
	for rows.Next() {
		var sem models.Semester
		var semesterNum sql.NullInt64
		err := rows.Scan(&sem.ID, &sem.RegulationID, &semesterNum, &sem.CardType)
		if err != nil {
			log.Println("Error scanning semester:", err)
			continue
		}
		if semesterNum.Valid {
			val := int(semesterNum.Int64)
			sem.SemesterNumber = &val
		}
		semesters = append(semesters, sem)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(semesters)
}

// CreateSemester creates a new semester for a regulation
func CreateSemester(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	vars := mux.Vars(r)
	regulationID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid regulation ID"})
		return
	}

	var sem models.Semester
	err = json.NewDecoder(r.Body).Decode(&sem)
	if err != nil {
		log.Println("Error decoding request body:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	sem.RegulationID = regulationID

	query := "INSERT INTO normal_cards (regulation_id, semester_number, card_type) VALUES (?, ?, ?)"
	result, err := db.DB.Exec(query, sem.RegulationID, sem.SemesterNumber, sem.CardType)
	if err != nil {
		log.Println("Error inserting semester:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create semester"})
		return
	}

	id, _ := result.LastInsertId()
	sem.ID = int(id)

	// Log the activity
	logName := ""
	if sem.SemesterNumber != nil {
		cardTypeLabel := sem.CardType
		if cardTypeLabel == "" {
			cardTypeLabel = "semester"
		}
		logName = strings.Title(cardTypeLabel) + " " + strconv.Itoa(*sem.SemesterNumber)
	} else {
		logName = "New Card"
	}
	LogCurriculumActivity(regulationID, "Card Added",
		"Added "+logName, "System")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sem)
}

// DeleteSemester deletes a semester
func DeleteSemester(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	vars := mux.Vars(r)
	semesterID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid semester ID"})
		return
	}

	// Delete the semester (cascade will handle related records)
	query := "DELETE FROM normal_cards WHERE id = ?"
	result, err := db.DB.Exec(query, semesterID)
	if err != nil {
		log.Println("Error deleting semester:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete semester"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Semester not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Semester deleted successfully"})
}

// GetSemesterCourses retrieves all courses for a specific semester
func GetSemesterCourses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	regulationID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid regulation ID"})
		return
	}

	semesterID, err := strconv.Atoi(vars["semId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid semester ID"})
		return
	}

	query := `
		SELECT c.course_id, c.course_code, c.course_name, c.course_type, c.category, c.credit, 
		       c.theory_hours, c.activity_hours, c.lecture_hours, c.tutorial_hours, c.practical_hours, 
		       c.cia_marks, c.see_marks, c.total_marks, c.total_hours,
		       rc.id as reg_course_id
		FROM courses c
		INNER JOIN curriculum_courses rc ON c.course_id = rc.course_id
		WHERE rc.regulation_id = ? AND rc.semester_id = ?
		ORDER BY c.course_code
	`

	rows, err := db.DB.Query(query, regulationID, semesterID)
	if err != nil {
		log.Println("Error querying courses:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch courses"})
		return
	}
	defer rows.Close()

	var courses []models.CourseWithDetails = make([]models.CourseWithDetails, 0)
	for rows.Next() {
		var course models.CourseWithDetails
		err := rows.Scan(&course.CourseID, &course.CourseCode, &course.CourseName, &course.CourseType, &course.Category, &course.Credit,
			&course.TheoryHours, &course.ActivityHours, &course.LectureHours, &course.TutorialHours, &course.PracticalHours,
			&course.CIAMarks, &course.SEEMarks, &course.TotalMarks, &course.TotalHours,
			&course.RegCourseID)
		if err != nil {
			log.Println("Error scanning course:", err)
			continue
		}
		courses = append(courses, course)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(courses)
}

// AddCourseToSemester adds a new course and links it to a semester
func AddCourseToSemester(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	vars := mux.Vars(r)
	regulationID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid regulation ID"})
		return
	}

	semesterID, err := strconv.Atoi(vars["semId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid semester ID"})
		return
	}

	var course models.Course
	err = json.NewDecoder(r.Body).Decode(&course)
	if err != nil {
		log.Println("Error decoding request body:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Get curriculum's max_credits
	var maxCredits int
	err = db.DB.QueryRow("SELECT max_credits FROM curriculum WHERE id = ?", regulationID).Scan(&maxCredits)
	if err != nil {
		log.Println("Error fetching curriculum max_credits:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to validate credits"})
		return
	}

	// Calculate current total credits for this curriculum (across all semesters)
	var currentCredits sql.NullInt64
	creditQuery := `SELECT SUM(c.credit) FROM courses c 
	                INNER JOIN curriculum_courses rc ON c.course_id = rc.course_id 
	                WHERE rc.regulation_id = ?`
	err = db.DB.QueryRow(creditQuery, regulationID).Scan(&currentCredits)
	if err != nil {
		log.Println("Error calculating current credits:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to validate credits"})
		return
	}

	totalCredits := 0
	if currentCredits.Valid {
		totalCredits = int(currentCredits.Int64)
	}

	// Check if adding this course would exceed curriculum's max_credits
	if totalCredits+course.Credit > maxCredits {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Total credits exceed the curriculum's maximum allowed limit"})
		return
	}

	// Check if course already exists
	var existingCourseID int
	checkQuery := "SELECT course_id FROM courses WHERE course_code = ?"
	err = db.DB.QueryRow(checkQuery, course.CourseCode).Scan(&existingCourseID)

	var courseID int
	if err == sql.ErrNoRows {
		// Insert new course (total_marks is auto-computed by database)
		insertCourseQuery := `INSERT INTO courses (course_code, course_name, course_type, category, credit, 
		                      theory_hours, activity_hours, lecture_hours, tutorial_hours, practical_hours, cia_marks, see_marks) 
		                      VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		result, err := db.DB.Exec(insertCourseQuery, course.CourseCode, course.CourseName, course.CourseType, course.Category, course.Credit,
			course.TheoryHours, course.ActivityHours, course.LectureHours, course.TutorialHours, course.PracticalHours, course.CIAMarks, course.SEEMarks)
		if err != nil {
			log.Println("Error inserting course:", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create course"})
			return
		}
		id, _ := result.LastInsertId()
		courseID = int(id)
	} else if err != nil {
		log.Println("Error checking existing course:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create course"})
		return
	} else {
		courseID = existingCourseID
	}

	// Link course to regulation and semester
	linkQuery := "INSERT INTO curriculum_courses (regulation_id, semester_id, course_id) VALUES (?, ?, ?)"
	result, err := db.DB.Exec(linkQuery, regulationID, semesterID, courseID)
	if err != nil {
		log.Println("Error linking course:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to add course to semester"})
		return
	}

	regCourseID, _ := result.LastInsertId()
	course.CourseID = courseID

	// Fetch the complete course details including computed total_marks
	var fullCourse models.CourseWithDetails
	fetchQuery := `SELECT course_id, course_code, course_name, course_type, category, credit, 
	               theory_hours, activity_hours, lecture_hours, tutorial_hours, practical_hours, cia_marks, see_marks, total_marks, total_hours 
	               FROM courses WHERE course_id = ?`
	err = db.DB.QueryRow(fetchQuery, courseID).Scan(&fullCourse.CourseID, &fullCourse.CourseCode, &fullCourse.CourseName,
		&fullCourse.CourseType, &fullCourse.Category, &fullCourse.Credit,
		&fullCourse.TheoryHours, &fullCourse.ActivityHours, &fullCourse.LectureHours, &fullCourse.TutorialHours, &fullCourse.PracticalHours,
		&fullCourse.CIAMarks, &fullCourse.SEEMarks, &fullCourse.TotalMarks)
	if err != nil {
		log.Println("Error fetching full course details:", err)
		// Fallback to sent values
		fullCourse = models.CourseWithDetails{
			CourseID:       courseID,
			CourseCode:     course.CourseCode,
			CourseName:     course.CourseName,
			CourseType:     course.CourseType,
			Category:       course.Category,
			Credit:         course.Credit,
			LectureHours:   course.LectureHours,
			TutorialHours:  course.TutorialHours,
			PracticalHours: course.PracticalHours,
			CIAMarks:       course.CIAMarks,
			SEEMarks:       course.SEEMarks,
			TotalMarks:     course.CIAMarks + course.SEEMarks,
		}
	}
	fullCourse.RegCourseID = int(regCourseID)

	// Log the activity
	LogCurriculumActivity(regulationID, "Course Added",
		"Added course "+course.CourseCode+" - "+course.CourseName+" to Semester "+strconv.Itoa(semesterID), "System")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(fullCourse)
}

// RemoveCourseFromSemester removes a course from a semester
func RemoveCourseFromSemester(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	vars := mux.Vars(r)
	regulationID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid regulation ID"})
		return
	}

	semesterID, err := strconv.Atoi(vars["semId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid semester ID"})
		return
	}

	courseID, err := strconv.Atoi(vars["courseId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid course ID"})
		return
	}

	query := "DELETE FROM curriculum_courses WHERE regulation_id = ? AND semester_id = ? AND course_id = ?"
	result, err := db.DB.Exec(query, regulationID, semesterID, courseID)
	if err != nil {
		log.Println("Error removing course:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to remove course"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Course not found in semester"})
		return
	}

	// Get course name for logging
	var courseName string
	db.DB.QueryRow("SELECT course_name FROM courses WHERE course_id = ?", courseID).Scan(&courseName)

	// Log the activity
	LogCurriculumActivity(regulationID, "Course Removed",
		"Removed course "+courseName+" from Semester "+strconv.Itoa(semesterID), "System")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Course removed successfully"})
}
