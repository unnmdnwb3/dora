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

var _ = Describe("services.metrics.deployment_frequency", func() {
	var (
		ctx        = context.Background()
		externalID = 40649465
	)

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

	var _ = When("CalculateDeploymentFrequency", func() {
		It("calculates the DeploymentFrequency for a given window.", func() {
			// create a new dataflow
			repository := models.Repository{
				IntegrationID:  primitive.NewObjectID(),
				ExternalID:     externalID,
				NamespacedName: "foobar/foobar",
				DefaultBranch:  "main",
			}
			pipeline := models.Pipeline{
				IntegrationID:  primitive.NewObjectID(),
				ExternalID:     externalID,
				NamespacedName: "foobar/foobar",
				DefaultBranch:  "main",
			}
			deployment := models.Deployment{
				IntegrationID: primitive.NewObjectID(),
				Query:         "job:http_total_requests:internal_server_error_percentage",
				Step:          300,
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

			// create pipeline runs per days
			pipelineRunsPerDays := []models.PipelineRunsPerDay{
				{
					PipelineID:        dataflow.Pipeline.ID,
					Date:              time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC),
					TotalPipelineRuns: 1,
				},
				{
					PipelineID:        dataflow.Pipeline.ID,
					Date:              time.Date(2022, 2, 2, 0, 0, 0, 0, time.UTC),
					TotalPipelineRuns: 2,
				},
				{
					PipelineID:        dataflow.Pipeline.ID,
					Date:              time.Date(2022, 2, 3, 0, 0, 0, 0, time.UTC),
					TotalPipelineRuns: 0,
				},
				{
					PipelineID:        dataflow.Pipeline.ID,
					Date:              time.Date(2022, 2, 4, 0, 0, 0, 0, time.UTC),
					TotalPipelineRuns: 1,
				},
				{
					PipelineID:        dataflow.Pipeline.ID,
					Date:              time.Date(2022, 2, 5, 0, 0, 0, 0, time.UTC),
					TotalPipelineRuns: 2,
				},
				{
					PipelineID:        dataflow.Pipeline.ID,
					Date:              time.Date(2022, 2, 6, 0, 0, 0, 0, time.UTC),
					TotalPipelineRuns: 0,
				},
				{
					PipelineID:        dataflow.Pipeline.ID,
					Date:              time.Date(2022, 2, 7, 0, 0, 0, 0, time.UTC),
					TotalPipelineRuns: 1,
				},
				{
					PipelineID:        dataflow.Pipeline.ID,
					Date:              time.Date(2022, 2, 8, 0, 0, 0, 0, time.UTC),
					TotalPipelineRuns: 2,
				},
				{
					PipelineID:        dataflow.Pipeline.ID,
					Date:              time.Date(2022, 2, 9, 0, 0, 0, 0, time.UTC),
					TotalPipelineRuns: 0,
				},
			}
			err = daos.CreatePipelineRunsPerDays(ctx, dataflow.Pipeline.ID, &pipelineRunsPerDays)
			Expect(err).To(BeNil())

			startDate := time.Date(2022, 2, 4, 0, 0, 0, 0, time.UTC)
			endDate := time.Date(2022, 2, 9, 0, 0, 0, 0, time.UTC)
			window := 3

			deploymentFrequency, err := metrics.DeploymentFrequency(ctx, dataflow.ID, startDate, endDate, window)
			Expect(err).To(BeNil())
			Expect(deploymentFrequency.DataflowID).To(Equal(dataflow.ID))
			Expect(deploymentFrequency.MovingAverages).To(Equal([]float64{1.0, 1.0, 1.0, 1.0, 1.0, 1.0}))
			Expect(deploymentFrequency.DailyPipelineRuns).To(Equal([]int{1, 2, 0, 1, 2, 0}))
		})
	})
})
