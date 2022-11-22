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

// DeployRun is the data access object for any deploy run
type DeployRun struct {
	ctx  *context.Context
	coll *mongo.Collection
}

// NewDeployRun creates a new data access object for any deploy run
func NewDeployRun(ctx *context.Context) (*DeployRun, error) {
	return &DeployRun{
		ctx:  ctx,
		coll: mongodb.DB.Collection("deploy_runs"),
	}, nil
}

// Create persists a new deploy run
func (d DeployRun) Create(deployRun *models.DeployRun) (string, error) {
	insertResult, err := d.coll.InsertOne(*d.ctx, deployRun)
	if err != nil {
		return "", err
	}

	oid := insertResult.InsertedID.(primitive.ObjectID)
	return oid.Hex(), nil
}

// CreateMany persists many new deploy runs
func (d DeployRun) CreateMany(deployRuns []models.DeployRun) (*[]string, error) {
	// convert structs to interface
	documents := make([]interface{}, len(deployRuns))
	for i, deployRun := range deployRuns {
		documents[i] = deployRun
	}

	insertManyResult, err := d.coll.InsertMany(*d.ctx, documents)
	if err != nil {
		return nil, err
	}

	result := []string{}
	for _, insertedID := range insertManyResult.InsertedIDs {
		oid := insertedID.(primitive.ObjectID)
		result = append(result, oid.String())
	}

	return &result, nil
}

// Read retrieves an deploy run
func (d DeployRun) Read(id string) (*models.DeployRun, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	readResult := d.coll.FindOne(*d.ctx, filter)
	if readResult.Err() != nil {
		return nil, readResult.Err()
	}

	var result models.DeployRun
	err = readResult.Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ReadAll retrieves all deploy runs
func (d DeployRun) ReadAll() (*[]models.DeployRun, error) {
	cursor, err := d.coll.Find(*d.ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	result := make([]models.DeployRun, cursor.RemainingBatchLength())
	pos := 0
	for cursor.Next(*d.ctx) {
		var cursorResult models.DeployRun
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

// Update persists changes to an alreay existing deploy run
func (d DeployRun) Update(deployRun *models.DeployRun) (*models.DeployRun, error) {
	objectID, err := primitive.ObjectIDFromHex(deployRun.ID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id": objectID,
	}
	update := bson.M{
		"$set": bson.M{
			"project_id": deployRun.Ref,
			"ref":        deployRun.Ref,
			"status":     deployRun.Status,
			"source":     deployRun.Source,
			"created_at": deployRun.CreatedAt,
			"updated_at": deployRun.UpdatedAt,
			"uri":        deployRun.URI,
		},
	}

	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	findResult := d.coll.FindOneAndUpdate(*d.ctx, filter, update, &opt)
	if findResult.Err() != nil {
		return nil, findResult.Err()
	}

	var result models.DeployRun
	err = findResult.Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Delete deletes an existing deploy run
func (d DeployRun) Delete(id string) (*models.DeployRun, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	deleteResult := d.coll.FindOneAndDelete(*d.ctx, filter)
	if deleteResult.Err() != nil {
		return nil, deleteResult.Err()
	}

	var result models.DeployRun
	err = deleteResult.Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
