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

// WorkflowRun is the data access object for any deploy run
type WorkflowRun struct {
	ctx  *context.Context
	coll *mongo.Collection
}

// NewWorkflowRun creates a new data access object for any deploy run
func NewWorkflowRun(ctx *context.Context) (*WorkflowRun, error) {
	return &WorkflowRun{
		ctx:  ctx,
		coll: mongodb.DB.Collection("deploy_runs"),
	}, nil
}

// Create persists a new deploy run
func (d WorkflowRun) Create(workflowRun *models.WorkflowRun) (string, error) {
	insertResult, err := d.coll.InsertOne(*d.ctx, workflowRun)
	if err != nil {
		return "", err
	}

	oid := insertResult.InsertedID.(primitive.ObjectID)
	return oid.Hex(), nil
}

// CreateMany persists many new deploy runs
func (d WorkflowRun) CreateMany(workflowRuns []models.WorkflowRun) (*[]string, error) {
	// convert structs to interface
	documents := make([]interface{}, len(workflowRuns))
	for i, workflowRun := range workflowRuns {
		documents[i] = workflowRun
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
func (d WorkflowRun) Read(id string) (*models.WorkflowRun, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	readResult := d.coll.FindOne(*d.ctx, filter)
	if readResult.Err() != nil {
		return nil, readResult.Err()
	}

	var result models.WorkflowRun
	err = readResult.Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ReadAll retrieves all deploy runs
func (d WorkflowRun) ReadAll(filter bson.M) (*[]models.WorkflowRun, error) {
	cursor, err := d.coll.Find(*d.ctx, filter)
	if err != nil {
		return nil, err
	}

	result := make([]models.WorkflowRun, cursor.RemainingBatchLength())
	pos := 0
	for cursor.Next(*d.ctx) {
		var cursorResult models.WorkflowRun
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
func (d WorkflowRun) Update(workflowRun *models.WorkflowRun) (*models.WorkflowRun, error) {
	objectID, err := primitive.ObjectIDFromHex(workflowRun.ID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id": objectID,
	}
	update := bson.M{
		"$set": bson.M{
			"project_id": workflowRun.Ref,
			"ref":        workflowRun.Ref,
			"status":     workflowRun.Status,
			"source":     workflowRun.Source,
			"created_at": workflowRun.CreatedAt,
			"updated_at": workflowRun.UpdatedAt,
			"uri":        workflowRun.URI,
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

	var result models.WorkflowRun
	err = findResult.Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Delete deletes an existing deploy run
func (d WorkflowRun) Delete(id string) (*models.WorkflowRun, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	deleteResult := d.coll.FindOneAndDelete(*d.ctx, filter)
	if deleteResult.Err() != nil {
		return nil, deleteResult.Err()
	}

	var result models.WorkflowRun
	err = deleteResult.Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
