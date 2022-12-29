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

var _ = Describe("daos.incidentsPerDays", func() {
	ctx := context.Background()

	var _ = When("CreateIncidentsPerDay", func() {
		It("creates creates a new pipelineRunsPerDay.", func() {
			deploymentID := primitive.NewObjectID()
			date, _ := time.Parse(time.RFC3339, "2020-02-04T14:29:50.092Z")
			incidentsPerDay := models.IncidentsPerDay{
				DeploymentID:          deploymentID,
				Date:                  date,
				TotalIncidents:        1,
				TotalIncidentDuration: 60,
			}
			err := daos.CreateIncidentsPerDay(ctx, deploymentID, &incidentsPerDay)
			Expect(err).To(BeNil())
			Expect(incidentsPerDay.ID).To(Not(BeEmpty()))
		})
	})

	var _ = When("CreateIncidentsPerDays", func() {
		It("creates creates many new IncidentsPerDays.", func() {
			deploymentID := primitive.NewObjectID()
			date1, _ := time.Parse(time.RFC3339, "2020-02-04T00:00:00.000Z")
			incidentsPerDay1 := models.IncidentsPerDay{
				DeploymentID:          deploymentID,
				Date:                  date1,
				TotalIncidents:        1,
				TotalIncidentDuration: 60,
			}

			date2, _ := time.Parse(time.RFC3339, "2020-02-05T00:00:00.000Z")
			incidentsPerDay2 := models.IncidentsPerDay{
				DeploymentID:          deploymentID,
				Date:                  date2,
				TotalIncidents:        2,
				TotalIncidentDuration: 180,
			}
			incidentsPerDays := []models.IncidentsPerDay{incidentsPerDay1, incidentsPerDay2}
			err := daos.CreateIncidentsPerDays(ctx, deploymentID, &incidentsPerDays)
			Expect(err).To(BeNil())
			Expect(incidentsPerDays[0].ID).To(Not(BeEmpty()))
			Expect(incidentsPerDays[1].ID).To(Not(BeEmpty()))
		})
	})

	var _ = When("GetIncidentsPerDay", func() {
		It("retrieves an IncidentsPerDay.", func() {
			deploymentID := primitive.NewObjectID()
			date, _ := time.Parse(time.RFC3339, "2020-02-04T14:29:50.092Z")
			incidentsPerDay := models.IncidentsPerDay{
				DeploymentID:          deploymentID,
				Date:                  date,
				TotalIncidents:        1,
				TotalIncidentDuration: 60,
			}
			err := daos.CreateIncidentsPerDay(ctx, deploymentID, &incidentsPerDay)
			Expect(err).To(BeNil())
			Expect(incidentsPerDay.ID).To(Not(BeEmpty()))

			var findIncidentsPerDay models.IncidentsPerDay
			err = daos.GetIncidentsPerDay(ctx, incidentsPerDay.ID, &findIncidentsPerDay)
			Expect(err).To(BeNil())
			Expect(findIncidentsPerDay.ID).To(Equal(incidentsPerDay.ID))
		})
	})

	var _ = When("ListIncidentsPerDays", func() {
		It("retrieves many IncidentsPerDays.", func() {
			deploymentID := primitive.NewObjectID()
			date1, _ := time.Parse(time.RFC3339, "2020-02-04T00:00:00.000Z")
			incidentsPerDay1 := models.IncidentsPerDay{
				DeploymentID:          deploymentID,
				Date:                  date1,
				TotalIncidents:        1,
				TotalIncidentDuration: 60,
			}

			date2, _ := time.Parse(time.RFC3339, "2020-02-05T00:00:00.000Z")
			incidentsPerDay2 := models.IncidentsPerDay{
				DeploymentID:          deploymentID,
				Date:                  date2,
				TotalIncidents:        2,
				TotalIncidentDuration: 180,
			}
			incidentsPerDays := []models.IncidentsPerDay{incidentsPerDay1, incidentsPerDay2}
			err := daos.CreateIncidentsPerDays(ctx, deploymentID, &incidentsPerDays)
			Expect(err).To(BeNil())

			var findIncidentsPerDays []models.IncidentsPerDay
			err = daos.ListIncidentsPerDays(ctx, deploymentID, &findIncidentsPerDays)
			Expect(err).To(BeNil())
			Expect(len(findIncidentsPerDays)).To(Equal(2))
		})
	})

	var _ = When("ListIncidentsPerDaysByFilter", func() {
		It("retrieves many IncidentsPerDays conforming to a filter.", func() {
			deploymentID := primitive.NewObjectID()
			date1, _ := time.Parse(time.RFC3339, "2020-02-04T00:00:00.000Z")
			incidentsPerDay1 := models.IncidentsPerDay{
				DeploymentID:          deploymentID,
				Date:                  date1,
				TotalIncidents:        1,
				TotalIncidentDuration: 60,
			}

			date2, _ := time.Parse(time.RFC3339, "2020-02-05T00:00:00.000Z")
			incidentsPerDay2 := models.IncidentsPerDay{
				DeploymentID:          deploymentID,
				Date:                  date2,
				TotalIncidents:        2,
				TotalIncidentDuration: 180,
			}
			incidentsPerDays := []models.IncidentsPerDay{incidentsPerDay1, incidentsPerDay2}
			err := daos.CreateIncidentsPerDays(ctx, deploymentID, &incidentsPerDays)
			Expect(err).To(BeNil())

			var findIncidentsPerDays []models.IncidentsPerDay
			filter := bson.M{"date": bson.M{"$gte": date2}}
			err = daos.ListIncidentsPerDaysByFilter(ctx, filter, &findIncidentsPerDays)
			Expect(err).To(BeNil())
			Expect(findIncidentsPerDays).To(HaveLen(1))
		})
	})

	var _ = When("UpdateIncidentsPerDay", func() {
		It("updates an IncidentsPerDay.", func() {
			deploymentID := primitive.NewObjectID()
			date, _ := time.Parse(time.RFC3339, "2020-02-04T14:29:50.092Z")
			incidentsPerDay := models.IncidentsPerDay{
				DeploymentID:          deploymentID,
				Date:                  date,
				TotalIncidents:        1,
				TotalIncidentDuration: 60,
			}
			err := daos.CreateIncidentsPerDay(ctx, deploymentID, &incidentsPerDay)
			Expect(err).To(BeNil())
			Expect(incidentsPerDay.ID).To(Not(BeEmpty()))

			updateIncidentsPerDay := models.IncidentsPerDay{
				DeploymentID:          deploymentID,
				Date:                  date,
				TotalIncidents:        2,
				TotalIncidentDuration: 180,
			}
			err = daos.UpdateIncidentsPerDay(ctx, incidentsPerDay.ID, &updateIncidentsPerDay)
			Expect(err).To(BeNil())
			Expect(updateIncidentsPerDay.TotalIncidents).To(Equal(2))
		})
	})

	var _ = When("DeleteOne", func() {
		It("deletes a document with ID in a collection", func() {
			deploymentID := primitive.NewObjectID()
			date, _ := time.Parse(time.RFC3339, "2020-02-04T14:29:50.092Z")
			incidentsPerDay := models.IncidentsPerDay{
				DeploymentID:          deploymentID,
				Date:                  date,
				TotalIncidents:        1,
				TotalIncidentDuration: 60,
			}
			err := daos.CreateIncidentsPerDay(ctx, deploymentID, &incidentsPerDay)
			Expect(err).To(BeNil())
			Expect(incidentsPerDay.ID).To(Not(BeEmpty()))

			err = daos.DeleteIncidentsPerDay(ctx, incidentsPerDay.ID)
			Expect(err).To(BeNil())

			var findIncidentsPerDay models.IncidentsPerDay
			err = daos.GetIncidentsPerDay(ctx, incidentsPerDay.ID, &findIncidentsPerDay)
			Expect(err).To(Not(BeNil()))
		})
	})
})
