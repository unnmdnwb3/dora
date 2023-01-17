package aggregate

import (
	"context"
	"time"

	"github.com/unnmdnwb3/dora/internal/daos"
	"github.com/unnmdnwb3/dora/internal/models"
	"github.com/unnmdnwb3/dora/internal/utils/times"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateChangesPerDays creates changes per days from commits and pipeline runs.
func CreateChangesPerDays(ctx context.Context, channel chan error, repositoryID primitive.ObjectID, pipelineID primitive.ObjectID) {
	changes := []models.Change{}
	err := daos.ListChanges(ctx, repositoryID, &changes)
	if err != nil {
		channel <- err
		return
	}

	changesPerDays, err := CalculateChangesPerDays(ctx, &changes)
	if err != nil {
		channel <- err
		return
	}

	err = daos.CreateChangesPerDays(ctx, repositoryID, pipelineID, changesPerDays)
	channel <- err
	return
}

// CalculateChangesPerDays calculates the changes per day.
// If no change is found for a date, no aggregate will be created for that date!
func CalculateChangesPerDays(ctx context.Context, changes *[]models.Change) (*[]models.ChangesPerDay, error) {
	changesPerDays := []models.ChangesPerDay{}

	date := (*changes)[0].FirstCommitDate
	var countPerDay int
	var durationPerDay time.Duration

	for index := 0; index < len(*changes); index++ {
		newDate := (*changes)[index].FirstCommitDate

		if !times.SameDay(date, newDate) {
			dayDate := times.Date(date)
			changesPerDay := models.ChangesPerDay{
				RepositoryID:  (*changes)[index].RepositoryID,
				Date:          dayDate,
				TotalChanges:  countPerDay,
				TotalLeadTime: durationPerDay.Seconds(),
			}

			changesPerDays = append(changesPerDays, changesPerDay)

			date = (*changes)[index].FirstCommitDate
			countPerDay = 0
			durationPerDay = 0
		}

		countPerDay++
		start := (*changes)[index].FirstCommitDate
		end := (*changes)[index].DeploymentDate
		durationPerDay += end.Sub(start)
	}

	changesPerDays = append(changesPerDays, models.ChangesPerDay{
		Date:          times.Date(date),
		TotalChanges:  countPerDay,
		TotalLeadTime: durationPerDay.Seconds(),
	})

	return &changesPerDays, nil
}
