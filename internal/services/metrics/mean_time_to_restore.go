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

// CalculateMeanTimeToRestore calculates the mean time to restore for a given dataflow.
func CalculateMeanTimeToRestore(ctx context.Context, dataflowID primitive.ObjectID, window int, endDate time.Time) (*models.MeanTimeToRestore, error) {
	var dataflow models.Dataflow
	err := daos.GetDataflow(ctx, dataflowID, &dataflow)
	if err != nil {
		return nil, err
	}

	offset := window - 1
	timeRange := offset * 2
	startDate := times.Date(endDate.AddDate(0, 0, -timeRange))

	var incidentsPerDays []models.IncidentsPerDay
	filter := bson.M{"deployment_id": dataflow.Deployment.ID, "date": bson.M{"$gte": startDate, "$lte": endDate}}
	err = daos.ListIncidentsPerDaysByFilter(ctx, filter, &incidentsPerDays)
	if err != nil {
		return nil, err
	}
	if len(incidentsPerDays) == 0 {
		return nil, fmt.Errorf("no incidents per days found for deployment with id: %s", dataflow.Deployment.ID.Hex())
	}

	dates, err := DatesBetween(startDate, endDate)
	if err != nil {
		return nil, err
	}

	dailyIncidents, dailyDurations, err := CompleteIncidentsPerDays(&incidentsPerDays, dates)
	if err != nil {
		return nil, err
	}

	movingAverages, err := CalculateMovingAverages(dailyDurations, window)
	if err != nil {
		return nil, err
	}

	return &models.MeanTimeToRestore{
		DataflowID:     dataflow.ID,
		Dates:          (*dates)[offset:],
		DailyIncidents: (*dailyIncidents)[offset:],
		DailyDurations: (*dailyDurations)[offset:],
		MovingAverages: *movingAverages,
	}, nil
}
