package db

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() error {
	var err error
	// Use only Aiven database - build DSN from env vars
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

	if host == "" || port == "" || user == "" || name == "" {
		return fmt.Errorf("missing required database environment variables (DB_HOST, DB_PORT, DB_USER, DB_NAME)")
	}

	// Setup TLS for Aiven
	rootCertPool := x509.NewCertPool()
	caPEM := os.Getenv("DB_CA_CERT")
	var dsn string

	if caPEM != "" {
		// Format the PEM certificate correctly by adding newlines
		// Replace escaped newlines if present
		caPEM = strings.ReplaceAll(caPEM, "\\n", "\n")

		// If certificate is on a single line, format it properly
		if !strings.Contains(caPEM, "\n") {
			// Add newline after BEGIN CERTIFICATE
			caPEM = strings.Replace(caPEM, "-----BEGIN CERTIFICATE-----", "-----BEGIN CERTIFICATE-----\n", 1)
			// Add newline before END CERTIFICATE
			caPEM = strings.Replace(caPEM, "-----END CERTIFICATE-----", "\n-----END CERTIFICATE-----", 1)
			// Split the certificate body into 64-character lines
			parts := strings.SplitN(caPEM, "\n", 2)
			if len(parts) == 2 {
				middle := strings.SplitN(parts[1], "\n", 2)
				if len(middle) == 2 {
					certBody := middle[0]
					var formattedBody strings.Builder
					for i := 0; i < len(certBody); i += 64 {
						end := i + 64
						if end > len(certBody) {
							end = len(certBody)
						}
						formattedBody.WriteString(certBody[i:end])
						formattedBody.WriteString("\n")
					}
					caPEM = parts[0] + "\n" + formattedBody.String() + middle[1]
				}
			}
		}

		pemBytes := []byte(caPEM)

		if rootCertPool.AppendCertsFromPEM(pemBytes) {
			cfg := &tls.Config{RootCAs: rootCertPool}
			tlsName := "aiven"
			if err := mysql.RegisterTLSConfig(tlsName, cfg); err != nil {
				log.Fatal("Failed to register TLS config:", err)
				return err
			}
			dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local&tls=%s", user, pass, host, port, name, tlsName)
		} else {
			return fmt.Errorf("failed to parse DB_CA_CERT")
		}
	} else {
		// Fallback without TLS (not recommended for production)
		log.Println("WARNING: Connecting without TLS certificate")
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", user, pass, host, port, name)
	}

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return err
	}

	// Test the connection
	err = DB.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
		return err
	}

	fmt.Println("Database connected successfully!")
	return nil
}

// CreateCurriculumTable creates the curriculum table if it doesn't exist
func CreateCurriculumTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS curriculum (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		academic_year VARCHAR(50) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)
	`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Failed to create curriculum table:", err)
		return err
	}

	fmt.Println("Curriculum table created/verified successfully!")
	return nil
}

// CreateCoursesTable creates the core courses table (using existing schema with course_id)
func CreateCoursesTable() error {
	// Check if table exists
	var tableExists bool
	err := DB.QueryRow("SELECT COUNT(*) > 0 FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'courses'").Scan(&tableExists)

	if err == nil && tableExists {
		// Table already exists, verify it has course_id column
		var courseIdExists bool
		err = DB.QueryRow("SELECT COUNT(*) > 0 FROM information_schema.columns WHERE table_schema = DATABASE() AND table_name = 'courses' AND column_name = 'course_id'").Scan(&courseIdExists)

		if courseIdExists {
			fmt.Println("Courses table verified successfully!")
			// Ensure all required columns exist
			_ = ensureColumnExists("courses", "course_code", "VARCHAR(50) NOT NULL")
			_ = ensureColumnExists("courses", "course_name", "VARCHAR(255) NOT NULL")
			_ = ensureColumnExists("courses", "course_type", "VARCHAR(50)")
			_ = ensureColumnExists("courses", "category", "VARCHAR(50)")
			_ = ensureColumnExists("courses", "credit", "INT")
			_ = ensureColumnExists("courses", "theory_hours", "INT")
			_ = ensureColumnExists("courses", "activity_hours", "INT")
			_ = ensureColumnExists("courses", "lecture_hours", "INT")
			_ = ensureColumnExists("courses", "tutorial_hours", "INT")
			_ = ensureColumnExists("courses", "practical_hours", "INT")
			_ = ensureColumnExists("courses", "cia_marks", "INT")
			_ = ensureColumnExists("courses", "see_marks", "INT")
			_ = ensureColumnExists("courses", "total_marks", "INT")
			_ = ensureColumnExists("courses", "total_hours", "INT")
			return nil
		}
	}

	// Create table if it doesn't exist (using course_id to match existing schema)
	query := `
	CREATE TABLE IF NOT EXISTS courses (
		course_id INT AUTO_INCREMENT PRIMARY KEY,
		course_code VARCHAR(50) NOT NULL,
		course_name VARCHAR(255) NOT NULL,
		course_type VARCHAR(50),
		category VARCHAR(50),
		credit INT,
		theory_hours INT,
		activity_hours INT,
		lecture_hours INT,
		tutorial_hours INT,
		practical_hours INT,
		cia_marks INT,
		see_marks INT,
		total_marks INT,
		total_hours INT,
		UNIQUE KEY unique_course_code (course_code)
	) ENGINE=InnoDB
	`

	_, err = DB.Exec(query)
	if err != nil {
		log.Println("Error creating courses table:", err)
		return err
	}

	fmt.Println("Courses table created/verified successfully!")
	return nil
}

// CreateCurriculumCoursesTable creates the junction table linking courses to curriculum semesters
func CreateCurriculumCoursesTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS curriculum_courses (
		id INT AUTO_INCREMENT PRIMARY KEY,
		regulation_id INT NOT NULL,
		semester_id INT NOT NULL,
		course_id INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (course_id) REFERENCES courses(course_id) ON DELETE CASCADE,
		UNIQUE KEY unique_course_semester (regulation_id, semester_id, course_id)
	) ENGINE=InnoDB
	`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Failed to create curriculum_courses table:", err)
		return err
	}

	fmt.Println("Curriculum courses table created/verified successfully!")
	return nil
}

// CreateCurriculumLogsTable creates the curriculum_logs table if it doesn't exist
func CreateCurriculumLogsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS curriculum_logs (
		id INT AUTO_INCREMENT PRIMARY KEY,
		curriculum_id INT NOT NULL,
		action VARCHAR(255) NOT NULL,
		description TEXT,
		changed_by VARCHAR(255) DEFAULT 'System',
		diff JSON,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (curriculum_id) REFERENCES curriculum(id) ON DELETE CASCADE
	)
	`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Failed to create curriculum_logs table:", err)
		return err
	}

	// Add diff column if it doesn't exist (for existing tables)
	alterQuery := `
	ALTER TABLE curriculum_logs 
	ADD COLUMN IF NOT EXISTS diff JSON AFTER changed_by
	`
	DB.Exec(alterQuery) // Ignore errors as column may already exist

	fmt.Println("Curriculum logs table created/verified successfully!")
	return nil
}

func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// CreateClusterTables creates cluster management tables
func CreateClusterTables() error {
	// Create clusters table
	clustersQuery := `
	CREATE TABLE IF NOT EXISTS clusters (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)
	`
	_, err := DB.Exec(clustersQuery)
	if err != nil {
		log.Fatal("Failed to create clusters table:", err)
		return err
	}

	// Create cluster_departments mapping table
	clusterDepartmentsQuery := `
	CREATE TABLE IF NOT EXISTS cluster_departments (
		id INT AUTO_INCREMENT PRIMARY KEY,
		cluster_id INT NOT NULL,
		department_id INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (cluster_id) REFERENCES clusters(id) ON DELETE CASCADE,
		UNIQUE KEY unique_department (department_id)
	)
	`
	_, err = DB.Exec(clusterDepartmentsQuery)
	if err != nil {
		log.Fatal("Failed to create cluster_departments table:", err)
		return err
	}

	fmt.Println("Cluster tables created/verified successfully!")
	return nil
}

// CreateDepartmentOverviewTables creates department overview and related tables
func CreateDepartmentOverviewTables() error {
	// Main department_overview table
	overviewQuery := `
	CREATE TABLE IF NOT EXISTS department_overview (
		id INT AUTO_INCREMENT PRIMARY KEY,
		regulation_id INT NOT NULL,
		vision TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE KEY unique_regulation (regulation_id)
	)
	`
	_, err := DB.Exec(overviewQuery)
	if err != nil {
		log.Fatal("Failed to create department_overview table:", err)
		return err
	}

	// Mission table
	missionQuery := `
	CREATE TABLE IF NOT EXISTS department_mission (
		id INT AUTO_INCREMENT PRIMARY KEY,
		department_id INT NOT NULL,
		mission_text TEXT NOT NULL,
		visibility ENUM('UNIQUE', 'CLUSTER') DEFAULT 'UNIQUE',
		position INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (department_id) REFERENCES department_overview(id) ON DELETE CASCADE
	)
	`
	_, err = DB.Exec(missionQuery)
	if err != nil {
		log.Fatal("Failed to create department_mission table:", err)
		return err
	}

	// PEOs table
	peosQuery := `
	CREATE TABLE IF NOT EXISTS department_peos (
		id INT AUTO_INCREMENT PRIMARY KEY,
		department_id INT NOT NULL,
		peo_text TEXT NOT NULL,
		visibility ENUM('UNIQUE', 'CLUSTER') DEFAULT 'UNIQUE',
		position INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (department_id) REFERENCES department_overview(id) ON DELETE CASCADE
	)
	`
	_, err = DB.Exec(peosQuery)
	if err != nil {
		log.Fatal("Failed to create department_peos table:", err)
		return err
	}

	// POs table
	posQuery := `
	CREATE TABLE IF NOT EXISTS department_pos (
		id INT AUTO_INCREMENT PRIMARY KEY,
		department_id INT NOT NULL,
		po_text TEXT NOT NULL,
		visibility ENUM('UNIQUE', 'CLUSTER') DEFAULT 'UNIQUE',
		position INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (department_id) REFERENCES department_overview(id) ON DELETE CASCADE
	)
	`
	_, err = DB.Exec(posQuery)
	if err != nil {
		log.Fatal("Failed to create department_pos table:", err)
		return err
	}

	// PSOs table
	psosQuery := `
	CREATE TABLE IF NOT EXISTS department_psos (
		id INT AUTO_INCREMENT PRIMARY KEY,
		department_id INT NOT NULL,
		pso_text TEXT NOT NULL,
		visibility ENUM('UNIQUE', 'CLUSTER') DEFAULT 'UNIQUE',
		position INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (department_id) REFERENCES department_overview(id) ON DELETE CASCADE
	)
	`
	_, err = DB.Exec(psosQuery)
	if err != nil {
		log.Fatal("Failed to create department_psos table:", err)
		return err
	}

	fmt.Println("Department overview tables created/verified successfully!")
	return nil
}

// AddVisibilityColumns adds visibility ENUM column to department data tables
func AddVisibilityColumns() error {
	tables := []string{"department_mission", "department_peos", "department_pos", "department_psos"}

	for _, table := range tables {
		err := ensureColumnExists(table, "visibility", "ENUM('UNIQUE', 'CLUSTER') DEFAULT 'UNIQUE'")
		if err != nil {
			log.Printf("Warning: Failed to add visibility column to %s: %v", table, err)
		}
	}

	fmt.Println("Visibility columns added/verified successfully!")
	return nil
}

// CreateCourseSyllabusTable - DEPRECATED: course_syllabus table removed in favor of course-centric design
// All syllabus data now references courses.id directly
func CreateCourseSyllabusTable() error {
	// Drop the legacy course_syllabus table if it exists
	_, _ = DB.Exec("DROP TABLE IF EXISTS course_syllabus")

	fmt.Println("Legacy course_syllabus table removed - using course-centric design")
	return nil
}

// CreateNormalizedSyllabusTables creates normalized tables for syllabus child data
func CreateNormalizedSyllabusTables() error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS course_objectives (
			id INT AUTO_INCREMENT PRIMARY KEY,
			course_id INT NOT NULL,
			objective TEXT NOT NULL,
			position INT NOT NULL,
			UNIQUE KEY unique_course_position (course_id, position),
			FOREIGN KEY (course_id) REFERENCES courses(course_id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS course_outcomes (
			id INT AUTO_INCREMENT PRIMARY KEY,
			course_id INT NOT NULL,
			outcome TEXT NOT NULL,
			position INT NOT NULL,
			UNIQUE KEY unique_course_position (course_id, position),
			FOREIGN KEY (course_id) REFERENCES courses(course_id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS course_references (
			id INT AUTO_INCREMENT PRIMARY KEY,
			course_id INT NOT NULL,
			reference_text TEXT NOT NULL,
			position INT NOT NULL,
			UNIQUE KEY unique_course_position (course_id, position),
			FOREIGN KEY (course_id) REFERENCES courses(course_id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS course_prerequisites (
			id INT AUTO_INCREMENT PRIMARY KEY,
			course_id INT NOT NULL,
			prerequisite TEXT NOT NULL,
			position INT NOT NULL,
			UNIQUE KEY unique_course_position (course_id, position),
			FOREIGN KEY (course_id) REFERENCES courses(course_id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS course_teamwork (
			course_id INT NOT NULL PRIMARY KEY,
			total_hours INT NOT NULL DEFAULT 0,
			FOREIGN KEY (course_id) REFERENCES courses(course_id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS course_teamwork_activities (
			id INT AUTO_INCREMENT PRIMARY KEY,
			course_id INT NOT NULL,
			activity TEXT NOT NULL,
			position INT NOT NULL,
			UNIQUE KEY unique_course_position (course_id, position),
			FOREIGN KEY (course_id) REFERENCES courses(course_id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS course_selflearning (
			course_id INT NOT NULL PRIMARY KEY,
			total_hours INT NOT NULL DEFAULT 0,
			FOREIGN KEY (course_id) REFERENCES courses(course_id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS course_selflearning_main (
			id INT AUTO_INCREMENT PRIMARY KEY,
			course_id INT NOT NULL,
			main_text TEXT NOT NULL,
			position INT NOT NULL,
			UNIQUE KEY unique_course_position (course_id, position),
			FOREIGN KEY (course_id) REFERENCES courses(course_id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS course_selflearning_internal (
			id INT AUTO_INCREMENT PRIMARY KEY,
			main_id INT NOT NULL,
			internal_text TEXT NOT NULL,
			position INT NOT NULL,
			UNIQUE KEY unique_main_position (main_id, position),
			FOREIGN KEY (main_id) REFERENCES course_selflearning_main(id) ON DELETE CASCADE
		)`,
	}

	for _, createSQL := range tables {
		if _, err := DB.Exec(createSQL); err != nil {
			log.Println("Error creating normalized table:", err)
			return err
		}
	}

	fmt.Println("Normalized syllabus tables created/verified successfully!")
	return nil
}

// CreateSyllabusRelationalTables creates models, titles, topics tables with cascades
// All tables now reference course_id directly (course-centric design)
func CreateSyllabusRelationalTables() error {
	// models - references course_id directly
	if _, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS syllabus_models (
			id INT AUTO_INCREMENT PRIMARY KEY,
			course_id INT NOT NULL,
			name VARCHAR(255) NOT NULL,
			model_name VARCHAR(255) NOT NULL,
			position INT DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT fk_models_course FOREIGN KEY (course_id) REFERENCES courses(course_id) ON DELETE CASCADE
		) ENGINE=InnoDB`); err != nil {
		return err
	}

	// Drop legacy syllabus_id column if exists
	_ = dropColumnIfExists("syllabus_models", "syllabus_id")

	// Ensure required columns exist for legacy schemas
	_ = ensureColumnExists("syllabus_models", "course_id", "INT")
	_ = ensureColumnExists("syllabus_models", "name", "VARCHAR(255) NOT NULL DEFAULT ''")
	_ = ensureColumnExists("syllabus_models", "model_name", "VARCHAR(255) NOT NULL DEFAULT ''")
	_ = ensureColumnExists("syllabus_models", "position", "INT DEFAULT 0")
	// Index for filtering by course_id
	_, _ = DB.Exec("CREATE INDEX IF NOT EXISTS idx_models_course ON syllabus_models(course_id)")
	// titles
	if _, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS syllabus_titles (
			id INT AUTO_INCREMENT PRIMARY KEY,
			model_id INT NOT NULL,
				title VARCHAR(512) NOT NULL,
				title_name VARCHAR(512) NOT NULL,
			hours INT DEFAULT 0,
			position INT DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT fk_titles_model FOREIGN KEY (model_id) REFERENCES syllabus_models(id) ON DELETE CASCADE
		) ENGINE=InnoDB`); err != nil {
		return err
	}
	_ = ensureColumnExists("syllabus_titles", "model_id", "INT")
	_ = ensureColumnExists("syllabus_titles", "title", "VARCHAR(512) NOT NULL")
	_ = ensureColumnExists("syllabus_titles", "title_name", "VARCHAR(512) NOT NULL DEFAULT ''")
	_ = ensureColumnExists("syllabus_titles", "hours", "INT DEFAULT 0")
	_ = ensureColumnExists("syllabus_titles", "position", "INT DEFAULT 0")
	_, _ = DB.Exec("CREATE INDEX IF NOT EXISTS idx_titles_model ON syllabus_titles(model_id)")
	// topics
	if _, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS syllabus_topics (
			id INT AUTO_INCREMENT PRIMARY KEY,
			title_id INT NOT NULL,
				topic VARCHAR(1024) NOT NULL,
				content TEXT NOT NULL,
			position INT DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT fk_topics_title FOREIGN KEY (title_id) REFERENCES syllabus_titles(id) ON DELETE CASCADE
		) ENGINE=InnoDB`); err != nil {
		return err
	}
	_ = ensureColumnExists("syllabus_topics", "title_id", "INT")
	_ = ensureColumnExists("syllabus_topics", "topic", "VARCHAR(1024) NOT NULL DEFAULT ''")
	_ = ensureColumnExists("syllabus_topics", "content", "TEXT NOT NULL")
	_ = ensureColumnExists("syllabus_topics", "position", "INT DEFAULT 0")
	_, _ = DB.Exec("CREATE INDEX IF NOT EXISTS idx_topics_title ON syllabus_topics(title_id)")
	fmt.Println("Syllabus relational tables created/verified successfully!")
	return nil
}
func ensureColumnExists(table, column, colType string) error {
	// First try IF NOT EXISTS (MySQL 8.0.29+)
	alter := fmt.Sprintf("ALTER TABLE %s ADD COLUMN IF NOT EXISTS %s %s", table, column, colType)
	if _, err := DB.Exec(alter); err == nil {
		fmt.Println("Ensured column", column, "on", table)
		return nil
	}
	// Check if column already exists
	exists, err := columnExists(table, column)
	if err != nil {
		return err
	}
	if exists {
		// already present
		return nil
	}
	// Fallback: add without IF NOT EXISTS
	alter2 := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", table, column, colType)
	if _, err := DB.Exec(alter2); err != nil {
		return err
	}
	fmt.Println("Added column", column, "to", table)
	return nil
}

func columnExists(table, column string) (bool, error) {
	// Some MySQL servers/drivers don't allow placeholders in SHOW statements
	q := fmt.Sprintf("SHOW COLUMNS FROM %s LIKE '%s'", table, column)
	rows, err := DB.Query(q)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}

// dropColumnIfExists drops a column only if it exists.
func dropColumnIfExists(table, column string) error {
	exists, err := columnExists(table, column)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}
	q := fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", table, column)
	if _, err := DB.Exec(q); err != nil {
		return err
	}
	fmt.Println("Dropped column", column, "from", table)
	return nil
}

// AddVisibilitySemestersCourses adds visibility columns to semesters and courses tables
func AddVisibilitySemestersCourses() error {
	fmt.Println("Adding visibility columns to semesters and courses tables...")
	
	// Add visibility to semesters table
	if err := ensureColumnExists("semesters", "visibility", "ENUM('UNIQUE', 'CLUSTER') DEFAULT 'UNIQUE'"); err != nil {
		return fmt.Errorf("failed to add visibility to semesters: %w", err)
	}
	
	// Add visibility to courses table
	if err := ensureColumnExists("courses", "visibility", "ENUM('UNIQUE', 'CLUSTER') DEFAULT 'UNIQUE'"); err != nil {
		return fmt.Errorf("failed to add visibility to courses: %w", err)
	}
	
	fmt.Println("Visibility columns added to semesters and courses successfully!")
	return nil
}

// AddSourceDepartmentColumns adds source_department_id to track shared item origins
func AddSourceDepartmentColumns() error {
	fmt.Println("Adding source_department_id columns to track shared items...")
	
	// Add to department tables
	tables := []string{"department_mission", "department_peos", "department_pos", "department_psos"}
	for _, table := range tables {
		if err := ensureColumnExists(table, "source_department_id", "INT DEFAULT NULL"); err != nil {
			return fmt.Errorf("failed to add source_department_id to %s: %w", table, err)
		}
	}
	
	// Add to semesters
	if err := ensureColumnExists("semesters", "source_department_id", "INT DEFAULT NULL"); err != nil {
		return fmt.Errorf("failed to add source_department_id to semesters: %w", err)
	}
	
	// Add to courses (source_regulation_id to track which regulation it came from)
	if err := ensureColumnExists("courses", "source_regulation_id", "INT DEFAULT NULL"); err != nil {
		return fmt.Errorf("failed to add source_regulation_id to courses: %w", err)
	}
	
	fmt.Println("Source tracking columns added successfully!")
	return nil
}
