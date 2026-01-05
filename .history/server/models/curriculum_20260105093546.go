package models

type Semester struct {
	ID             int `json:"id"`
	RegulationID   int `json:"regulation_id"`
	SemesterNumber int `json:"semester_number"`
}

type Course struct {
	CourseID       int    `json:"id"`
	CourseCode     string `json:"course_code"`
	CourseName     string `json:"course_name"`
	CourseType     string `json:"course_type"`
	Category       string `json:"category"`
	Credit         int    `json:"credit"`
	TheoryHours    int    `json:"theory_hours"`
	ActivityHours  int    `json:"activity_hours"`
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
	CourseID       int    `json:"id"`
	CourseCode     string `json:"course_code"`
	CourseName     string `json:"course_name"`
	CourseType     string `json:"course_type"`
	Category       string `json:"category"`
	Credit         int    `json:"credit"`
	TheoryHours    int    `json:"theory_hours"`
	ActivityHours  int    `json:"activity_hours"`
	LectureHours   int    `json:"lecture_hours"`
	TutorialHours  int    `json:"tutorial_hours"`
	PracticalHours int    `json:"practical_hours"`
	TotalHours     int    `json:"total_hours"`
	CIAMarks       int    `json:"cia_marks"`
	SEEMarks       int    `json:"see_marks"`
	TotalMarks     int    `json:"total_marks"`
	RegCourseID    int    `json:"reg_course_id"`
}

type HonourCard struct {
	ID             int    `json:"id"`
	RegulationID   int    `json:"regulation_id"`
	Title          string `json:"title"`
	SemesterNumber int    `json:"semester_number"`
}

type HonourVertical struct {
	ID           int    `json:"id"`
	HonourCardID int    `json:"honour_card_id"`
	Name         string `json:"name"`
}

type HonourVerticalWithCourses struct {
	ID           int                 `json:"id"`
	HonourCardID int                 `json:"honour_card_id"`
	Name         string              `json:"name"`
	Courses      []CourseWithDetails `json:"courses"`
}

type HonourCardWithVerticals struct {
	ID             int                         `json:"id"`
	RegulationID   int                         `json:"regulation_id"`
	Title          string                      `json:"title"`
	SemesterNumber int                         `json:"semester_number"`
	Verticals      []HonourVerticalWithCourses `json:"verticals"`
}
