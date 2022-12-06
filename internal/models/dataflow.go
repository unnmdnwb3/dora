package models

// Dataflow describes a toolchain including source-control, cicd and incident management
type Dataflow struct {
	ID         string      `bson:"_id,omitempty"`
	Repository *Repository `bson:"repository" json:"repository"`
	Pipeline   *Pipeline   `bson:"pipeline" json:"pipeline"`
	Deployment *Deployment `bson:"deployment" json:"deployment"`
}
