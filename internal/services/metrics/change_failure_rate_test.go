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

var _ = Describe("services.metrics.change_failure_rate", func() {
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

	var _ = When("CalculateChangeFailureRate", func() {
		It("calculates the ChangeFailureRate for a given window.", func() {
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

			// create incidents per days
			incidentsPerDays := []models.IncidentsPerDay{
				{
					DeploymentID:   dataflow.Deployment.ID,
					Date:           time.Date(2022, 12, 24, 0, 0, 0, 0, time.UTC),
					TotalIncidents: 1,
					TotalDuration:  600,
				},
				{
					DeploymentID:   dataflow.Deployment.ID,
					Date:           time.Date(2022, 12, 25, 0, 0, 0, 0, time.UTC),
					TotalIncidents: 0,
					TotalDuration:  900,
				},
				{
					DeploymentID:   dataflow.Deployment.ID,
					Date:           time.Date(2022, 12, 26, 0, 0, 0, 0, time.UTC),
					TotalIncidents: 2,
					TotalDuration:  600,
				},
				{
					DeploymentID:   dataflow.Deployment.ID,
					Date:           time.Date(2022, 12, 27, 0, 0, 0, 0, time.UTC),
					TotalIncidents: 1,
					TotalDuration:  6000,
				},
				{
					DeploymentID:   dataflow.Deployment.ID,
					Date:           time.Date(2022, 12, 28, 0, 0, 0, 0, time.UTC),
					TotalIncidents: 0,
					TotalDuration:  1200,
				},
				{
					DeploymentID:   dataflow.Deployment.ID,
					Date:           time.Date(2022, 12, 29, 0, 0, 0, 0, time.UTC),
					TotalIncidents: 2,
					TotalDuration:  1600,
				},
				{
					DeploymentID:   dataflow.Deployment.ID,
					Date:           time.Date(2022, 12, 30, 0, 0, 0, 0, time.UTC),
					TotalIncidents: 3,
					TotalDuration:  600,
				},
			}

			err = daos.CreateIncidentsPerDays(ctx, dataflow.Deployment.ID, &incidentsPerDays)
			Expect(err).To(BeNil())

			// create pipeline runs per days
			pipelineRunsPerDays := []models.PipelineRunsPerDay{
				{
					PipelineID:        dataflow.Pipeline.ID,
					Date:              time.Date(2022, 12, 24, 0, 0, 0, 0, time.UTC),
					TotalPipelineRuns: 5,
				},
				{
					PipelineID:        dataflow.Pipeline.ID,
					Date:              time.Date(2022, 12, 25, 0, 0, 0, 0, time.UTC),
					TotalPipelineRuns: 4,
				},
				{
					PipelineID:        dataflow.Pipeline.ID,
					Date:              time.Date(2022, 12, 26, 0, 0, 0, 0, time.UTC),
					TotalPipelineRuns: 6,
				},
				{
					PipelineID:        dataflow.Pipeline.ID,
					Date:              time.Date(2022, 12, 27, 0, 0, 0, 0, time.UTC),
					TotalPipelineRuns: 2,
				},
				{
					PipelineID:        dataflow.Pipeline.ID,
					Date:              time.Date(2022, 12, 28, 0, 0, 0, 0, time.UTC),
					TotalPipelineRuns: 8,
				},
				{
					PipelineID:        dataflow.Pipeline.ID,
					Date:              time.Date(2022, 12, 29, 0, 0, 0, 0, time.UTC),
					TotalPipelineRuns: 5,
				},
				{
					PipelineID:        dataflow.Pipeline.ID,
					Date:              time.Date(2022, 12, 30, 0, 0, 0, 0, time.UTC),
					TotalPipelineRuns: 10,
				},
			}

			err = daos.CreatePipelineRunsPerDays(ctx, dataflow.Pipeline.ID, &pipelineRunsPerDays)
			Expect(err).To(BeNil())

			// calculate change failure rate
			endDate := time.Date(2022, 12, 29, 23, 59, 59, 0, time.UTC)
			window := 3
			changeFailureRate, err := metrics.CalculateChangeFailureRate(ctx, dataflow.ID, window, endDate)
			Expect(err).To(BeNil())
			Expect(changeFailureRate.DailyDeployments).To(Equal([]int{2, 8, 5}))
			Expect(changeFailureRate.DailyIncidents).To(Equal([]int{1, 0, 2}))
			Expect(changeFailureRate.MovingAverages).To(Equal([]float64{25, 18.75, 20}))
		})
	})
})
