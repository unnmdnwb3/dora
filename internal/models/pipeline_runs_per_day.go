package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PipelineRunsPerDay represents the daily pipeline runs.
type PipelineRunsPerDay struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`
	PipelineID        primitive.ObjectID `bson:"pipeline_id" json:"pipeline_id"`
	Date              time.Time          `bson:"date" json:"date"`
	TotalPipelineRuns int                `bson:"total_pipeline_runs" json:"total_pipeline_runs"`
}
