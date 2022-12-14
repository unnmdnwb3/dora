package trigger

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

// GetFirstCommits returns the first commit of a change.
func GetFirstCommits(ctx context.Context, repositoryID primitive.ObjectID, pipelineRuns *[]models.PipelineRun) (*[]models.Commit, error) {
	// get the merge commits
	mergeCommitShas := []string{}
	for _, pipelineRun := range *pipelineRuns {
		mergeCommitShas = append(mergeCommitShas, pipelineRun.Sha)
	}

	var mergeCommits []models.Commit
	filter := bson.M{
		"repository_id": repositoryID,
		"sha":           bson.M{"$in": mergeCommitShas},
	}

	err := daos.ListCommitsByFilter(ctx, filter, &mergeCommits)
	if err != nil {
		return nil, err
	}

	// get the second-parent commits
	secondParentShas := []string{}
	for _, mergeCommit := range mergeCommits {
		secondParentShas = append(secondParentShas, mergeCommit.ParentShas[1])
	}

	// get all commits between the first second-parent commit and the last merge commit
	firstSecondParentSha := secondParentShas[0]
	lastMergeSha := mergeCommitShas[len(mergeCommitShas)-1]

	var boundaryCommits []models.Commit
	filter = bson.M{
		"repository_id": repositoryID,
		"sha":           bson.M{"$in": []string{firstSecondParentSha, lastMergeSha}},
	}
	err = daos.ListCommitsByFilter(ctx, filter, &boundaryCommits)
	if err != nil {
		return nil, err
	}

	firstCommit := boundaryCommits[0]
	lastCommit := boundaryCommits[1]

	var commits []models.Commit
	filter = bson.M{
		"repository_id": repositoryID,
		"created_at":    bson.M{"$gte": firstCommit.CreatedAt, "$lte": lastCommit.CreatedAt},
	}
	err = daos.ListCommitsByFilter(ctx, filter, &commits)
	if err != nil {
		return nil, err
	}

	// get the first commit of each pipeline run
	firstCommits := []models.Commit{}
	for index := 0; index < len(secondParentShas); index++ {
		pointer := 0
		for pointer < len(commits) {
			if commits[pointer].Sha == secondParentShas[index] {
				firstCommits = append(firstCommits, commits[pointer+1])
				pointer = pointer + 2
				break
			}
			pointer++
		}
	}

	return &firstCommits, nil
}

// CalculateChanges calculates the changes from commits and pipeline runs.
func CalculateChanges(ctx context.Context, commits *[]models.Commit, pipelineRuns *[]models.PipelineRun) (*[]models.Change, error) {
	if len(*commits) == 0 || len(*pipelineRuns) == 0 {
		return nil, fmt.Errorf("no commits or pipeline runs found")
	}
	if len(*commits) != len(*pipelineRuns) {
		return nil, fmt.Errorf("number of commits and pipeline runs do not match")
	}

	repositoryID := (*commits)[0].RepositoryID
	commits, err := GetFirstCommits(ctx, repositoryID, pipelineRuns)
	if err != nil {
		return nil, err
	}

	var changes []models.Change
	for index := 0; index < len(*pipelineRuns); index++ {
		start := (*commits)[index].CreatedAt
		end := (*pipelineRuns)[index].UpdatedAt

		change := models.Change{
			RepositoryID:    (*commits)[index].RepositoryID,
			PipelineID:      (*pipelineRuns)[index].PipelineID,
			FirstCommitDate: start,
			DeploymentDate:  end,
			LeadTime:        end.Sub(start).Seconds(),
		}
		changes = append(changes, change)
	}

	return &changes, nil
}

// CreateChanges creates changes from commits and pipeline runs.
func CreateChanges(ctx context.Context, repositoryID primitive.ObjectID, pipelineID primitive.ObjectID) error {
	var pipelineRuns []models.PipelineRun
	err := daos.ListPipelineRuns(ctx, pipelineID, &pipelineRuns)
	if err != nil {
		return err
	}

	commits, err := GetFirstCommits(ctx, repositoryID, &pipelineRuns)
	if err != nil {
		return err
	}

	changes, err := CalculateChanges(ctx, commits, &pipelineRuns)
	if err != nil {
		return err
	}

	err = daos.CreateChanges(ctx, repositoryID, changes)
	if err != nil {
		return err
	}

	return nil
}
