package daos

import (
	"context"
	"os"

	"github.com/unnmdnwb3/dora/internal/database/mongodb"
	"github.com/unnmdnwb3/dora/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

// default pipelineRunsAggregateCollection
const pipelineRunsAggregateCollection = "pipeline_runs_aggregates"

// CreatePipelineRunsAggregate creates a new PipelineRunsAggregate.
func CreatePipelineRunsAggregate(ctx context.Context, pipelineRunsAggregate *models.PipelineRunsAggregate) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.InsertOne(ctx, pipelineRunsAggregateCollection, pipelineRunsAggregate)
	return err
}

// CreatePipelineRunsAggregates creates many new PipelineRunsAggregates.
func CreatePipelineRunsAggregates(ctx context.Context, pipelineRunsAggregates *[]models.PipelineRunsAggregate) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	for index, pipelineRunsAggregate := range *pipelineRunsAggregates {
		err = service.InsertOne(ctx, pipelineRunsAggregateCollection, &pipelineRunsAggregate)
		if err != nil {
			return err
		}
		(*pipelineRunsAggregates)[index] = pipelineRunsAggregate
	}
	return nil
}

// GetPipelineRunsAggregate retrieves an PipelineRunsAggregate.
func GetPipelineRunsAggregate(ctx context.Context, ID string, pipelineRunsAggregate *models.PipelineRunsAggregate) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.FindOneByID(ctx, pipelineRunsAggregateCollection, ID, pipelineRunsAggregate)
	return err
}

// ListPipelineRunsAggregates retrieves many PipelineRunsAggregates.
func ListPipelineRunsAggregates(ctx context.Context, pipelineID string, pipelineRunsAggregates *[]models.PipelineRunsAggregate) error {
	filter := bson.M{"pipeline_id": pipelineID}
	err := ListPipelineRunsAggregatesByFilter(ctx, filter, pipelineRunsAggregates)
	return err
}

// ListPipelineRunsAggregatesByFilter retrieves many PipelineRunsAggregates conforming to a filter.
func ListPipelineRunsAggregatesByFilter(ctx context.Context, filter bson.M, pipelineRunsAggregates *[]models.PipelineRunsAggregate) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.Find(ctx, pipelineRunsAggregateCollection, filter, pipelineRunsAggregates)
	return err
}

// UpdatePipelineRunsAggregate updates an PipelineRunsAggregate.
func UpdatePipelineRunsAggregate(ctx context.Context, ID string, pipelineRunsAggregate *models.PipelineRunsAggregate) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.UpdateOne(ctx, pipelineRunsAggregateCollection, ID, &pipelineRunsAggregate)
	if err != nil {
		return err
	}

	pipelineRunsAggregate.ID = ID
	return nil
}

// DeletePipelineRunsAggregate deletes an PipelineRunsAggregate.
func DeletePipelineRunsAggregate(ctx context.Context, ID string) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.DeleteOne(ctx, pipelineRunsAggregateCollection, ID)
	return err
}
