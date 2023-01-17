package aggregate

import (
	"context"

	"github.com/unnmdnwb3/dora/internal/daos"
	"github.com/unnmdnwb3/dora/internal/models"
	"github.com/unnmdnwb3/dora/internal/utils/times"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreatePipelineRunsPerDays calculates and creates the pipeline runs for each day.
func CreatePipelineRunsPerDays(ctx context.Context, channel chan error, pipelineID primitive.ObjectID) {
	var pipelineRuns []models.PipelineRun
	err := daos.ListPipelineRuns(ctx, pipelineID, &pipelineRuns)
	if err != nil {
		channel <- err
		return
	}

	pipelineRunsPerDays, err := CalculatePipelineRunsPerDays(ctx, &pipelineRuns)
	if err != nil {
		channel <- err
		return
	}

	err = daos.CreatePipelineRunsPerDays(ctx, pipelineID, pipelineRunsPerDays)
	channel <- err
	return
}

// CalculatePipelineRunsPerDays calculates the pipeline runs per day.
// If no pipeline run is found for a date, no aggregate will be created for that date!
func CalculatePipelineRunsPerDays(ctx context.Context, pipelineRuns *[]models.PipelineRun) (*[]models.PipelineRunsPerDay, error) {
	pipelineRunsPerDays := []models.PipelineRunsPerDay{}

	date := (*pipelineRuns)[0].UpdatedAt
	countPerDay := 0

	for index := 0; index < len(*pipelineRuns); index++ {
		if !times.SameDay(date, (*pipelineRuns)[index].UpdatedAt) {
			dayDate := times.Date(date)
			pipelineRunsPerDay := models.PipelineRunsPerDay{
				PipelineID:        (*pipelineRuns)[index].PipelineID,
				Date:              dayDate,
				TotalPipelineRuns: countPerDay,
			}

			pipelineRunsPerDays = append(pipelineRunsPerDays, pipelineRunsPerDay)

			date = (*pipelineRuns)[index].UpdatedAt
			countPerDay = 0
		}

		countPerDay++
	}

	pipelineRunsPerDays = append(pipelineRunsPerDays, models.PipelineRunsPerDay{
		Date:              times.Date(date),
		TotalPipelineRuns: countPerDay,
	})

	return &pipelineRunsPerDays, nil
}
