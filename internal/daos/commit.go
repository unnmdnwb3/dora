package daos

import (
	"context"
	"os"

	"github.com/unnmdnwb3/dora/internal/database/mongodb"
	"github.com/unnmdnwb3/dora/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// default commitCollection
const commitCollection = "commits"

// CreateCommit creates a new Commit.
func CreateCommit(ctx context.Context, commit *models.Commit) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.InsertOne(ctx, commitCollection, commit)
	return err
}

// CreateCommits creates many new Commits.
func CreateCommits(ctx context.Context, commits *[]models.Commit) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	for index, commit := range *commits {
		err = service.InsertOne(ctx, commitCollection, &commit)
		if err != nil {
			return err
		}
		(*commits)[index] = commit
	}
	return nil
}

// GetCommit retrieves an Commit.
func GetCommit(ctx context.Context, commitID primitive.ObjectID, commit *models.Commit) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.FindOneByID(ctx, commitCollection, commitID, commit)
	return err
}

// ListCommits retrieves many Commits.
func ListCommits(ctx context.Context, pipelineID primitive.ObjectID, commits *[]models.Commit) error {
	filter := bson.M{"pipeline_id": pipelineID}
	err := ListCommitsByFilter(ctx, filter, commits)
	return err
}

// ListCommitsByFilter retrieves many Commits conforming to a filter.
func ListCommitsByFilter(ctx context.Context, filter bson.M, commits *[]models.Commit) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.Find(ctx, commitCollection, filter, commits)
	return err
}

// UpdateCommit updates an Commit.
func UpdateCommit(ctx context.Context, commitID primitive.ObjectID, commit *models.Commit) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.UpdateOne(ctx, commitCollection, commitID, &commit)
	if err != nil {
		return err
	}

	commit.ID = commitID
	return nil
}

// DeleteCommit deletes an Commit.
func DeleteCommit(ctx context.Context, commitID primitive.ObjectID) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.DeleteOne(ctx, commitCollection, commitID)
	return err
}
