package daos

import (
	"context"
	"os"

	"github.com/unnmdnwb3/dora/internal/database/mongodb"
	"github.com/unnmdnwb3/dora/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// default incidentsPerDayCollection
const incidentsPerDayCollection = "incidents_per_days"

// CreateIncidentsPerDay creates a new IncidentsPerDay.
func CreateIncidentsPerDay(ctx context.Context, deploymentID primitive.ObjectID, incidentsPerDay *models.IncidentsPerDay) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	incidentsPerDay.DeploymentID = deploymentID

	err = service.InsertOne(ctx, incidentsPerDayCollection, incidentsPerDay)
	return err
}

// CreateIncidentsPerDays creates many new IncidentsPerDay.
func CreateIncidentsPerDays(ctx context.Context, deploymentID primitive.ObjectID, incidentsPerDays *[]models.IncidentsPerDay) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	for index, incidentsPerDay := range *incidentsPerDays {
		incidentsPerDay.DeploymentID = deploymentID

		err = service.InsertOne(ctx, incidentsPerDayCollection, &incidentsPerDay)
		if err != nil {
			return err
		}
		(*incidentsPerDays)[index] = incidentsPerDay
	}
	return nil
}

// GetIncidentsPerDay retrieves a IncidentsPerDay.
func GetIncidentsPerDay(ctx context.Context, incidentsPerDayID primitive.ObjectID, incidentsPerDay *models.IncidentsPerDay) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.FindOneByID(ctx, incidentsPerDayCollection, incidentsPerDayID, incidentsPerDay)
	return err
}

// ListIncidentsPerDays retrieves many IncidentsPerDay.
func ListIncidentsPerDays(ctx context.Context, deploymentID primitive.ObjectID, incidentsPerDay *[]models.IncidentsPerDay) error {
	filter := bson.M{"deployment_id": deploymentID}
	err := ListIncidentsPerDaysByFilter(ctx, filter, incidentsPerDay)
	return err
}

// ListIncidentsPerDaysByFilter retrieves many IncidentsPerDay conforming to a filter.
func ListIncidentsPerDaysByFilter(ctx context.Context, filter bson.M, incidentsPerDay *[]models.IncidentsPerDay) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.Find(ctx, incidentsPerDayCollection, filter, incidentsPerDay)
	return err
}

// UpdateIncidentsPerDay updates a IncidentsPerDay.
func UpdateIncidentsPerDay(ctx context.Context, incidentsPerDayID primitive.ObjectID, incidentsPerDay *models.IncidentsPerDay) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.UpdateOne(ctx, incidentsPerDayCollection, incidentsPerDayID, &incidentsPerDay)
	if err != nil {
		return err
	}

	incidentsPerDay.ID = incidentsPerDayID
	return nil
}

// DeleteIncidentsPerDay deletes a IncidentsPerDay.
func DeleteIncidentsPerDay(ctx context.Context, incidentsPerDayID primitive.ObjectID) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.DeleteOne(ctx, incidentsPerDayCollection, incidentsPerDayID)
	return err
}
