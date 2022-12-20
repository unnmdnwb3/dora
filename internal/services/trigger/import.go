package trigger

import (
	"context"

	"github.com/unnmdnwb3/dora/internal/connectors/gitlab"
	"github.com/unnmdnwb3/dora/internal/daos"
	"github.com/unnmdnwb3/dora/internal/models"
)

// ImportData gets and persists historical data for each defined source in a Dataflow.
func ImportData(ctx context.Context, dataflow *models.Dataflow) error {
	commitsChannel := make(chan error)
	defer close(commitsChannel)
	go ImportCommits(ctx, commitsChannel, &dataflow.Repository)

	pipelineRunsChannel := make(chan error)
	defer close(pipelineRunsChannel)
	go ImportPipelineRuns(ctx, pipelineRunsChannel, &dataflow.Pipeline)

	incidentsChannel := make(chan error)
	defer close(incidentsChannel)
	go ImportIncidents(ctx, incidentsChannel, &dataflow.Deployment)

	err := <-commitsChannel
	if err != nil {
		return err
	}

	err = <-pipelineRunsChannel
	if err != nil {
		return err
	}

	err = <-incidentsChannel
	return err

}

// ImportCommits gets and persists historical data for each commit in a repository.
func ImportCommits(ctx context.Context, channel chan error, repository *models.Repository) {
	// TODO implement
	channel <- nil
}

// ImportPipelineRuns gets and persists historical data for each run of a pipeline.
func ImportPipelineRuns(ctx context.Context, channel chan error, pipeline *models.Pipeline) {
	client, err := gitlab.NewClient()
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

// ImportIncidents gets and persists historical data for each incident of a deployment.
func ImportIncidents(ctx context.Context, channel chan error, deployment *models.Deployment) {
	// TODO implement
	channel <- nil
}

// CreateWebhooks creates a Webhook for each defined source in a Dataflow.
func CreateWebhooks(ctx context.Context, dataflow *models.Dataflow) error {
	err := CreateWebhook(ctx, dataflow.Repository)
	if err != nil {
		return err
	}

	err = CreateWebhook(ctx, dataflow.Pipeline)
	if err != nil {
		return err
	}

	err = CreateWebhook(ctx, dataflow.Deployment)
	return err
}

// CreateWebhook creates a Webhook for a defined source.
func CreateWebhook(ctx context.Context, v any) error {
	// TODO implement
	return nil
}
