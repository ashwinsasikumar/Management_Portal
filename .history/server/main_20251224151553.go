package main

import (
	"fmt"
	"log"
	"net/http"

	"server/db"
	"server/handlers"
	"server/middleware"

	"github.com/gorilla/mux"
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

	// Ensure course syllabus table and prerequisites column
	if err := db.CreateCourseSyllabusTable(); err != nil {
		log.Fatal("Failed to create course_syllabus table:", err)
	}


	// Relational syllabus tables (models, titles, topics)
	if err := db.CreateSyllabusRelationalTables(); err != nil {
		log.Fatal("Failed to create relational syllabus tables:", err)
	}

	// Setup routes with Gorilla Mux for path parameters
	router := mux.NewRouter()

	// Curriculum routes
	router.HandleFunc("/api/curriculum", handlers.GetRegulations).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/curriculum/create", handlers.CreateRegulation).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/curriculum/delete", handlers.DeleteRegulation).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/curriculum/{id}", handlers.UpdateCurriculum).Methods("PUT", "OPTIONS")

	// Department Overview routes
	router.HandleFunc("/api/regulation/{id}/overview", handlers.GetDepartmentOverview).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/regulation/{id}/overview", handlers.SaveDepartmentOverview).Methods("POST", "OPTIONS")

	// Curriculum routes
	router.HandleFunc("/api/regulation/{id}/semesters", handlers.GetSemesters).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/regulation/{id}/semester", handlers.CreateSemester).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/semester/{id}", handlers.UpdateSemester).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/regulation/{id}/semester/{semId}/courses", handlers.GetSemesterCourses).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/regulation/{id}/semester/{semId}/course", handlers.AddCourseToSemester).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/course/{id}", handlers.GetCourse).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/course/{id}", handlers.UpdateCourse).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/regulation/{id}/semester/{semId}/course/{courseId}", handlers.RemoveCourseFromSemester).Methods("DELETE", "OPTIONS")

	// Syllabus routes
	// Return nested syllabus (header + models/titles/topics)
	router.HandleFunc("/api/course/{courseId}/syllabus", handlers.GetCourseSyllabusNested).Methods("GET", "OPTIONS")
	// Save header-only fields (outcomes, resources, prerequisites)
	router.HandleFunc("/api/course/{courseId}/syllabus", handlers.SaveCourseSyllabus).Methods("POST", "OPTIONS")

	// Relational CRUD
	router.HandleFunc("/api/course/{courseId}/syllabus/model", handlers.CreateModel).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/syllabus/model/{modelId}", handlers.UpdateModel).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/syllabus/model/{modelId}", handlers.DeleteModel).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/syllabus/model/{modelId}/title", handlers.CreateTitle).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/syllabus/title/{titleId}", handlers.UpdateTitle).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/syllabus/title/{titleId}", handlers.DeleteTitle).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/syllabus/title/{titleId}/topic", handlers.CreateTopic).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/syllabus/topic/{topicId}", handlers.UpdateTopic).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/syllabus/topic/{topicId}", handlers.DeleteTopic).Methods("DELETE", "OPTIONS")

	// CO-PO and CO-PSO Mapping routes
	router.HandleFunc("/api/course/{courseId}/mapping", handlers.GetCourseMapping).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/course/{courseId}/mapping", handlers.SaveCourseMapping).Methods("POST", "OPTIONS")

	// PEO-PO Mapping routes
	router.HandleFunc("/api/regulation/{id}/peo-po-mapping", handlers.GetPEOPOMapping).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/regulation/{id}/peo-po-mapping", handlers.SavePEOPOMapping).Methods("POST", "OPTIONS")

	// Curriculum Logs routes
	router.HandleFunc("/api/curriculum/{id}/log", handlers.CreateCurriculumLog).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/curriculum/{id}/logs", handlers.GetCurriculumLogs).Methods("GET", "OPTIONS")

	// PDF Generation route
	router.HandleFunc("/api/regulation/{id}/pdf", handlers.GenerateRegulationPDF).Methods("GET", "OPTIONS")

	// Wrap with CORS middleware
	handler := middleware.CORSMiddleware(router)

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
