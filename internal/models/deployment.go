package models

// Deployment describes a deployed and monitored runnable
type Deployment struct {
	IntegrationID string `bson:"integration_id" json:"integration_id"`
	TargetURI     string `bson:"target_uri" json:"target_uri"`
}
