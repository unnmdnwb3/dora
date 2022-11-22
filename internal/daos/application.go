package daos

import (
	"context"

	"github.com/unnmdnwb3/dora/internal/database/mongodb"
	"github.com/unnmdnwb3/dora/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Application is the data access object for any application
type Application struct {
	ctx  *context.Context
	coll *mongo.Collection
}

// NewApplication creates a new data access object for any application
func NewApplication(ctx *context.Context) (Application, error) {
	return Application{
		ctx:  ctx,
		coll: mongodb.DB.Collection("applications"),
	}, nil
}

// Create persists a new application
func (i *Application) Create(integration *models.Application) (*models.Application, error) {
	insertResult, err := i.coll.InsertOne(*i.ctx, integration)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": insertResult.InsertedID}
	integrationResult := i.coll.FindOne(*i.ctx, filter)
	if integrationResult.Err() != nil {
		return nil, err
	}

	var result models.Application
	err = integrationResult.Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Read retrieves an application
func (i *Application) Read(id string) (*models.Application, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	readResult := i.coll.FindOne(*i.ctx, filter)
	if readResult.Err() != nil {
		return nil, readResult.Err()
	}

	var result models.Application
	err = readResult.Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ReadAll retrieves all applications
func (i *Application) ReadAll() (*[]models.Application, error) {
	cursor, err := i.coll.Find(*i.ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	result := make([]models.Application, cursor.RemainingBatchLength())
	pos := 0
	for cursor.Next(*i.ctx) {
		var cursorResult models.Application
		err := cursor.Decode(&cursorResult)
		if err != nil {
			return nil, err
		}
		result[pos] = cursorResult
		pos++
	}

	err = cursor.Err()
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Update persists changes to an alreay existing application
func (i *Application) Update(integration *models.Application) (*models.Application, error) {
	objectID, err := primitive.ObjectIDFromHex(integration.ID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id": objectID,
	}
	update := bson.M{
		"$set": bson.M{
			"auth": integration.Auth,
			"type": integration.Type,
			"uri":  integration.URI,
		},
	}

	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	findResult := i.coll.FindOneAndUpdate(*i.ctx, filter, update, &opt)
	if findResult.Err() != nil {
		return nil, findResult.Err()
	}

	var result models.Application
	err = findResult.Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Delete deletes an existing application
func (i *Application) Delete(id string) (*models.Application, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	deleteResult := i.coll.FindOneAndDelete(*i.ctx, filter)
	if deleteResult.Err() != nil {
		return nil, deleteResult.Err()
	}

	var result models.Application
	err = deleteResult.Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
