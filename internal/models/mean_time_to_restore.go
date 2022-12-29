package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MeanTimeToRestore represents the mean time to restore
type MeanTimeToRestore struct {
	DataflowID     primitive.ObjectID `bson:"dataflow_id" json:"dataflow_id"`
	Dates          []time.Time        `bson:"date" json:"date"`
	DailyIncidents []int              `bson:"daily_incidents" json:"daily_incidents"`
	DailyDurations []int              `bson:"daily_durations" json:"daily_durations"`
	MovingAverages []float64          `bson:"moving_average" json:"moving_average"`
}

// MeanTimeToRestoreRequest represents the request to calculate the mean time to restore
type MeanTimeToRestoreRequest struct {
	DataflowID primitive.ObjectID `bson:"dataflow_id" json:"dataflow_id"`
	EndDate    time.Time          `bson:"end_date" json:"end_date"`
	Window     int                `bson:"window" json:"window"`
}
