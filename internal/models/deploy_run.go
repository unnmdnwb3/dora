package models

import "time"

// DeployRun describes a deployment run
type DeployRun struct {
	ID        int       `json:"id"`
	ProjectID int       `json:"project_id"`
	Ref       string    `json:"ref"`
	Status    string    `json:"status"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	URI       string    `json:"web_url"`
}
