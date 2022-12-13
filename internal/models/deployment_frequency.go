package models

import "time"

// DeploymentFrequency represents the frequency of deployments
type DeploymentFrequency struct {
	ID                 string       `bson:"_id,omitempty"`
	DataflowID         string       `bson:"dataflow_id"`
	Date               time.Weekday `bson:"date" json:"date"`
	PipelineRunsPerDay int          `bson:"pipeline_runs_per_day" json:"pipeline_runs_per_day"`
	MovingAverage      float64      `bson:"moving_average" json:"moving_average"`
}
