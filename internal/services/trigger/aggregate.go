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
// If no incident is found for a date, no aggregate will be created for that date!
func CalculateIncidentsPerDays(ctx context.Context, incidents *[]models.Incident) (*[]models.IncidentsPerDay, error) {
	incidentsPerDays := []models.IncidentsPerDay{}

	date := (*incidents)[0].StartDate
	var countPerDay int
	var durationPerDay time.Duration

	for index := 0; index < len(*incidents); index++ {
		newDate := (*incidents)[index].StartDate

		if !times.SameDay(date, newDate) {
			dayDate := times.Date(date)
			incidentsPerDay := models.IncidentsPerDay{
				DeploymentID:   (*incidents)[index].DeploymentID,
				Date:           dayDate,
				TotalIncidents: countPerDay,
				TotalDuration:  durationPerDay.Seconds(),
			}

			incidentsPerDays = append(incidentsPerDays, incidentsPerDay)

			date = (*incidents)[index].StartDate
			countPerDay = 0
			durationPerDay = 0
		}

		countPerDay++
		start := (*incidents)[index].StartDate
		end := (*incidents)[index].EndDate
		durationPerDay += end.Sub(start)
	}

	incidentsPerDays = append(incidentsPerDays, models.IncidentsPerDay{
		Date:           times.Date(date),
		TotalIncidents: countPerDay,
		TotalDuration:  durationPerDay.Seconds(),
	})

	return &incidentsPerDays, nil
}

// CreateIncidents calculates and creates the incidents for a given deployment.
func CreateIncidents(ctx context.Context, deployment *models.Deployment, monitoringDataPoints *[]models.MonitoringDataPoint) error {
	incidents, err := CalculateIncidents(ctx, deployment, monitoringDataPoints)
	if err != nil {
		return err
	}

	err = daos.CreateIncidents(ctx, incidents)
	return err
}

// CalculateIncidents calculates the incidents for a given deployment.
func CalculateIncidents(ctx context.Context, deployment *models.Deployment, monitoringDataPoints *[]models.MonitoringDataPoint) (*[]models.Incident, error) {
	var incidents []models.Incident

	index := FirstNonIncident(deployment.Relation, deployment.Threshold, monitoringDataPoints)
	if index == -1 {
		return &[]models.Incident{
			{
				DeploymentID: deployment.ID,
				StartDate:    (*monitoringDataPoints)[0].CreatedAt,
				EndDate:      (*monitoringDataPoints)[len(*monitoringDataPoints)-1].CreatedAt,
			},
		}, nil
	}

	// cut slice to the first non-incident point
	*monitoringDataPoints = (*monitoringDataPoints)[index:]
	isIncidentPrev := false
	slow := 0
	for fast := 1; fast < len(*monitoringDataPoints); fast++ {
		isIncident := IsIncident(deployment.Relation, deployment.Threshold, (*monitoringDataPoints)[fast])

		if !isIncident {
			if isIncidentPrev {
				incident := models.Incident{
					DeploymentID: deployment.ID,
					StartDate:    (*monitoringDataPoints)[slow].CreatedAt,
					EndDate:      (*monitoringDataPoints)[fast-1].CreatedAt,
				}
				incidents = append(incidents, incident)
			}
		}

		if isIncident {
			step, err := times.Duration(deployment.Step)
			if err != nil {
				return nil, err
			}
			isContinuation := IsContinuation((*monitoringDataPoints)[fast-1], (*monitoringDataPoints)[fast], step)

			if !isIncidentPrev {
				slow = fast
			}

			if isIncidentPrev && !isContinuation {
				incident := models.Incident{
					DeploymentID: deployment.ID,
					StartDate:    (*monitoringDataPoints)[slow].CreatedAt,
					EndDate:      (*monitoringDataPoints)[fast-1].CreatedAt,
				}
				incidents = append(incidents, incident)
				slow = fast
			}
		}

		isIncidentPrev = isIncident
	}

	// add last incident if still open
	if isIncidentPrev {
		incident := models.Incident{
			DeploymentID: deployment.ID,
			StartDate:    (*monitoringDataPoints)[slow].CreatedAt,
			EndDate:      (*monitoringDataPoints)[len(*monitoringDataPoints)-1].CreatedAt,
		}
		incidents = append(incidents, incident)
	}

	return &incidents, nil
}

// IsIncident checks if a given monitoring data point is part of an incident.
func IsIncident(relation string, threshold float64, monitoringDataPoint models.MonitoringDataPoint) bool {
	if relation == "gt" {
		return monitoringDataPoint.Value > threshold
	}

	return monitoringDataPoint.Value < threshold
}

// FirstNonIncident finds the first non-incident data point.
func FirstNonIncident(relation string, threshold float64, monitoringDataPoints *[]models.MonitoringDataPoint) int {
	for index := 0; index < len(*monitoringDataPoints); index++ {
		if !IsIncident(relation, threshold, (*monitoringDataPoints)[index]) {
			return index
		}
	}
	return -1
}

// IsContinuation checks if the given monitoring data points are part of the same incident.
func IsContinuation(prev models.MonitoringDataPoint, curr models.MonitoringDataPoint, step time.Duration) bool {
	elapsedTime := curr.CreatedAt.Sub(prev.CreatedAt)
	// timesteps in Prometheus time series are not exact, thus, we need to implement some tolerance
	tolerance := elapsedTime / 2
	return step-tolerance <= elapsedTime && elapsedTime <= step+tolerance
}
