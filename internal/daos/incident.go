package daos

import (
	"context"
	"os"

	"github.com/unnmdnwb3/dora/internal/database/mongodb"
	"github.com/unnmdnwb3/dora/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// default incidentCollection
const incidentCollection = "incidents"

// CreateIncident creates a new Incident.
func CreateIncident(ctx context.Context, incident *models.Incident) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.InsertOne(ctx, incidentCollection, incident)
	return err
}

// CreateIncidents creates many new Incidents.
func CreateIncidents(ctx context.Context, incidents *[]models.Incident) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	for index, incident := range *incidents {
		err = service.InsertOne(ctx, incidentCollection, &incident)
		if err != nil {
			return err
		}
		(*incidents)[index] = incident
	}
	return nil
}

// GetIncident retrieves an Incident.
func GetIncident(ctx context.Context, incidentID primitive.ObjectID, incident *models.Incident) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.FindOneByID(ctx, incidentCollection, incidentID, incident)
	return err
}

// ListIncidents retrieves many Incidents.
func ListIncidents(ctx context.Context, deploymentID primitive.ObjectID, incidents *[]models.Incident) error {
	filter := bson.M{"deployment_id": deploymentID}
	err := ListIncidentsByFilter(ctx, filter, incidents)
	return err
}

// ListIncidentsByFilter retrieves many Incidents conforming to a filter.
func ListIncidentsByFilter(ctx context.Context, filter bson.M, incidents *[]models.Incident) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.Find(ctx, incidentCollection, filter, incidents)
	return err
}

// UpdateIncident updates an Incident.
func UpdateIncident(ctx context.Context, incidentID primitive.ObjectID, incident *models.Incident) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.UpdateOne(ctx, incidentCollection, incidentID, &incident)
	if err != nil {
		return err
	}

	incident.ID = incidentID
	return nil
}

// DeleteIncident deletes an Incident.
func DeleteIncident(ctx context.Context, incidentID primitive.ObjectID) error {
	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		return err
	}
	defer service.Disconnect(ctx)

	err = service.DeleteOne(ctx, incidentCollection, incidentID)
	return err
}
