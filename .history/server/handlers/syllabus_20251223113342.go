package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"server/db"
	"server/models"
	"strconv"

	"github.com/gorilla/mux"
)

// GetCourseSyllabus handles GET /course/:courseId/syllabus
// Fetches data from normalized tables and returns in the same JSON format as before
func GetCourseSyllabus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	courseID, err := strconv.Atoi(vars["courseId"])
	if err != nil {
		http.Error(w, "Invalid course ID", http.StatusBadRequest)
		return
	}

	// Ensure syllabus record exists
	syllabusID, err := ensureSyllabusRecord(courseID)
	if err != nil {
		log.Println("Error ensuring syllabus record:", err)
		http.Error(w, "Failed to fetch syllabus", http.StatusInternalServerError)
		return
	}

	// Fetch all data from normalized tables
	syllabus := models.Syllabus{
		ID:       syllabusID,
		CourseID: courseID,
	}

	syllabus.Objectives, _ = fetchObjectives(courseID)
	syllabus.Outcomes, _ = fetchOutcomes(courseID)
	syllabus.ReferenceList, _ = fetchReferences(courseID)
	syllabus.Prerequisites, _ = fetchPrerequisites(courseID)
	syllabus.Teamwork, _ = fetchTeamwork(courseID)
	syllabus.SelfLearning, _ = fetchSelfLearning(courseID)

	json.NewEncoder(w).Encode(syllabus)
}

// SaveCourseSyllabus handles POST /course/:courseId/syllabus
// Saves data to normalized tables while maintaining the same API interface
func SaveCourseSyllabus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	courseID, err := strconv.Atoi(vars["courseId"])
	if err != nil {
		http.Error(w, "Invalid course ID", http.StatusBadRequest)
		return
	}

	var syllabus models.Syllabus
	if err := json.NewDecoder(r.Body).Decode(&syllabus); err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	syllabus.CourseID = courseID

	// Ensure syllabus record exists
	syllabusID, err := ensureSyllabusRecord(courseID)
	if err != nil {
		log.Println("Error ensuring syllabus record:", err)
		http.Error(w, "Failed to save syllabus", http.StatusInternalServerError)
		return
	}

	syllabus.ID = syllabusID

	// Save all data to normalized tables
	if err := saveObjectives(courseID, syllabus.Objectives); err != nil {
		log.Println("Error saving objectives:", err)
		http.Error(w, "Failed to save objectives", http.StatusInternalServerError)
		return
	}

	if err := saveOutcomes(courseID, syllabus.Outcomes); err != nil {
		log.Println("Error saving outcomes:", err)
		http.Error(w, "Failed to save outcomes", http.StatusInternalServerError)
		return
	}

	if err := saveReferences(courseID, syllabus.ReferenceList); err != nil {
		log.Println("Error saving references:", err)
		http.Error(w, "Failed to save references", http.StatusInternalServerError)
		return
	}

	if err := savePrerequisites(courseID, syllabus.Prerequisites); err != nil {
		log.Println("Error saving prerequisites:", err)
		http.Error(w, "Failed to save prerequisites", http.StatusInternalServerError)
		return
	}

	if err := saveTeamwork(courseID, syllabus.Teamwork); err != nil {
		log.Println("Error saving teamwork:", err)
		http.Error(w, "Failed to save teamwork", http.StatusInternalServerError)
		return
	}

	if err := saveSelfLearning(courseID, syllabus.SelfLearning); err != nil {
		log.Println("Error saving self-learning:", err)
		http.Error(w, "Failed to save self-learning", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(syllabus)
}
