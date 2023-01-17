package ingest

import (
	"context"

	"github.com/unnmdnwb3/dora/internal/models"
)

// All gets and persists historical data for each defined source in a Dataflow.
func All(ctx context.Context, dataflow *models.Dataflow) error {
	err := Raw(ctx, dataflow)
	if err != nil {
		return err
	}

	err = Advanced(ctx, dataflow)
	return err
}

// Raw gets and persists raw historical data for each defined source in a Dataflow.
func Raw(ctx context.Context, dataflow *models.Dataflow) error {
	commitsChannel := make(chan error)
	defer close(commitsChannel)
	go ImportCommits(ctx, commitsChannel, &dataflow.Repository)

	pipelineRunsChannel := make(chan error)
	defer close(pipelineRunsChannel)
	go ImportPipelineRuns(ctx, pipelineRunsChannel, &dataflow.Pipeline)

	// need to wait for each channel get a message, otherwise send on a closed channel will panic
	commitsErr := <-commitsChannel
	pipelineRunsErr := <-pipelineRunsChannel

	if commitsErr != nil {
		return commitsErr
	}

	return pipelineRunsErr
}

// Advanced gets and persists advanced historical data based on raw data for each defined source in a Dataflow.
func Advanced(ctx context.Context, dataflow *models.Dataflow) error {
	changesChannel := make(chan error)
	defer close(changesChannel)
	go ImportChanges(ctx, changesChannel, dataflow.Repository.ID, dataflow.Pipeline.ID)

	incidentsChannel := make(chan error)
	defer close(incidentsChannel)
	go ImportIncidents(ctx, incidentsChannel, &dataflow.Deployment)

	// need to wait for each channel get a message, otherwise send on a closed channel will panic
	changesErr := <-changesChannel
	incidentsErr := <-incidentsChannel

	if changesErr != nil {
		return changesErr
	}
	return incidentsErr
}
