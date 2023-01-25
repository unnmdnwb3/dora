package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MeanTimeToRestore represents the mean time to restore after an incident for a specific dataflow.
type MeanTimeToRestore struct {
	DataflowID     primitive.ObjectID `bson:"dataflow_id" json:"dataflow_id"`
	StartDate      time.Time          `bson:"start_date" json:"start_date"`
	EndDate        time.Time          `bson:"end_date" json:"end_date"`
	Window         int                `bson:"window" json:"window"`
	Dates          []time.Time        `bson:"date" json:"date"`
	DailyIncidents []int              `bson:"daily_incidents" json:"daily_incidents"`
	DailyDurations []int              `bson:"daily_durations" json:"daily_durations"`
	MovingAverages []float64          `bson:"moving_averages" json:"moving_averages"`
}

// GeneralMeanTimeToRestore represents the general mean time to restore after an incident.
type GeneralMeanTimeToRestore struct {
	StartDate      time.Time   `bson:"start_date" json:"start_date"`
	EndDate        time.Time   `bson:"end_date" json:"end_date"`
	Window         int         `bson:"window" json:"window"`
	Dates          []time.Time `bson:"date" json:"date"`
	DailyIncidents []int       `bson:"daily_incidents" json:"daily_incidents"`
	DailyDurations []int       `bson:"daily_durations" json:"daily_durations"`
	MovingAverages []float64   `bson:"moving_averages" json:"moving_averages"`
}
