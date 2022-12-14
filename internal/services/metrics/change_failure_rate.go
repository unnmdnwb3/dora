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

// CalculateChangeFailureRate calculates the change failure rate for a given dataflow.
func CalculateChangeFailureRate(ctx context.Context, dataflowID primitive.ObjectID, window int, endDate time.Time) (*models.ChangeFailureRate, error) {
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

	var pipelineRunsPerDays []models.PipelineRunsPerDay
	filter = bson.M{"pipeline_id": dataflow.Pipeline.ID, "date": bson.M{"$gte": startDate, "$lte": endDate}}
	err = daos.ListPipelineRunsPerDaysByFilter(ctx, filter, &pipelineRunsPerDays)
	if err != nil {
		return nil, err
	}
	if len(pipelineRunsPerDays) == 0 {
		return nil, fmt.Errorf("no pipline runs per days found for pipeline with id: %s", dataflow.Pipeline.ID.Hex())
	}

	dates, err := DatesBetween(startDate, endDate)
	if err != nil {
		return nil, err
	}

	dailyIncidents, _, err := CompleteIncidentsPerDays(&incidentsPerDays, dates)
	if err != nil {
		return nil, err
	}

	dailyDeployments, err := CompletePipelineRunsPerDays(&pipelineRunsPerDays, dates)
	if err != nil {
		return nil, err
	}

	movingAverages, err := CalculateMovingAveragesRatio(dailyIncidents, dailyDeployments, window)
	if err != nil {
		return nil, err
	}

	return &models.ChangeFailureRate{
		DataflowID:       dataflow.ID,
		Dates:            (*dates)[offset:],
		DailyIncidents:   (*dailyIncidents)[offset:],
		DailyDeployments: (*dailyDeployments)[offset:],
		MovingAverages:   *movingAverages,
	}, nil
}
