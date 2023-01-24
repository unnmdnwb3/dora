package ingest

import (
	"context"
	"fmt"
	"log"
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
	start := times.Date(end.AddDate(0, -1, 0))

	client := prometheus.NewClient(integration.URI, integration.BearerToken, deployment.Query, start, end, deployment.Step)
	monitoringDataPoints, err := client.GetMonitoringDataPoints()
	if err != nil {
		return nil, err
	}

	log.Printf("Created %d monitoring data points", len(*monitoringDataPoints))

	return monitoringDataPoints, nil
}

// CreateIncidents calculates and creates the incidents for a given deployment.
func CreateIncidents(ctx context.Context, deployment *models.Deployment, monitoringDataPoints *[]models.MonitoringDataPoint) error {
	incidents, err := CalculateIncidents(ctx, deployment, monitoringDataPoints)
	if err != nil {
		return err
	}

	err = daos.CreateIncidents(ctx, incidents)
	if err != nil {
		return err
	}

	log.Printf("Created %d incidents", len(*incidents))

	return nil
}

// CalculateIncidents calculates the incidents for a given deployment.
// Ongoing incidents at the start and the end of the monitoring data points are used to create incidents with the information we got.
func CalculateIncidents(ctx context.Context, deployment *models.Deployment, monitoringDataPoints *[]models.MonitoringDataPoint) (*[]models.Incident, error) {
	if len(*monitoringDataPoints) == 0 {
		return &[]models.Incident{}, fmt.Errorf("no monitoring data points provided")
	}

	incidents := []models.Incident{}
	prevIsIncident := IsIncident(deployment.Relation, deployment.Threshold, (*monitoringDataPoints)[0])

	for slow, fast := 0, 1; fast < len(*monitoringDataPoints); fast++ {
		curr := (*monitoringDataPoints)[fast]
		currIsIncident := IsIncident(deployment.Relation, deployment.Threshold, curr)

		// if the current data point is an incident and the previous was one as well, we do not need to do anything
		// same applies if the current data point is not an incident and the previous was not one as well

		if currIsIncident && !prevIsIncident {
			slow = fast
		}

		if !currIsIncident && prevIsIncident {
			incidents = append(incidents, models.Incident{
				DeploymentID: deployment.ID,
				StartDate:    (*monitoringDataPoints)[slow].CreatedAt,
				EndDate:      (*monitoringDataPoints)[fast-1].CreatedAt,
			})
		}

		prevIsIncident = currIsIncident
	}

	// if the last data point part of an incident, we need to create an incident for it
	len := len(*monitoringDataPoints)
	prevIsIncident = IsIncident(deployment.Relation, deployment.Threshold, (*monitoringDataPoints)[len-2])
	currIsIncident := IsIncident(deployment.Relation, deployment.Threshold, (*monitoringDataPoints)[len-1])

	if currIsIncident && prevIsIncident {
		incidents = append(incidents, models.Incident{
			DeploymentID: deployment.ID,
			StartDate:    (*monitoringDataPoints)[len-2].CreatedAt,
			EndDate:      (*monitoringDataPoints)[len-1].CreatedAt,
		})
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
