package daos

import (
	"context"
	"os"

	"github.com/unnmdnwb3/dora/internal/database/mongodb"
	"github.com/unnmdnwb3/dora/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// default changesPerDayCollection
const changesPerDayCollection = "changes_per_days"

// CreateChangesPerDay creates a new ChangesPerDay.
func CreateChangesPerDay(ctx context.Context, repositoryID primitive.ObjectID, changesPerDay *models.ChangesPerDay) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	changesPerDay.RepositoryID = repositoryID

	err = service.InsertOne(ctx, changesPerDayCollection, changesPerDay)
	return err
}

// CreateChangesPerDays creates many new ChangesPerDay.
func CreateChangesPerDays(ctx context.Context, repositoryID primitive.ObjectID, changesPerDays *[]models.ChangesPerDay) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	for index, changesPerDay := range *changesPerDays {
		changesPerDay.RepositoryID = repositoryID

		err = service.InsertOne(ctx, changesPerDayCollection, &changesPerDay)
		if err != nil {
			return err
		}
		(*changesPerDays)[index] = changesPerDay
	}
	return nil
}

// GetChangesPerDay retrieves a ChangesPerDay.
func GetChangesPerDay(ctx context.Context, changesPerDayID primitive.ObjectID, changesPerDay *models.ChangesPerDay) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.FindOneByID(ctx, changesPerDayCollection, changesPerDayID, changesPerDay)
	return err
}

// ListChangesPerDays retrieves many ChangesPerDay.
func ListChangesPerDays(ctx context.Context, repositoryID primitive.ObjectID, changesPerDay *[]models.ChangesPerDay) error {
	filter := bson.M{"repository_id": repositoryID}
	err := ListChangesPerDaysByFilter(ctx, filter, changesPerDay)
	return err
}

// ListChangesPerDaysByFilter retrieves many ChangesPerDay conforming to a filter.
func ListChangesPerDaysByFilter(ctx context.Context, filter bson.M, changesPerDay *[]models.ChangesPerDay) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.Find(ctx, changesPerDayCollection, filter, changesPerDay)
	return err
}

// UpdateChangesPerDay updates a ChangesPerDay.
func UpdateChangesPerDay(ctx context.Context, changesPerDayID primitive.ObjectID, changesPerDay *models.ChangesPerDay) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.UpdateOne(ctx, changesPerDayCollection, changesPerDayID, &changesPerDay)
	if err != nil {
		return err
	}

	changesPerDay.ID = changesPerDayID
	return nil
}

// DeleteChangesPerDay deletes a ChangesPerDay.
func DeleteChangesPerDay(ctx context.Context, changesPerDayID primitive.ObjectID) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.DeleteOne(ctx, changesPerDayCollection, changesPerDayID)
	return err
}
