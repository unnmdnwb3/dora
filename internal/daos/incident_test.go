package daos_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/unnmdnwb3/dora/internal/daos"
	"github.com/unnmdnwb3/dora/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ = Describe("daos.incident", func() {
	ctx := context.Background()

	var _ = When("CreateIncident", func() {
		It("creates a new Incident.", func() {
			incident := models.Incident{
				DeploymentID: primitive.NewObjectID(),
				StartDate:    time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
				EndDate:      time.Date(2022, 12, 27, 13, 21, 42, 0, time.UTC),
			}
			err := daos.CreateIncident(ctx, &incident)
			Expect(err).To(BeNil())
			Expect(incident.ID).To(Not(BeEmpty()))
		})
	})

	var _ = When("CreatePipelineRuns", func() {
		It("creates creates many new PipelineRuns.", func() {
			deploymentID := primitive.NewObjectID()
			incident1 := models.Incident{
				DeploymentID: deploymentID,
				StartDate:    time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
				EndDate:      time.Date(2022, 12, 27, 13, 21, 42, 0, time.UTC),
			}

			incident2 := models.Incident{
				DeploymentID: deploymentID,
				StartDate:    time.Date(2022, 12, 27, 14, 16, 42, 0, time.UTC),
				EndDate:      time.Date(2022, 12, 27, 14, 21, 42, 0, time.UTC),
			}
			incidents := []models.Incident{incident1, incident2}
			err := daos.CreateIncidents(ctx, &incidents)
			Expect(err).To(BeNil())
			Expect(incidents[0].ID).To(Not(BeEmpty()))
			Expect(incidents[1].ID).To(Not(BeEmpty()))
		})
	})

	var _ = When("GetIncident", func() {
		It("retrieves an Incident.", func() {
			incident := models.Incident{
				DeploymentID: primitive.NewObjectID(),
				StartDate:    time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
				EndDate:      time.Date(2022, 12, 27, 13, 21, 42, 0, time.UTC),
			}
			err := daos.CreateIncident(ctx, &incident)
			Expect(err).To(BeNil())
			Expect(incident.ID).To(Not(BeEmpty()))

			var findIncident models.Incident
			err = daos.GetIncident(ctx, incident.ID, &findIncident)
			Expect(err).To(BeNil())
			Expect(findIncident.ID).To(Equal(incident.ID))
		})
	})

	var _ = When("ListIncidents", func() {
		It("retrieves many Incidents.", func() {
			deploymentID := primitive.NewObjectID()
			incident1 := models.Incident{
				DeploymentID: deploymentID,
				StartDate:    time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
				EndDate:      time.Date(2022, 12, 27, 13, 21, 42, 0, time.UTC),
			}
			incident2 := models.Incident{
				DeploymentID: deploymentID,
				StartDate:    time.Date(2022, 12, 27, 14, 51, 21, 0, time.UTC),
				EndDate:      time.Date(2022, 12, 27, 14, 59, 34, 0, time.UTC),
			}
			incident3 := models.Incident{
				DeploymentID: deploymentID,
				StartDate:    time.Date(2022, 12, 28, 21, 27, 40, 0, time.UTC),
				EndDate:      time.Date(2022, 12, 28, 21, 45, 46, 0, time.UTC),
			}
			_ = daos.CreateIncident(ctx, &incident1)
			_ = daos.CreateIncident(ctx, &incident2)
			_ = daos.CreateIncident(ctx, &incident3)
			Expect(incident1.ID).To(Not(BeNil()))
			Expect(incident2.ID).To(Not(BeNil()))
			Expect(incident3.ID).To(Not(BeNil()))

			var findIncidents []models.Incident
			err := daos.ListIncidents(ctx, deploymentID, &findIncidents)
			Expect(err).To(BeNil())
			Expect(findIncidents).To(HaveLen(3))
		})
	})

	var _ = When("ListIncidentsByFilter", func() {
		It("retrieves many Incidents conforming to a filter.", func() {
			deploymentID := primitive.NewObjectID()
			incident1 := models.Incident{
				DeploymentID: deploymentID,
				StartDate:    time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
				EndDate:      time.Date(2022, 12, 27, 13, 21, 42, 0, time.UTC),
			}
			incident2 := models.Incident{
				DeploymentID: deploymentID,
				StartDate:    time.Date(2022, 12, 27, 14, 51, 21, 0, time.UTC),
				EndDate:      time.Date(2022, 12, 27, 14, 59, 34, 0, time.UTC),
			}
			incident3 := models.Incident{
				DeploymentID: deploymentID,
				StartDate:    time.Date(2022, 12, 28, 21, 27, 40, 0, time.UTC),
				EndDate:      time.Date(2022, 12, 28, 21, 45, 46, 0, time.UTC),
			}
			_ = daos.CreateIncident(ctx, &incident1)
			_ = daos.CreateIncident(ctx, &incident2)
			_ = daos.CreateIncident(ctx, &incident3)
			Expect(incident1.ID).To(Not(BeNil()))
			Expect(incident2.ID).To(Not(BeNil()))
			Expect(incident3.ID).To(Not(BeNil()))

			var findIncidents []models.Incident
			filter := bson.M{"deployment_id": deploymentID}
			err := daos.ListIncidentsByFilter(ctx, filter, &findIncidents)
			Expect(err).To(BeNil())
			Expect(findIncidents).To(HaveLen(3))
		})
	})

	var _ = When("UpdateIncident", func() {
		It("updates an Incident.", func() {
			incident := models.Incident{
				DeploymentID: primitive.NewObjectID(),
				StartDate:    time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
				EndDate:      time.Date(2022, 12, 27, 13, 21, 42, 0, time.UTC),
			}
			err := daos.CreateIncident(ctx, &incident)
			Expect(err).To(BeNil())
			Expect(incident.ID).To(Not(BeEmpty()))

			updateIncident := models.Incident{
				DeploymentID: primitive.NewObjectID(),
				StartDate:    time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
				EndDate:      time.Date(2022, 12, 27, 13, 26, 42, 0, time.UTC),
			}
			err = daos.UpdateIncident(ctx, incident.ID, &updateIncident)
			Expect(err).To(BeNil())
			Expect(updateIncident.EndDate).To(Equal(time.Date(2022, 12, 27, 13, 26, 42, 0, time.UTC)))
		})
	})

	var _ = When("DeleteOne", func() {
		It("deletes a document with ID in a collection", func() {
			incident := models.Incident{
				DeploymentID: primitive.NewObjectID(),
				StartDate:    time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
				EndDate:      time.Date(2022, 12, 27, 13, 21, 42, 0, time.UTC),
			}
			err := daos.CreateIncident(ctx, &incident)
			Expect(err).To(BeNil())
			Expect(incident.ID).To(Not(BeEmpty()))

			err = daos.DeleteIncident(ctx, incident.ID)
			Expect(err).To(BeNil())

			var findIncident models.Incident
			err = daos.GetIncident(ctx, incident.ID, &findIncident)
			Expect(err).To(Not(BeNil()))
		})
	})
})
