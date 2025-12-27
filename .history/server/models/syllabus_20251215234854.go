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
	Unit1         []string      `json:"unit1"`
	Unit2         []string      `json:"unit2"`
	Unit3         []string      `json:"unit3"`
	Unit4         []string      `json:"unit4"`
	Unit5         []string      `json:"unit5"`
	Textbooks     []string      `json:"textbooks"`
	ReferenceList []string      `json:"reference_list"`
	Prerequisites []string      `json:"prerequisites"`
	Teamwork      *Teamwork     `json:"teamwork,omitempty"`
	SelfLearning  *SelfLearning `json:"selflearning,omitempty"`
}
