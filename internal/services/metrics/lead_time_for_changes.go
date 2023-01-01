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

// CalculateLeadTimeForChanges calculates the lead time for changes.
func CalculateLeadTimeForChanges(ctx context.Context, dataflowID primitive.ObjectID, window int, endDate time.Time) (*models.LeadTimeForChanges, error) {
	var dataflow models.Dataflow
	err := daos.GetDataflow(ctx, dataflowID, &dataflow)
	if err != nil {
		return nil, err
	}

	offset := window - 1
	timeRange := offset * 2
	startDate := times.Date(endDate.AddDate(0, 0, -timeRange))

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

	movingAverages, err := CalculateMovingAverages(dailyLeadTimes, window)
	if err != nil {
		return nil, err
	}

	return &models.LeadTimeForChanges{
		DataflowID:     dataflow.ID,
		Dates:          *dates,
		DailyChanges:   (*dailyChanges)[offset:],
		DailyLeadTimes: (*dailyLeadTimes)[offset:],
		MovingAverages: *movingAverages,
	}, nil
}
