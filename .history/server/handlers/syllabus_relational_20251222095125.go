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

// GetCourseSyllabusNested returns the complete syllabus structure:
// Header (course_syllabus) + Models (syllabus_models) → Titles (syllabus_titles) → Topics (syllabus_topics)
func GetCourseSyllabusNested(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	courseID, err := strconv.Atoi(vars["courseId"])
	if err != nil {
		http.Error(w, "Invalid course ID", http.StatusBadRequest)
		return
	}

	var resp models.CourseSyllabusResponse

	// 1. Get header from course_syllabus
	var objectivesJSON, outcomesJSON, referenceJSON, prereqJSON, teamworkJSON, selflearningJSON []byte
	err = db.DB.QueryRow(`
		 SELECT id, course_id, 
			 COALESCE(objectives, JSON_ARRAY()), 
			 COALESCE(outcomes, JSON_ARRAY()), 
			 COALESCE(reference_list, JSON_ARRAY()), 
			 COALESCE(prerequisites, JSON_ARRAY()),
			 COALESCE(teamwork, '{}'),
			 COALESCE(selflearning, '{}')
		FROM course_syllabus 
		WHERE course_id = ?`, courseID).
		Scan(&resp.Header.ID, &resp.Header.CourseID,
			&objectivesJSON, &outcomesJSON, &referenceJSON, &prereqJSON, &teamworkJSON, &selflearningJSON)

	if err == sql.ErrNoRows {
		// No syllabus exists yet, return empty structure
		resp.Header.CourseID = courseID
		resp.Header.Objectives = []string{}
		resp.Header.Outcomes = []string{}
		resp.Header.ReferenceList = []string{}
		resp.Header.Prerequisites = []string{}
		resp.Models = []models.SyllabusModel{}
		json.NewEncoder(w).Encode(resp)
		return
	} else if err != nil {
		log.Println("Error fetching syllabus header:", err)
		http.Error(w, "Failed to fetch syllabus", http.StatusInternalServerError)
		return
	}

	// Parse JSON arrays
	json.Unmarshal(objectivesJSON, &resp.Header.Objectives)
	json.Unmarshal(outcomesJSON, &resp.Header.Outcomes)
	json.Unmarshal(referenceJSON, &resp.Header.ReferenceList)
	json.Unmarshal(prereqJSON, &resp.Header.Prerequisites)

	// Parse teamwork and selflearning
	if len(teamworkJSON) > 0 && string(teamworkJSON) != "{}" {
		var teamwork models.Teamwork
		if err := json.Unmarshal(teamworkJSON, &teamwork); err == nil {
			resp.Header.Teamwork = &teamwork
		}
	}
	if len(selflearningJSON) > 0 && string(selflearningJSON) != "{}" {
		var selflearning models.SelfLearning
		if err := json.Unmarshal(selflearningJSON, &selflearning); err == nil {
			resp.Header.SelfLearning = &selflearning
		}
	}

	// 2. Get models linked via syllabus_id
	modelRows, err := db.DB.Query(`
		SELECT id, syllabus_id, model_name, position 
		FROM syllabus_models 
		WHERE syllabus_id = ? 
		ORDER BY position, id`, resp.Header.ID)

	if err != nil {
		log.Println("Error fetching models:", err)
		http.Error(w, "Failed to fetch models", http.StatusInternalServerError)
		return
	}
	defer modelRows.Close()

	modelsList := []models.SyllabusModel{}
	for modelRows.Next() {
		var model models.SyllabusModel
		if err := modelRows.Scan(&model.ID, &model.SyllabusID, &model.ModelName, &model.Position); err != nil {
			log.Println("Error scanning model:", err)
			continue
		}

		// 3. Get titles for this model
		titleRows, err := db.DB.Query(`
			SELECT id, model_id, title_name, hours, position 
			FROM syllabus_titles 
			WHERE model_id = ? 
			ORDER BY position, id`, model.ID)

		if err != nil {
			log.Println("Error fetching titles for model", model.ID, ":", err)
			model.Titles = []models.SyllabusTitle{}
			modelsList = append(modelsList, model)
			continue
		}

		titlesList := []models.SyllabusTitle{}
		for titleRows.Next() {
			var title models.SyllabusTitle
			if err := titleRows.Scan(&title.ID, &title.ModelID, &title.TitleName, &title.Hours, &title.Position); err != nil {
				log.Println("Error scanning title:", err)
				continue
			}

			// 4. Get topics for this title
			topicRows, err := db.DB.Query(`
				SELECT id, title_id, topic, position 
				FROM syllabus_topics 
				WHERE title_id = ? 
				ORDER BY position, id`, title.ID)

			if err != nil {
				log.Println("Error fetching topics for title", title.ID, ":", err)
				title.Topics = []models.SyllabusTopic{}
				titlesList = append(titlesList, title)
				continue
			}

			topicsList := []models.SyllabusTopic{}
			for topicRows.Next() {
				var topic models.SyllabusTopic
				if err := topicRows.Scan(&topic.ID, &topic.TitleID, &topic.Topic, &topic.Position); err != nil {
					log.Println("Error scanning topic:", err)
					continue
				}
				topicsList = append(topicsList, topic)
			}
			topicRows.Close()

			title.Topics = topicsList
			titlesList = append(titlesList, title)
		}
		titleRows.Close()

		model.Titles = titlesList
		modelsList = append(modelsList, model)
	}

	resp.Models = modelsList
	json.NewEncoder(w).Encode(resp)
}

// ============================================================================
// MODEL CRUD OPERATIONS
// ============================================================================

// CreateModel creates a new module under a course syllabus
func CreateModel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID, err := strconv.Atoi(vars["courseId"])
	if err != nil {
		http.Error(w, "Invalid course ID", http.StatusBadRequest)
		return
	}

	var body struct {
		ModelName string `json:"model_name"`
		Position  int    `json:"position"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Ensure course_syllabus record exists
	var syllabusID int
	err = db.DB.QueryRow("SELECT id FROM course_syllabus WHERE course_id = ?", courseID).Scan(&syllabusID)

	if err == sql.ErrNoRows {
		// Create new syllabus header
		result, err := db.DB.Exec(`
			INSERT INTO course_syllabus (course_id, objectives, outcomes, reference_list, prerequisites) 
			VALUES (?, JSON_ARRAY(), JSON_ARRAY(), JSON_ARRAY(), JSON_ARRAY())`, courseID)
		if err != nil {
			log.Println("Error creating syllabus header:", err)
			http.Error(w, "Failed to create syllabus", http.StatusInternalServerError)
			return
		}
		id, _ := result.LastInsertId()
		syllabusID = int(id)
	} else if err != nil {
		log.Println("Error checking syllabus:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Insert model with both syllabus_id and course_id for data integrity
	result, err := db.DB.Exec(`
		INSERT INTO syllabus_models (syllabus_id, model_name, name, position, course_id) 
		VALUES (?, ?, ?, ?, ?)`, syllabusID, body.ModelName, body.ModelName, body.Position, courseID)

	if err != nil {
		log.Println("CreateModel error:", err)
		http.Error(w, "Failed to create model", http.StatusInternalServerError)
		return
	}

	modelID, _ := result.LastInsertId()
	json.NewEncoder(w).Encode(map[string]int{"id": int(modelID)})
}

// UpdateModel updates a model's name and position
func UpdateModel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	modelID, err := strconv.Atoi(vars["modelId"])
	if err != nil {
		http.Error(w, "Invalid model ID", http.StatusBadRequest)
		return
	}

	var body struct {
		ModelName string `json:"model_name"`
		Position  int    `json:"position"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec(`
		UPDATE syllabus_models 
		SET model_name = ?, name = ?, position = ? 
		WHERE id = ?`, body.ModelName, body.ModelName, body.Position, modelID)

	if err != nil {
		log.Println("UpdateModel error:", err)
		http.Error(w, "Failed to update model", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteModel deletes a model (cascades to titles and topics)
func DeleteModel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	modelID, err := strconv.Atoi(vars["modelId"])
	if err != nil {
		http.Error(w, "Invalid model ID", http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec("DELETE FROM syllabus_models WHERE id = ?", modelID)
	if err != nil {
		log.Println("DeleteModel error:", err)
		http.Error(w, "Failed to delete model", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ============================================================================
// TITLE CRUD OPERATIONS
// ============================================================================

// CreateTitle creates a new title under a model
func CreateTitle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	modelID, err := strconv.Atoi(vars["modelId"])
	if err != nil {
		http.Error(w, "Invalid model ID", http.StatusBadRequest)
		return
	}

	var body struct {
		TitleName string `json:"title_name"`
		Hours     int    `json:"hours"`
		Position  int    `json:"position"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := db.DB.Exec(`
		INSERT INTO syllabus_titles (model_id, title_name, title, hours, position) 
		VALUES (?, ?, ?, ?, ?)`, modelID, body.TitleName, body.TitleName, body.Hours, body.Position)

	if err != nil {
		log.Println("CreateTitle error:", err)
		http.Error(w, "Failed to create title", http.StatusInternalServerError)
		return
	}

	titleID, _ := result.LastInsertId()
	json.NewEncoder(w).Encode(map[string]int{"id": int(titleID)})
}

// UpdateTitle updates a title's name, hours, and position
func UpdateTitle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	titleID, err := strconv.Atoi(vars["titleId"])
	if err != nil {
		http.Error(w, "Invalid title ID", http.StatusBadRequest)
		return
	}

	var body struct {
		TitleName string `json:"title_name"`
		Hours     int    `json:"hours"`
		Position  int    `json:"position"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec(`
		UPDATE syllabus_titles 
		SET title_name = ?, title = ?, hours = ?, position = ? 
		WHERE id = ?`, body.TitleName, body.TitleName, body.Hours, body.Position, titleID)

	if err != nil {
		log.Println("UpdateTitle error:", err)
		http.Error(w, "Failed to update title", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteTitle deletes a title (cascades to topics)
func DeleteTitle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	titleID, err := strconv.Atoi(vars["titleId"])
	if err != nil {
		http.Error(w, "Invalid title ID", http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec("DELETE FROM syllabus_titles WHERE id = ?", titleID)
	if err != nil {
		log.Println("DeleteTitle error:", err)
		http.Error(w, "Failed to delete title", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ============================================================================
// TOPIC CRUD OPERATIONS
// ============================================================================

// CreateTopic creates a new topic under a title
func CreateTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	titleID, err := strconv.Atoi(vars["titleId"])
	if err != nil {
		http.Error(w, "Invalid title ID", http.StatusBadRequest)
		return
	}

	var body struct {
		Topic    string `json:"topic"`
		Position int    `json:"position"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := db.DB.Exec(`
		INSERT INTO syllabus_topics (title_id, topic, content, position) 
		VALUES (?, ?, ?, ?)`, titleID, body.Topic, body.Topic, body.Position)

	if err != nil {
		log.Println("CreateTopic error:", err)
		http.Error(w, "Failed to create topic", http.StatusInternalServerError)
		return
	}

	topicID, _ := result.LastInsertId()
	json.NewEncoder(w).Encode(map[string]int{"id": int(topicID)})
}

// UpdateTopic updates a topic's content and position
func UpdateTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicID, err := strconv.Atoi(vars["topicId"])
	if err != nil {
		http.Error(w, "Invalid topic ID", http.StatusBadRequest)
		return
	}

	var body struct {
		Topic    string `json:"topic"`
		Position int    `json:"position"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec(`
		UPDATE syllabus_topics 
		SET topic = ?, content = ?, position = ? 
		WHERE id = ?`, body.Topic, body.Topic, body.Position, topicID)

	if err != nil {
		log.Println("UpdateTopic error:", err)
		http.Error(w, "Failed to update topic", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteTopic deletes a topic
func DeleteTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicID, err := strconv.Atoi(vars["topicId"])
	if err != nil {
		http.Error(w, "Invalid topic ID", http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec("DELETE FROM syllabus_topics WHERE id = ?", topicID)
	if err != nil {
		log.Println("DeleteTopic error:", err)
		http.Error(w, "Failed to delete topic", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
