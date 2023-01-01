package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ChangeFailureRate represents the ratio between the number of incidents and the number of successful pipeline runs
type ChangeFailureRate struct {
	DataflowID       primitive.ObjectID `bson:"dataflow_id" json:"dataflow_id"`
	Dates            []time.Time        `bson:"date" json:"date"`
	DailyIncidents   []int              `bson:"daily_incidents" json:"daily_incidents"`
	DailyDeployments []int              `bson:"daily_deployments" json:"daily_deployments"`
	MovingAverages   []float64          `bson:"moving_averages" json:"moving_averages"`
}

// ChangeFailureRateRequest represents the request to calculate the change failure rate
type ChangeFailureRateRequest struct {
	DataflowID primitive.ObjectID `bson:"dataflow_id" json:"dataflow_id"`
	EndDate    time.Time          `bson:"end_date" json:"end_date"`
	Window     int                `bson:"window" json:"window"`
}
