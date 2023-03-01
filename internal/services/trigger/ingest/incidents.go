package ingest

import (
	"context"
	"log"
	"time"

	"github.com/unnmdnwb3/dora/internal/connectors/prometheus"
	"github.com/unnmdnwb3/dora/internal/daos"
	"github.com/unnmdnwb3/dora/internal/models"
)

// ImportIncidents gets and persists historical data for each incident of a deployment.
// However, this functions does not persist the raw alerts, but rather aggregates them already to incidents.
// This is because the raw alerts are not relevant for the user, but the incidents are.
func ImportIncidents(ctx context.Context, channel chan error, deployment *models.Deployment) {
	alerts, err := ImportAlerts(ctx, deployment)
	if err != nil {
		channel <- err
		return
	}

	err = CreateIncidents(ctx, deployment, alerts)
	channel <- err
	return
}

// ImportAlerts gets the historical raw alert data.
func ImportAlerts(ctx context.Context, deployment *models.Deployment) (*[]models.Alert, error) {
	var integration models.Integration
	err := daos.GetIntegration(ctx, deployment.IntegrationID, &integration)
	if err != nil {
		return nil, err
	}

	client := prometheus.NewClient(integration.URI, integration.BearerToken, deployment.Query)
	alerts, err := client.GetAlerts()
	if err != nil {
		return nil, err
	}

	log.Printf("Imported %d alerts", len(*alerts))

	return alerts, nil
}

// CreateIncidents calculates and creates the incidents for a given deployment.
func CreateIncidents(ctx context.Context, deployment *models.Deployment, alerts *[]models.Alert) error {
	incidents, err := CalculateIncidents(ctx, deployment, alerts)
	if err != nil {
		return err
	}

	if incidents == nil {
		return nil
	}

	err = daos.CreateIncidents(ctx, incidents)
	if err != nil {
		return err
	}

	log.Printf("Created %d incidents", len(*incidents))

	return nil
}

// CalculateIncidents calculates the incidents for a given deployment.
func CalculateIncidents(ctx context.Context, deployment *models.Deployment, alerts *[]models.Alert) (*[]models.Incident, error) {
	if len(*alerts) == 0 {
		return nil, nil
	}

	incidents := []models.Incident{}

	// we assume each incidents to have at least one alert
	// if the next alert is more than 120 seconds away, we assume a new incident

	start := (*alerts)[0]
	for i := 1; i < len(*alerts); i++ {
		prev := (*alerts)[i-1]
		curr := (*alerts)[i]

		diff := curr.CreatedAt.Sub(prev.CreatedAt).Seconds()
		if diff > 120 {
			incidents = append(incidents, models.Incident{
				DeploymentID: deployment.ID,
				StartDate:    start.CreatedAt,
				EndDate:      prev.CreatedAt.Add(30 * time.Second), // add 30 seconds
			})

			start = curr
		}
	}

	// add the last incident
	incidents = append(incidents, models.Incident{
		DeploymentID: deployment.ID,
		StartDate:    start.CreatedAt,
		EndDate:      (*alerts)[len(*alerts)-1].CreatedAt.Add(60 * time.Second), // add 30 seconds
	})

	return &incidents, nil
}
