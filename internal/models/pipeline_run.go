package models

import "time"

// PipelineRun describes a single run of a CICD pipeline
type PipelineRun struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	ProjectID   string    `json:"project_id" bson:"project_id"`
	Ref         string    `json:"ref" bson:"ref"`
	Status      string    `json:"status" bson:"status"`
	EventSource string    `json:"source" bson:"event_source"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
	URI         string    `json:"web_url" bson:"uri"`
}
