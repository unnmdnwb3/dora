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

// LeadTimeForChanges calculates the lead time for changes for a specific dataflow.
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
		return nil, fmt.Errorf("error listing changes per days: %w", err)
	}

	dates, err := DatesBetween(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error getting dates between %s and %s: %w", startDate, endDate, err)
	}

	dailyChanges, dailyLeadTimes, err := CompleteChangesPerDays(&changesPerDay, dates)
	if err != nil {
		return nil, fmt.Errorf("error completing changes per days: %w", err)
	}

	movingAverages, err := MovingAverages(dailyLeadTimes, window)
	if err != nil {
		return nil, fmt.Errorf("error calculating moving averages: %w", err)
	}

	return &models.LeadTimeForChanges{
		DataflowID:     dataflow.ID,
		Dates:          (*dates)[offset:],
		StartDate:      startDate,
		EndDate:        endDate,
		Window:         window,
		DailyChanges:   (*dailyChanges)[offset:],
		DailyLeadTimes: (*dailyLeadTimes)[offset:],
		MovingAverages: *movingAverages,
	}, nil
}

// GeneralLeadTimeForChanges calculates the general lead time for changes over all dataflows.
func GeneralLeadTimeForChanges(ctx context.Context, startDate time.Time, endDate time.Time, window int) (*models.GeneralLeadTimeForChanges, error) {
	if window < 1 {
		return nil, fmt.Errorf("window must be greater than 0")
	}

	if startDate.After(endDate) {
		return nil, fmt.Errorf("start date must be before end date")
	}

	offset := window - 1
	startDate = times.Date(startDate.AddDate(0, 0, -offset))

	var changesPerDay []models.ChangesPerDay
	filter := bson.M{"date": bson.M{"$gte": startDate, "$lte": endDate}}
	err := daos.ListChangesPerDaysByFilter(ctx, filter, &changesPerDay)
	if err != nil {
		return nil, fmt.Errorf("error listing changes per days: %w", err)
	}

	dates, err := DatesBetween(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error getting dates between %s and %s: %w", startDate, endDate, err)
	}

	dailyChanges, dailyLeadTimes, err := CompleteChangesPerDays(&changesPerDay, dates)
	if err != nil {
		return nil, fmt.Errorf("error completing changes per days: %w", err)
	}

	movingAverages, err := MovingAverages(dailyLeadTimes, window)
	if err != nil {
		return nil, fmt.Errorf("error calculating moving averages: %w", err)
	}

	return &models.GeneralLeadTimeForChanges{
		Dates:          (*dates)[offset:],
		StartDate:      startDate,
		EndDate:        endDate,
		Window:         window,
		DailyChanges:   (*dailyChanges)[offset:],
		DailyLeadTimes: (*dailyLeadTimes)[offset:],
		MovingAverages: *movingAverages,
	}, nil
}
