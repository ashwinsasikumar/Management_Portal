package models

import "time"

type CurriculumLog struct {
	ID           int       `json:"id"`
	CurriculumID int       `json:"curriculum_id"`
	Action       string    `json:"action"`
	Description  string    `json:"description"`
	ChangedBy    string    `json:"changed_by"`
	CreatedAt    time.Time `json:"created_at"`
}
