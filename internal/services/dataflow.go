package services

import (
	"context"

	"github.com/unnmdnwb3/dora/internal/connectors/gitlab"
	"github.com/unnmdnwb3/dora/internal/daos"
	"github.com/unnmdnwb3/dora/internal/models"
)

// InitializeDataflow creates a new Dataflow, gets the historical data
// and sets up webhooks for new events from the provided sources
func InitializeDataflow(ctx context.Context, dataflow *models.Dataflow) (*models.Dataflow, error) {
	dataflow, err := CreateDataflow(ctx, dataflow)
	if err != nil {
		return nil, err
	}

	err = ImportHistoricalData(ctx, dataflow)
	err = CreateWebhooks(ctx, dataflow)
	return dataflow, err
}

// CreateDataflow creates a new Dataflow.
func CreateDataflow(ctx context.Context, dataflow *models.Dataflow) (*models.Dataflow, error) {
	err := daos.CreateDataflow(ctx, dataflow)
	return dataflow, err
}

// ImportHistoricalData gets and persists historical data for each defined source in a Dataflow.
func ImportHistoricalData(ctx context.Context, dataflow *models.Dataflow) error {
	commitsChannel := make(chan error)
	defer close(commitsChannel)
	go ImportCommits(ctx, commitsChannel, dataflow.Repository)

	runsChannel := make(chan error)
	defer close(runsChannel)
	go ImportPipelineRuns(ctx, runsChannel, dataflow.Pipeline)

	incidentsChannel := make(chan error)
	defer close(incidentsChannel)
	go ImportIncidents(ctx, incidentsChannel, dataflow.Deployment)

	err := <-commitsChannel
	if err != nil {
		return err
	}

	err = <-runsChannel
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

	pipelines, err := client.GetPipelineRuns(pipeline.ExternalID, pipeline.DefaultBranch)
	if err != nil {
		channel <- err
		return
	}

	err = daos.CreatePipelineRuns(ctx, pipelines)
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
