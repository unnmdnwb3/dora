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

type Integration struct {
	ctx    *context.Context
	client *mongo.Client
	coll   *mongo.Collection
}

func NewIntegration(ctx *context.Context) (Integration, error) {
	client, err := mongodb.NewClient(ctx)
	if err != nil {
		return Integration{}, err
	}
	coll := client.Database(mongodb.Database).Collection("integrations")

	return Integration{
		ctx:    ctx,
		client: client,
		coll:   coll,
	}, nil
}

func (i *Integration) Create(integration *models.Integration) (*models.Integration, error) {
	insertResult, err := i.coll.InsertOne(*i.ctx, integration)
	if err != nil {
		return nil, err
	}

	filter := bson.M{ "_id": insertResult.InsertedID }
	integrationResult := i.coll.FindOne(*i.ctx, filter)
	if integrationResult.Err() != nil {
		return nil, err
	}

	var result models.Integration
	err = integrationResult.Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (i *Integration) Read(id string) (*models.Integration, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{ "_id": objectId }
	readResult := i.coll.FindOne(*i.ctx, filter)
	if readResult.Err() != nil {
		return nil, readResult.Err()
	}

	var result models.Integration
	err = readResult.Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (i *Integration) ReadAll() (*[]models.Integration, error) {
	cursor, err := i.coll.Find(*i.ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	result := make([]models.Integration, cursor.RemainingBatchLength())
	pos := 0
	for cursor.Next(*i.ctx) {
		var cursorResult models.Integration
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

func (i *Integration) Update(integration *models.Integration) (*models.Integration, error) {
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

	var result models.Integration
	err = findResult.Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (i *Integration) Delete(id string) (*models.Integration, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{ "_id": objectId }
	deleteResult := i.coll.FindOneAndDelete(*i.ctx, filter)
	if deleteResult.Err() != nil {
		return nil, deleteResult.Err()
	}

	var result models.Integration
	err = deleteResult.Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}