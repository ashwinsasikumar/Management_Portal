package models

type Teamwork struct {
	Activities []string `json:"activities"`
	Hours      int      `json:"hours"`
}

type SelfLearningInternal struct {
	Main     string   `json:"main"`
	Internal []string `json:"internal"`
}

type SelfLearning struct {
	MainInputs []SelfLearningInternal `json:"main_inputs"`
	Hours      int                    `json:"hours"`
}

type Syllabus struct {
	ID            int           `json:"id"`
	CourseID      int           `json:"course_id"`
	Objectives    []string      `json:"objectives"`
	Outcomes      []string      `json:"outcomes"`
	ReferenceList []string      `json:"reference_list"`
	Prerequisites []string      `json:"prerequisites"`
	Teamwork      *Teamwork     `json:"teamwork,omitempty"`
	SelfLearning  *SelfLearning `json:"selflearning,omitempty"`
}
