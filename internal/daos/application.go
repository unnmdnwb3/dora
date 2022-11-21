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

type Application struct {
	ctx    *context.Context
	client *mongo.Client
	coll   *mongo.Collection
}

func NewApplication(ctx *context.Context) (Application, error) {
	client, err := mongodb.NewClient(ctx)
	if err != nil {
		return Application{}, err
	}
	coll := client.Database(mongodb.Database).Collection("integrations")

	return Application{
		ctx:    ctx,
		client: client,
		coll:   coll,
	}, nil
}

func (i *Application) Create(integration *models.Application) (*models.Application, error) {
	insertResult, err := i.coll.InsertOne(*i.ctx, integration)
	if err != nil {
		return nil, err
	}

	filter := bson.M{ "_id": insertResult.InsertedID }
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

func (i *Application) Read(id string) (*models.Application, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{ "_id": objectId }
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

	err = cursor.Err(); 
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (i *Application) Update(integration *models.Application) (*models.Application, error) {
	objectId, err := primitive.ObjectIDFromHex(integration.Id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id": objectId,
	}
	update := bson.M{
		"$set": bson.M{
			"auth": integration.Auth,
			"type": integration.Type,
			"uri":  integration.Uri,
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

func (i *Application) Delete(id string) (*models.Application, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{ "_id": objectId }
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