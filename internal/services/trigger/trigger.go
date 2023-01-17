package trigger

import (
	"context"

	"github.com/unnmdnwb3/dora/internal/models"
	"github.com/unnmdnwb3/dora/internal/services/trigger/aggregate"
	"github.com/unnmdnwb3/dora/internal/services/trigger/ingest"
)

// OnNewDataflow gets the historical data, creates the necessary aggregates
// and sets up webhooks for new events from the provided sources
func OnNewDataflow(ctx context.Context, dataflow *models.Dataflow) error {
	err := ingest.All(ctx, dataflow)
	if err != nil {
		return err
	}

	err = aggregate.All(ctx, dataflow)
	return err
}
