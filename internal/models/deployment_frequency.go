package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DeploymentFrequency represents the frequency of deployments
type DeploymentFrequency struct {
	DataflowID        primitive.ObjectID `bson:"dataflow_id" json:"dataflow_id"`
	StartDate         time.Time          `bson:"start_date" json:"start_date"`
	EndDate           time.Time          `bson:"end_date" json:"end_date"`
	Window            int                `bson:"window" json:"window"`
	Dates             []time.Time        `bson:"date" json:"date"`
	DailyPipelineRuns []int              `bson:"daily_pipeline_runs" json:"daily_pipeline_runs"`
	MovingAverages    []float64          `bson:"moving_averages" json:"moving_averages"`
}

// DeploymentFrequencyRequest represents the request to calculate the deployment frequency
type DeploymentFrequencyRequest struct {
	DataflowID primitive.ObjectID `bson:"dataflow_id" json:"dataflow_id"`
	StartDate  time.Time          `bson:"start_date" json:"start_date"`
	EndDate    time.Time          `bson:"end_date" json:"end_date"`
	Window     int                `bson:"window" json:"window"`
}
