package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ChangesPerDay represents the daily changes.
type ChangesPerDay struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	RepositoryID  primitive.ObjectID `bson:"repository_id" json:"repository_id"`
	Date          time.Time          `bson:"date" json:"date"`
	TotalChanges  int                `bson:"total_changes" json:"total_changes"`
	TotalLeadTime float64            `bson:"total_lead_time" json:"total_lead_time"`
}
