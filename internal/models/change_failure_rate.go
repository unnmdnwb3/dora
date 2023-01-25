package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ChangeFailureRate represents the ratio between the number of incidents and the number of successful pipeline runs of a specific dataflow.
type ChangeFailureRate struct {
	DataflowID       primitive.ObjectID `bson:"dataflow_id" json:"dataflow_id"`
	StartDate        time.Time          `bson:"start_date" json:"start_date"`
	EndDate          time.Time          `bson:"end_date" json:"end_date"`
	Window           int                `bson:"window" json:"window"`
	Dates            []time.Time        `bson:"date" json:"date"`
	DailyIncidents   []int              `bson:"daily_incidents" json:"daily_incidents"`
	DailyDeployments []int              `bson:"daily_deployments" json:"daily_deployments"`
	MovingAverages   []float64          `bson:"moving_averages" json:"moving_averages"`
}

// GeneralChangeFailureRate represents the general ratio between the number of incidents and the number of successful pipeline runs.
type GeneralChangeFailureRate struct {
	StartDate        time.Time   `bson:"start_date" json:"start_date"`
	EndDate          time.Time   `bson:"end_date" json:"end_date"`
	Window           int         `bson:"window" json:"window"`
	Dates            []time.Time `bson:"date" json:"date"`
	DailyIncidents   []int       `bson:"daily_incidents" json:"daily_incidents"`
	DailyDeployments []int       `bson:"daily_deployments" json:"daily_deployments"`
	MovingAverages   []float64   `bson:"moving_averages" json:"moving_averages"`
}
