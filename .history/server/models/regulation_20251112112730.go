package models

import "time"

type Regulation struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	AcademicYear string    `json:"academic_year"`
	MaxCredits   int       `json:"max_credits"`
	CreatedAt    time.Time `json:"created_at"`
}
