package models

// Repository described the body of version control
type Repository struct {
	ID             string `json:"id"`
	CreatedAt      string `json:"created_at"`
	DefaultBranch  string `json:"default_branch"`
	NamespacedName string `json:"path_with_namespace"`
	URI            string `json:"web_url"`
}
