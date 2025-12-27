package main

import (
	"fmt"
	"log"
	"net/http"

	"server/db"
	"server/middleware"
	"server/routes"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env in server directory
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found or failed to load; using environment defaults")
	}
	// Initialize database
	err := db.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.CloseDB()

	

	// Relational syllabus tables (models, titles, topics)
	if err := db.CreateSyllabusRelationalTables(); err != nil {
		log.Fatal("Failed to create relational syllabus tables:", err)
	}

	// Setup routes
	router := routes.SetupRoutes()

	// Wrap with CORS middleware
	handler := middleware.CORSMiddleware(router)

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
