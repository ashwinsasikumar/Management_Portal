package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"crypto/tls"
	"crypto/x509"
	"strings"

	"github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	caPEM := os.Getenv("DB_CA_CERT")

	// Setup TLS
	rootCertPool := x509.NewCertPool()
	caPEM = strings.ReplaceAll(caPEM, "\\n", "\n")
	if ok := rootCertPool.AppendCertsFromPEM([]byte(caPEM)); !ok {
		log.Fatal("Failed to parse CA certificate")
	}

	mysql.RegisterTLSConfig("custom", &tls.Config{
		RootCAs: rootCertPool,
	})

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?tls=custom&parseTime=true", user, pass, host, port, name)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer db.Close()

	tables := []string{
		"cluster_psos",
		"cluster_pos",
		"cluster_peos",
		"cluster_mission",
		"cluster_departments",
		"clusters",
	}

	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table))
		if err != nil {
			log.Printf("Failed to drop table %s: %v", table, err)
		} else {
			fmt.Printf("Dropped table: %s\n", table)
		}
	}

	fmt.Println("All cluster tables dropped successfully!")
}
