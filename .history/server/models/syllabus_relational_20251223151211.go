package models

// SyllabusModel represents a module in the syllabus (e.g., Module I, Module II)
// Links to courses via course_id (course-centric design)
type SyllabusModel struct {
	ID        int             `json:"id"`
	CourseID  int             `json:"course_id"`
	ModelName string          `json:"model_name"`
	Position  int             `json:"position"`
	Titles    []SyllabusTitle `json:"titles,omitempty"`
}

// SyllabusTitle represents a title under a model with hours
// Links to syllabus_models via model_id
type SyllabusTitle struct {
	ID        int             `json:"id"`
	ModelID   int             `json:"model_id"`
	TitleName string          `json:"title_name"`
	Hours     int             `json:"hours"`
	Position  int             `json:"position"`
	Topics    []SyllabusTopic `json:"topics,omitempty"`
}

// SyllabusTopic represents individual topic lines under a title
// Links to syllabus_titles via title_id
type SyllabusTopic struct {
	ID       int    `json:"id"`
	TitleID  int    `json:"title_id"`
	Topic    string `json:"topic"`
	Position int    `json:"position"`
}

// CourseSyllabusResponse is the complete nested response structure
type CourseSyllabusResponse struct {
	Header Syllabus        `json:"header"`
	Models []SyllabusModel `json:"models"`
}
