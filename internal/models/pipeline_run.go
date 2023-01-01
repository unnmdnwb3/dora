package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PipelineRun describes a single run of a CICD pipeline
type PipelineRun struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	PipelineID  primitive.ObjectID `bson:"pipeline_id" json:"pipeline_id"`
	ExternalID  int                `bson:"external_id" json:"id"`
	Sha         string             `json:"sha" bson:"sha"`
	Ref         string             `json:"ref" bson:"ref"`
	Status      string             `json:"status" bson:"status"`
	EventSource string             `json:"source" bson:"event_source"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
	URI         string             `json:"web_url" bson:"uri"`
}
