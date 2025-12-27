package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"server/db"
	"server/models"
	"strconv"

	"github.com/gorilla/mux"
)

// CreateCurriculumLog handles POST /curriculum/:id/log
func CreateCurriculumLog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	curriculumID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid curriculum ID", http.StatusBadRequest)
		return
	}

	var logEntry models.CurriculumLog
	if err := json.NewDecoder(r.Body).Decode(&logEntry); err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	logEntry.CurriculumID = curriculumID

	// Insert log entry
	result, err := db.DB.Exec(`
		INSERT INTO curriculum_logs (curriculum_id, action, description, changed_by)
		VALUES (?, ?, ?, ?)
	`, logEntry.CurriculumID, logEntry.Action, logEntry.Description, logEntry.ChangedBy)

	if err != nil {
		log.Println("Error creating log entry:", err)
		http.Error(w, "Failed to create log entry", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	logEntry.ID = int(id)

	json.NewEncoder(w).Encode(logEntry)
}

// GetCurriculumLogs handles GET /curriculum/:id/logs
func GetCurriculumLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	curriculumID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid curriculum ID", http.StatusBadRequest)
		return
	}

	rows, err := db.DB.Query(`
		SELECT id, curriculum_id, action, description, changed_by, diff, created_at
		FROM curriculum_logs
		WHERE curriculum_id = ?
		ORDER BY created_at DESC
	`, curriculumID)

	if err != nil {
		log.Println("Error fetching logs:", err)
		http.Error(w, "Failed to fetch logs", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var logs []models.CurriculumLog
	for rows.Next() {
		var logEntry models.CurriculumLog
		var diffData []byte
		err := rows.Scan(&logEntry.ID, &logEntry.CurriculumID, &logEntry.Action,
			&logEntry.Description, &logEntry.ChangedBy, &diffData, &logEntry.CreatedAt)
		if err != nil {
			log.Println("Error scanning log entry:", err)
			continue
		}
		if len(diffData) > 0 {
			logEntry.Diff = json.RawMessage(diffData)
		}
		logs = append(logs, logEntry)
	}

	if logs == nil {
		logs = []models.CurriculumLog{}
	}

	json.NewEncoder(w).Encode(logs)
}

// Helper function to create log entries (non-blocking)
func LogCurriculumActivity(curriculumID int, action, description, changedBy string) {
	go func() {
		if changedBy == "" {
			changedBy = "System"
		}
		_, err := db.DB.Exec(`
			INSERT INTO curriculum_logs (curriculum_id, action, description, changed_by)
			VALUES (?, ?, ?, ?)
		`, curriculumID, action, description, changedBy)
		if err != nil {
			log.Printf("Warning: Failed to log activity for curriculum %d: %v", curriculumID, err)
		}
	}()
}
