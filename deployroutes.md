# local

package routes

import (
	"net/http"
	"server/handlers"

	"github.com/gorilla/mux"
)

// SetupRoutes configures all application routes
func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Curriculum routes
	router.HandleFunc("/curriculum", handlers.GetRegulations).Methods("GET", "OPTIONS")
	router.HandleFunc("/curriculum/create", handlers.CreateRegulation).Methods("POST", "OPTIONS")
	router.HandleFunc("/curriculum/delete", handlers.DeleteRegulation).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/curriculum/{id}", handlers.UpdateCurriculum).Methods("PUT", "OPTIONS")

	// NEW Regulation Management routes (isolated from curriculum)
	router.HandleFunc("/regulations", handlers.GetRegulationsNew).Methods("GET", "OPTIONS")
	router.HandleFunc("/regulations", handlers.CreateRegulationNew).Methods("POST", "OPTIONS")
	router.HandleFunc("/regulations/{id}", handlers.GetRegulationByID).Methods("GET", "OPTIONS")
	router.HandleFunc("/regulations/{id}", handlers.UpdateRegulationNew).Methods("PUT", "OPTIONS")
	router.HandleFunc("/regulations/{id}", handlers.DeleteRegulationNew).Methods("DELETE", "OPTIONS")

	// Regulation Clauses routes
	router.HandleFunc("/regulations/{id}/clauses", handlers.GetRegulationClauses).Methods("GET", "OPTIONS")
	router.HandleFunc("/regulations/{id}/clauses", handlers.CreateRegulationClause).Methods("POST", "OPTIONS")
	router.HandleFunc("/regulations/clauses/{clauseId}", handlers.UpdateRegulationClause).Methods("PUT", "OPTIONS")
	router.HandleFunc("/regulations/clauses/{clauseId}", handlers.DeleteRegulationClause).Methods("DELETE", "OPTIONS")

	// Regulation Editor routes (structured editing)
	router.HandleFunc("/regulations/{id}/structure", handlers.GetRegulationStructure).Methods("GET", "OPTIONS")

	// Section management
	router.HandleFunc("/regulations/{id}/sections", handlers.CreateSection).Methods("POST", "OPTIONS")
	router.HandleFunc("/regulations/sections/{sectionId}", handlers.UpdateSection).Methods("PUT", "OPTIONS")
	router.HandleFunc("/regulations/sections/{sectionId}", handlers.DeleteSection).Methods("DELETE", "OPTIONS")

	// Clause management
	router.HandleFunc("/regulations/sections/{sectionId}/clauses", handlers.CreateClause).Methods("POST", "OPTIONS")
	router.HandleFunc("/regulations/clauses/{clauseId}", handlers.UpdateClause).Methods("PUT", "OPTIONS")
	router.HandleFunc("/regulations/clauses/{clauseId}", handlers.DeleteClause).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/regulations/clauses/{clauseId}/history", handlers.GetClauseHistory).Methods("GET", "OPTIONS")

	// Department Overview routes
	router.HandleFunc("/regulation/{id}/overview", handlers.GetDepartmentOverview).Methods("GET", "OPTIONS")
	router.HandleFunc("/regulation/{id}/overview", handlers.SaveDepartmentOverview).Methods("POST", "OPTIONS")

	// Semester routes
	router.HandleFunc("/regulation/{id}/semesters", handlers.GetSemesters).Methods("GET", "OPTIONS")
	router.HandleFunc("/regulation/{id}/semester", handlers.CreateSemester).Methods("POST", "OPTIONS")
	router.HandleFunc("/semester/{id}", handlers.UpdateSemester).Methods("PUT", "OPTIONS")
	router.HandleFunc("/semester/{id}", handlers.DeleteSemester).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/regulation/{id}/semester/{semId}/courses", handlers.GetSemesterCourses).Methods("GET", "OPTIONS")
	router.HandleFunc("/regulation/{id}/semester/{semId}/course", handlers.AddCourseToSemester).Methods("POST", "OPTIONS")
	router.HandleFunc("/course/{id}", handlers.GetCourse).Methods("GET", "OPTIONS")
	router.HandleFunc("/course/{id}", handlers.UpdateCourse).Methods("PUT", "OPTIONS")
	router.HandleFunc("/regulation/{id}/semester/{semId}/course/{courseId}", handlers.RemoveCourseFromSemester).Methods("DELETE", "OPTIONS")

	// Honour Card routes
	router.HandleFunc("/regulation/{id}/honour-cards", handlers.GetHonourCards).Methods("GET", "OPTIONS")
	router.HandleFunc("/regulation/{id}/honour-card", handlers.CreateHonourCard).Methods("POST", "OPTIONS")
	router.HandleFunc("/honour-card/{cardId}", handlers.DeleteHonourCard).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/honour-card/{cardId}/vertical", handlers.CreateHonourVertical).Methods("POST", "OPTIONS")
	router.HandleFunc("/honour-vertical/{verticalId}", handlers.DeleteHonourVertical).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/honour-vertical/{verticalId}/course", handlers.AddCourseToVertical).Methods("POST", "OPTIONS")
	router.HandleFunc("/honour-vertical/{verticalId}/course/{courseId}", handlers.RemoveCourseFromVertical).Methods("DELETE", "OPTIONS")

	// Syllabus routes
	// Return nested syllabus (header + models/titles/topics)
	router.HandleFunc("/course/{courseId}/syllabus", handlers.GetCourseSyllabusNested).Methods("GET", "OPTIONS")
	// Save header-only fields (outcomes, resources, prerequisites)
	router.HandleFunc("/course/{courseId}/syllabus", handlers.SaveCourseSyllabus).Methods("POST", "OPTIONS")

	// Relational CRUD
	router.HandleFunc("/course/{courseId}/syllabus/model", handlers.CreateModel).Methods("POST", "OPTIONS")
	router.HandleFunc("/syllabus/model/{modelId}", handlers.UpdateModel).Methods("PUT", "OPTIONS")
	router.HandleFunc("/syllabus/model/{modelId}", handlers.DeleteModel).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/syllabus/model/{modelId}/title", handlers.CreateTitle).Methods("POST", "OPTIONS")
	router.HandleFunc("/syllabus/title/{titleId}", handlers.UpdateTitle).Methods("PUT", "OPTIONS")
	router.HandleFunc("/syllabus/title/{titleId}", handlers.DeleteTitle).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/syllabus/title/{titleId}/topic", handlers.CreateTopic).Methods("POST", "OPTIONS")
	router.HandleFunc("/syllabus/topic/{topicId}", handlers.UpdateTopic).Methods("PUT", "OPTIONS")
	router.HandleFunc("/syllabus/topic/{topicId}", handlers.DeleteTopic).Methods("DELETE", "OPTIONS")

	// CO-PO and CO-PSO Mapping routes
	router.HandleFunc("/course/{courseId}/mapping", handlers.GetCourseMapping).Methods("GET", "OPTIONS")
	router.HandleFunc("/course/{courseId}/mapping", handlers.SaveCourseMapping).Methods("POST", "OPTIONS")

	// PEO-PO Mapping routes
	router.HandleFunc("/regulation/{id}/peo-po-mapping", handlers.GetPEOPOMapping).Methods("GET", "OPTIONS")
	router.HandleFunc("/regulation/{id}/peo-po-mapping", handlers.SavePEOPOMapping).Methods("POST", "OPTIONS")

	// Experiments routes (2022 template)
	router.HandleFunc("/course/{courseId}/experiments", handlers.GetCourseExperiments).Methods("GET", "OPTIONS")
	router.HandleFunc("/course/{courseId}/experiments", handlers.CreateExperiment).Methods("POST", "OPTIONS")
	router.HandleFunc("/experiments/{expId}", handlers.UpdateExperiment).Methods("PUT", "OPTIONS")
	router.HandleFunc("/experiments/{expId}", handlers.DeleteExperiment).Methods("DELETE", "OPTIONS")

	// Curriculum Logs routes
	router.HandleFunc("/curriculum/{id}/log", handlers.CreateCurriculumLog).Methods("POST", "OPTIONS")
	router.HandleFunc("/curriculum/{id}/logs", handlers.GetCurriculumLogs).Methods("GET", "OPTIONS")

	// PDF Generation route
	router.HandleFunc("/regulation/{id}/pdf", handlers.GenerateRegulationPDFHTML).Methods("GET", "OPTIONS")

	// Cluster Management routes
	router.HandleFunc("/clusters", handlers.GetClusters).Methods("GET", "OPTIONS")
	router.HandleFunc("/clusters", handlers.CreateCluster).Methods("POST", "OPTIONS")
	router.HandleFunc("/clusters/available-departments", handlers.GetAvailableDepartments).Methods("GET", "OPTIONS")
	router.HandleFunc("/cluster/{id}", handlers.DeleteCluster).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/cluster/{id}/departments", handlers.GetClusterDepartments).Methods("GET", "OPTIONS")
	router.HandleFunc("/cluster/{id}/department", handlers.AddDepartmentToCluster).Methods("POST", "OPTIONS")
	router.HandleFunc("/cluster/{id}/department/{deptId}", handlers.RemoveDepartmentFromCluster).Methods("DELETE", "OPTIONS")

	// Sharing Management routes
	router.HandleFunc("/regulation/{id}/sharing", handlers.GetDepartmentSharingInfo).Methods("GET", "OPTIONS")
	router.HandleFunc("/sharing/visibility", handlers.UpdateItemVisibility).Methods("PUT", "OPTIONS")
	router.HandleFunc("/sharing/{item_type}/{item_id}/departments", handlers.GetItemSharedDepartments).Methods("GET", "OPTIONS")
	router.HandleFunc("/cluster/{id}/shared-content", handlers.GetClusterSharedContent).Methods("GET", "OPTIONS")

	// Authentication routes
	router.HandleFunc("/auth/login", handlers.Login).Methods("POST", "OPTIONS")

	// User Management routes
	router.HandleFunc("/users", handlers.GetUsers).Methods("GET", "OPTIONS")
	router.HandleFunc("/users", handlers.CreateUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/users/{id}", handlers.GetUser).Methods("GET", "OPTIONS")
	router.HandleFunc("/users/{id}", handlers.UpdateUser).Methods("PUT", "OPTIONS")
	router.HandleFunc("/users/{id}", handlers.DeleteUser).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/users/{id}/password", handlers.ChangePassword).Methods("PUT", "OPTIONS")

	return router
}







# deploy

package routes

import (
	"net/http"
	"server/handlers"

	"github.com/gorilla/mux"
)

// SetupRoutes configures all application routes
func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Health check endpoint
	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Curriculum routes
	router.HandleFunc("/api/curriculum", handlers.GetRegulations).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/curriculum/create", handlers.CreateRegulation).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/curriculum/delete", handlers.DeleteRegulation).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/curriculum/{id}", handlers.UpdateCurriculum).Methods("PUT", "OPTIONS")

	// Regulation Management routes
	router.HandleFunc("/api/regulations", handlers.GetRegulationsNew).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/regulations", handlers.CreateRegulationNew).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/regulations/{id}", handlers.GetRegulationByID).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/regulations/{id}", handlers.UpdateRegulationNew).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/regulations/{id}", handlers.DeleteRegulationNew).Methods("DELETE", "OPTIONS")

	// Regulation Clauses routes
	router.HandleFunc("/api/regulations/{id}/clauses", handlers.GetRegulationClauses).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/regulations/{id}/clauses", handlers.CreateRegulationClause).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/regulations/clauses/{clauseId}", handlers.UpdateRegulationClause).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/regulations/clauses/{clauseId}", handlers.DeleteRegulationClause).Methods("DELETE", "OPTIONS")

	// Regulation Editor routes
	router.HandleFunc("/api/regulations/{id}/structure", handlers.GetRegulationStructure).Methods("GET", "OPTIONS")

	// Section management
	router.HandleFunc("/api/regulations/{id}/sections", handlers.CreateSection).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/regulations/sections/{sectionId}", handlers.UpdateSection).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/regulations/sections/{sectionId}", handlers.DeleteSection).Methods("DELETE", "OPTIONS")

	// Clause management
	router.HandleFunc("/api/regulations/sections/{sectionId}/clauses", handlers.CreateClause).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/regulations/clauses/{clauseId}", handlers.UpdateClause).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/regulations/clauses/{clauseId}", handlers.DeleteClause).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/regulations/clauses/{clauseId}/history", handlers.GetClauseHistory).Methods("GET", "OPTIONS")

	// Department Overview routes
	router.HandleFunc("/api/regulation/{id}/overview", handlers.GetDepartmentOverview).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/regulation/{id}/overview", handlers.SaveDepartmentOverview).Methods("POST", "OPTIONS")

	// Semester routes
	router.HandleFunc("/api/regulation/{id}/semesters", handlers.GetSemesters).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/regulation/{id}/semester", handlers.CreateSemester).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/semester/{id}", handlers.UpdateSemester).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/semester/{id}", handlers.DeleteSemester).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/regulation/{id}/semester/{semId}/courses", handlers.GetSemesterCourses).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/regulation/{id}/semester/{semId}/course", handlers.AddCourseToSemester).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/course/{id}", handlers.GetCourse).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/course/{id}", handlers.UpdateCourse).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/regulation/{id}/semester/{semId}/course/{courseId}", handlers.RemoveCourseFromSemester).Methods("DELETE", "OPTIONS")

	// Honour Card routes
	router.HandleFunc("/api/regulation/{id}/honour-cards", handlers.GetHonourCards).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/regulation/{id}/honour-card", handlers.CreateHonourCard).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/honour-card/{cardId}", handlers.DeleteHonourCard).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/honour-card/{cardId}/vertical", handlers.CreateHonourVertical).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/honour-vertical/{verticalId}", handlers.DeleteHonourVertical).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/honour-vertical/{verticalId}/course", handlers.AddCourseToVertical).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/honour-vertical/{verticalId}/course/{courseId}", handlers.RemoveCourseFromVertical).Methods("DELETE", "OPTIONS")

	// Syllabus routes
	router.HandleFunc("/api/course/{courseId}/syllabus", handlers.GetCourseSyllabusNested).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/course/{courseId}/syllabus", handlers.SaveCourseSyllabus).Methods("POST", "OPTIONS")

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

	// Experiments routes
	router.HandleFunc("/api/course/{courseId}/experiments", handlers.GetCourseExperiments).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/course/{courseId}/experiments", handlers.CreateExperiment).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/experiments/{expId}", handlers.UpdateExperiment).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/experiments/{expId}", handlers.DeleteExperiment).Methods("DELETE", "OPTIONS")

	// Curriculum Logs routes
	router.HandleFunc("/api/curriculum/{id}/log", handlers.CreateCurriculumLog).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/curriculum/{id}/logs", handlers.GetCurriculumLogs).Methods("GET", "OPTIONS")

	// PDF Generation
	router.HandleFunc("/api/regulation/{id}/pdf", handlers.GenerateRegulationPDFHTML).Methods("GET", "OPTIONS")

	// Cluster Management routes
	router.HandleFunc("/api/clusters", handlers.GetClusters).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/clusters", handlers.CreateCluster).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/clusters/available-departments", handlers.GetAvailableDepartments).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/cluster/{id}", handlers.DeleteCluster).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/cluster/{id}/departments", handlers.GetClusterDepartments).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/cluster/{id}/department", handlers.AddDepartmentToCluster).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/cluster/{id}/department/{deptId}", handlers.RemoveDepartmentFromCluster).Methods("DELETE", "OPTIONS")

	// Sharing Management routes
	router.HandleFunc("/api/regulation/{id}/sharing", handlers.GetDepartmentSharingInfo).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/sharing/visibility", handlers.UpdateItemVisibility).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/sharing/{item_type}/{item_id}/departments", handlers.GetItemSharedDepartments).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/cluster/{id}/shared-content", handlers.GetClusterSharedContent).Methods("GET", "OPTIONS")

	// Authentication routes
	router.HandleFunc("/api/auth/login", handlers.Login).Methods("POST", "OPTIONS")

	// User Management routes
	router.HandleFunc("/api/users", handlers.GetUsers).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/users", handlers.CreateUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/users/{id}", handlers.GetUser).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/users/{id}", handlers.UpdateUser).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/users/{id}", handlers.DeleteUser).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/users/{id}/password", handlers.ChangePassword).Methods("PUT", "OPTIONS")

	return router
}


# FOR LOCAL USE :
// API Configuration
const API_BASE_URL =
  process.env.REACT_APP_API_URL ||
  "http://localhost:5000";

export { API_BASE_URL };

# FOR DEPLOY USE :

// API Configuration
const API_BASE_URL =
  process.env.REACT_APP_API_URL 

export { API_BASE_URL };

# FOR LOCAL USE :
port : 5001

# for deploy use :
port : 5000