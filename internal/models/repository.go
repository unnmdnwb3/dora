package models

// Repository described the body of version control
type Repository struct {
	ID             string `bson:"_id,omitempty"`
	IntegrationID  string `bson:"integration_id"`
	ExternalID     string `bson:"external_id"`
	NamespacedName string `json:"path_with_namespace"`
	DefaultBranch  string `json:"default_branch"`
	URI            string `json:"web_url"`
}
