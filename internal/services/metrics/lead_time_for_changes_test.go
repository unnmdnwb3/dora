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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ = Describe("services.metrics.lead_time_for_changes", func() {
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

	var _ = When("CalculateLeadTimeForChanges", func() {
		It("calculates the CalculateLeadTime for a given window.", func() {
			// create a new dataflow
			repository := models.Repository{
				IntegrationID:  primitive.NewObjectID(),
				ExternalID:     40649465,
				NamespacedName: "foobar/foobar",
				DefaultBranch:  "main",
				URI:            "https://gitlab.com/foobar/foobar",
			}

			pipeline := models.Pipeline{
				IntegrationID:  primitive.NewObjectID(),
				ExternalID:     40649465,
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

			// create changes per days
			changesPerDays := []models.ChangesPerDay{
				{
					RepositoryID:  dataflow.Repository.ID,
					PipelineID:    dataflow.Pipeline.ID,
					Date:          time.Date(2022, 12, 24, 0, 0, 0, 0, time.UTC),
					TotalChanges:  1,
					TotalLeadTime: 600,
				},
				{
					RepositoryID:  dataflow.Repository.ID,
					PipelineID:    dataflow.Pipeline.ID,
					Date:          time.Date(2022, 12, 25, 0, 0, 0, 0, time.UTC),
					TotalChanges:  0,
					TotalLeadTime: 900,
				},
				{
					RepositoryID:  dataflow.Repository.ID,
					PipelineID:    dataflow.Pipeline.ID,
					Date:          time.Date(2022, 12, 26, 0, 0, 0, 0, time.UTC),
					TotalChanges:  2,
					TotalLeadTime: 600,
				},
				{
					RepositoryID:  dataflow.Repository.ID,
					PipelineID:    dataflow.Pipeline.ID,
					Date:          time.Date(2022, 12, 27, 0, 0, 0, 0, time.UTC),
					TotalChanges:  1,
					TotalLeadTime: 6000,
				},
				{
					RepositoryID:  dataflow.Repository.ID,
					PipelineID:    dataflow.Pipeline.ID,
					Date:          time.Date(2022, 12, 28, 0, 0, 0, 0, time.UTC),
					TotalChanges:  0,
					TotalLeadTime: 1200,
				},
				{
					RepositoryID:  dataflow.Repository.ID,
					PipelineID:    dataflow.Pipeline.ID,
					Date:          time.Date(2022, 12, 29, 0, 0, 0, 0, time.UTC),
					TotalChanges:  2,
					TotalLeadTime: 1800,
				},
				{
					RepositoryID:  dataflow.Repository.ID,
					PipelineID:    dataflow.Pipeline.ID,
					Date:          time.Date(2022, 12, 30, 0, 0, 0, 0, time.UTC),
					TotalChanges:  3,
					TotalLeadTime: 600,
				},
			}

			err = daos.CreateChangesPerDays(ctx, dataflow.Repository.ID, dataflow.Pipeline.ID, &changesPerDays)
			Expect(err).To(BeNil())

			// calculate lead time for changes rate
			endDate := time.Date(2022, 12, 29, 23, 59, 59, 0, time.UTC)
			window := 3
			leadTimeForChanges, err := metrics.CalculateLeadTimeForChanges(ctx, dataflow.ID, window, endDate)
			Expect(err).To(BeNil())
			Expect(leadTimeForChanges.DailyChanges).To(Equal([]int{1, 0, 2}))
			Expect(leadTimeForChanges.MovingAverages).To(Equal([]float64{2500, 2600, 3000}))
		})
	})
})
