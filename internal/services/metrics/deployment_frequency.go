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

// CalculateDeploymentFrequency calculates the deployment frequency for a given dataflow.
func CalculateDeploymentFrequency(ctx context.Context, dataflowID primitive.ObjectID, window int, endDate time.Time) (*models.DeploymentFrequency, error) {
	var dataflow models.Dataflow
	err := daos.GetDataflow(ctx, dataflowID, &dataflow)
	if err != nil {
		return nil, err
	}

	offset := window - 1
	timeRange := offset * 2
	startDate := times.Date(endDate.AddDate(0, 0, -timeRange))

	var pipelineRunsPerDay []models.PipelineRunsPerDay
	filter := bson.M{"pipeline_id": dataflow.Pipeline.ID, "date": bson.M{"$gte": startDate}}
	err = daos.ListPipelineRunsPerDaysByFilter(ctx, filter, &pipelineRunsPerDay)
	if err != nil {
		return nil, err
	}
	if len(pipelineRunsPerDay) == 0 {
		return nil, fmt.Errorf("no pipeline runs per days found for pipeline with id: %s", dataflow.Pipeline.ID.Hex())
	}

	dates, err := DatesBetween(startDate, endDate)
	if err != nil {
		return nil, err
	}

	dailyPipelineRuns, err := CompletePipelineRunsPerDays(&pipelineRunsPerDay, dates)
	if err != nil {
		return nil, err
	}

	movingAverages, err := CalculateMovingAverages(dailyPipelineRuns, window)
	if err != nil {
		return nil, err
	}

	return &models.DeploymentFrequency{
		DataflowID:        dataflow.ID,
		Dates:             (*dates)[offset:],
		DailyPipelineRuns: (*dailyPipelineRuns)[offset:],
		MovingAverages:    *movingAverages,
	}, nil
}

// CompletePipelineRunsPerDays returns a slice of the number of pipeline runs per day,
// since provided PipelineRunsPerDays only account for the dates that any pipeline runs were found.
func CompletePipelineRunsPerDays(pipelineRunsPerDays *[]models.PipelineRunsPerDay, dates *[]time.Time) (*[]int, error) {
	if len(*pipelineRunsPerDays) == 0 {
		return nil, fmt.Errorf("no pipeline runs aggregates provided")
	}
	if len(*dates) == 0 {
		return nil, fmt.Errorf("no dates provided")
	}
	if len(*dates) <= len(*pipelineRunsPerDays) {
		return nil, fmt.Errorf("more pipeline runs per day than dates provided")
	}

	dailyPipelineRuns := make([]int, len(*dates))
	curr := 0
	for index, date := range *dates {
		if curr < len(*pipelineRunsPerDays) && date == (*pipelineRunsPerDays)[curr].Date {
			dailyPipelineRuns[index] = (*pipelineRunsPerDays)[curr].TotalPipelineRuns
			curr++
		} else {
			dailyPipelineRuns[index] = 0
		}
	}

	return &dailyPipelineRuns, nil
}
