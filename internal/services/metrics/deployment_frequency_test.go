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

			// create pipeline runs
			createdAt1, _ := time.Parse(time.RFC3339, "2020-02-03T14:29:50.092Z")
			updatedAt1, _ := time.Parse(time.RFC3339, "2020-02-03T14:35:51.459Z")
			pipelineRun1 := models.PipelineRun{
				PipelineID:  dataflow.Pipeline.ID,
				ExternalID:  713437220,
				Sha:         "1cfffa2ae16528e36115ece8b1f2601bcf74414e",
				Ref:         "main",
				Status:      "success",
				EventSource: "push",
				CreatedAt:   createdAt1,
				UpdatedAt:   updatedAt1,
				URI:         "https://gitlab.com/foobar/foobar/-/pipelines/114884002",
			}

			createdAt2, _ := time.Parse(time.RFC3339, "2020-02-04T14:29:50.092Z")
			updatedAt2, _ := time.Parse(time.RFC3339, "2020-02-04T14:35:51.459Z")
			pipelineRun2 := models.PipelineRun{
				PipelineID:  dataflow.Pipeline.ID,
				ExternalID:  713437220,
				Sha:         "345207c839e94a939aebdc86835ae2e2a6c85acb",
				Ref:         "main",
				Status:      "success",
				EventSource: "push",
				CreatedAt:   createdAt2,
				UpdatedAt:   updatedAt2,
				URI:         "https://gitlab.com/foobar/foobar/-/pipelines/114884002",
			}

			createdAt3, _ := time.Parse(time.RFC3339, "2020-02-04T15:29:50.092Z")
			updatedAt3, _ := time.Parse(time.RFC3339, "2020-02-04T15:35:51.459Z")
			pipelineRun3 := models.PipelineRun{
				PipelineID:  dataflow.Pipeline.ID,
				ExternalID:  713437221,
				Sha:         "dcc7ef44dc6a376854c5f2cc42b0b24aa3a9ed10",
				Ref:         "main",
				Status:      "success",
				EventSource: "push",
				CreatedAt:   createdAt3,
				UpdatedAt:   updatedAt3,
				URI:         "https://gitlab.com/foobar/foobar/-/pipelines/114884003",
			}

			createdAt4, _ := time.Parse(time.RFC3339, "2020-02-06T14:29:50.092Z")
			updatedAt4, _ := time.Parse(time.RFC3339, "2020-02-06T14:35:51.459Z")
			pipelineRun4 := models.PipelineRun{
				PipelineID:  dataflow.Pipeline.ID,
				ExternalID:  713437222,
				Sha:         "05358053edd401e1fba272ecfe6f8839e9d8fdec",
				Ref:         "main",
				Status:      "success",
				EventSource: "push",
				CreatedAt:   createdAt4,
				UpdatedAt:   updatedAt4,
				URI:         "https://gitlab.com/foobar/foobar/-/pipelines/114884004",
			}

			createdAt5, _ := time.Parse(time.RFC3339, "2020-02-07T14:29:50.092Z")
			updatedAt5, _ := time.Parse(time.RFC3339, "2020-02-07T14:35:51.459Z")
			pipelineRun5 := models.PipelineRun{
				PipelineID:  dataflow.Pipeline.ID,
				ExternalID:  713437223,
				Sha:         "e6425994236ddd01a8ed24852d5f0ec19e70d701",
				Ref:         "main",
				Status:      "success",
				EventSource: "push",
				CreatedAt:   createdAt5,
				UpdatedAt:   updatedAt5,
				URI:         "https://gitlab.com/foobar/foobar/-/pipelines/114884005",
			}

			createdAt6, _ := time.Parse(time.RFC3339, "2020-02-07T15:29:50.092Z")
			updatedAt6, _ := time.Parse(time.RFC3339, "2020-02-07T15:35:51.459Z")
			pipelineRun6 := models.PipelineRun{
				PipelineID:  dataflow.Pipeline.ID,
				ExternalID:  713437224,
				Sha:         "19950955997854e605aca6b27f3b7bec6063069a",
				Ref:         "main",
				Status:      "success",
				EventSource: "push",
				CreatedAt:   createdAt6,
				UpdatedAt:   updatedAt6,
				URI:         "https://gitlab.com/foobar/foobar/-/pipelines/114884006",
			}
			pipelineRuns := []models.PipelineRun{pipelineRun1, pipelineRun2, pipelineRun3, pipelineRun4, pipelineRun5, pipelineRun6}
			err = daos.CreatePipelineRuns(ctx, dataflow.Pipeline.ID, &pipelineRuns)
			Expect(err).To(BeNil())

			// create pipeline runs aggregates
			pipelineRunsPerDays, err := trigger.CalculatePipelineRunsPerDays(ctx, &pipelineRuns)
			Expect(err).To(BeNil())
			err = daos.CreatePipelineRunsPerDays(ctx, dataflow.Pipeline.ID, pipelineRunsPerDays)
			Expect(err).To(BeNil())

			// actions above are assumed to be carried out by the trigger, when a new dataflow is created
			deploymentFrequency, err := metrics.CalculateDeploymentFrequency(ctx, dataflow.ID, 3, createdAt5)
			Expect(err).To(BeNil())
			Expect(deploymentFrequency.DataflowID).To(Equal(dataflow.ID))
			Expect(deploymentFrequency.MovingAverages).To(Equal([]float64{1.0, 1.0, 1.0}))
			Expect(deploymentFrequency.DailyPipelineRuns).To(Equal([]int{0, 1, 2}))
		})
	})
})
