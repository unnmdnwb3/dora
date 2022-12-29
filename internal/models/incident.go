package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Incident describes a single incident.
type Incident struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	DeploymentID primitive.ObjectID `bson:"deployment_id" json:"deployment_id"`
	StartDate    time.Time          `json:"start_date" bson:"start_date"`
	EndDate      time.Time          `json:"end_date" bson:"end_date"`
}
