package models

type Semester struct {
	ID             int `json:"id"`
	RegulationID   int `json:"regulation_id"`
	SemesterNumber int `json:"semester_number"`
}

type Course struct {
	ID             int    `json:"id"`
	CourseCode     string `json:"course_code"`
	CourseName     string `json:"course_name"`
	CourseType     string `json:"course_type"`
	Category       string `json:"category"`
	Credit         int    `json:"credit"`
	LectureHours   int    `json:"lecture_hours"`
	TutorialHours  int    `json:"tutorial_hours"`
	PracticalHours int    `json:"practical_hours"`
	TotalHours     int    `json:"total_hours"`
	CIAMarks       int    `json:"cia_marks"`
	SEEMarks       int    `json:"see_marks"`
	TotalMarks     int    `json:"total_marks"`
}

type RegulationCourse struct {
	ID           int `json:"id"`
	RegulationID int `json:"regulation_id"`
	SemesterID   int `json:"semester_id"`
	CourseID     int `json:"course_id"`
}

type CourseWithDetails struct {
	ID             int    `json:"id"`
	CourseCode     string `json:"course_code"`
	CourseName     string `json:"course_name"`
	CourseType     string `json:"course_type"`
	Category       string `json:"category"`
	Credit         int    `json:"credit"`
	LectureHours   int    `json:"lecture_hours"`
	TutorialHours  int    `json:"tutorial_hours"`
	PracticalHours int    `json:"practical_hours"`
	TotalHours     int    `json:"total_hours"`
	CIAMarks       int    `json:"cia_marks"`
	SEEMarks       int    `json:"see_marks"`
	TotalMarks     int    `json:"total_marks"`
	RegCourseID    int    `json:"reg_course_id"`
}
