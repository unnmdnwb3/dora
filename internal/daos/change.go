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

// default changeCollection
const changeCollection = "changes"

// CreateChange creates a new Change.
func CreateChange(ctx context.Context, repositoryID primitive.ObjectID, change *models.Change) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	change.RepositoryID = repositoryID
	err = service.InsertOne(ctx, changeCollection, change)
	return err
}

// CreateChanges creates many new Changes.
func CreateChanges(ctx context.Context, repositoryID primitive.ObjectID, changes *[]models.Change) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	for index, change := range *changes {
		change.RepositoryID = repositoryID

		err = service.InsertOne(ctx, changeCollection, &change)
		if err != nil {
			return err
		}
		(*changes)[index] = change
	}
	return nil
}

// GetChange retrieves an Change.
func GetChange(ctx context.Context, changeID primitive.ObjectID, change *models.Change) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.FindOneByID(ctx, changeCollection, changeID, change)
	return err
}

// ListChanges retrieves many Changes.
func ListChanges(ctx context.Context, repositoryID primitive.ObjectID, changes *[]models.Change) error {
	filter := bson.M{"repository_id": repositoryID}
	err := ListChangesByFilter(ctx, filter, changes)
	return err
}

// ListChangesByFilter retrieves many Changes conforming to a filter.
func ListChangesByFilter(ctx context.Context, filter bson.M, changes *[]models.Change) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	ops := options.Find().SetSort(bson.M{"first_commit_date": 1})
	err = service.Find(ctx, changeCollection, filter, changes, ops)
	return err
}

// UpdateChange updates an Change.
func UpdateChange(ctx context.Context, changeID primitive.ObjectID, change *models.Change) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.UpdateOne(ctx, changeCollection, changeID, &change)
	if err != nil {
		return err
	}

	change.ID = changeID
	return nil
}

// DeleteChange deletes an Change.
func DeleteChange(ctx context.Context, changeID primitive.ObjectID) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.DeleteOne(ctx, changeCollection, changeID)
	return err
}
