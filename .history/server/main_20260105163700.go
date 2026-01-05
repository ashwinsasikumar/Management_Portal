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

	// Create curriculum table
	if err := db.CreateCurriculumTable(); err != nil {
		log.Fatal("Failed to create curriculum table:", err)
	}

	// Create courses table (must happen before AddVisibilitySemestersCourses)
	if err := db.CreateCoursesTable(); err != nil {
		log.Fatal("Failed to create courses table:", err)
	}

	// Create curriculum courses junction table
	if err := db.CreateCurriculumCoursesTable(); err != nil {
		log.Fatal("Failed to create curriculum_courses table:", err)
	}

	// Create normal_cards table (must happen before any normal_cards migrations)
	if err := db.CreateNormalCardsTable(); err != nil {
		log.Fatal("Failed to create normal_cards table:", err)
	}

	// Rename semesters table to normal_cards (must happen after CreateNormalCardsTable)
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

	// Remove name column from normal_cards table (card_type is now used as the display name)
	if err := db.RemoveNameColumnFromNormalCards(); err != nil {
		log.Fatal("Failed to remove name column from normal_cards:", err)
	}

	// Setup routes
	router := routes.SetupRoutes()

	// Wrap with CORS middleware
	handler := middleware.CORSMiddleware(router)

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
