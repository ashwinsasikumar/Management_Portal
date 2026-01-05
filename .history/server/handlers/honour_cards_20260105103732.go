package handlers

import (
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
	regulationID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid regulation ID"})
		return
	}

	query := "SELECT id, regulation_id, title, semester_number FROM honour_cards WHERE regulation_id = ? ORDER BY semester_number"
	rows, err := db.DB.Query(query, regulationID)
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
		err := rows.Scan(&card.ID, &card.RegulationID, &card.Title, &card.SemesterNumber)
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
	query := "SELECT id, honour_card_id, name FROM honour_verticals WHERE honour_card_id = ? ORDER BY id"
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
		       c.credit, c.theory_hours, c.activity_hours, c.lecture_hours, c.tutorial_hours, 
		       c.practical_hours, c.total_hours, c.cia_marks, c.see_marks, c.total_marks,
		       hvc.id as honour_vertical_course_id
		FROM courses c
		INNER JOIN honour_vertical_courses hvc ON c.course_id = hvc.course_id
		WHERE hvc.honour_vertical_id = ?
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
			&course.Category, &course.Credit, &course.TheoryHours, &course.ActivityHours,
			&course.LectureHours, &course.TutorialHours, &course.PracticalHours,
			&course.TotalHours, &course.CIAMarks, &course.SEEMarks, &course.TotalMarks,
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
	regulationID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid regulation ID"})
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

	card.RegulationID = regulationID

	query := "INSERT INTO honour_cards (regulation_id, title, semester_number) VALUES (?, ?, ?)"
	result, err := db.DB.Exec(query, card.RegulationID, card.Title, card.SemesterNumber)
	if err != nil {
		log.Println("Error inserting honour card:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create honour card"})
		return
	}

	id, _ := result.LastInsertId()
	card.ID = int(id)

	// Log the activity
	LogCurriculumActivity(regulationID, "Honour Card Added",
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

	var reqBody struct {
		CourseID int `json:"course_id"`
	}
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		log.Println("Error decoding request body:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Check if course exists
	var exists bool
	checkQuery := "SELECT EXISTS(SELECT 1 FROM courses WHERE course_id = ?)"
	err = db.DB.QueryRow(checkQuery, reqBody.CourseID).Scan(&exists)
	if err != nil || !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Course not found"})
		return
	}

	query := "INSERT INTO honour_vertical_courses (honour_vertical_id, course_id) VALUES (?, ?)"
	result, err := db.DB.Exec(query, verticalID, reqBody.CourseID)
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
		"course_id":          reqBody.CourseID,
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

	query := "DELETE FROM honour_vertical_courses WHERE honour_vertical_id = ? AND course_id = ?"
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

	query := "DELETE FROM honour_verticals WHERE id = ?"
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

	query := "DELETE FROM honour_cards WHERE id = ?"
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
