package routes

import (
	"net/http"
	curriculum "server/handlers/curriculum"
	studentteacher "server/handlers/student-teacher_entry"

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
	router.HandleFunc("/api/curriculum", curriculum.GetRegulations).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/curriculum/create", curriculum.CreateRegulation).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/curriculum/delete", curriculum.DeleteRegulation).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/curriculum/{id}", curriculum.UpdateCurriculum).Methods("PUT", "OPTIONS")

	// NEW Regulation Management routes (isolated from curriculum)
	router.HandleFunc("/api/regulations", curriculum.GetRegulationsNew).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/regulations", curriculum.CreateRegulationNew).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/regulations/{id}", curriculum.GetRegulationByID).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/regulations/{id}", curriculum.UpdateRegulationNew).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/regulations/{id}", curriculum.DeleteRegulationNew).Methods("DELETE", "OPTIONS")

	// Regulation Clauses routes
	router.HandleFunc("/api/regulations/{id}/clauses", curriculum.GetRegulationClauses).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/regulations/{id}/clauses", curriculum.CreateRegulationClause).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/regulations/clauses/{clauseId}", curriculum.UpdateRegulationClause).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/regulations/clauses/{clauseId}", curriculum.DeleteRegulationClause).Methods("DELETE", "OPTIONS")

	// Regulation Editor routes (structured editing)
	router.HandleFunc("/api/regulations/{id}/structure", curriculum.GetRegulationStructure).Methods("GET", "OPTIONS")

	// Section management
	router.HandleFunc("/api/regulations/{id}/sections", curriculum.CreateSection).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/regulations/sections/{sectionId}", curriculum.UpdateSection).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/regulations/sections/{sectionId}", curriculum.DeleteSection).Methods("DELETE", "OPTIONS")

	// Clause management
	router.HandleFunc("/api/regulations/sections/{sectionId}/clauses", curriculum.CreateClause).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/regulations/clauses/{clauseId}", curriculum.UpdateClause).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/regulations/clauses/{clauseId}", curriculum.DeleteClause).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/regulations/clauses/{clauseId}/history", curriculum.GetClauseHistory).Methods("GET", "OPTIONS")

	// Department Overview routes
	router.HandleFunc("/api/curriculum/{id}/overview", curriculum.GetDepartmentOverview).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/curriculum/{id}/overview", curriculum.SaveDepartmentOverview).Methods("POST", "OPTIONS")

	// Semester routes
	router.HandleFunc("/api/curriculum/{id}/semesters", curriculum.GetSemesters).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/curriculum/{id}/semester", curriculum.CreateSemester).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/semester/{id}", curriculum.UpdateSemester).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/semester/{id}", curriculum.DeleteSemester).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/curriculum/{id}/semester/{semId}/courses", curriculum.GetSemesterCourses).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/curriculum/{id}/semester/{semId}/course", curriculum.AddCourseToSemester).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/course/{id}", curriculum.GetCourse).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/course/{id}", curriculum.UpdateCourse).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/curriculum/{id}/semester/{semId}/course/{courseId}", curriculum.RemoveCourseFromSemester).Methods("DELETE", "OPTIONS")

	// Honour Card routes
	router.HandleFunc("/api/curriculum/{id}/honour-cards", curriculum.GetHonourCards).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/curriculum/{id}/honour-card", curriculum.CreateHonourCard).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/honour-card/{cardId}", curriculum.DeleteHonourCard).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/honour-card/{cardId}/vertical", curriculum.CreateHonourVertical).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/honour-vertical/{verticalId}", curriculum.DeleteHonourVertical).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/honour-vertical/{verticalId}/course", curriculum.AddCourseToVertical).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/honour-vertical/{verticalId}/course/{courseId}", curriculum.RemoveCourseFromVertical).Methods("DELETE", "OPTIONS")

	// Syllabus routes
	// Return nested syllabus (header + models/titles/topics)
	router.HandleFunc("/api/course/{courseId}/syllabus", curriculum.GetCourseSyllabusNested).Methods("GET", "OPTIONS")
	// Save header-only fields (outcomes, resources, prerequisites)
	router.HandleFunc("/api/course/{courseId}/syllabus", curriculum.SaveCourseSyllabus).Methods("POST", "OPTIONS")

	// Relational CRUD
	router.HandleFunc("/api/course/{courseId}/syllabus/model", curriculum.CreateModel).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/syllabus/model/{modelId}", curriculum.UpdateModel).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/syllabus/model/{modelId}", curriculum.DeleteModel).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/syllabus/model/{modelId}/title", curriculum.CreateTitle).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/syllabus/title/{titleId}", curriculum.UpdateTitle).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/syllabus/title/{titleId}", curriculum.DeleteTitle).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/syllabus/title/{titleId}/topic", curriculum.CreateTopic).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/syllabus/topic/{topicId}", curriculum.UpdateTopic).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/syllabus/topic/{topicId}", curriculum.DeleteTopic).Methods("DELETE", "OPTIONS")

	// CO-PO and CO-PSO Mapping routes
	router.HandleFunc("/api/course/{courseId}/mapping", curriculum.GetCourseMapping).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/course/{courseId}/mapping", curriculum.SaveCourseMapping).Methods("POST", "OPTIONS")

	// PEO-PO Mapping routes
	router.HandleFunc("/api/curriculum/{id}/peo-po-mapping", curriculum.GetPEOPOMapping).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/curriculum/{id}/peo-po-mapping", curriculum.SavePEOPOMapping).Methods("POST", "OPTIONS")

	// Experiments routes (2022 template)
	router.HandleFunc("/api/course/{courseId}/experiments", curriculum.GetCourseExperiments).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/course/{courseId}/experiments", curriculum.CreateExperiment).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/experiments/{expId}", curriculum.UpdateExperiment).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/experiments/{expId}", curriculum.DeleteExperiment).Methods("DELETE", "OPTIONS")

	// Curriculum Logs routes
	router.HandleFunc("/api/curriculum/{id}/log", curriculum.CreateCurriculumLog).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/curriculum/{id}/logs", curriculum.GetCurriculumLogs).Methods("GET", "OPTIONS")

	// PDF Generation route
	router.HandleFunc("/api/curriculum/{id}/pdf", curriculum.GenerateRegulationPDFHTML).Methods("GET", "OPTIONS")

	// Cluster Management routes
	router.HandleFunc("/api/clusters", curriculum.GetClusters).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/clusters", curriculum.CreateCluster).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/clusters/available-departments", curriculum.GetAvailableDepartments).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/cluster/{id}", curriculum.DeleteCluster).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/cluster/{id}/departments", curriculum.GetClusterDepartments).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/cluster/{id}/department", curriculum.AddDepartmentToCluster).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/cluster/{id}/department/{deptId}", curriculum.RemoveDepartmentFromCluster).Methods("DELETE", "OPTIONS")

	// Sharing Management routes
	router.HandleFunc("/api/curriculum/{id}/sharing", curriculum.GetDepartmentSharingInfo).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/sharing/visibility", curriculum.UpdateItemVisibility).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/sharing/{item_type}/{item_id}/departments", curriculum.GetItemSharedDepartments).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/cluster/{id}/shared-content", curriculum.GetClusterSharedContent).Methods("GET", "OPTIONS")

	// Authentication routes
	router.HandleFunc("/api/auth/login", curriculum.Login).Methods("POST", "OPTIONS")

	// User Management routes
	router.HandleFunc("/api/users", curriculum.GetUsers).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/users", curriculum.CreateUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/users/{id}", curriculum.GetUser).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/users/{id}", curriculum.UpdateUser).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/users/{id}", curriculum.DeleteUser).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/users/{id}/password", curriculum.ChangePassword).Methods("PUT", "OPTIONS")

	// Student-Teacher Entry routes
	router.HandleFunc("/api/students", studentteacher.GetStudents).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/students/{id}", studentteacher.GetStudent).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/students", studentteacher.CreateStudent).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/students/{id}", studentteacher.UpdateStudent).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/students/{id}", studentteacher.DeleteStudent).Methods("DELETE", "OPTIONS")

	return router
}
