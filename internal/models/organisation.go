package models

// Organisation describes the organisational subject of a shared repository
type Organisation []struct {
	ID   string `json:"id"`
	URI  string `json:"web_url"`
	Name string `json:"full_path"`
}
