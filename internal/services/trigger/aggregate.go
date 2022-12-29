package trigger

import (
	"context"
	"time"

	"github.com/unnmdnwb3/dora/internal/daos"
	"github.com/unnmdnwb3/dora/internal/models"
	"github.com/unnmdnwb3/dora/internal/utils/times"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreatePipelineRunsPerDays calculates and creates the pipeline runs for each day.
func CreatePipelineRunsPerDays(ctx context.Context, pipelineID primitive.ObjectID) error {
	var pipelineRuns []models.PipelineRun
	err := daos.ListPipelineRuns(ctx, pipelineID, &pipelineRuns)
	if err != nil {
		return err
	}

	pipelineRunsPerDays, err := CalculatePipelineRunsPerDays(ctx, &pipelineRuns)
	if err != nil {
		return err
	}

	err = daos.CreatePipelineRunsPerDays(ctx, pipelineID, pipelineRunsPerDays)
	return err
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

// CreateIncidentsPerDays calculates and creates the incidents for each day.
func CreateIncidentsPerDays(ctx context.Context, deploymentID primitive.ObjectID) error {
	var incidents []models.Incident
	err := daos.ListIncidents(ctx, deploymentID, &incidents)
	if err != nil {
		return err
	}

	incidentsPerDays, err := CalculateIncidentsPerDays(ctx, &incidents)
	if err != nil {
		return err
	}

	err = daos.CreateIncidentsPerDays(ctx, deploymentID, incidentsPerDays)
	return err
}

// CalculateIncidentsPerDays calculates the incidents per day.
func CalculateIncidentsPerDays(ctx context.Context, incidents *[]models.Incident) (*[]models.IncidentsPerDay, error) {
	incidentsPerDays := []models.IncidentsPerDay{}

	date := (*incidents)[0].StartDate
	var countPerDay int
	var durationPerDay time.Duration

	// there is always only a single incident active
	// but, altough counted as a single incident for a single day, incidents can span across multiple days
	// so we need to cut the last incident to the end of the day if it would span across multiple days
	for index := 0; index < len(*incidents); index++ {
		if !times.SameDay(date, (*incidents)[index].StartDate) {
			dayDate := times.Date(date)
			incidentsPerDay := models.IncidentsPerDay{
				DeploymentID:          (*incidents)[index].DeploymentID,
				Date:                  dayDate,
				TotalIncidents:        countPerDay,
				TotalIncidentDuration: durationPerDay,
			}

			incidentsPerDays = append(incidentsPerDays, incidentsPerDay)

			date = (*incidents)[index].StartDate
			countPerDay = 0
			durationPerDay = 0
		}

		incidentDuration := (*incidents)[index].EndDate.Sub((*incidents)[index].StartDate)

		countPerDay++
		durationPerDay += incidentDuration
	}

	incidentsPerDays = append(incidentsPerDays, models.IncidentsPerDay{
		Date:                  times.Date(date),
		TotalIncidents:        countPerDay,
		TotalIncidentDuration: durationPerDay,
	})

	return &incidentsPerDays, nil
}
