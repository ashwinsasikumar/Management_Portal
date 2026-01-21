package models

type PEOPOMapping struct {
	ID           int `json:"id"`
	CurriculumID int `json:"curriculum_id"`
	PEOIndex     int `json:"peo_index"`
	POIndex      int `json:"po_index"`
	MappingValue int `json:"mapping_value"`
}

type PEOPOMappingResponse struct {
	Matrix map[string]int `json:"matrix"` // key: "peo_index-po_index"
}

type PEOPOMappingRequest struct {
	Mappings []PEOPOMapping `json:"mappings"`
}
