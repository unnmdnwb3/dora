package ingest

import (
	"context"
	"fmt"

	"github.com/unnmdnwb3/dora/internal/daos"
	"github.com/unnmdnwb3/dora/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ImportChanges gets and persists advanced historical data based on raw data for each defined source in a Dataflow.
func ImportChanges(ctx context.Context, channel chan error, repositoryID primitive.ObjectID, pipelineID primitive.ObjectID) {
	err := CreateChanges(ctx, repositoryID, pipelineID)
	channel <- err
	return
}

// CreateChanges creates changes from commits and pipeline runs.
func CreateChanges(ctx context.Context, repositoryID primitive.ObjectID, pipelineID primitive.ObjectID) error {
	var pipelineRuns []models.PipelineRun
	err := daos.ListPipelineRuns(ctx, pipelineID, &pipelineRuns)
	if err != nil {
		return err
	}

	firstCommits, err := GetFirstCommits(ctx, repositoryID, &pipelineRuns)
	if err != nil {
		return err
	}

	changes, err := CalculateChanges(ctx, firstCommits, &pipelineRuns)
	if err != nil {
		return err
	}

	err = daos.CreateChanges(ctx, repositoryID, changes)
	if err != nil {
		return err
	}

	return nil
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

	// get the first-parent commits
	firstParentShas := []string{}
	for _, mergeCommit := range mergeCommits {
		firstParentShas = append(firstParentShas, mergeCommit.ParentShas[0])
	}

	// get all commits between the first second-parent commit and the last merge commit
	firstSha := firstParentShas[0]
	lastSha := mergeCommitShas[len(mergeCommitShas)-1]

	var boundaryCommits []models.Commit
	filter = bson.M{
		"repository_id": repositoryID,
		"sha":           bson.M{"$in": []string{firstSha, lastSha}},
	}
	err = daos.ListCommitsByFilter(ctx, filter, &boundaryCommits)
	if err != nil {
		return nil, err
	}

	firstCommit, lastCommit := boundaryCommits[0], boundaryCommits[1]

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
	for index := 0; index < len(firstParentShas); index++ {
		pointer := 0
		for pointer < len(commits) {
			if commits[pointer].Sha == firstParentShas[index] {
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

	var changes []models.Change
	repositoryID := (*commits)[0].RepositoryID
	pipelineID := (*pipelineRuns)[0].PipelineID

	for index := 0; index < len(*pipelineRuns); index++ {
		start := (*commits)[index].CreatedAt
		end := (*pipelineRuns)[index].UpdatedAt
		leadTime := end.Sub(start).Seconds()

		change := models.Change{
			RepositoryID:    repositoryID,
			PipelineID:      pipelineID,
			FirstCommitDate: start,
			DeploymentDate:  end,
			LeadTime:        leadTime,
		}
		changes = append(changes, change)
	}

	return &changes, nil
}
