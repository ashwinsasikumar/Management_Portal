package db

import (
	"fmt"
	"log"
)

// CreateNormalizedSchema creates the complete normalized schema with zero JSON fields
// This replaces all JSON columns with properly normalized relational tables
func CreateNormalizedSchema() error {
	log.Println("Creating normalized schema (3NF compliant, zero JSON fields)...")

	schemas := []struct {
		name string
		ddl  string
	}{
		// Core hierarchy tables
		{"curriculum", createCurriculumTableDDL()},
		{"regulations", createRegulationsTableDDL()},
		{"semesters", createSemestersTableDDL()},
		{"courses", createCoursesTableDDL()},
		{"regulation_courses", createRegulationCoursesTableDDL()},

		// Syllabus header (normalized)
		{"course_syllabus_normalized", createCourseSyllabusNormalizedDDL()},

		// Syllabus components (replacing JSON arrays)
		{"course_objectives", createCourseObjectivesDDL()},
		{"course_outcomes", createCourseOutcomesDDL()},
		{"course_textbooks", createCourseTextbooksDDL()},
		{"course_references", createCourseReferencesDDL()},
		{"course_prerequisites", createCoursePrerequisitesDDL()},

		// Teamwork (replacing teamwork JSON)
		{"teamwork_activities", createTeamworkActivitiesDDL()},

		// Self-learning (replacing selflearning nested JSON)
		{"self_learning_main_topics", createSelfLearningMainTopicsDDL()},
		{"self_learning_resources", createSelfLearningResourcesDDL()},

		// Syllabus modules (enhanced)
		{"syllabus_modules_normalized", createSyllabusModulesNormalizedDDL()},
		{"syllabus_titles_normalized", createSyllabusTitlesNormalizedDDL()},
		{"syllabus_topics_normalized", createSyllabusTopicsNormalizedDDL()},

		// Mappings
		{"co_po_mappings", createCOPOMappingsDDL()},
		{"co_pso_mappings", createCOPSOMappingsDDL()},
		{"peo_po_mappings", createPEOPOMappingsDDL()},

		// Audit trail (normalized, no diff JSON)
		{"curriculum_logs_normalized", createCurriculumLogsNormalizedDDL()},
		{"curriculum_log_changes", createCurriculumLogChangesDDL()},

		// Department overview (normalized)
		{"department_overview_normalized", createDepartmentOverviewNormalizedDDL()},
		{"department_peos", createDepartmentPEOsDDL()},
		{"department_pos", createDepartmentPOsDDL()},
		{"department_psos", createDepartmentPSOsDDL()},
	}

	for _, schema := range schemas {
		log.Printf("Creating table: %s", schema.name)
		if _, err := DB.Exec(schema.ddl); err != nil {
			log.Printf("Warning: Failed to create %s: %v", schema.name, err)
			// Continue with other tables
		}
	}

	log.Println("Normalized schema creation completed!")
	return nil
}

func createCurriculumTableDDL() string {
	return `
	CREATE TABLE IF NOT EXISTS curriculum (
		id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		academic_year VARCHAR(50) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		UNIQUE KEY uk_curriculum_name_year (name, academic_year),
		INDEX idx_academic_year (academic_year)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`
}

func createRegulationsTableDDL() string {
	return `
	CREATE TABLE IF NOT EXISTS regulations (
		id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		curriculum_id INT UNSIGNED NOT NULL,
		regulation_code VARCHAR(100) NOT NULL,
		regulation_name VARCHAR(255) NOT NULL,
		effective_year YEAR NOT NULL,
		is_active BOOLEAN NOT NULL DEFAULT TRUE,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		CONSTRAINT fk_regulations_curriculum 
			FOREIGN KEY (curriculum_id) REFERENCES curriculum(id) 
			ON DELETE RESTRICT ON UPDATE CASCADE,
		UNIQUE KEY uk_regulation_code (regulation_code),
		INDEX idx_curriculum_id (curriculum_id),
		INDEX idx_effective_year (effective_year),
		CHECK (effective_year >= 2000 AND effective_year <= 2100)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`
}

func createSemestersTableDDL() string {
	return `
	CREATE TABLE IF NOT EXISTS semesters (
		id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		regulation_id INT UNSIGNED NOT NULL,
		semester_number TINYINT UNSIGNED NOT NULL,
		semester_name VARCHAR(100) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT fk_semesters_regulation 
			FOREIGN KEY (regulation_id) REFERENCES regulations(id) 
			ON DELETE CASCADE ON UPDATE CASCADE,
		UNIQUE KEY uk_regulation_semester (regulation_id, semester_number),
		INDEX idx_regulation_id (regulation_id),
		CHECK (semester_number BETWEEN 1 AND 12)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`
}

func createCoursesTableDDL() string {
	return `
	CREATE TABLE IF NOT EXISTS courses (
		id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		course_code VARCHAR(50) NOT NULL,
		course_name VARCHAR(255) NOT NULL,
		course_type ENUM('Theory', 'Practical', 'Lab', 'Project', 'Seminar', 'Internship') NOT NULL,
		category ENUM('Professional Core', 'Professional Elective', 'Open Elective', 'Mandatory', 'Audit') NOT NULL,
		credit DECIMAL(4,2) UNSIGNED NOT NULL,
		lecture_hours TINYINT UNSIGNED NOT NULL DEFAULT 0,
		tutorial_hours TINYINT UNSIGNED NOT NULL DEFAULT 0,
		practical_hours TINYINT UNSIGNED NOT NULL DEFAULT 0,
		cia_marks SMALLINT UNSIGNED NOT NULL DEFAULT 0,
		see_marks SMALLINT UNSIGNED NOT NULL DEFAULT 0,
		total_marks SMALLINT UNSIGNED GENERATED ALWAYS AS (cia_marks + see_marks) STORED,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		UNIQUE KEY uk_course_code (course_code),
		INDEX idx_course_type (course_type),
		INDEX idx_category (category),
		CHECK (credit >= 0 AND credit <= 10),
		CHECK (lecture_hours + tutorial_hours + practical_hours > 0),
		CHECK (cia_marks >= 0 AND see_marks >= 0)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`
}

func createRegulationCoursesTableDDL() string {
	return `
	CREATE TABLE IF NOT EXISTS regulation_courses (
		id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		regulation_id INT UNSIGNED NOT NULL,
		semester_id INT UNSIGNED NOT NULL,
		course_id INT UNSIGNED NOT NULL,
		display_order SMALLINT UNSIGNED NOT NULL DEFAULT 0,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT fk_regulation_courses_regulation 
			FOREIGN KEY (regulation_id) REFERENCES regulations(id) 
			ON DELETE CASCADE ON UPDATE CASCADE,
		CONSTRAINT fk_regulation_courses_semester 
			FOREIGN KEY (semester_id) REFERENCES semesters(id) 
			ON DELETE CASCADE ON UPDATE CASCADE,
		CONSTRAINT fk_regulation_courses_course 
			FOREIGN KEY (course_id) REFERENCES courses(id) 
			ON DELETE RESTRICT ON UPDATE CASCADE,
		UNIQUE KEY uk_regulation_semester_course (regulation_id, semester_id, course_id),
		INDEX idx_regulation_id (regulation_id),
		INDEX idx_semester_id (semester_id),
		INDEX idx_course_id (course_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`
}

func createCourseSyllabusNormalizedDDL() string {
	return `
	CREATE TABLE IF NOT EXISTS course_syllabus_normalized (
		id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		course_id INT UNSIGNED NOT NULL,
		teamwork_total_hours TINYINT UNSIGNED NOT NULL DEFAULT 0,
		self_learning_total_hours TINYINT UNSIGNED NOT NULL DEFAULT 0,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		CONSTRAINT fk_course_syllabus_norm_course 
			FOREIGN KEY (course_id) REFERENCES courses(id) 
			ON DELETE CASCADE ON UPDATE CASCADE,
		UNIQUE KEY uk_course_syllabus_norm (course_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`
}

func createCourseObjectivesDDL() string {
	return `
	CREATE TABLE IF NOT EXISTS course_objectives (
		id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		syllabus_id INT UNSIGNED NOT NULL,
		objective_text TEXT NOT NULL,
		display_order SMALLINT UNSIGNED NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT fk_course_objectives_syllabus 
			FOREIGN KEY (syllabus_id) REFERENCES course_syllabus_normalized(id) 
			ON DELETE CASCADE ON UPDATE CASCADE,
		UNIQUE KEY uk_syllabus_order (syllabus_id, display_order),
		INDEX idx_syllabus_id (syllabus_id),
		CHECK (CHAR_LENGTH(objective_text) >= 10)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`
}

func createCourseOutcomesDDL() string {
	return `
	CREATE TABLE IF NOT EXISTS course_outcomes (
		id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		syllabus_id INT UNSIGNED NOT NULL,
		outcome_code VARCHAR(20) NOT NULL,
		outcome_text TEXT NOT NULL,
		display_order SMALLINT UNSIGNED NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT fk_course_outcomes_syllabus 
			FOREIGN KEY (syllabus_id) REFERENCES course_syllabus_normalized(id) 
			ON DELETE CASCADE ON UPDATE CASCADE,
		UNIQUE KEY uk_syllabus_outcome_code (syllabus_id, outcome_code),
		UNIQUE KEY uk_syllabus_order (syllabus_id, display_order),
		INDEX idx_syllabus_id (syllabus_id),
		CHECK (CHAR_LENGTH(outcome_text) >= 10)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`
}

func createCourseTextbooksDDL() string {
	return `
	CREATE TABLE IF NOT EXISTS course_textbooks (
		id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		syllabus_id INT UNSIGNED NOT NULL,
		title VARCHAR(500) NOT NULL,
		authors VARCHAR(500) NOT NULL,
		publisher VARCHAR(255) DEFAULT NULL,
		edition VARCHAR(50) DEFAULT NULL,
		year YEAR DEFAULT NULL,
		isbn VARCHAR(20) DEFAULT NULL,
		display_order SMALLINT UNSIGNED NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT fk_course_textbooks_syllabus 
			FOREIGN KEY (syllabus_id) REFERENCES course_syllabus_normalized(id) 
			ON DELETE CASCADE ON UPDATE CASCADE,
		UNIQUE KEY uk_syllabus_order (syllabus_id, display_order),
		INDEX idx_syllabus_id (syllabus_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`
}

func createCourseReferencesDDL() string {
	return `
	CREATE TABLE IF NOT EXISTS course_references (
		id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		syllabus_id INT UNSIGNED NOT NULL,
		title VARCHAR(500) NOT NULL,
		authors VARCHAR(500) NOT NULL,
		publisher VARCHAR(255) DEFAULT NULL,
		edition VARCHAR(50) DEFAULT NULL,
		year YEAR DEFAULT NULL,
		isbn VARCHAR(20) DEFAULT NULL,
		url VARCHAR(1000) DEFAULT NULL,
		display_order SMALLINT UNSIGNED NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT fk_course_references_syllabus 
			FOREIGN KEY (syllabus_id) REFERENCES course_syllabus_normalized(id) 
			ON DELETE CASCADE ON UPDATE CASCADE,
		UNIQUE KEY uk_syllabus_order (syllabus_id, display_order),
		INDEX idx_syllabus_id (syllabus_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`
}

func createCoursePrerequisitesDDL() string {
	return `
	CREATE TABLE IF NOT EXISTS course_prerequisites (
		id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		syllabus_id INT UNSIGNED NOT NULL,
		prerequisite_course_id INT UNSIGNED NOT NULL,
		prerequisite_type ENUM('Mandatory', 'Recommended', 'Co-requisite') NOT NULL DEFAULT 'Mandatory',
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT fk_course_prerequisites_syllabus 
			FOREIGN KEY (syllabus_id) REFERENCES course_syllabus_normalized(id) 
			ON DELETE CASCADE ON UPDATE CASCADE,
		CONSTRAINT fk_course_prerequisites_course 
			FOREIGN KEY (prerequisite_course_id) REFERENCES courses(id) 
			ON DELETE RESTRICT ON UPDATE CASCADE,
		UNIQUE KEY uk_syllabus_prerequisite (syllabus_id, prerequisite_course_id),
		INDEX idx_syllabus_id (syllabus_id),
		INDEX idx_prerequisite_course_id (prerequisite_course_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`
}

func createTeamworkActivitiesDDL() string {
	return `
	CREATE TABLE IF NOT EXISTS teamwork_activities (
		id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		syllabus_id INT UNSIGNED NOT NULL,
		activity_name VARCHAR(500) NOT NULL,
		display_order SMALLINT UNSIGNED NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		CONSTRAINT fk_teamwork_activities_syllabus 
			FOREIGN KEY (syllabus_id) REFERENCES course_syllabus_normalized(id) 
			ON DELETE CASCADE ON UPDATE CASCADE,
		UNIQUE KEY uk_syllabus_order (syllabus_id, display_order),
		INDEX idx_syllabus_id (syllabus_id),
		CHECK (CHAR_LENGTH(activity_name) >= 5)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`
}

func createSelfLearningMainTopicsDDL() string {
	return `
	CREATE TABLE IF NOT EXISTS self_learning_main_topics (
		id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		syllabus_id INT UNSIGNED NOT NULL,
		main_topic VARCHAR(500) NOT NULL,
		display_order SMALLINT UNSIGNED NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		CONSTRAINT fk_self_learning_main_syllabus 
			FOREIGN KEY (syllabus_id) REFERENCES course_syllabus_normalized(id) 
			ON DELETE CASCADE ON UPDATE CASCADE,
		UNIQUE KEY uk_syllabus_order (syllabus_id, display_order),
		INDEX idx_syllabus_id (syllabus_id),
		CHECK (CHAR_LENGTH(main_topic) >= 5)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`
}

func createSelfLearningResourcesDDL() string {
	return `
	CREATE TABLE IF NOT EXISTS self_learning_resources (
		id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		main_topic_id INT UNSIGNED NOT NULL,
		resource_text VARCHAR(1000) NOT NULL,
		resource_type ENUM('Text', 'Link', 'Document', 'Video') NOT NULL DEFAULT 'Text',
		display_order SMALLINT UNSIGNED NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT fk_self_learning_resources_main 
			FOREIGN KEY (main_topic_id) REFERENCES self_learning_main_topics(id) 
			ON DELETE CASCADE ON UPDATE CASCADE,
		UNIQUE KEY uk_main_topic_order (main_topic_id, display_order),
		INDEX idx_main_topic_id (main_topic_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`
}

func createSyllabusModulesNormalizedDDL() string {
	return `
	CREATE TABLE IF NOT EXISTS syllabus_modules_normalized (
		id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		syllabus_id INT UNSIGNED NOT NULL,
		module_name VARCHAR(100) NOT NULL,
		module_number TINYINT UNSIGNED NOT NULL,
		display_order SMALLINT UNSIGNED NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT fk_syllabus_modules_norm_syllabus 
			FOREIGN KEY (syllabus_id) REFERENCES course_syllabus_normalized(id) 
			ON DELETE CASCADE ON UPDATE CASCADE,
		UNIQUE KEY uk_syllabus_module_number (syllabus_id, module_number),
		UNIQUE KEY uk_syllabus_order (syllabus_id, display_order),
		INDEX idx_syllabus_id (syllabus_id),
		CHECK (module_number BETWEEN 1 AND 50)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`
}

func createSyllabusTitlesNormalizedDDL() string {
	return `
	CREATE TABLE IF NOT EXISTS syllabus_titles_normalized (
		id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		module_id INT UNSIGNED NOT NULL,
		title_name VARCHAR(255) NOT NULL,
		hours TINYINT UNSIGNED NOT NULL DEFAULT 0,
		display_order SMALLINT UNSIGNED NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT fk_syllabus_titles_norm_module 
			FOREIGN KEY (module_id) REFERENCES syllabus_modules_normalized(id) 
			ON DELETE CASCADE ON UPDATE CASCADE,
		UNIQUE KEY uk_module_order (module_id, display_order),
		INDEX idx_module_id (module_id),
		CHECK (hours <= 100)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`
}

func createSyllabusTopicsNormalizedDDL() string {
	return `
	CREATE TABLE IF NOT EXISTS syllabus_topics_normalized (
		id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		title_id INT UNSIGNED NOT NULL,
		topic_text TEXT NOT NULL,
		display_order SMALLINT UNSIGNED NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT fk_syllabus_topics_norm_title 
			FOREIGN KEY (title_id) REFERENCES syllabus_titles_normalized(id) 
			ON DELETE CASCADE ON UPDATE CASCADE,
		UNIQUE KEY uk_title_order (title_id, display_order),
		INDEX idx_title_id (title_id),
		CHECK (CHAR_LENGTH(topic_text) >= 3)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`
}

func createCOPOMappingsDDL() string {
	return `
	CREATE TABLE IF NOT EXISTS co_po_mappings (
		id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		course_id INT UNSIGNED NOT NULL,
		co_code VARCHAR(20) NOT NULL,
		po_number TINYINT UNSIGNED NOT NULL,
		mapping_level ENUM('Low', 'Medium', 'High') NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		CONSTRAINT fk_co_po_course 
			FOREIGN KEY (course_id) REFERENCES courses(id) 
			ON DELETE CASCADE ON UPDATE CASCADE,
		UNIQUE KEY uk_course_co_po (course_id, co_code, po_number),
		INDEX idx_course_id (course_id),
		CHECK (po_number BETWEEN 1 AND 20)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`
}

func createCOPSOMappingsDDL() string {
	return `
	CREATE TABLE IF NOT EXISTS co_pso_mappings (
		id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		course_id INT UNSIGNED NOT NULL,
		co_code VARCHAR(20) NOT NULL,
		pso_number TINYINT UNSIGNED NOT NULL,
		mapping_level ENUM('Low', 'Medium', 'High') NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		CONSTRAINT fk_co_pso_course 
			FOREIGN KEY (course_id) REFERENCES courses(id) 
			ON DELETE CASCADE ON UPDATE CASCADE,
		UNIQUE KEY uk_course_co_pso (course_id, co_code, pso_number),
		INDEX idx_course_id (course_id),
		CHECK (pso_number BETWEEN 1 AND 20)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`
}

func createPEOPOMappingsDDL() string {
	return `
	CREATE TABLE IF NOT EXISTS peo_po_mappings (
		id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		regulation_id INT UNSIGNED NOT NULL,
		peo_number TINYINT UNSIGNED NOT NULL,
		po_number TINYINT UNSIGNED NOT NULL,
		mapping_level ENUM('Low', 'Medium', 'High') NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		CONSTRAINT fk_peo_po_regulation 
			FOREIGN KEY (regulation_id) REFERENCES regulations(id) 
			ON DELETE CASCADE ON UPDATE CASCADE,
		UNIQUE KEY uk_regulation_peo_po (regulation_id, peo_number, po_number),
		INDEX idx_regulation_id (regulation_id),
		CHECK (peo_number BETWEEN 1 AND 10),
		CHECK (po_number BETWEEN 1 AND 20)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`
}

func createCurriculumLogsNormalizedDDL() string {
	return `
	CREATE TABLE IF NOT EXISTS curriculum_logs_normalized (
		id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		curriculum_id INT UNSIGNED NOT NULL,
		action ENUM('CREATE', 'UPDATE', 'DELETE', 'IMPORT', 'EXPORT', 'APPROVE', 'REJECT') NOT NULL,
		entity_type ENUM('Regulation', 'Semester', 'Course', 'Syllabus', 'Mapping', 'Other') NOT NULL,
		entity_id INT UNSIGNED DEFAULT NULL,
		description TEXT NOT NULL,
		changed_by VARCHAR(255) NOT NULL DEFAULT 'System',
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT fk_curriculum_logs_norm_curriculum 
			FOREIGN KEY (curriculum_id) REFERENCES curriculum(id) 
			ON DELETE CASCADE ON UPDATE CASCADE,
		INDEX idx_curriculum_id (curriculum_id),
		INDEX idx_created_at (created_at),
		INDEX idx_action (action),
		INDEX idx_entity (entity_type, entity_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`
}

func createCurriculumLogChangesDDL() string {
	return `
	CREATE TABLE IF NOT EXISTS curriculum_log_changes (
		id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		log_id INT UNSIGNED NOT NULL,
		field_name VARCHAR(100) NOT NULL,
		old_value TEXT DEFAULT NULL,
		new_value TEXT DEFAULT NULL,
		CONSTRAINT fk_curriculum_log_changes_log 
			FOREIGN KEY (log_id) REFERENCES curriculum_logs_normalized(id) 
			ON DELETE CASCADE ON UPDATE CASCADE,
		INDEX idx_log_id (log_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`
}

func createDepartmentOverviewNormalizedDDL() string {
	return `
	CREATE TABLE IF NOT EXISTS department_overview_normalized (
		id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		regulation_id INT UNSIGNED NOT NULL,
		vision TEXT NOT NULL,
		mission TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		CONSTRAINT fk_department_overview_norm_regulation 
			FOREIGN KEY (regulation_id) REFERENCES regulations(id) 
			ON DELETE CASCADE ON UPDATE CASCADE,
		UNIQUE KEY uk_regulation (regulation_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`
}

func createDepartmentPEOsDDL() string {
	return `
	CREATE TABLE IF NOT EXISTS department_peos (
		id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		department_overview_id INT UNSIGNED NOT NULL,
		peo_number TINYINT UNSIGNED NOT NULL,
		peo_text TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT fk_department_peos_overview 
			FOREIGN KEY (department_overview_id) REFERENCES department_overview_normalized(id) 
			ON DELETE CASCADE ON UPDATE CASCADE,
		UNIQUE KEY uk_overview_peo (department_overview_id, peo_number),
		INDEX idx_department_overview_id (department_overview_id),
		CHECK (peo_number BETWEEN 1 AND 10)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`
}

func createDepartmentPOsDDL() string {
	return `
	CREATE TABLE IF NOT EXISTS department_pos (
		id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		department_overview_id INT UNSIGNED NOT NULL,
		po_number TINYINT UNSIGNED NOT NULL,
		po_text TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT fk_department_pos_overview 
			FOREIGN KEY (department_overview_id) REFERENCES department_overview_normalized(id) 
			ON DELETE CASCADE ON UPDATE CASCADE,
		UNIQUE KEY uk_overview_po (department_overview_id, po_number),
		INDEX idx_department_overview_id (department_overview_id),
		CHECK (po_number BETWEEN 1 AND 20)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`
}

func createDepartmentPSOsDDL() string {
	return `
	CREATE TABLE IF NOT EXISTS department_psos (
		id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		department_overview_id INT UNSIGNED NOT NULL,
		pso_number TINYINT UNSIGNED NOT NULL,
		pso_text TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT fk_department_psos_overview 
			FOREIGN KEY (department_overview_id) REFERENCES department_overview_normalized(id) 
			ON DELETE CASCADE ON UPDATE CASCADE,
		UNIQUE KEY uk_overview_pso (department_overview_id, pso_number),
		INDEX idx_department_overview_id (department_overview_id),
		CHECK (pso_number BETWEEN 1 AND 20)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`
}
