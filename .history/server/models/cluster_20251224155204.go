package models

import "time"

// Cluster represents a group of departments that share common objects
type Cluster struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// ClusterDepartment maps departments to clusters
type ClusterDepartment struct {
	ID           int `json:"id"`
	ClusterID    int `json:"cluster_id"`
	DepartmentID int `json:"department_id"`
}

// ClusterObjects holds cluster-level shared objects
type ClusterObjects struct {
	ClusterID int      `json:"cluster_id"`
	Mission   []string `json:"mission"`
	PEOs      []string `json:"peos"`
	POs       []string `json:"pos"`
	PSOs      []string `json:"psos"`
}
