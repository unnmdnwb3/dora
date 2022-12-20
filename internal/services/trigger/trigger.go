package trigger

import (
	"context"

	"github.com/unnmdnwb3/dora/internal/models"
)

// OnNewDataflow gets the historical data, creates the necessary aggregates
// and sets up webhooks for new events from the provided sources
func OnNewDataflow(ctx context.Context, dataflow *models.Dataflow) error {
	err := ImportData(ctx, dataflow)
	if err != nil {
		return err
	}

	err = CreatePipelineRunsPerDays(ctx, dataflow.Pipeline.ID)
	if err != nil {
		return err
	}

	err = CreateWebhooks(ctx, dataflow)
	return err
}
