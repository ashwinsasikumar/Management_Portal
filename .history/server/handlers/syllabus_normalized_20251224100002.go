package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"server/db"
	"server/models"
)

// Helper functions for normalized syllabus data access
// All tables reference course_id directly (course-centric design)

// fetchObjectives retrieves all objectives for a course ordered by position
func fetchObjectives(courseID int) ([]string, error) {
	rows, err := db.DB.Query(`
		SELECT objective 
		FROM course_objectives 
		WHERE course_id = ? 
		ORDER BY position`, courseID)
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()

	objectives := []string{}
	for rows.Next() {
		var text string
		if err := rows.Scan(&text); err == nil {
			objectives = append(objectives, text)
		}
	}
	return objectives, nil
}

// fetchOutcomes retrieves all outcomes for a course ordered by position
func fetchOutcomes(courseID int) ([]string, error) {
	rows, err := db.DB.Query(`
		SELECT outcome 
		FROM course_outcomes 
		WHERE course_id = ? 
		ORDER BY position`, courseID)
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()

	outcomes := []string{}
	for rows.Next() {
		var text string
		if err := rows.Scan(&text); err == nil {
			outcomes = append(outcomes, text)
		}
	}
	return outcomes, nil
}

// fetchReferences retrieves all references for a course ordered by position
func fetchReferences(courseID int) ([]string, error) {
	rows, err := db.DB.Query(`
		SELECT reference 
		FROM course_references 
		WHERE course_id = ? 
		ORDER BY position`, courseID)
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()

	references := []string{}
	for rows.Next() {
		var text string
		if err := rows.Scan(&text); err == nil {
			references = append(references, text)
		}
	}
	return references, nil
}

// fetchPrerequisites retrieves all prerequisites for a course ordered by position
func fetchPrerequisites(courseID int) ([]string, error) {
	rows, err := db.DB.Query(`
		SELECT prerequisite 
		FROM course_prerequisites 
		WHERE course_id = ? 
		ORDER BY position`, courseID)
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()

	prerequisites := []string{}
	for rows.Next() {
		var text string
		if err := rows.Scan(&text); err == nil {
			prerequisites = append(prerequisites, text)
		}
	}
	return prerequisites, nil
}

// fetchTeamwork retrieves teamwork data for a course
func fetchTeamwork(courseID int) (*models.Teamwork, error) {
	// Get hours from first row
	var hours int
	var activityText string
	err := db.DB.QueryRow(`
		SELECT hours, activity 
		FROM course_teamwork 
		WHERE course_id = ? 
		ORDER BY position 
		LIMIT 1`, courseID).Scan(&hours, &activityText)

	if err == sql.ErrNoRows {
		return nil, nil // No teamwork data
	}
	if err != nil {
		return nil, err
	}

	// Fetch all activities
	rows, err := db.DB.Query(`
		SELECT activity 
		FROM course_teamwork 
		WHERE course_id = ? 
		ORDER BY position`, courseID)
	if err != nil {
		return &models.Teamwork{Hours: hours, Activities: []string{}}, nil
	}
	defer rows.Close()

	activities := []string{}
	for rows.Next() {
		var text string
		if err := rows.Scan(&text); err == nil {
			activities = append(activities, text)
		}
	}

	return &models.Teamwork{
		Hours:      hours,
		Activities: activities,
	}, nil
}

// fetchSelfLearning retrieves self-learning data for a course
func fetchSelfLearning(courseID int) (*models.SelfLearning, error) {
	var hours int
	err := db.DB.QueryRow(`
		SELECT hours 
		FROM course_selflearning 
		WHERE course_id = ?`, courseID).Scan(&hours)

	if err == sql.ErrNoRows {
		return nil, nil // No self-learning data
	}
	if err != nil {
		return nil, err
	}

	// Fetch main topics
	rows, err := db.DB.Query(`
		SELECT id, main 
		FROM course_selflearning_main 
		WHERE course_id = ? 
		ORDER BY position`, courseID)
	if err != nil {
		return &models.SelfLearning{Hours: hours, MainInputs: []models.SelfLearningInternal{}}, nil
	}
	defer rows.Close()

	mainInputs := []models.SelfLearningInternal{}
	for rows.Next() {
		var mainID int
		var mainText string
		if err := rows.Scan(&mainID, &mainText); err == nil {
			// Fetch internal resources for this main topic
			internalRows, err := db.DB.Query(`
				SELECT internal 
				FROM course_selflearning_internal 
				WHERE main_id = ? 
				ORDER BY position`, mainID)

			internal := []string{}
			if err == nil {
				defer internalRows.Close()
				for internalRows.Next() {
					var text string
					if err := internalRows.Scan(&text); err == nil {
						internal = append(internal, text)
					}
				}
			}

			mainInputs = append(mainInputs, models.SelfLearningInternal{
				Main:     mainText,
				Internal: internal,
			})
		}
	}

	return &models.SelfLearning{
		Hours:      hours,
		MainInputs: mainInputs,
	}, nil
}

// saveObjectives saves objectives for a course, replacing existing ones
func saveObjectives(courseID int, objectives []string) error {
	// Delete existing
	_, err := db.DB.Exec("DELETE FROM course_objectives WHERE course_id = ?", courseID)
	if err != nil {
		return err
	}

	// Insert new ones with position
	for i, text := range objectives {
		if text == "" {
			continue
		}
		_, err := db.DB.Exec(`
			INSERT INTO course_objectives (course_id, objective, position) 
			VALUES (?, ?, ?)`, courseID, text, i)
		if err != nil {
			return err
		}
	}
	return nil
}

// saveOutcomes saves outcomes for a course, replacing existing ones
func saveOutcomes(courseID int, outcomes []string) error {
	// Delete existing
	_, err := db.DB.Exec("DELETE FROM course_outcomes WHERE course_id = ?", courseID)
	if err != nil {
		return err
	}

	// Insert new ones with position
	for i, text := range outcomes {
		if text == "" {
			continue
		}
		_, err := db.DB.Exec(`
			INSERT INTO course_outcomes (course_id, outcome, position) 
			VALUES (?, ?, ?)`, courseID, text, i)
		if err != nil {
			return err
		}
	}
	return nil
}

// saveReferences saves references for a course, replacing existing ones
func saveReferences(courseID int, references []string) error {
	// Delete existing
	_, err := db.DB.Exec("DELETE FROM course_references WHERE course_id = ?", courseID)
	if err != nil {
		return err
	}

	// Insert new ones with position
	for i, text := range references {
		if text == "" {
			continue
		}
		_, err := db.DB.Exec(`
			INSERT INTO course_references (course_id, reference, position) 
			VALUES (?, ?, ?)`, courseID, text, i)
		if err != nil {
			return err
		}
	}
	return nil
}

// savePrerequisites saves prerequisites for a course, replacing existing ones
func savePrerequisites(courseID int, prerequisites []string) error {
	// Delete existing
	_, err := db.DB.Exec("DELETE FROM course_prerequisites WHERE course_id = ?", courseID)
	if err != nil {
		return err
	}

	// Insert new ones with position
	for i, text := range prerequisites {
		if text == "" {
			continue
		}
		_, err := db.DB.Exec(`
			INSERT INTO course_prerequisites (course_id, prerequisite, position) 
			VALUES (?, ?, ?)`, courseID, text, i)
		if err != nil {
			return err
		}
	}
	return nil
}

// saveTeamwork saves teamwork data for a course
func saveTeamwork(courseID int, teamwork *models.Teamwork) error {
	if teamwork == nil {
		log.Printf("DEBUG: saveTeamwork called with nil teamwork for courseID %d", courseID)
		// Delete if nil
		db.DB.Exec("DELETE FROM course_teamwork WHERE course_id = ?", courseID)
		return nil
	}

	log.Printf("DEBUG: saveTeamwork called for courseID %d with hours=%d, activities=%v", courseID, teamwork.Hours, teamwork.Activities)

	// Delete existing teamwork data
	_, err := db.DB.Exec("DELETE FROM course_teamwork WHERE course_id = ?", courseID)
	if err != nil {
		log.Printf("ERROR: Failed to delete existing teamwork: %v", err)
		return err
	}

	// Insert new activities with hours and position
	log.Printf("DEBUG: Inserting %d teamwork activities", len(teamwork.Activities))
	for i, text := range teamwork.Activities {
		if text == "" {
			log.Printf("DEBUG: Skipping empty activity at position %d", i)
			continue
		}
		result, err := db.DB.Exec(`
			INSERT INTO course_teamwork (course_id, hours, activity, position) 
			VALUES (?, ?, ?, ?)`, courseID, teamwork.Hours, text, i)
		if err != nil {
			log.Printf("ERROR: Failed to insert teamwork activity at position %d: %v", i, err)
			return err
		}
		rowsAffected, _ := result.RowsAffected()
		log.Printf("DEBUG: Inserted teamwork activity at position %d, affected %d rows", i, rowsAffected)
	}
	log.Printf("DEBUG: Successfully saved teamwork data for courseID %d", courseID)
	return nil
}

// saveSelfLearning saves self-learning data for a course
func saveSelfLearning(courseID int, selflearning *models.SelfLearning) error {
	if selflearning == nil {
		// Delete if nil
		db.DB.Exec("DELETE FROM course_selflearning_internal WHERE main_id IN (SELECT id FROM course_selflearning_main WHERE course_id = ?)", courseID)
		db.DB.Exec("DELETE FROM course_selflearning_main WHERE course_id = ?", courseID)
		db.DB.Exec("DELETE FROM course_selflearning WHERE course_id = ?", courseID)
		return nil
	}

	// Upsert self-learning hours
	_, err := db.DB.Exec(`
		INSERT INTO course_selflearning (course_id, hours) 
		VALUES (?, ?) 
		ON DUPLICATE KEY UPDATE hours = ?`,
		courseID, selflearning.Hours, selflearning.Hours)
	if err != nil {
		return err
	}

	// Delete existing main topics and their internals
	db.DB.Exec("DELETE FROM course_selflearning_internal WHERE main_id IN (SELECT id FROM course_selflearning_main WHERE course_id = ?)", courseID)
	db.DB.Exec("DELETE FROM course_selflearning_main WHERE course_id = ?", courseID)

	// Insert new main topics
	for i, mainInput := range selflearning.MainInputs {
		if mainInput.Main == "" {
			continue
		}

		result, err := db.DB.Exec(`
			INSERT INTO course_selflearning_main (course_id, main, position) 
			VALUES (?, ?, ?)`, courseID, mainInput.Main, i)
		if err != nil {
			return err
		}

		mainID, _ := result.LastInsertId()

		// Insert internal resources for this main topic
		for j, text := range mainInput.Internal {
			if text == "" {
				continue
			}
			_, err := db.DB.Exec(`
				INSERT INTO course_selflearning_internal (main_id, internal, position) 
				VALUES (?, ?, ?)`, mainID, text, j)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
