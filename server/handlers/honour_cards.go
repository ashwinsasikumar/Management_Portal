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

// GetHonourCards retrieves all honour cards for a regulation
func GetHonourCards(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	curriculumID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid curriculum ID"})
		return
	}

	query := "SELECT id, curriculum_id, title FROM honour_cards WHERE curriculum_id = ? AND status = 1 ORDER BY id"
	rows, err := db.DB.Query(query, curriculumID)
	if err != nil {
		log.Println("Error querying honour cards:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch honour cards"})
		return
	}
	defer rows.Close()

	var honourCards []models.HonourCardWithVerticals = make([]models.HonourCardWithVerticals, 0)
	for rows.Next() {
		var card models.HonourCardWithVerticals
		err := rows.Scan(&card.ID, &card.CurriculumID, &card.Title)
		if err != nil {
			log.Println("Error scanning honour card:", err)
			continue
		}

		// Fetch verticals for this honour card
		card.Verticals = fetchVerticalsForCard(card.ID)
		honourCards = append(honourCards, card)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(honourCards)
}

// fetchVerticalsForCard retrieves all verticals and their courses for a given honour card
func fetchVerticalsForCard(honourCardID int) []models.HonourVerticalWithCourses {
	query := "SELECT id, honour_card_id, name FROM honour_verticals WHERE honour_card_id = ? AND status = 1 ORDER BY id"
	rows, err := db.DB.Query(query, honourCardID)
	if err != nil {
		log.Println("Error querying verticals:", err)
		return []models.HonourVerticalWithCourses{}
	}
	defer rows.Close()

	var verticals []models.HonourVerticalWithCourses = make([]models.HonourVerticalWithCourses, 0)
	for rows.Next() {
		var vertical models.HonourVerticalWithCourses
		err := rows.Scan(&vertical.ID, &vertical.HonourCardID, &vertical.Name)
		if err != nil {
			log.Println("Error scanning vertical:", err)
			continue
		}

		// Fetch courses for this vertical
		vertical.Courses = fetchCoursesForVertical(vertical.ID)
		verticals = append(verticals, vertical)
	}

	return verticals
}

// fetchCoursesForVertical retrieves all courses for a given vertical
func fetchCoursesForVertical(verticalID int) []models.CourseWithDetails {
	query := `
		SELECT c.course_id, c.course_code, c.course_name, c.course_type, c.category, 
		       c.credit, c.lecture_hrs, c.tutorial_hrs, c.practical_hrs, c.activity_hrs, COALESCE(c.` + "`tw/sl`" + `, 0) as tw_sl,
		       COALESCE(c.theory_total_hrs, 0), COALESCE(c.tutorial_total_hrs, 0), COALESCE(c.practical_total_hrs, 0), COALESCE(c.activity_total_hrs, 0), COALESCE(c.total_hrs, 0),
		       c.cia_marks, c.see_marks, c.total_marks,
		       hvc.id as honour_vertical_course_id
		FROM courses c
		INNER JOIN honour_vertical_courses hvc ON c.course_id = hvc.course_id
		WHERE hvc.honour_vertical_id = ? AND hvc.status = 1 AND c.status = 1
		ORDER BY c.course_code
	`
	rows, err := db.DB.Query(query, verticalID)
	if err != nil {
		log.Println("Error querying courses for vertical:", err)
		return []models.CourseWithDetails{}
	}
	defer rows.Close()

	var courses []models.CourseWithDetails = make([]models.CourseWithDetails, 0)
	for rows.Next() {
		var course models.CourseWithDetails
		err := rows.Scan(
			&course.CourseID, &course.CourseCode, &course.CourseName, &course.CourseType,
			&course.Category, &course.Credit, &course.LectureHrs, &course.TutorialHrs, &course.PracticalHrs, &course.ActivityHrs, &course.TwSlHrs,
			&course.TheoryTotalHrs, &course.TutorialTotalHrs, &course.PracticalTotalHrs, &course.ActivityTotalHrs, &course.TotalHrs,
			&course.CIAMarks, &course.SEEMarks, &course.TotalMarks,
			&course.RegCourseID,
		)
		if err != nil {
			log.Println("Error scanning course:", err)
			continue
		}
		courses = append(courses, course)
	}

	return courses
}

// CreateHonourCard creates a new honour card for a regulation
func CreateHonourCard(w http.ResponseWriter, r *http.Request) {
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

	var card models.HonourCard
	err = json.NewDecoder(r.Body).Decode(&card)
	if err != nil {
		log.Println("Error decoding request body:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	card.CurriculumID = curriculumID

	query := "INSERT INTO honour_cards (curriculum_id, title) VALUES (?, ?)"
	result, err := db.DB.Exec(query, card.CurriculumID, card.Title)
	if err != nil {
		log.Println("Error inserting honour card:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create honour card"})
		return
	}

	id, _ := result.LastInsertId()
	card.ID = int(id)

	// Log the activity
	LogCurriculumActivity(curriculumID, "Honour Card Added",
		"Added Honour Card: "+card.Title, "System")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(card)
}

// CreateHonourVertical creates a new vertical within an honour card
func CreateHonourVertical(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	vars := mux.Vars(r)
	honourCardID, err := strconv.Atoi(vars["cardId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid honour card ID"})
		return
	}

	var vertical models.HonourVertical
	err = json.NewDecoder(r.Body).Decode(&vertical)
	if err != nil {
		log.Println("Error decoding request body:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	vertical.HonourCardID = honourCardID

	query := "INSERT INTO honour_verticals (honour_card_id, name) VALUES (?, ?)"
	result, err := db.DB.Exec(query, vertical.HonourCardID, vertical.Name)
	if err != nil {
		log.Println("Error inserting vertical:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create vertical"})
		return
	}

	id, _ := result.LastInsertId()
	vertical.ID = int(id)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(vertical)
}

// AddCourseToVertical adds a course to a vertical
func AddCourseToVertical(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	vars := mux.Vars(r)
	verticalID, err := strconv.Atoi(vars["verticalId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid vertical ID"})
		return
	}

	// Support two ways of adding a course:
	// 1) By existing course_id (legacy behaviour)
	// 2) By full course details (same as normal card flow)
	var payload struct {
		CourseID          *int   `json:"course_id,omitempty"`
		CourseCode        string `json:"course_code,omitempty"`
		CourseName        string `json:"course_name,omitempty"`
		CourseType        string `json:"course_type,omitempty"`
		Category          string `json:"category,omitempty"`
		Credit            int    `json:"credit,omitempty"`
		LectureHrs        int    `json:"lecture_hrs,omitempty"`
		TutorialHrs       int    `json:"tutorial_hrs,omitempty"`
		PracticalHrs      int    `json:"practical_hrs,omitempty"`
		ActivityHrs       int    `json:"activity_hrs,omitempty"`
		TwSlHrs           int    `json:"tw_sl_hrs,omitempty"`
		TheoryTotalHrs    int    `json:"theory_total_hrs,omitempty"`
		TutorialTotalHrs  int    `json:"tutorial_total_hrs,omitempty"`
		PracticalTotalHrs int    `json:"practical_total_hrs,omitempty"`
		ActivityTotalHrs  int    `json:"activity_total_hrs,omitempty"`
		CIAMarks          int    `json:"cia_marks,omitempty"`
		SEEMarks          int    `json:"see_marks,omitempty"`
	}

	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Println("Error decoding request body:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Get curriculum ID and template from vertical
	var curriculumID int
	var curriculumTemplate string
	err = db.DB.QueryRow(`
		SELECT hc.curriculum_id, c.curriculum_template 
		FROM honour_verticals hv
		INNER JOIN honour_cards hc ON hv.honour_card_id = hc.id
		INNER JOIN curriculum c ON hc.curriculum_id = c.id
		WHERE hv.id = ?`, verticalID).Scan(&curriculumID, &curriculumTemplate)
	if err != nil {
		log.Println("Error fetching curriculum template for vertical:", err)
		// Don't fail, just use default calculation
		curriculumTemplate = "2022"
	}

	var courseID int

	if payload.CourseID != nil && *payload.CourseID > 0 {
		// Legacy path: link an existing course by ID
		var exists bool
		checkQuery := "SELECT EXISTS(SELECT 1 FROM courses WHERE course_id = ?)"
		err = db.DB.QueryRow(checkQuery, *payload.CourseID).Scan(&exists)
		if err != nil || !exists {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Course not found"})
			return
		}
		courseID = *payload.CourseID
	} else {
		// New path: create or reuse a course based on course_code (similar to AddCourseToSemester)
		if payload.CourseCode == "" || payload.CourseName == "" || payload.CourseType == "" || payload.Category == "" || payload.Credit < 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Missing required course fields"})
			return
		}

		// Use total hours from frontend (already calculated)
		theoryTotal := payload.TheoryTotalHrs
		tutorialTotal := payload.TutorialTotalHrs
		practicalTotal := payload.PracticalTotalHrs
		activityTotal := payload.ActivityTotalHrs

		// Check if course already exists by course_code
		var existingCourseID int
		checkQuery := "SELECT course_id FROM courses WHERE course_code = ?"
		err = db.DB.QueryRow(checkQuery, payload.CourseCode).Scan(&existingCourseID)
		if err == sql.ErrNoRows {
			// Insert new course
			insertCourseQuery := `INSERT INTO courses (course_code, course_name, course_type, category, credit,
				lecture_hrs, tutorial_hrs, practical_hrs, activity_hrs, ` + "`tw/sl`" + `,
				theory_total_hrs, tutorial_total_hrs, practical_total_hrs, activity_total_hrs,
				cia_marks, see_marks)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
			result, err := db.DB.Exec(insertCourseQuery,
				payload.CourseCode,
				payload.CourseName,
				payload.CourseType,
				payload.Category,
				payload.Credit,
				payload.LectureHrs,
				payload.TutorialHrs,
				payload.PracticalHrs,
				payload.ActivityHrs,
				payload.TwSlHrs,
				theoryTotal,
				tutorialTotal,
				practicalTotal,
				activityTotal,
				payload.CIAMarks,
				payload.SEEMarks,
			)
			if err != nil {
				log.Println("Error inserting course for honour vertical:", err)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create course"})
				return
			}
			id, _ := result.LastInsertId()
			courseID = int(id)
		} else if err != nil {
			log.Println("Error checking existing course for honour vertical:", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create course"})
			return
		} else {
			courseID = existingCourseID
		}
	}

	query := "INSERT INTO honour_vertical_courses (honour_vertical_id, course_id) VALUES (?, ?)"
	result, err := db.DB.Exec(query, verticalID, courseID)
	if err != nil {
		log.Println("Error adding course to vertical:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to add course to vertical"})
		return
	}

	id, _ := result.LastInsertId()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":                 id,
		"honour_vertical_id": verticalID,
		"course_id":          courseID,
	})
}

// RemoveCourseFromVertical removes a course from a vertical
func RemoveCourseFromVertical(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	vars := mux.Vars(r)
	verticalID, err := strconv.Atoi(vars["verticalId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid vertical ID"})
		return
	}

	courseID, err := strconv.Atoi(vars["courseId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid course ID"})
		return
	}

	query := "UPDATE honour_vertical_courses SET status = 0 WHERE honour_vertical_id = ? AND course_id = ? AND status = 1"
	result, err := db.DB.Exec(query, verticalID, courseID)
	if err != nil {
		log.Println("Error removing course from vertical:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to remove course"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Course not found in vertical"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Course removed successfully"})
}

// DeleteHonourVertical deletes a vertical
func DeleteHonourVertical(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	vars := mux.Vars(r)
	verticalID, err := strconv.Atoi(vars["verticalId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid vertical ID"})
		return
	}

	query := "UPDATE honour_verticals SET status = 0 WHERE id = ? AND status = 1"
	result, err := db.DB.Exec(query, verticalID)
	if err != nil {
		log.Println("Error deleting vertical:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete vertical"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Vertical not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Vertical deleted successfully"})
}

// DeleteHonourCard deletes an honour card and all its verticals
func DeleteHonourCard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	vars := mux.Vars(r)
	cardID, err := strconv.Atoi(vars["cardId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid card ID"})
		return
	}

	query := "UPDATE honour_cards SET status = 0 WHERE id = ? AND status = 1"
	result, err := db.DB.Exec(query, cardID)
	if err != nil {
		log.Println("Error deleting honour card:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete honour card"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Honour card not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Honour card deleted successfully"})
}
