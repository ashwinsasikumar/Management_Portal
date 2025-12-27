package models

// RegulationPDF represents all data needed for PDF generation
type RegulationPDF struct {
	RegulationID   int                `json:"regulation_id"`
	RegulationName string             `json:"regulation_name"`
	AcademicYear   string             `json:"academic_year"`
	Overview       DepartmentOverview `json:"overview"`
	Semesters      []SemesterPDF      `json:"semesters"`
	PEOPOMapping   map[string]int     `json:"peo_po_mapping"`
}

type SemesterPDF struct {
	SemesterNumber int         `json:"semester_number"`
	Courses        []CoursePDF `json:"courses"`
}

type CoursePDF struct {
	ID             int            `json:"id"`
	CourseCode     string         `json:"course_code"`
	CourseName     string         `json:"course_name"`
	CourseType     string         `json:"course_type"`
	Category       string         `json:"category"`
	Credit         int            `json:"credit"`
	LectureHours   int            `json:"lecture_hours"`
	TutorialHours  int            `json:"tutorial_hours"`
	PracticalHours int            `json:"practical_hours"`
	CIAMarks       int            `json:"cia_marks"`
	SEEMarks       int            `json:"see_marks"`
	TotalMarks     int            `json:"total_marks"`
	Syllabus       SyllabusPDF    `json:"syllabus"`
	COPOMapping    map[string]int `json:"co_po_mapping"`
	COPSOMapping   map[string]int `json:"co_pso_mapping"`
}

type SyllabusPDF struct {
	Objectives    []string `json:"objectives"`
	Outcomes      []string `json:"outcomes"`
	ReferenceList []string `json:"reference_list"`
}
