package ingest

import (
	"context"
	"time"

	"github.com/unnmdnwb3/dora/internal/connectors/prometheus"
	"github.com/unnmdnwb3/dora/internal/daos"
	"github.com/unnmdnwb3/dora/internal/models"
	"github.com/unnmdnwb3/dora/internal/utils/times"
)

// ImportIncidents gets and persists historical data for each incident of a deployment.
// However, this functions does not persist the raw monitoring data points, but rather aggregates them already to incidents.
// This is because the raw data points are not relevant for the user, but the incidents are.
// Additionally, persisting the raw data points would be too expensive, especially considering a small step size.
func ImportIncidents(ctx context.Context, channel chan error, deployment *models.Deployment) {
	monitoringDataPoints, err := ImportMonitoringDataPoints(ctx, deployment)
	if err != nil {
		channel <- err
		return
	}

	err = CreateIncidents(ctx, deployment, monitoringDataPoints)
	channel <- err
	return
}

// ImportMonitoringDataPoints gets historical monitoring data.
func ImportMonitoringDataPoints(ctx context.Context, deployment *models.Deployment) (*[]models.MonitoringDataPoint, error) {
	var integration models.Integration
	err := daos.GetIntegration(ctx, deployment.IntegrationID, &integration)
	if err != nil {
		return nil, err
	}

	// TODO standardize the time range for data imports
	end := time.Now()
	start := times.Date(end.AddDate(0, 0, -90))

	client := prometheus.NewClient(integration.URI, integration.BearerToken, deployment.Query, start, end, deployment.Step)
	monitoringDataPoints, err := client.GetMonitoringDataPoints()
	if err != nil {
		return nil, err
	}

	return monitoringDataPoints, nil
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
