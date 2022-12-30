package metrics_test

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
	"github.com/unnmdnwb3/dora/internal/services/metrics"
	"github.com/unnmdnwb3/dora/internal/services/trigger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ = Describe("services.metrics.mean_time_to_restore", func() {
	ctx := context.Background()

	var _ = BeforeEach(func() {
		_ = godotenv.Load("./../../../test/.env")
	})

	var _ = AfterEach(func() {
		ctx := context.Background()
		service := mongodb.NewService()
		service.Connect(ctx, os.Getenv("MONGODB_DATABASE"))
		service.DB.Drop(ctx)
		defer service.Disconnect(ctx)

		os.Remove("MONGODB_URI")
		os.Remove("MONGODB_PORT")
		os.Remove("MONGODB_USER")
		os.Remove("MONGODB_PASSWORD")
	})

	var _ = When("CalculateMeanTimeToRestore", func() {
		It("calculates the MeanTimeToRestore for a given window.", func() {
			// create a new dataflow
			externalID := 40649465
			repository := models.Repository{
				IntegrationID:  primitive.NewObjectID(),
				ExternalID:     externalID,
				NamespacedName: "foobar/foobar",
				DefaultBranch:  "main",
				URI:            "https://gitlab.com/foobar/foobar",
			}
			pipeline := models.Pipeline{
				IntegrationID:  primitive.NewObjectID(),
				ExternalID:     externalID,
				NamespacedName: "foobar/foobar",
				DefaultBranch:  "main",
				URI:            "https://gitlab.com/foobar/foobar/-/pipelines",
			}
			deployment := models.Deployment{
				IntegrationID: primitive.NewObjectID(),
				Query:         "job:http_total_requests:internal_server_error_percentage",
				Step:          "5m",
				Relation:      "gt",
				Threshold:     0.2,
			}
			dataflow := models.Dataflow{
				Repository: repository,
				Pipeline:   pipeline,
				Deployment: deployment,
			}
			err := daos.CreateDataflow(ctx, &dataflow)
			Expect(err).To(BeNil())

			// create incidents
			incidents := []models.Incident{
				{
					DeploymentID: dataflow.Deployment.ID,
					StartDate:    time.Date(2022, 12, 24, 13, 20, 56, 0, time.UTC),
					EndDate:      time.Date(2022, 12, 24, 13, 30, 56, 0, time.UTC),
				},
				{
					DeploymentID: dataflow.Deployment.ID,
					StartDate:    time.Date(2022, 12, 26, 8, 20, 56, 0, time.UTC),
					EndDate:      time.Date(2022, 12, 26, 8, 35, 56, 0, time.UTC),
				},
				{
					DeploymentID: dataflow.Deployment.ID,
					StartDate:    time.Date(2022, 12, 26, 10, 35, 56, 0, time.UTC),
					EndDate:      time.Date(2022, 12, 26, 10, 40, 56, 0, time.UTC),
				},
				{
					DeploymentID: dataflow.Deployment.ID,
					StartDate:    time.Date(2022, 12, 27, 9, 10, 56, 0, time.UTC),
					EndDate:      time.Date(2022, 12, 27, 9, 20, 56, 0, time.UTC),
				},
				{
					DeploymentID: dataflow.Deployment.ID,
					StartDate:    time.Date(2022, 12, 29, 9, 10, 56, 0, time.UTC),
					EndDate:      time.Date(2022, 12, 29, 9, 20, 56, 0, time.UTC),
				},
				{
					DeploymentID: dataflow.Deployment.ID,
					StartDate:    time.Date(2022, 12, 29, 9, 40, 56, 0, time.UTC),
					EndDate:      time.Date(2022, 12, 29, 9, 50, 56, 0, time.UTC),
				},
			}

			err = daos.CreateIncidents(ctx, &incidents)
			Expect(err).To(BeNil())

			// create incidents per days
			err = trigger.CreateIncidentsPerDays(ctx, dataflow.Deployment.ID)
			Expect(err).To(BeNil())

			// calculate mean time to restore
			endDate := time.Date(2022, 12, 29, 23, 59, 59, 0, time.UTC)
			window := 3
			meanTimeToRestore, err := metrics.CalculateMeanTimeToRestore(ctx, dataflow.ID, window, endDate)
			Expect(err).To(BeNil())
			Expect(meanTimeToRestore.DailyIncidents).To(Equal([]int{1, 0, 2}))
			Expect(meanTimeToRestore.DailyDurations).To(Equal([]int{600, 0, 1200}))
			Expect(meanTimeToRestore.MovingAverages).To(Equal([]float64{600, 600, 600}))
		})
	})
})
