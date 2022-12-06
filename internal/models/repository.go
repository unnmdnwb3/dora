package models

// Repository described the body of version control
type Repository struct {
	IntegrationID  string `bson:"integration_id" json:"integration_id"`
	ExternalID     string `bson:"external_id" json:"external_id"`
	NamespacedName string `bson:"namespaced_name" json:"namespaced_name"`
	DefaultBranch  string `bson:"default_branch" json:"default_branch"`
	URI            string `bson:"uri" json:"uri"`
}
