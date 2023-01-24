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

// MeanTimeToRestore calculates the mean time to restore for a given dataflow.
func MeanTimeToRestore(ctx context.Context, dataflowID primitive.ObjectID, startDate time.Time, endDate time.Time, window int) (*models.MeanTimeToRestore, error) {
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

	var incidentsPerDays []models.IncidentsPerDay
	filter := bson.M{"deployment_id": dataflow.Deployment.ID, "date": bson.M{"$gte": startDate, "$lte": endDate}}
	err = daos.ListIncidentsPerDaysByFilter(ctx, filter, &incidentsPerDays)
	if err != nil {
		return nil, fmt.Errorf("error getting incidents per days: %w", err)
	}

	dates, err := DatesBetween(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error getting dates between %s and %s: %w", startDate, endDate, err)
	}

	dailyIncidents, dailyDurations, err := CompleteIncidentsPerDays(&incidentsPerDays, dates)
	if err != nil {
		return nil, fmt.Errorf("error completing incidents per days: %w", err)
	}

	movingAverages, err := MovingAverages(dailyDurations, window)
	if err != nil {
		return nil, fmt.Errorf("error calculating moving averages: %w", err)
	}

	return &models.MeanTimeToRestore{
		DataflowID:     dataflow.ID,
		StartDate:      startDate,
		EndDate:        endDate,
		Window:         window,
		Dates:          (*dates)[offset:],
		DailyIncidents: (*dailyIncidents)[offset:],
		DailyDurations: (*dailyDurations)[offset:],
		MovingAverages: *movingAverages,
	}, nil
}
