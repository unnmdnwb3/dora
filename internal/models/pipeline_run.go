package models

import "time"

// PipelineRun describes a single run of a CICD pipeline
type PipelineRun struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	PipelineID  string    `json:"pipeline_id" bson:"pipeline_id"`
	ExternalID  string    `bson:"external_id" json:"external_id"`
	Ref         string    `json:"ref" bson:"ref"`
	Status      string    `json:"status" bson:"status"`
	EventSource string    `json:"source" bson:"event_source"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
	URI         string    `json:"web_url" bson:"uri"`
}
