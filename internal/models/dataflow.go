package models

// Dataflow describes a toolchain including source-control, cicd and incident management
type Dataflow struct {
	ID         string      `bson:"_id,omitempty"`
	repository *Repository `bson:"repository"`
	pipeline   *Pipeline   `bson:"pipeline"`
	instance   *Deployment `bson:"deployment"`
}
