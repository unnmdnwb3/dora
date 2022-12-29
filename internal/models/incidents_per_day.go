package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// IncidentsPerDay represents the daily pipeline runs.
type IncidentsPerDay struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	DeploymentID   primitive.ObjectID `bson:"deployment_id" json:"deployment_id"`
	Date           time.Time          `bson:"date" json:"date"`
	TotalIncidents int                `bson:"total_incidents" json:"total_incidents"`
	TotalDuration  float64            `bson:"total_duration" json:"total_duration"`
}
