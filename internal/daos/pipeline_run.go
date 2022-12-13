package daos

import (
	"context"
	"os"

	"github.com/unnmdnwb3/dora/internal/database/mongodb"
	"github.com/unnmdnwb3/dora/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

// default pipelineRunCollection
const pipelineRunCollection = "pipeline_runs"

// CreatePipelineRun creates a new PipelineRun.
func CreatePipelineRun(ctx context.Context, pipelineRun *models.PipelineRun) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.InsertOne(ctx, pipelineRunCollection, pipelineRun)
	return err
}

// CreatePipelineRuns creates many new PipelineRuns.
func CreatePipelineRuns(ctx context.Context, pipelineRuns *[]models.PipelineRun) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	for index, pipelineRun := range *pipelineRuns {
		err = service.InsertOne(ctx, pipelineRunCollection, &pipelineRun)
		if err != nil {
			return err
		}
		(*pipelineRuns)[index] = pipelineRun
	}
	return nil
}

// GetPipelineRun retrieves an PipelineRun.
func GetPipelineRun(ctx context.Context, ID string, pipelineRun *models.PipelineRun) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.FindOneByID(ctx, pipelineRunCollection, ID, pipelineRun)
	return err
}

// ListPipelineRuns retrieves many PipelineRuns.
func ListPipelineRuns(ctx context.Context, pipelineID string, pipelineRuns *[]models.PipelineRun) error {
	filter := bson.M{"pipeline_id": pipelineID}
	err := ListPipelineRunsByFilter(ctx, filter, pipelineRuns)
	return err
}

// ListPipelineRunsByFilter retrieves many PipelineRuns conforming to a filter.
// TODO change to pass a struct instead of bson.M
func ListPipelineRunsByFilter(ctx context.Context, filter bson.M, pipelineRuns *[]models.PipelineRun) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.Find(ctx, pipelineRunCollection, filter, pipelineRuns)
	return err
}

// UpdatePipelineRun updates an PipelineRun.
func UpdatePipelineRun(ctx context.Context, ID string, pipelineRun *models.PipelineRun) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.UpdateOne(ctx, pipelineRunCollection, ID, &pipelineRun)
	if err != nil {
		return err
	}

	pipelineRun.ID = ID
	return nil
}

// DeletePipelineRun deletes an PipelineRun.
func DeletePipelineRun(ctx context.Context, ID string) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.DeleteOne(ctx, pipelineRunCollection, ID)
	return err
}
