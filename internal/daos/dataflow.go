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

	dataflow.Repository.ID = primitive.NewObjectID()
	dataflow.Pipeline.ID = primitive.NewObjectID()
	dataflow.Deployment.ID = primitive.NewObjectID()

	err = service.InsertOne(ctx, dataflowCollection, dataflow)
	return err
}

// GetDataflow retrieves an Dataflow.
func GetDataflow(ctx context.Context, objectID primitive.ObjectID, dataflow *models.Dataflow) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.FindOneByID(ctx, dataflowCollection, objectID, dataflow)
	return err
}

// ListDataflows retrieves many Dataflows.
func ListDataflows(ctx context.Context, dataflows *[]models.Dataflow) error {
	err := ListDataflowsByFilter(ctx, bson.M{}, dataflows)
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

	ops := options.Find().SetSort(bson.M{"_id": 1})
	err = service.Find(ctx, dataflowCollection, filter, dataflows, ops)
	return err
}

// UpdateDataflow updates an Dataflow.
func UpdateDataflow(ctx context.Context, objectID primitive.ObjectID, dataflow *models.Dataflow) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.UpdateOne(ctx, dataflowCollection, objectID, &dataflow)
	if err != nil {
		return err
	}

	dataflow.ID = objectID
	return nil
}

// DeleteDataflow deletes an Dataflow.
func DeleteDataflow(ctx context.Context, objectID primitive.ObjectID) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.DeleteOne(ctx, dataflowCollection, objectID)
	return err
}
