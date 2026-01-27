package models

import "time"

// CourseAllocation represents a teacher assigned to a course
type CourseAllocation struct {
	ID           int       `json:"id"`
	CourseID     int       `json:"course_id"`
	TeacherID    int       `json:"teacher_id"`
	TeacherName  string    `json:"teacher_name,omitempty"`
	AcademicYear string    `json:"academic_year"`
	Semester     int       `json:"semester"`
	Section      string    `json:"section"`
	Role         string    `json:"role"` // Primary, Assistant
	Status       int       `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CourseWithAllocations represents a course along with its assigned faculty
type CourseWithAllocations struct {
	CourseID    int                `json:"id"`
	CourseCode  string             `json:"course_code"`
	CourseName  string             `json:"course_name"`
	CourseType  string             `json:"course_type"`
	Credit      int                `json:"credit"`
	Allocations []CourseAllocation `json:"allocations"`
}
