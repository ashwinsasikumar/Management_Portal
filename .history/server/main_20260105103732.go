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

	// Create department overview tables
	if err := db.CreateDepartmentOverviewTables(); err != nil {
		log.Fatal("Failed to create department overview tables:", err)
	}

	// Create cluster tables
	if err := db.CreateClusterTables(); err != nil {
		log.Fatal("Failed to create cluster tables:", err)
	}

	// Add visibility columns to department data tables
	if err := db.AddVisibilityColumns(); err != nil {
		log.Fatal("Failed to add visibility columns:", err)
	}

	// Rename semesters table to normal_cards (must happen before AddVisibilitySemestersCourses)
	if err := db.RenameSemestersToNormalCards(); err != nil {
		log.Fatal("Failed to rename semesters to normal_cards:", err)
	}

	// Add visibility columns to normal_cards and courses tables
	if err := db.AddVisibilitySemestersCourses(); err != nil {
		log.Fatal("Failed to add visibility columns to semesters/courses:", err)
	}

	// Add source tracking columns for shared items
	if err := db.AddSourceDepartmentColumns(); err != nil {
		log.Fatal("Failed to add source tracking columns:", err)
	}

	// Create sharing tracking table
	if err := db.CreateSharingTrackingTable(); err != nil {
		log.Fatal("Failed to create sharing tracking table:", err)
	}

	// Create regulation management tables (PHASE 1 - isolated, zero breakage)
	if err := db.CreateRegulationTables(); err != nil {
		log.Fatal("Failed to create regulation tables:", err)
	}

	// Add regulation reference columns (PHASE 2 - shadow links, nullable)
	if err := db.AddRegulationRefColumns(); err != nil {
		log.Fatal("Failed to add regulation reference columns:", err)
	}

	// Create honour card tables
	if err := db.CreateHonourCardTables(); err != nil {
		log.Fatal("Failed to create honour card tables:", err)
	}

	// Add name column to normal_cards table
	if err := db.AddSemesterNameColumn(); err != nil {
		log.Fatal("Failed to add semester name column:", err)
	}

	// Setup routes
	router := routes.SetupRoutes()

	// Wrap with CORS middleware
	handler := middleware.CORSMiddleware(router)

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
