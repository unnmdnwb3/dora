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

// DeploymentFrequency calculates the deployment frequency for a specific dataflow.
func DeploymentFrequency(ctx context.Context, dataflowID primitive.ObjectID, startDate time.Time, endDate time.Time, window int) (*models.DeploymentFrequency, error) {
	if window < 1 {
		return nil, fmt.Errorf("window must be greater than 0")
	}

	if startDate.After(endDate) {
		return nil, fmt.Errorf("start date must be before end date")
	}

	var dataflow models.Dataflow
	err := daos.GetDataflow(ctx, dataflowID, &dataflow)
	if err != nil {
		return nil, fmt.Errorf("error getting dataflow: %w", err)
	}

	offset := window - 1
	startDate = times.Date(startDate.AddDate(0, 0, -offset))

	var pipelineRunsPerDay []models.PipelineRunsPerDay
	filter := bson.M{"pipeline_id": dataflow.Pipeline.ID, "date": bson.M{"$gte": startDate, "$lte": endDate}}
	err = daos.ListPipelineRunsPerDaysByFilter(ctx, filter, &pipelineRunsPerDay)
	if err != nil {
		return nil, fmt.Errorf("error getting pipeline runs per days: %w", err)
	}

	dates, err := DatesBetween(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error getting dates between %s and %s: %w", startDate, endDate, err)
	}

	dailyPipelineRuns, err := CompletePipelineRunsPerDays(&pipelineRunsPerDay, dates)
	if err != nil {
		return nil, fmt.Errorf("error completing pipeline runs per days: %w", err)
	}

	movingAverages, err := MovingAverages(dailyPipelineRuns, window)
	if err != nil {
		return nil, fmt.Errorf("error calculating moving averages: %w", err)
	}

	return &models.DeploymentFrequency{
		DataflowID:        dataflow.ID,
		StartDate:         startDate,
		EndDate:           endDate,
		Window:            window,
		Dates:             (*dates)[offset:],
		DailyPipelineRuns: (*dailyPipelineRuns)[offset:],
		MovingAverages:    (*movingAverages),
	}, nil
}

// GeneralDeploymentFrequency calculates the general deployment frequency over all dataflows.
func GeneralDeploymentFrequency(ctx context.Context, startDate time.Time, endDate time.Time, window int) (*models.GeneralDeploymentFrequency, error) {
	if window < 1 {
		return nil, fmt.Errorf("window must be greater than 0")
	}

	if startDate.After(endDate) {
		return nil, fmt.Errorf("start date must be before end date")
	}

	offset := window - 1
	startDate = times.Date(startDate.AddDate(0, 0, -offset))

	var pipelineRunsPerDay []models.PipelineRunsPerDay
	filter := bson.M{"date": bson.M{"$gte": startDate, "$lte": endDate}}
	err := daos.ListPipelineRunsPerDaysByFilter(ctx, filter, &pipelineRunsPerDay)
	if err != nil {
		return nil, fmt.Errorf("error getting pipeline runs per days: %w", err)
	}

	dates, err := DatesBetween(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error getting dates between %s and %s: %w", startDate, endDate, err)
	}

	dailyPipelineRuns, err := CompletePipelineRunsPerDays(&pipelineRunsPerDay, dates)
	if err != nil {
		return nil, fmt.Errorf("error completing pipeline runs per days: %w", err)
	}

	movingAverages, err := MovingAverages(dailyPipelineRuns, window)
	if err != nil {
		return nil, fmt.Errorf("error calculating moving averages: %w", err)
	}

	return &models.GeneralDeploymentFrequency{
		StartDate:         startDate,
		EndDate:           endDate,
		Window:            window,
		Dates:             (*dates)[offset:],
		DailyPipelineRuns: (*dailyPipelineRuns)[offset:],
		MovingAverages:    (*movingAverages),
	}, nil
}
