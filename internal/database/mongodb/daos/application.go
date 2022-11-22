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
func (a *Application) Create(application *models.Application) (*models.Application, error) {
	insertResult, err := a.coll.InsertOne(*a.ctx, application)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": insertResult.InsertedID}
	applicationResult := a.coll.FindOne(*a.ctx, filter)
	if applicationResult.Err() != nil {
		return nil, err
	}

	var result models.Application
	err = applicationResult.Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Read retrieves an application
func (a *Application) Read(id string) (*models.Application, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	readResult := a.coll.FindOne(*a.ctx, filter)
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
func (a *Application) ReadAll() (*[]models.Application, error) {
	cursor, err := a.coll.Find(*a.ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	result := make([]models.Application, cursor.RemainingBatchLength())
	pos := 0
	for cursor.Next(*a.ctx) {
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
func (a *Application) Update(application *models.Application) (*models.Application, error) {
	objectID, err := primitive.ObjectIDFromHex(application.ID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id": objectID,
	}
	update := bson.M{
		"$set": bson.M{
			"auth": application.Auth,
			"type": application.Type,
			"uri":  application.URI,
		},
	}

	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	findResult := a.coll.FindOneAndUpdate(*a.ctx, filter, update, &opt)
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
func (a *Application) Delete(id string) (*models.Application, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	deleteResult := a.coll.FindOneAndDelete(*a.ctx, filter)
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
