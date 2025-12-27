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

// UpdateCurriculum updates curriculum name and max_credits
func UpdateCurriculum(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPut {
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

	var updateData struct {
		Name       string `json:"name"`
		MaxCredits int    `json:"max_credits"`
	}

	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		log.Println("Error decoding request body:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Fetch old values for diff
	var oldName string
	var oldMaxCredits int
	err = db.DB.QueryRow("SELECT name, max_credits FROM curriculum WHERE id = ?", curriculumID).Scan(&oldName, &oldMaxCredits)
	if err != nil {
		log.Println("Error fetching old curriculum data:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch curriculum data"})
		return
	}

	// Update curriculum
	query := "UPDATE curriculum SET name = ?, max_credits = ? WHERE id = ?"
	_, err = db.DB.Exec(query, updateData.Name, updateData.MaxCredits, curriculumID)
	if err != nil {
		log.Println("Error updating curriculum:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update curriculum"})
		return
	}

	// Generate diff and log changes
	diff := make(map[string]map[string]interface{})
	if oldName != updateData.Name {
		diff["name"] = map[string]interface{}{"old": oldName, "new": updateData.Name}
	}
	if oldMaxCredits != updateData.MaxCredits {
		diff["max_credits"] = map[string]interface{}{"old": oldMaxCredits, "new": updateData.MaxCredits}
	}

	if len(diff) > 0 {
		LogCurriculumActivityWithDiff(curriculumID, "Curriculum Updated",
			fmt.Sprintf("Updated curriculum details"), "System", diff)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Curriculum updated successfully"})
}

// UpdateSemester updates semester name/number
func UpdateSemester(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPut {
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

	var updateData struct {
		SemesterNumber int `json:"semester_number"`
	}

	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		log.Println("Error decoding request body:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Fetch old values and regulation_id for diff
	var oldSemesterNumber int
	var regulationID int
	err = db.DB.QueryRow("SELECT semester_number, regulation_id FROM semesters WHERE id = ?", semesterID).Scan(&oldSemesterNumber, &regulationID)
	if err != nil {
		log.Println("Error fetching old semester data:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch semester data"})
		return
	}

	// Update semester
	query := "UPDATE semesters SET semester_number = ? WHERE id = ?"
	_, err = db.DB.Exec(query, updateData.SemesterNumber, semesterID)
	if err != nil {
		log.Println("Error updating semester:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update semester"})
		return
	}

	// Generate diff and log changes
	if oldSemesterNumber != updateData.SemesterNumber {
		diff := map[string]map[string]interface{}{
			"semester_number": {"old": oldSemesterNumber, "new": updateData.SemesterNumber},
		}
		LogCurriculumActivityWithDiff(regulationID, "Semester Updated",
			fmt.Sprintf("Updated Semester %d to Semester %d", oldSemesterNumber, updateData.SemesterNumber), "System", diff)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Semester updated successfully"})
}

// UpdateCourse updates course details
func UpdateCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	vars := mux.Vars(r)
	courseID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid course ID"})
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

	// Fetch old course data for diff
	var oldCourse models.Course
	err = db.DB.QueryRow(`SELECT course_code, course_name, course_type, category, credit, 
		lecture_hours, tutorial_hours, practical_hours, cia_marks, see_marks 
		FROM courses WHERE id = ?`, courseID).Scan(
		&oldCourse.CourseCode, &oldCourse.CourseName, &oldCourse.CourseType, &oldCourse.Category,
		&oldCourse.Credit, &oldCourse.LectureHours, &oldCourse.TutorialHours, &oldCourse.PracticalHours,
		&oldCourse.CIAMarks, &oldCourse.SEEMarks)
	if err != nil {
		log.Println("Error fetching old course data:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch course data"})
		return
	}

	// Get regulation_id from curriculum_courses
	var regulationID int
	err = db.DB.QueryRow("SELECT regulation_id FROM curriculum_courses WHERE course_id = ? LIMIT 1", courseID).Scan(&regulationID)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Error fetching regulation_id:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch regulation ID"})
		return
	}

	// Update course (do not set generated column total_marks)
	course.TotalMarks = course.CIAMarks + course.SEEMarks
	query := `UPDATE courses SET course_code = ?, course_name = ?, course_type = ?, category = ?, 
		credit = ?, lecture_hours = ?, tutorial_hours = ?, practical_hours = ?, 
		cia_marks = ?, see_marks = ? WHERE id = ?`

	_, err = db.DB.Exec(query, course.CourseCode, course.CourseName, course.CourseType, course.Category,
		course.Credit, course.LectureHours, course.TutorialHours, course.PracticalHours,
		course.CIAMarks, course.SEEMarks, courseID)
	if err != nil {
		log.Println("Error updating course:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update course"})
		return
	}

	// Generate diff and log changes
	diff := make(map[string]map[string]interface{})
	if oldCourse.CourseCode != course.CourseCode {
		diff["course_code"] = map[string]interface{}{"old": oldCourse.CourseCode, "new": course.CourseCode}
	}
	if oldCourse.CourseName != course.CourseName {
		diff["course_name"] = map[string]interface{}{"old": oldCourse.CourseName, "new": course.CourseName}
	}
	if oldCourse.CourseType != course.CourseType {
		diff["course_type"] = map[string]interface{}{"old": oldCourse.CourseType, "new": course.CourseType}
	}
	if oldCourse.Category != course.Category {
		diff["category"] = map[string]interface{}{"old": oldCourse.Category, "new": course.Category}
	}
	if oldCourse.Credit != course.Credit {
		diff["credit"] = map[string]interface{}{"old": oldCourse.Credit, "new": course.Credit}
	}
	if oldCourse.LectureHours != course.LectureHours {
		diff["lecture_hours"] = map[string]interface{}{"old": oldCourse.LectureHours, "new": course.LectureHours}
	}
	if oldCourse.TutorialHours != course.TutorialHours {
		diff["tutorial_hours"] = map[string]interface{}{"old": oldCourse.TutorialHours, "new": course.TutorialHours}
	}
	if oldCourse.PracticalHours != course.PracticalHours {
		diff["practical_hours"] = map[string]interface{}{"old": oldCourse.PracticalHours, "new": course.PracticalHours}
	}
	// Computed total_hours change for logging
	if (oldCourse.LectureHours+oldCourse.TutorialHours+oldCourse.PracticalHours) != (course.LectureHours+course.TutorialHours+course.PracticalHours) {
		diff["total_hours"] = map[string]interface{}{"old": oldCourse.LectureHours + oldCourse.TutorialHours + oldCourse.PracticalHours, "new": course.LectureHours + course.TutorialHours + course.PracticalHours}
	}
	if oldCourse.CIAMarks != course.CIAMarks {
		diff["cia_marks"] = map[string]interface{}{"old": oldCourse.CIAMarks, "new": course.CIAMarks}
	}
	if oldCourse.SEEMarks != course.SEEMarks {
		diff["see_marks"] = map[string]interface{}{"old": oldCourse.SEEMarks, "new": course.SEEMarks}
	}

	if len(diff) > 0 && regulationID > 0 {
		LogCurriculumActivityWithDiff(regulationID, "Course Updated",
			fmt.Sprintf("Updated course: %s - %s", course.CourseCode, course.CourseName), "System", diff)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Course updated successfully"})
}

// GetCourse retrieves a single course by ID
func GetCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	courseID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid course ID"})
		return
	}

	var course models.Course
	err = db.DB.QueryRow(`
		SELECT id, course_code, course_name, course_type, category, credit, 
		       lecture_hours, tutorial_hours, practical_hours,
		       (lecture_hours + tutorial_hours + practical_hours) AS total_hours,
		       cia_marks, see_marks, total_marks
		FROM courses 
		WHERE id = ?`, courseID).
		Scan(&course.ID, &course.CourseCode, &course.CourseName, &course.CourseType,
			&course.Category, &course.Credit, &course.LectureHours, &course.TutorialHours,
			&course.PracticalHours, &course.TotalHours, &course.CIAMarks, &course.SEEMarks, &course.TotalMarks)

	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Course not found"})
		return
	} else if err != nil {
		log.Println("Error fetching course:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch course"})
		return
	}

	json.NewEncoder(w).Encode(course)
}
