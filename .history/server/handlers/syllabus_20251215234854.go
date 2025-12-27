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
func GetCourseSyllabus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	courseID, err := strconv.Atoi(vars["courseId"])
	if err != nil {
		http.Error(w, "Invalid course ID", http.StatusBadRequest)
		return
	}

	var syllabus models.Syllabus
	var objectivesJSON, outcomesJSON, textbooksJSON, referenceListJSON, prerequisitesJSON, teamworkJSON, selflearningJSON []byte

	err = db.DB.QueryRow(`
		SELECT id, course_id, objectives, outcomes, textbooks, reference_list, prerequisites, 
		       COALESCE(teamwork, '{}'), COALESCE(selflearning, '{}')
		FROM course_syllabus
		WHERE course_id = ?
	`, courseID).Scan(
		&syllabus.ID,
		&syllabus.CourseID,
		&objectivesJSON,
		&outcomesJSON,
		&textbooksJSON,
		&referenceListJSON,
		&prerequisitesJSON,
		&teamworkJSON,
		&selflearningJSON,
	)

	if err == sql.ErrNoRows {
		// Return empty/default values if no syllabus exists
		syllabus = models.Syllabus{
			CourseID:      courseID,
			Objectives:    []string{},
			Outcomes:      []string{},
			Textbooks:     []string{},
			ReferenceList: []string{},
			Prerequisites: []string{},
		}
		json.NewEncoder(w).Encode(syllabus)
		return
	}

	if err != nil {
		log.Println("Error fetching syllabus:", err)
		http.Error(w, "Failed to fetch syllabus", http.StatusInternalServerError)
		return
	}

	// Parse JSON arrays
	if err := json.Unmarshal(objectivesJSON, &syllabus.Objectives); err != nil {
		syllabus.Objectives = []string{}
	}
	if err := json.Unmarshal(outcomesJSON, &syllabus.Outcomes); err != nil {
		syllabus.Outcomes = []string{}
	}
	if err := json.Unmarshal(textbooksJSON, &syllabus.Textbooks); err != nil {
		syllabus.Textbooks = []string{}
	}
	if err := json.Unmarshal(referenceListJSON, &syllabus.ReferenceList); err != nil {
		syllabus.ReferenceList = []string{}
	}
	if err := json.Unmarshal(prerequisitesJSON, &syllabus.Prerequisites); err != nil {
		syllabus.Prerequisites = []string{}
	}
	if len(teamworkJSON) > 0 && string(teamworkJSON) != "{}" {
		var teamwork models.Teamwork
		if err := json.Unmarshal(teamworkJSON, &teamwork); err == nil {
			syllabus.Teamwork = &teamwork
		}
	}
	if len(selflearningJSON) > 0 && string(selflearningJSON) != "{}" {
		var selflearning models.SelfLearning
		if err := json.Unmarshal(selflearningJSON, &selflearning); err == nil {
			syllabus.SelfLearning = &selflearning
		}
	}

	json.NewEncoder(w).Encode(syllabus)
}

// SaveCourseSyllabus handles POST /course/:courseId/syllabus
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

	// Convert arrays to JSON (header only)
	objectivesJSON, _ := json.Marshal(syllabus.Objectives)
	outcomesJSON, _ := json.Marshal(syllabus.Outcomes)
	textbooksJSON, _ := json.Marshal(syllabus.Textbooks)
	referenceListJSON, _ := json.Marshal(syllabus.ReferenceList)
	prerequisitesJSON, _ := json.Marshal(syllabus.Prerequisites)
	teamworkJSON, _ := json.Marshal(syllabus.Teamwork)
	selflearningJSON, _ := json.Marshal(syllabus.SelfLearning)

	// Column existence is handled during startup migration; avoid runtime ALTERs for broader MySQL compatibility

	// Check if syllabus exists
	var existingID int
	err = db.DB.QueryRow("SELECT id FROM course_syllabus WHERE course_id = ?", courseID).Scan(&existingID)

	if err == sql.ErrNoRows {
		// Insert new syllabus (header only)
		result, err := db.DB.Exec(`
			INSERT INTO course_syllabus (course_id, objectives, outcomes, textbooks, reference_list, prerequisites, teamwork, selflearning)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`, courseID, objectivesJSON, outcomesJSON, textbooksJSON, referenceListJSON, prerequisitesJSON, teamworkJSON, selflearningJSON)

		if err != nil {
			log.Println("Error inserting syllabus:", err)
			http.Error(w, "Failed to create syllabus", http.StatusInternalServerError)
			return
		}

		id, _ := result.LastInsertId()
		syllabus.ID = int(id)
	} else if err != nil {
		log.Println("Error checking existing syllabus:", err)
		http.Error(w, "Failed to save syllabus", http.StatusInternalServerError)
		return
	} else {
		// Update existing syllabus (header only)
		_, err := db.DB.Exec(`
			UPDATE course_syllabus 
			SET objectives = ?, outcomes = ?, textbooks = ?, reference_list = ?, prerequisites = ?, teamwork = ?, selflearning = ?
			WHERE course_id = ?
		`, objectivesJSON, outcomesJSON, textbooksJSON, referenceListJSON, prerequisitesJSON, teamworkJSON, selflearningJSON, courseID)

		if err != nil {
			log.Println("Error updating syllabus:", err)
			http.Error(w, "Failed to update syllabus", http.StatusInternalServerError)
			return
		}

		syllabus.ID = existingID
	}

	json.NewEncoder(w).Encode(syllabus)
}
