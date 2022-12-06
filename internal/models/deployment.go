package models

// Deployment describes a deployed and monitored runnable
type Deployment struct {
	ID            string `bson:"_id,omitempty"`
	IntegrationID string `bson:"integration_id"`
	TargetURI     string `bson:"target_uri"`
}
