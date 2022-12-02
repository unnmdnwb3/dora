package models

import "time"

// Commit describes a single commit of a repository
type Commit struct {
	ID        string    `json:"id"`
	ShortID   string    `json:"short_id"`
	CreatedAt time.Time `json:"created_at"`
	ParentIds []string  `json:"parent_ids"`
	Title     string    `json:"title"`
}
