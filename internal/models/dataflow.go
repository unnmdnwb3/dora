package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Dataflow represents a complete dataflow, from repository, to pipeline, to deployment
type Dataflow struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Repository Repository         `bson:"repository" json:"repository"`
	Pipeline   Pipeline           `bson:"pipeline" json:"pipeline"`
	Deployment Deployment         `bson:"deployment" json:"deployment"`
}

// Repository represents a repository used for version control
type Repository struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	IntegrationID  primitive.ObjectID `bson:"integration_id,omitempty" json:"integration_id"`
	ExternalID     int                `bson:"external_id" json:"id"`
	NamespacedName string             `bson:"namespaced_name" json:"namespaced_name"`
	DefaultBranch  string             `bson:"default_branch" json:"default_branch"`
	URI            string             `bson:"uri" json:"uri"`
}

// Pipeline represents a pipeline used for CI/CD
type Pipeline struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	IntegrationID  primitive.ObjectID `bson:"integration_id,omitempty" json:"integration_id"`
	ExternalID     int                `bson:"external_id" json:"external_id"`
	NamespacedName string             `bson:"namespaced_name" json:"namespaced_name"`
	DefaultBranch  string             `bson:"default_branch" json:"default_branch"`
	URI            string             `bson:"uri" json:"uri"`
}

// Deployment represents a running deployment
type Deployment struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	IntegrationID primitive.ObjectID `bson:"integration_id,omitempty" json:"integration_id"`
	TargetURI     string             `bson:"target_uri" json:"target_uri"`
}
