package metrics

import (
	"context"
	"fmt"
	"time"

	"github.com/unnmdnwb3/dora/internal/daos"
	"github.com/unnmdnwb3/dora/internal/models"
	"github.com/unnmdnwb3/dora/internal/utils/times"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// LeadTimeForChanges calculates the lead time for changes.
func LeadTimeForChanges(ctx context.Context, dataflowID primitive.ObjectID, startDate time.Time, endDate time.Time, window int) (*models.LeadTimeForChanges, error) {
	if window < 1 {
		return nil, fmt.Errorf("window must be greater than 0")
	}
	if startDate.After(endDate) {
		return nil, fmt.Errorf("start date must be before end date")
	}

	var dataflow models.Dataflow
	err := daos.GetDataflow(ctx, dataflowID, &dataflow)
	if err != nil {
		return nil, err
	}

	offset := window - 1
	startDate = times.Date(startDate.AddDate(0, 0, -offset))

	var changesPerDay []models.ChangesPerDay
	filter := bson.M{"repository_id": dataflow.Repository.ID, "date": bson.M{"$gte": startDate, "$lte": endDate}}
	err = daos.ListChangesPerDaysByFilter(ctx, filter, &changesPerDay)
	if err != nil {
		return nil, err
	}
	if len(changesPerDay) == 0 {
		return nil, fmt.Errorf("no changes per days found for repository with id: %s", dataflow.Repository.ID.Hex())
	}

	dates, err := DatesBetween(startDate, endDate)
	if err != nil {
		return nil, err
	}

	dailyChanges, dailyLeadTimes, err := CompleteChangesPerDays(&changesPerDay, dates)
	if err != nil {
		return nil, err
	}

	movingAverages, err := MovingAverages(dailyLeadTimes, window)
	if err != nil {
		return nil, err
	}

	return &models.LeadTimeForChanges{
		DataflowID:     dataflow.ID,
		Dates:          (*dates)[offset:],
		DailyChanges:   (*dailyChanges)[offset:],
		DailyLeadTimes: (*dailyLeadTimes)[offset:],
		MovingAverages: *movingAverages,
	}, nil
}
