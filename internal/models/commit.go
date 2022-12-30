package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Commit describes a single commit of a repository
type Commit struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	RepositoryID primitive.ObjectID `json:"repository_id" bson:"repository_id"`
	Sha          string             `json:"id" bson:"sha"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	ParentIds    []string           `json:"parent_ids" bson:"parent_ids"`
}
