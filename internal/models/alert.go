package models

import "time"

// Alert describes a single alert.
type Alert struct {
	CreatedAt time.Time `json:"created_at"`
}
