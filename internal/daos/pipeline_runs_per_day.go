package daos

import (
	"context"
	"os"

	"github.com/unnmdnwb3/dora/internal/database/mongodb"
	"github.com/unnmdnwb3/dora/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// default pipelineRunsPerDayCollection
const pipelineRunsPerDayCollection = "pipeline_runs_per_days"

// CreatePipelineRunsPerDay creates a new PipelineRunsPerDay.
func CreatePipelineRunsPerDay(ctx context.Context, pipelineID primitive.ObjectID, pipelineRunsPerDay *models.PipelineRunsPerDay) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	pipelineRunsPerDay.PipelineID = pipelineID

	err = service.InsertOne(ctx, pipelineRunsPerDayCollection, pipelineRunsPerDay)
	return err
}

// CreatePipelineRunsPerDays creates many new PipelineRunsPerDay.
func CreatePipelineRunsPerDays(ctx context.Context, pipelineID primitive.ObjectID, pipelineRunsPerDays *[]models.PipelineRunsPerDay) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	for index, pipelineRunsPerDay := range *pipelineRunsPerDays {
		pipelineRunsPerDay.PipelineID = pipelineID

		err = service.InsertOne(ctx, pipelineRunsPerDayCollection, &pipelineRunsPerDay)
		if err != nil {
			return err
		}
		(*pipelineRunsPerDays)[index] = pipelineRunsPerDay
	}
	return nil
}

// GetPipelineRunsPerDay retrieves a PipelineRunsPerDay.
func GetPipelineRunsPerDay(ctx context.Context, pipelineRunsPerDayID primitive.ObjectID, pipelineRunsPerDay *models.PipelineRunsPerDay) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.FindOneByID(ctx, pipelineRunsPerDayCollection, pipelineRunsPerDayID, pipelineRunsPerDay)
	return err
}

// ListPipelineRunsPerDays retrieves many PipelineRunsPerDay.
func ListPipelineRunsPerDays(ctx context.Context, pipelineID primitive.ObjectID, pipelineRunsPerDay *[]models.PipelineRunsPerDay) error {
	filter := bson.M{"pipeline_id": pipelineID}
	err := ListPipelineRunsPerDaysByFilter(ctx, filter, pipelineRunsPerDay)
	return err
}

// ListPipelineRunsPerDaysByFilter retrieves many PipelineRunsPerDay conforming to a filter.
func ListPipelineRunsPerDaysByFilter(ctx context.Context, filter bson.M, pipelineRunsPerDay *[]models.PipelineRunsPerDay) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.Find(ctx, pipelineRunsPerDayCollection, filter, pipelineRunsPerDay)
	return err
}

// UpdatePipelineRunsPerDay updates a PipelineRunsPerDay.
func UpdatePipelineRunsPerDay(ctx context.Context, pipelineRunsPerDayID primitive.ObjectID, pipelineRunsPerDay *models.PipelineRunsPerDay) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.UpdateOne(ctx, pipelineRunsPerDayCollection, pipelineRunsPerDayID, &pipelineRunsPerDay)
	if err != nil {
		return err
	}

	pipelineRunsPerDay.ID = pipelineRunsPerDayID
	return nil
}

// DeletePipelineRunsPerDay deletes a PipelineRunsPerDay.
func DeletePipelineRunsPerDay(ctx context.Context, pipelineRunsPerDayID primitive.ObjectID) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.DeleteOne(ctx, pipelineRunsPerDayCollection, pipelineRunsPerDayID)
	return err
}
