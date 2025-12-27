package routes

import (
	"server/handlers"

	"github.com/gorilla/mux"
)

// SetupRoutes configures all application routes
func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Curriculum routes
	router.HandleFunc("/api/curriculum", handlers.GetRegulations).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/curriculum/create", handlers.CreateRegulation).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/curriculum/delete", handlers.DeleteRegulation).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/curriculum/{id}", handlers.UpdateCurriculum).Methods("PUT", "OPTIONS")

	// Department Overview routes
	router.HandleFunc("/api/regulation/{id}/overview", handlers.GetDepartmentOverview).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/regulation/{id}/overview", handlers.SaveDepartmentOverview).Methods("POST", "OPTIONS")

	// Semester routes
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

	// Cluster Management routes
	router.HandleFunc("/api/clusters", handlers.GetClusters).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/clusters", handlers.CreateCluster).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/cluster/{id}/departments", handlers.GetClusterDepartments).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/cluster/{id}/department", handlers.AddDepartmentToCluster).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/cluster/{id}/department/{deptId}", handlers.RemoveDepartmentFromCluster).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/cluster/{id}/objects", handlers.GetClusterObjects).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/cluster/{id}/objects", handlers.SaveClusterObjects).Methods("POST", "OPTIONS")

	return router
}
