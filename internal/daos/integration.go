package daos

import (
	"context"
	"os"

	"github.com/unnmdnwb3/dora/internal/database/mongodb"
	"github.com/unnmdnwb3/dora/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// default integrationCollection
const integrationCollection = "integrations"

// CreateIntegration creates a new Integration.
func CreateIntegration(ctx context.Context, integration *models.Integration) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.InsertOne(ctx, integrationCollection, integration)
	return err
}

// GetIntegration retrieves an Integration.
func GetIntegration(ctx context.Context, objectID primitive.ObjectID, integration *models.Integration) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.FindOneByID(ctx, integrationCollection, objectID, integration)
	return err
}

// ListIntegrations retrieves many Integrations.
func ListIntegrations(ctx context.Context, integrations *[]models.Integration) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.Find(ctx, integrationCollection, bson.M{}, integrations)
	return err
}

// ListIntegrationsByFilter retrieves many Integrations conforming to a filter.
// TODO change to pass a struct instead of bson.M
func ListIntegrationsByFilter(ctx context.Context, filter bson.M, integrations *[]models.Integration) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.Find(ctx, integrationCollection, filter, integrations)
	return err
}

// UpdateIntegration updates an Integration.
func UpdateIntegration(ctx context.Context, objectID primitive.ObjectID, integration *models.Integration) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.UpdateOne(ctx, integrationCollection, objectID, &integration)
	if err != nil {
		return err
	}

	integration.ID = objectID
	return nil
}

// DeleteIntegration deletes an Integration.
func DeleteIntegration(ctx context.Context, objectID primitive.ObjectID) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.DeleteOne(ctx, integrationCollection, objectID)
	return err
}
