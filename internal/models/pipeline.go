package models

// Pipeline describes a workflow tool for CICD
type Pipeline struct {
	ID             string `bson:"_id,omitempty"`
	IntegrationID  string `bson:"integration_id"`
	ExternalID     string `bson:"external_id"`
	DefaultBranch  string `bson:"default_branch"`
	NamespacedName string `bson:"path_with_namespace"`
	URI            string `bson:"web_url"`
}
