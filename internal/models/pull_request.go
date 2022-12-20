package models

import "time"

// PullRequest describes a successfully merged pull-request
type PullRequest struct {
	ID               int       `json:"id"`
	ProjectID        int       `json:"project_id"`
	Title            string    `json:"title"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	TargetBranch     string    `json:"target_branch"`
	SourceBranch     string    `json:"source_branch"`
	PreCommitTailSha string    `json:"sha"`
	MergeCommitSha   string    `json:"merge_commit_sha"` // could also be a squashed commit
	Reference        string    `json:"reference"`
	WebURL           string    `json:"web_url"`
}
