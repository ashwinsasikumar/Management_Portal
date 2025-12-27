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
	// Prefer full DSN if provided
	// Example: MYSQL_DSN="user:pass@tcp(host:port)/db?parseTime=true&loc=Local"
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		// Build DSN from individual env vars when MYSQL_DSN not set
		host := getenvDefault("DB_HOST", "localhost")
		port := getenvDefault("DB_PORT", "3306")
		user := getenvDefault("DB_USER", "root")
		pass := getenvDefault("DB_PASSWORD", "")
		name := getenvDefault("DB_NAME", "Management_Portal")

		// Optional SSL
		sslMode := strings.ToLower(getenvDefault("DB_SSL_MODE", "false"))
		useTLS := sslMode == "true" || sslMode == "1" || sslMode == "on"
		tlsName := ""
		if useTLS {
			rootCertPool := x509.NewCertPool()
			// Try path first, then inline cert, then default file
			caPath := os.Getenv("DB_CA_CERT_PATH")
			caPEM := os.Getenv("DB_CA_CERT")
			var pemBytes []byte
			if caPath != "" {
				if b, rerr := os.ReadFile(caPath); rerr == nil {
					pemBytes = b
				} else {
					log.Println("Warning: failed to read DB_CA_CERT_PATH:", rerr)
				}
			}
			if len(pemBytes) == 0 && caPEM != "" {
				// Allow \n escaped newlines in env
				caPEM = strings.ReplaceAll(caPEM, "\\n", "\n")
				pemBytes = []byte(caPEM)
			}
			if len(pemBytes) == 0 {
				if b, rerr := os.ReadFile("./ca.pem"); rerr == nil {
					pemBytes = b
				}
			}
			if len(pemBytes) > 0 && rootCertPool.AppendCertsFromPEM(pemBytes) {
				cfg := &tls.Config{RootCAs: rootCertPool}
				tlsName = "custom"
				if err := mysql.RegisterTLSConfig(tlsName, cfg); err != nil {
					log.Fatal("Failed to register TLS config:", err)
					return err
				}
			} else if useTLS {
				return fmt.Errorf("DB_SSL_MODE=true but no valid CA certificate was provided")
			}

			if tlsName != "" {
				dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local&tls=%s", user, pass, host, port, name, tlsName)
			} else {
				dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", user, pass, host, port, name)
			}
		} else {
			dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", user, pass, host, port, name)
		}
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

func getenvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
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

// CreateCourseSyllabusTable ensures course_syllabus exists and includes prerequisites column
func CreateCourseSyllabusTable() error {
	create := `
	CREATE TABLE IF NOT EXISTS course_syllabus (
		id INT AUTO_INCREMENT PRIMARY KEY,
		course_id INT NOT NULL,
		objectives JSON,
		outcomes JSON,
		unit1 JSON,
		unit2 JSON,
		unit3 JSON,
		unit4 JSON,
		unit5 JSON,
		textbooks JSON,
		reference_list JSON,
		prerequisites JSON,
		teamwork JSON,
		selflearning JSON,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	if _, err := DB.Exec(create); err != nil {
		log.Println("Failed to ensure course_syllabus table:", err)
		return err
	}

	// Ensure prerequisites column exists for older schemas
	if err := ensureColumnExists("course_syllabus", "prerequisites", "JSON"); err != nil {
		log.Println("Warning: could not ensure prerequisites column:", err)
	}

	// Ensure teamwork and selflearning columns exist
	if err := ensureColumnExists("course_syllabus", "teamwork", "JSON"); err != nil {
		log.Println("Warning: could not ensure teamwork column:", err)
	}
	if err := ensureColumnExists("course_syllabus", "selflearning", "JSON"); err != nil {
		log.Println("Warning: could not ensure selflearning column:", err)
	}

	fmt.Println("Course syllabus table created/verified successfully!")
	return nil
}

// CreateSyllabusRelationalTables creates models, titles, topics tables with cascades
func CreateSyllabusRelationalTables() error {
	// models
	if _, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS syllabus_models (
			id INT AUTO_INCREMENT PRIMARY KEY,
			course_id INT NOT NULL,
			name VARCHAR(255) NOT NULL,
			position INT DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT fk_models_course FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
		) ENGINE=InnoDB`); err != nil {
		return err
	}
	// Ensure required columns exist for legacy schemas
	_ = ensureColumnExists("syllabus_models", "course_id", "INT")
	_ = ensureColumnExists("syllabus_models", "name", "VARCHAR(255) NOT NULL DEFAULT ''")
	_ = ensureColumnExists("syllabus_models", "position", "INT DEFAULT 0")
	// Optional index for filtering by course_id
	_, _ = DB.Exec("CREATE INDEX IF NOT EXISTS idx_models_course ON syllabus_models(course_id)")
	// titles
	if _, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS syllabus_titles (
			id INT AUTO_INCREMENT PRIMARY KEY,
			model_id INT NOT NULL,
			title VARCHAR(512) NOT NULL,
			hours INT DEFAULT 0,
			position INT DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT fk_titles_model FOREIGN KEY (model_id) REFERENCES syllabus_models(id) ON DELETE CASCADE
		) ENGINE=InnoDB`); err != nil {
		return err
	}
	_ = ensureColumnExists("syllabus_titles", "model_id", "INT")
	_ = ensureColumnExists("syllabus_titles", "title", "VARCHAR(512) NOT NULL")
	_ = ensureColumnExists("syllabus_titles", "hours", "INT DEFAULT 0")
	_ = ensureColumnExists("syllabus_titles", "position", "INT DEFAULT 0")
	_, _ = DB.Exec("CREATE INDEX IF NOT EXISTS idx_titles_model ON syllabus_titles(model_id)")
	// topics
	if _, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS syllabus_topics (
			id INT AUTO_INCREMENT PRIMARY KEY,
			title_id INT NOT NULL,
			content TEXT NOT NULL,
			position INT DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT fk_topics_title FOREIGN KEY (title_id) REFERENCES syllabus_titles(id) ON DELETE CASCADE
		) ENGINE=InnoDB`); err != nil {
		return err
	}
	_ = ensureColumnExists("syllabus_topics", "title_id", "INT")
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
