package models

import "time"

// WorkflowRun describes a single run of a CICD workflow
type WorkflowRun struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	ProjectID string    `json:"project_id" bson:"project_id"`
	Ref       string    `json:"ref" bson:"ref"`
	Status    string    `json:"status" bson:"status"`
	Source    string    `json:"source" bson:"source"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	URI       string    `json:"web_url" bson:"uri"`
}
