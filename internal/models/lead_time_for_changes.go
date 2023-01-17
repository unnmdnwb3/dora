package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// LeadTimeForChanges represents the lead time for changes.
type LeadTimeForChanges struct {
	DataflowID     primitive.ObjectID `bson:"dataflow_id" json:"dataflow_id"`
	StartDate      time.Time          `bson:"start_date" json:"start_date"`
	EndDate        time.Time          `bson:"end_date" json:"end_date"`
	Window         int                `bson:"window" json:"window"`
	Dates          []time.Time        `bson:"date" json:"date"`
	DailyChanges   []int              `bson:"daily_changes" json:"daily_changes"`
	DailyLeadTimes []int              `bson:"daily_lead_times" json:"daily_lead_times"`
	MovingAverages []float64          `bson:"moving_averages" json:"moving_averages"`
}

// LeadTimeForChangesRequest represents the request to calculate the lead time for changes.
type LeadTimeForChangesRequest struct {
	DataflowID primitive.ObjectID `bson:"dataflow_id" json:"dataflow_id"`
	StartDate  time.Time          `bson:"start_date" json:"start_date"`
	EndDate    time.Time          `bson:"end_date" json:"end_date"`
	Window     int                `bson:"window" json:"window"`
}
