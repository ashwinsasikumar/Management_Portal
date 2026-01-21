package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/db"
	"server/models"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// GetSemesters retrieves all semesters for a regulation
func GetSemesters(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	curriculumID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid curriculum ID"})
		return
	}

	query := "SELECT id, curriculum_id, semester_number, COALESCE(card_type, 'semester') as card_type FROM normal_cards WHERE curriculum_id = ? ORDER BY COALESCE(semester_number, 999), id"
	rows, err := db.DB.Query(query, curriculumID)
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
		err := rows.Scan(&sem.ID, &sem.CurriculumID, &semesterNum, &sem.CardType)
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
	curriculumID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid curriculum ID"})
		return
	}

	var card models.Semester
	err = json.NewDecoder(r.Body).Decode(&card)
	if err != nil {
		log.Println("Error decoding request body:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	card.CurriculumID = curriculumID

	// Set default card type if not provided
	if card.CardType == "" {
		card.CardType = "semester"
	}

	// Check if the number already exists within the same card type
	// Semester numbers must be unique among semester cards
	// Vertical numbers must be unique among vertical cards
	if card.SemesterNumber != nil {
		var existingCount int
		err = db.DB.QueryRow("SELECT COUNT(*) FROM normal_cards WHERE curriculum_id = ? AND semester_number = ? AND card_type = ?",
			curriculumID, *card.SemesterNumber, card.CardType).Scan(&existingCount)
		if err != nil {
			log.Println("Error checking for duplicate number:", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to validate number"})
			return
		}
		if existingCount > 0 {
			if card.CardType == "vertical" {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Vertical %d already exists in this curriculum", *card.SemesterNumber)})
			} else {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Semester %d already exists in this curriculum", *card.SemesterNumber)})
			}
			return
		}
	}

	query := "INSERT INTO normal_cards (curriculum_id, semester_number, card_type) VALUES (?, ?, ?)"
	result, err := db.DB.Exec(query, card.CurriculumID, card.SemesterNumber, card.CardType)
	if err != nil {
		log.Println("Error inserting semester:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create semester"})
		return
	}

	id, _ := result.LastInsertId()
	card.ID = int(id)

	// Log the activity
	logName := ""
	if card.SemesterNumber != nil {
		cardTypeLabel := card.CardType
		if cardTypeLabel == "" {
			cardTypeLabel = "semester"
		}
		logName = strings.Title(cardTypeLabel) + " " + strconv.Itoa(*card.SemesterNumber)
	} else {
		logName = "New Card"
	}
	LogCurriculumActivity(curriculumID, "Card Added",
		"Added "+logName, "System")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(card)
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
	curriculumID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid curriculum ID"})
		return
	}

	semesterID, err := strconv.Atoi(vars["semId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid semester ID"})
		return
	}

	curriculumTemplate := getCurriculumTemplateByRegulation(curriculumID)

	query := `
		SELECT c.course_id, c.course_code, c.course_name, c.course_type, c.category, c.credit, 
		       c.lecture_hrs, c.tutorial_hrs, c.practical_hrs, c.activity_hrs, COALESCE(c.` + "`tw/sl`" + `, 0) as tw_sl,
		       COALESCE(c.theory_total_hrs, 0), COALESCE(c.tutorial_total_hrs, 0), COALESCE(c.practical_total_hrs, 0), COALESCE(c.activity_total_hrs, 0), COALESCE(c.total_hrs, 0),
		       c.cia_marks, c.see_marks, c.total_marks,
		       rc.id as reg_course_id
		FROM courses c
		INNER JOIN curriculum_courses rc ON c.course_id = rc.course_id
		WHERE rc.curriculum_id = ? AND rc.semester_id = ? AND rc.status = 1 AND c.status = 1
		ORDER BY c.course_code
	`

	rows, err := db.DB.Query(query, curriculumID, semesterID)
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
			&course.LectureHrs, &course.TutorialHrs, &course.PracticalHrs, &course.ActivityHrs, &course.TwSlHrs,
			&course.TheoryTotalHrs, &course.TutorialTotalHrs, &course.PracticalTotalHrs, &course.ActivityTotalHrs, &course.TotalHrs,
			&course.CIAMarks, &course.SEEMarks, &course.TotalMarks,
			&course.RegCourseID)
		if err != nil {
			log.Println("Error scanning course:", err)
			continue
		}
		course.CurriculumTemplate = curriculumTemplate
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
	curriculumID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid curriculum ID"})
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
	err = db.DB.QueryRow("SELECT max_credits FROM curriculum WHERE id = ? AND status = 1", curriculumID).Scan(&maxCredits)
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
	                WHERE rc.curriculum_id = ?`
	err = db.DB.QueryRow(creditQuery, curriculumID).Scan(&currentCredits)
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

	// Get curriculum template
	var curriculumTemplate string
	err = db.DB.QueryRow("SELECT curriculum_template FROM curriculum WHERE id = ? AND status = 1", curriculumID).Scan(&curriculumTemplate)
	if err != nil {
		log.Println("Error fetching curriculum template:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch curriculum template"})
		return
	}

	// Use total hours from frontend (already calculated)
	theoryTotal := course.TheoryTotalHrs
	tutorialTotal := course.TutorialTotalHrs
	practicalTotal := course.PracticalTotalHrs
	activityTotal := course.ActivityTotalHrs

	// Check if course already exists
	var existingCourseID int
	checkQuery := "SELECT course_id FROM courses WHERE course_code = ?"
	err = db.DB.QueryRow(checkQuery, course.CourseCode).Scan(&existingCourseID)

	var courseID int
	if err == sql.ErrNoRows {
		// Insert new course - calculate total hours (total_hrs and total_marks are GENERATED columns)
		insertCourseQuery := `INSERT INTO courses (course_code, course_name, course_type, category, credit, 
		                      lecture_hrs, tutorial_hrs, practical_hrs, activity_hrs, ` + "`tw/sl`" + `,
		                      theory_total_hrs, tutorial_total_hrs, practical_total_hrs, activity_total_hrs,
		                      cia_marks, see_marks) 
		                      VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		result, err := db.DB.Exec(insertCourseQuery, course.CourseCode, course.CourseName, course.CourseType, course.Category, course.Credit,
			course.LectureHrs, course.TutorialHrs, course.PracticalHrs, course.ActivityHrs, course.TwSlHrs,
			theoryTotal, tutorialTotal, practicalTotal, activityTotal,
			course.CIAMarks, course.SEEMarks)
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
		// Course code already exists - return error
		log.Printf("Course with code %s already exists (ID: %d)", course.CourseCode, existingCourseID)
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"error": "A course with this course code already exists. Please use a different course code."})
		return
	}

	// Link course to curriculum and semester
	linkQuery := "INSERT INTO curriculum_courses (curriculum_id, semester_id, course_id) VALUES (?, ?, ?)"
	result, err := db.DB.Exec(linkQuery, curriculumID, semesterID, courseID)
	if err != nil {
		log.Println("Error linking course:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to add course to semester"})
		return
	}

	regCourseID, _ := result.LastInsertId()
	course.CourseID = courseID

	// Fetch the complete course details including computed fields
	var fullCourse models.CourseWithDetails
	fetchQuery := `SELECT course_id, course_code, course_name, course_type, category, credit, 
	               lecture_hrs, tutorial_hrs, practical_hrs, activity_hrs, COALESCE(` + "`tw/sl`" + `, 0) as tw_sl,
	               COALESCE(theory_total_hrs, 0), COALESCE(tutorial_total_hrs, 0), COALESCE(practical_total_hrs, 0), COALESCE(activity_total_hrs, 0), COALESCE(total_hrs, 0),
	               cia_marks, see_marks, total_marks 
	               FROM courses WHERE course_id = ?`
	err = db.DB.QueryRow(fetchQuery, courseID).Scan(&fullCourse.CourseID, &fullCourse.CourseCode, &fullCourse.CourseName,
		&fullCourse.CourseType, &fullCourse.Category, &fullCourse.Credit,
		&fullCourse.LectureHrs, &fullCourse.TutorialHrs, &fullCourse.PracticalHrs, &fullCourse.ActivityHrs, &fullCourse.TwSlHrs,
		&fullCourse.TheoryTotalHrs, &fullCourse.TutorialTotalHrs, &fullCourse.PracticalTotalHrs, &fullCourse.ActivityTotalHrs, &fullCourse.TotalHrs,
		&fullCourse.CIAMarks, &fullCourse.SEEMarks, &fullCourse.TotalMarks)
	if err != nil {
		log.Println("Error fetching full course details:", err)
		// Fallback to sent values
		fullCourse = models.CourseWithDetails{
			CourseID:     courseID,
			CourseCode:   course.CourseCode,
			CourseName:   course.CourseName,
			CourseType:   course.CourseType,
			Category:     course.Category,
			Credit:       course.Credit,
			LectureHrs:   course.LectureHrs,
			TutorialHrs:  course.TutorialHrs,
			PracticalHrs: course.PracticalHrs,
			ActivityHrs:  course.ActivityHrs,
			CIAMarks:     course.CIAMarks,
			SEEMarks:     course.SEEMarks,
			TotalMarks:   course.CIAMarks + course.SEEMarks,
		}
	}
	fullCourse.RegCourseID = int(regCourseID)

	// Log the activity
	LogCurriculumActivity(curriculumID, "Course Added",
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
	curriculumID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid curriculum ID"})
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

	query := "UPDATE curriculum_courses SET status = 0 WHERE curriculum_id = ? AND semester_id = ? AND course_id = ? AND status = 1"
	result, err := db.DB.Exec(query, curriculumID, semesterID, courseID)
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
	LogCurriculumActivity(curriculumID, "Course Removed",
		"Removed course "+courseName+" from Semester "+strconv.Itoa(semesterID), "System")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Course removed successfully"})
}
