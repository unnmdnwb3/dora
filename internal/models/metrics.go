package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MetricsRequest represents a generic metrics request body for a specific dataflow.
type MetricsRequest struct {
	DataflowID primitive.ObjectID `bson:"dataflow_id" json:"dataflow_id"`
	StartDate  time.Time          `bson:"start_date" json:"start_date"`
	EndDate    time.Time          `bson:"end_date" json:"end_date"`
	Window     int                `bson:"window" json:"window"`
}

// GeneralMetricsRequest represents a general generic metrics request body.
type GeneralMetricsRequest struct {
	StartDate time.Time `bson:"start_date" json:"start_date"`
	EndDate   time.Time `bson:"end_date" json:"end_date"`
	Window    int       `bson:"window" json:"window"`
}
