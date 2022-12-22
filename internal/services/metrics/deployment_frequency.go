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

// // CalculateDeploymentFrequency calculates the deployment frequency for a given dataflow.
// func CalculateDeploymentFrequency(ctx context.Context, dataflowID string, window int) (*models.DeploymentFrequency, error) {
// 	var dataflow models.Dataflow
// 	err := daos.GetDataflow(ctx, dataflowID, &dataflow)
// 	if err != nil {
// 		return nil, err
// 	}

// 	pipelineID := dataflow.Pipeline.ID
// 	timeRange := (window - 1) * 2
// 	today := times.Date(time.Now())
// 	startDate := today.AddDate(0, 0, -timeRange)

// 	var pipelineRunsAggregate []models.PipelineRunsPerDay
// 	filter := bson.M{"pipeline_id": pipelineID, "date": bson.M{"$gte": startDate}}
// 	err = daos.ListPipelineRunsPerDaysByFilter(ctx, filter, &pipelineRunsAggregate)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(pipelineRunsAggregate) == 0 {
// 		return nil, fmt.Errorf("no pipeline runs found for pipeline %s", pipelineID)
// 	}

// 	dates, err := GetDates(&pipelineRunsAggregate)
// 	if err != nil {
// 		return nil, err
// 	}

// 	dailyPipelineRuns, err := GetDailyPipelineRuns(&pipelineRunsAggregate)
// 	if err != nil {
// 		return nil, err
// 	}

// 	movingAverages, err := CalculateMovingAverages(dailyPipelineRuns, window)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &models.DeploymentFrequency{
// 		DataflowID:        pipelineID,
// 		Dates:             *dates,
// 		DailyPipelineRuns: *dailyPipelineRuns,
// 		MovingAverages:    *movingAverages,
// 	}, nil
// }

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

// DatesBetween returns a slice of dates between the start and end dates.
func DatesBetween(startDate time.Time, endDate time.Time) (*[]time.Time, error) {
	startDate = times.Date(startDate)
	endDate = times.Date(endDate)

	if startDate.After(endDate) {
		return nil, fmt.Errorf("start date is after end date")
	}

	dates := []time.Time{}
	for date := startDate; date.Before(endDate); date = date.AddDate(0, 0, 1) {
		dates = append(dates, times.Date(date))
	}
	dates = append(dates, times.Date(endDate))

	return &dates, nil
}

// CompletePipelineRunsPerDays returns a slice of the number of pipeline runs per day,
// since provided PipelineRunsPerDays only account for the dates that any pipeline runs were found.
func CompletePipelineRunsPerDays(pipelineRunsPerDays *[]models.PipelineRunsPerDay, dates *[]time.Time) (*[]int, error) {
	if len(*pipelineRunsPerDays) == 0 {
		return nil, fmt.Errorf("no pipeline runs aggregates provided")
	}
	if len(*dates) == 0 {
		return nil, fmt.Errorf("no dates provided pipeline runs aggregates")
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

// CalculateMovingAverages calculates the moving averages for a given slice of deployments per day.
func CalculateMovingAverages(pipelineRunsPerDays *[]int, window int) (*[]float64, error) {
	if len(*pipelineRunsPerDays) == 0 {
		return nil, fmt.Errorf("no deployments per day provided to calculate moving averages")
	}
	if len(*pipelineRunsPerDays) != (window*2)-1 {
		return nil, fmt.Errorf("number of pipeline runs per day provided to calculate moving averages does not match window")
	}

	offset := window - 1
	totalDeploymentsInWindow := 0
	for index := 0; index < offset; index++ {
		totalDeploymentsInWindow += (*pipelineRunsPerDays)[index]
	}

	movingAverages := make([]float64, len(*pipelineRunsPerDays)-offset)
	for index := offset; index < len(*pipelineRunsPerDays); index++ {
		totalDeploymentsInWindow += (*pipelineRunsPerDays)[index]
		movingAverages[index-offset] = float64(totalDeploymentsInWindow) / float64(window)
		totalDeploymentsInWindow -= (*pipelineRunsPerDays)[index-offset]
	}

	return &movingAverages, nil
}
