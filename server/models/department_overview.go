package models

type DepartmentOverview struct {
	ID           int                  `json:"id"`
	CurriculumID int                  `json:"curriculum_id"`
	Vision       string               `json:"vision"`
	Mission      []DepartmentListItem `json:"mission"`
	PEOs         []DepartmentListItem `json:"peos"`
	POs          []DepartmentListItem `json:"pos"`
	PSOs         []DepartmentListItem `json:"psos"`
}

type DepartmentListItem struct {
	ID                 int    `json:"id,omitempty"`
	Text               string `json:"text"`
	Visibility         string `json:"visibility"` // "UNIQUE" or "CLUSTER"
	SourceDepartmentID int    `json:"source_department_id,omitempty"`
}

type Cluster struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
}

type ClusterDepartment struct {
	ID           int    `json:"id"`
	ClusterID    int    `json:"cluster_id"`
	DepartmentID int    `json:"department_id"`
	CreatedAt    string `json:"created_at,omitempty"`
}
