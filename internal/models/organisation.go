package models

// Organisation describes the organisational subject of a shared repository
type Organisation struct {
	ExternalID     int    `json:"id"`
	WebURL         string `json:"web_url"`
	NamespacedName string `json:"full_path"`
}
