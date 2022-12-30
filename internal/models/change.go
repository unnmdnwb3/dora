package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Change describes a single change from first commit to deployment.
type Change struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	RepositoryID    primitive.ObjectID `json:"repository_id" bson:"repository_id"`
	FirstCommitDate time.Time          `json:"first_commit_date" bson:"first_commit_date"`
	DeploymentDate  time.Time          `json:"deployment_date" bson:"deployment_date"`
}
