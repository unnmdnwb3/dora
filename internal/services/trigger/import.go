package trigger

import (
	"context"
	"time"

	"github.com/unnmdnwb3/dora/internal/connectors/gitlab"
	"github.com/unnmdnwb3/dora/internal/connectors/prometheus"
	"github.com/unnmdnwb3/dora/internal/daos"
	"github.com/unnmdnwb3/dora/internal/models"
	"github.com/unnmdnwb3/dora/internal/utils/times"
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

	// need to wait for each channel get a message, otherwise send on a closed channel will panic
	commitsErr := <-commitsChannel
	pipelineRunsErr := <-pipelineRunsChannel
	incidentsErr := <-incidentsChannel

	if commitsErr != nil {
		return commitsErr
	}
	if pipelineRunsErr != nil {
		return pipelineRunsErr
	}
	return incidentsErr

}

// ImportCommits gets and persists historical data for each commit in a repository.
func ImportCommits(ctx context.Context, channel chan error, repository *models.Repository) {
	// TODO implement
	channel <- nil
}

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

// ImportIncidents gets and persists historical data for each incident of a deployment.
// However, this functions does not persist the raw monitoring data points, but rather aggregates them already to incidents.
// This is because the raw data points are not relevant for the user, but the incidents are.
// Additionally, persisting the raw data points would be too expensive, especially considering a small step size.
func ImportIncidents(ctx context.Context, channel chan error, deployment *models.Deployment) {
	monitoringDataPoints, err := ImportMonitoringDataPoints(ctx, deployment)
	if err != nil {
		channel <- err
		return
	}

	err = CreateIncidents(ctx, deployment, monitoringDataPoints)
	channel <- err
	return
}

// ImportMonitoringDataPoints gets historical monitoring data.
func ImportMonitoringDataPoints(ctx context.Context, deployment *models.Deployment) (*[]models.MonitoringDataPoint, error) {
	var integration models.Integration
	err := daos.GetIntegration(ctx, deployment.IntegrationID, &integration)
	if err != nil {
		return nil, err
	}

	// TODO standardize the time range for data imports
	end := time.Now()
	start := times.Date(end.AddDate(0, 0, -90))

	client := prometheus.NewClient(integration.URI, integration.BearerToken, deployment.Query, start, end, deployment.Step)
	monitoringDataPoints, err := client.GetMonitoringDataPoints()
	if err != nil {
		return nil, err
	}

	return monitoringDataPoints, nil
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
