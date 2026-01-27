package curriculum

import (
	"database/sql"
	"log"

	"server/db"
)

// getCurriculumTemplateByRegulation returns the curriculum_template for a curriculum with a safe default.
func getCurriculumTemplateByRegulation(curriculumID int) string {
	var tmpl sql.NullString
	err := db.DB.QueryRow(`SELECT curriculum_template FROM curriculum WHERE id = ?`, curriculumID).Scan(&tmpl)
	if err != nil {
		log.Println("template lookup failed, defaulting to 2026:", err)
		return "2026"
	}
	if tmpl.Valid && tmpl.String != "" {
		return tmpl.String
	}
	return "2026"
}

// getCurriculumTemplateForCourse finds the template for a course via curriculum or honour linkages.
func getCurriculumTemplateForCourse(courseID int) string {
	var tmpl sql.NullString

	// Try curriculum_courses linkage
	err := db.DB.QueryRow(`
        SELECT c.curriculum_template
        FROM curriculum_courses cc
        INNER JOIN curriculum c ON c.id = cc.curriculum_id
        WHERE cc.course_id = ?
        LIMIT 1`, courseID).Scan(&tmpl)
	if err == nil && tmpl.Valid && tmpl.String != "" {
		return tmpl.String
	}

	// Try honour vertical linkage
	err = db.DB.QueryRow(`
        SELECT c.curriculum_template
        FROM honour_vertical_courses hvc
        INNER JOIN honour_verticals hv ON hv.id = hvc.honour_vertical_id
        INNER JOIN honour_cards hc ON hc.id = hv.honour_card_id
        INNER JOIN curriculum c ON c.id = hc.curriculum_id
        WHERE hvc.course_id = ?
        LIMIT 1`, courseID).Scan(&tmpl)
	if err == nil && tmpl.Valid && tmpl.String != "" {
		return tmpl.String
	}

	return "2026"
}
