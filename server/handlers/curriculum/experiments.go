package curriculum

import (
	"encoding/json"
	"log"
	"net/http"
	"server/db"
	"server/models"
	"strconv"

	"github.com/gorilla/mux"
)

// GetCourseExperiments returns all experiments for a course (2022 template)
func GetCourseExperiments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	courseID, err := strconv.Atoi(vars["courseId"])
	if err != nil {
		http.Error(w, "Invalid course ID", http.StatusBadRequest)
		return
	}

	experiments := []models.Experiment{}
	rows, err := db.DB.Query(`
		SELECT id, course_id, experiment_number, experiment_name, hours
		FROM course_experiments
		WHERE course_id = ?
		ORDER BY experiment_number`, courseID)
	if err != nil {
		log.Println("Error fetching experiments:", err)
		http.Error(w, "Failed to fetch experiments", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var exp models.Experiment
		if err := rows.Scan(&exp.ID, &exp.CourseID, &exp.ExperimentNumber, &exp.ExperimentName, &exp.Hours); err != nil {
			continue
		}

		// Fetch topics for this experiment
		topicRows, err := db.DB.Query(`
			SELECT topic_text
			FROM course_experiment_topics
			WHERE experiment_id = ?
			ORDER BY topic_order`, exp.ID)
		if err != nil {
			exp.Topics = []string{}
		} else {
			topics := []string{}
			for topicRows.Next() {
				var topic string
				if err := topicRows.Scan(&topic); err == nil {
					topics = append(topics, topic)
				}
			}
			topicRows.Close()
			exp.Topics = topics
		}

		experiments = append(experiments, exp)
	}

	json.NewEncoder(w).Encode(experiments)
}

// CreateExperiment creates a new experiment for a course
func CreateExperiment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	courseID, err := strconv.Atoi(vars["courseId"])
	if err != nil {
		http.Error(w, "Invalid course ID", http.StatusBadRequest)
		return
	}

	var req struct {
		ExperimentNumber int      `json:"experiment_number"`
		ExperimentName   string   `json:"experiment_name"`
		Hours            int      `json:"hours"`
		Topics           []string `json:"topics"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := db.DB.Exec(`
		INSERT INTO course_experiments (course_id, experiment_number, experiment_name, hours)
		VALUES (?, ?, ?, ?)`, courseID, req.ExperimentNumber, req.ExperimentName, req.Hours)
	if err != nil {
		log.Println("Error creating experiment:", err)
		http.Error(w, "Failed to create experiment", http.StatusInternalServerError)
		return
	}

	expID, _ := result.LastInsertId()

	// Insert topics
	for i, topic := range req.Topics {
		_, _ = db.DB.Exec(`
			INSERT INTO course_experiment_topics (experiment_id, topic_text, topic_order)
			VALUES (?, ?, ?)`, expID, topic, i)
	}

	json.NewEncoder(w).Encode(map[string]int64{"id": expID})
}

// UpdateExperiment updates an experiment and its topics
func UpdateExperiment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	expID, err := strconv.Atoi(vars["expId"])
	if err != nil {
		http.Error(w, "Invalid experiment ID", http.StatusBadRequest)
		return
	}

	var req struct {
		ExperimentNumber int      `json:"experiment_number"`
		ExperimentName   string   `json:"experiment_name"`
		Hours            int      `json:"hours"`
		Topics           []string `json:"topics"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec(`
		UPDATE course_experiments
		SET experiment_number = ?, experiment_name = ?, hours = ?
		WHERE id = ?`, req.ExperimentNumber, req.ExperimentName, req.Hours, expID)
	if err != nil {
		log.Println("Error updating experiment:", err)
		http.Error(w, "Failed to update experiment", http.StatusInternalServerError)
		return
	}

	// Replace topics
	_, _ = db.DB.Exec("DELETE FROM course_experiment_topics WHERE experiment_id = ?", expID)
	for i, topic := range req.Topics {
		_, _ = db.DB.Exec(`
			INSERT INTO course_experiment_topics (experiment_id, topic_text, topic_order)
			VALUES (?, ?, ?)`, expID, topic, i)
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteExperiment deletes an experiment and its topics
func DeleteExperiment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	expID, err := strconv.Atoi(vars["expId"])
	if err != nil {
		http.Error(w, "Invalid experiment ID", http.StatusBadRequest)
		return
	}

	// Get course_id before deleting
	var courseID int
	db.DB.QueryRow("SELECT course_id FROM course_experiments WHERE id = ?", expID).Scan(&courseID)

	_, err = db.DB.Exec("DELETE FROM course_experiments WHERE id = ?", expID)
	if err != nil {
		log.Println("Error deleting experiment:", err)
		http.Error(w, "Failed to delete experiment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
