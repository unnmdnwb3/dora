package ingest

import (
	"context"

	"github.com/unnmdnwb3/dora/internal/connectors/gitlab"
	"github.com/unnmdnwb3/dora/internal/daos"
	"github.com/unnmdnwb3/dora/internal/models"
)

// ImportPipelineRuns gets and persists historical data for each run of a pipeline.
func ImportPipelineRuns(ctx context.Context, channel chan error, pipeline *models.Pipeline) {
	var integration models.Integration
	err := daos.GetIntegration(ctx, pipeline.IntegrationID, &integration)
	if err != nil {
		channel <- err
		return
	}

	client := gitlab.NewClient(integration.URI, integration.BearerToken)
	if err != nil {
		channel <- err
		return
	}

	pipelineRuns, err := client.GetPipelineRuns(pipeline.ExternalID, pipeline.DefaultBranch)
	if err != nil {
		channel <- err
		return
	}

	err = daos.CreatePipelineRuns(ctx, pipeline.ID, pipelineRuns)
	if err != nil {
		channel <- err
		return
	}

	channel <- nil
	return
}
