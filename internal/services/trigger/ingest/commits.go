package ingest

import (
	"context"
	"log"

	"github.com/unnmdnwb3/dora/internal/connectors/gitlab"
	"github.com/unnmdnwb3/dora/internal/daos"
	"github.com/unnmdnwb3/dora/internal/models"
)

// ImportCommits gets and persists historical data for each commit in a repository.
func ImportCommits(ctx context.Context, channel chan error, repository *models.Repository) {
	var integration models.Integration
	err := daos.GetIntegration(ctx, repository.IntegrationID, &integration)
	if err != nil {
		channel <- err
		return
	}

	client := gitlab.NewClient(integration.URI, integration.BearerToken)
	if err != nil {
		channel <- err
		return
	}

	commits, err := client.GetCommits(repository.ExternalID, repository.DefaultBranch)
	if err != nil {
		channel <- err
		return
	}

	err = daos.CreateCommits(ctx, repository.ID, commits)
	if err != nil {
		channel <- err
		return
	}

	log.Printf("Created %d commits for repository %s", len(*commits), repository.NamespacedName)

	channel <- nil
}
