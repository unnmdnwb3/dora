package aggregate

import (
	"context"
	"time"

	"github.com/unnmdnwb3/dora/internal/daos"
	"github.com/unnmdnwb3/dora/internal/models"
	"github.com/unnmdnwb3/dora/internal/utils/times"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateIncidentsPerDays calculates and creates the incidents for each day.
func CreateIncidentsPerDays(ctx context.Context, channel chan error, deploymentID primitive.ObjectID) {
	var incidents []models.Incident
	err := daos.ListIncidents(ctx, deploymentID, &incidents)
	if err != nil {
		channel <- err
		return
	}

	incidentsPerDays, err := CalculateIncidentsPerDays(ctx, &incidents)
	if err != nil {
		channel <- err
		return
	}

	err = daos.CreateIncidentsPerDays(ctx, deploymentID, incidentsPerDays)
	channel <- err
	return
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
