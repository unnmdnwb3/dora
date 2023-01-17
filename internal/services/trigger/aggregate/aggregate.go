package aggregate

import (
	"context"

	"github.com/unnmdnwb3/dora/internal/models"
)

// All aggregates the data per day.
func All(ctx context.Context, dataflow *models.Dataflow) error {
	cpdChannel := make(chan error)
	defer close(cpdChannel)
	go CreateChangesPerDays(ctx, cpdChannel, dataflow.Repository.ID, dataflow.Pipeline.ID)

	ipdChannel := make(chan error)
	defer close(ipdChannel)
	go CreateIncidentsPerDays(ctx, ipdChannel, dataflow.Deployment.ID)

	prpdChannel := make(chan error)
	defer close(prpdChannel)
	go CreatePipelineRunsPerDays(ctx, prpdChannel, dataflow.Pipeline.ID)

	// need to wait for each channel get a message, otherwise send on a closed channel will panic
	cpdErr := <-cpdChannel
	ipdErr := <-ipdChannel
	prpdErr := <-prpdChannel

	if cpdErr != nil {
		return cpdErr
	}

	if ipdErr != nil {
		return ipdErr
	}

	if prpdErr != nil {
		return prpdErr
	}

	return nil
}
