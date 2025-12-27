package models

type COPOMapping struct {
	ID           int `json:"id"`
	CourseID     int `json:"course_id"`
	COIndex      int `json:"co_index"`
	POIndex      int `json:"po_index"`
	MappingValue int `json:"mapping_value"`
}

type COPSOMapping struct {
	ID           int `json:"id"`
	CourseID     int `json:"course_id"`
	COIndex      int `json:"co_index"`
	PSOIndex     int `json:"pso_index"`
	MappingValue int `json:"mapping_value"`
}

type MappingResponse struct {
	COs         []string       `json:"cos"`
	COPOMatrix  map[string]int `json:"co_po_matrix"`  // key: "co_index-po_index"
	COPSOMatrix map[string]int `json:"co_pso_matrix"` // key: "co_index-pso_index"
}

type MappingRequest struct {
	COPOMatrix  []COPOMapping  `json:"co_po_matrix"`
	COPSOMatrix []COPSOMapping `json:"co_pso_matrix"`
}
