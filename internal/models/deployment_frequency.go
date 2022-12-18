package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// DeploymentFrequency represents the frequency of deployments
type DeploymentFrequency struct {
	DataflowID        primitive.ObjectID `bson:"dataflow_id" json:"dataflow_id"`
	Dates             []string           `bson:"date" json:"date"`
	DailyPipelineRuns []int              `bson:"daily_pipeline_runs" json:"daily_pipeline_runs"`
	MovingAverages    []float64          `bson:"moving_average" json:"moving_average"`
}
