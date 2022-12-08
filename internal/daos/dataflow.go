package daos

import (
	"context"
	"os"

	"github.com/unnmdnwb3/dora/internal/database/mongodb"
	"github.com/unnmdnwb3/dora/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

// default dataflowCollection
const dataflowCollection = "dataflows"

// CreateDataflow creates a new Dataflow.
func CreateDataflow(ctx context.Context, dataflow *models.Dataflow) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.InsertOne(ctx, dataflowCollection, dataflow)
	return err
}

// GetDataflow retrieves an Dataflow.
func GetDataflow(ctx context.Context, ID string, dataflow *models.Dataflow) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.FindOneByID(ctx, dataflowCollection, ID, dataflow)
	return err
}

// ListDataflows retrieves many Dataflows.
func ListDataflows(ctx context.Context, dataflows *[]models.Dataflow) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.Find(ctx, dataflowCollection, bson.M{}, dataflows)
	return err
}

// ListDataflowsByFilter retrieves many Dataflows conforming to a filter.
// TODO change to pass a struct instead of bson.M
func ListDataflowsByFilter(ctx context.Context, filter bson.M, dataflows *[]models.Dataflow) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.Find(ctx, dataflowCollection, filter, dataflows)
	return err
}

// UpdateDataflow updates an Dataflow.
func UpdateDataflow(ctx context.Context, ID string, dataflow *models.Dataflow) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.UpdateOne(ctx, dataflowCollection, ID, &dataflow)
	if err != nil {
		return err
	}

	dataflow.ID = ID
	return nil
}

// DeleteDataflow deletes an Dataflow.
func DeleteDataflow(ctx context.Context, ID string) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.DeleteOne(ctx, dataflowCollection, ID)
	return err
}
