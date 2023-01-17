package aggregate_test

import (
	"context"
	"os"
	"time"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/unnmdnwb3/dora/internal/daos"
	"github.com/unnmdnwb3/dora/internal/database/mongodb"
	"github.com/unnmdnwb3/dora/internal/models"
	"github.com/unnmdnwb3/dora/internal/services/trigger/aggregate"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ = Describe("services.trigger.aggregate", func() {
	ctx := context.Background()

	var _ = BeforeEach(func() {
		_ = godotenv.Load("./../../../../test/.env")
	})

	var _ = AfterEach(func() {
		service := mongodb.NewService()
		service.Connect(ctx, os.Getenv("MONGODB_DATABASE"))
		service.DB.Drop(ctx)
		defer service.Disconnect(ctx)

		os.Remove("MONGODB_URI")
		os.Remove("MONGODB_PORT")
		os.Remove("MONGODB_USER")
		os.Remove("MONGODB_PASSWORD")
	})

	var _ = When("CalculateIncidentsPerDays", func() {
		It("calculates IncidentsPerDays based on Incidents.", func() {
			deploymentID := primitive.NewObjectID()
			incidents := []models.Incident{
				{
					DeploymentID: deploymentID,
					StartDate:    time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
					EndDate:      time.Date(2022, 12, 27, 14, 16, 42, 0, time.UTC),
				},
				{
					DeploymentID: deploymentID,
					StartDate:    time.Date(2022, 12, 27, 17, 39, 21, 0, time.UTC),
					EndDate:      time.Date(2022, 12, 27, 17, 54, 21, 0, time.UTC),
				},
				{
					DeploymentID: deploymentID,
					StartDate:    time.Date(2022, 12, 29, 02, 21, 42, 0, time.UTC),
					EndDate:      time.Date(2022, 12, 29, 02, 41, 42, 0, time.UTC),
				},
			}

			err := daos.CreateIncidents(ctx, &incidents)
			Expect(err).To(BeNil())

			incidentsPerDays, err := aggregate.CalculateIncidentsPerDays(ctx, &incidents)
			Expect(err).To(BeNil())
			Expect(len(*incidentsPerDays)).To(Equal(2))
			Expect((*incidentsPerDays)[0].TotalIncidents).To(Equal(2))
			Expect((*incidentsPerDays)[0].TotalDuration).To(Equal(float64(4500)))
			Expect((*incidentsPerDays)[1].TotalIncidents).To(Equal(1))
			Expect((*incidentsPerDays)[1].TotalDuration).To(Equal(float64(1200)))
		})
	})

	var _ = When("CreateIncidentsPerDays", func() {
		It("creates IncidentsPerDays based on Incidents.", func() {
			deploymentID := primitive.NewObjectID()
			incidents := []models.Incident{
				{
					DeploymentID: deploymentID,
					StartDate:    time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
					EndDate:      time.Date(2022, 12, 27, 14, 16, 42, 0, time.UTC),
				},
				{
					DeploymentID: deploymentID,
					StartDate:    time.Date(2022, 12, 27, 17, 39, 21, 0, time.UTC),
					EndDate:      time.Date(2022, 12, 27, 17, 54, 21, 0, time.UTC),
				},
				{
					DeploymentID: deploymentID,
					StartDate:    time.Date(2022, 12, 29, 02, 21, 42, 0, time.UTC),
					EndDate:      time.Date(2022, 12, 29, 02, 41, 42, 0, time.UTC),
				},
			}

			err := daos.CreateIncidents(ctx, &incidents)
			Expect(err).To(BeNil())

			channel := make(chan error)
			defer close(channel)

			go aggregate.CreateIncidentsPerDays(ctx, channel, deploymentID)
			err = <-channel
			Expect(err).To(BeNil())

			var incidentsPerDays []models.IncidentsPerDay
			err = daos.ListIncidentsPerDays(ctx, deploymentID, &incidentsPerDays)
			Expect(err).To(BeNil())
			Expect(incidentsPerDays).To(HaveLen(2))
			Expect(incidentsPerDays[0].TotalIncidents).To(Equal(2))
			Expect(incidentsPerDays[0].TotalDuration).To(Equal(float64(4500)))
			Expect(incidentsPerDays[1].TotalIncidents).To(Equal(1))
			Expect(incidentsPerDays[1].TotalDuration).To(Equal(float64(1200)))
		})
	})
})
