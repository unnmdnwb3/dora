package daos

import (
	"context"
	"os"

	"github.com/unnmdnwb3/dora/internal/database/mongodb"
	"github.com/unnmdnwb3/dora/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// default pipelineRunCollection
const pipelineRunCollection = "pipeline_runs"

// CreatePipelineRun creates a new PipelineRun.
func CreatePipelineRun(ctx context.Context, pipelineID primitive.ObjectID, pipelineRun *models.PipelineRun) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	pipelineRun.PipelineID = pipelineID

	err = service.InsertOne(ctx, pipelineRunCollection, pipelineRun)
	return err
}

// CreatePipelineRuns creates many new PipelineRuns.
func CreatePipelineRuns(ctx context.Context, pipelineID primitive.ObjectID, pipelineRuns *[]models.PipelineRun) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	for index, pipelineRun := range *pipelineRuns {
		pipelineRun.PipelineID = pipelineID

		err = service.InsertOne(ctx, pipelineRunCollection, &pipelineRun)
		if err != nil {
			return err
		}
		(*pipelineRuns)[index] = pipelineRun
	}
	return nil
}

// GetPipelineRun retrieves an PipelineRun.
func GetPipelineRun(ctx context.Context, pipelineRunID primitive.ObjectID, pipelineRun *models.PipelineRun) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.FindOneByID(ctx, pipelineRunCollection, pipelineRunID, pipelineRun)
	return err
}

// ListPipelineRuns retrieves many PipelineRuns.
func ListPipelineRuns(ctx context.Context, pipelineID primitive.ObjectID, pipelineRuns *[]models.PipelineRun) error {
	filter := bson.M{"pipeline_id": pipelineID}
	err := ListPipelineRunsByFilter(ctx, filter, pipelineRuns)
	return err
}

// ListPipelineRunsByFilter retrieves many PipelineRuns conforming to a filter.
func ListPipelineRunsByFilter(ctx context.Context, filter bson.M, pipelineRuns *[]models.PipelineRun) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	ops := options.Find().SetSort(bson.M{"created_at": 1})
	err = service.Find(ctx, pipelineRunCollection, filter, pipelineRuns, ops)
	return err
}

// UpdatePipelineRun updates an PipelineRun.
func UpdatePipelineRun(ctx context.Context, pipelineRunID primitive.ObjectID, pipelineRun *models.PipelineRun) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.UpdateOne(ctx, pipelineRunCollection, pipelineRunID, &pipelineRun)
	if err != nil {
		return err
	}

	pipelineRun.ID = pipelineRunID
	return nil
}

// DeletePipelineRun deletes an PipelineRun.
func DeletePipelineRun(ctx context.Context, pipelineRunID primitive.ObjectID) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.DeleteOne(ctx, pipelineRunCollection, pipelineRunID)
	return err
}
