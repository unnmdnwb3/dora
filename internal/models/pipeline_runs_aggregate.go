package models

import "time"

// PipelineRunsAggregate represents the pipeline runs aggregate.
type PipelineRunsAggregate struct {
	ID                string    `bson:"_id,omitempty"`
	PipelineID        string    `bson:"pipeline_id" json:"pipeline_id"`
	Date              time.Time `bson:"date" json:"date"`
	TotalPipelineRuns int       `bson:"total_pipeline_runs" json:"total_pipeline_runs"`
}
