package trigger_test

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
	"github.com/unnmdnwb3/dora/internal/services/trigger"
	"github.com/unnmdnwb3/dora/internal/utils/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ = Describe("services.trigger.aggregate", func() {
	ctx := context.Background()
	externalID := 40649465

	var _ = BeforeEach(func() {
		_ = godotenv.Load("./../../../test/.env")
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

	var _ = When("CalculatePipelineRunsPerDays", func() {
		It("calculates the pipeline runs per day.", func() {
			pipelineID, _ := types.StringToObjectID("638e00b85edd5bef25e5e9e1")
			createdAt1, _ := time.Parse(time.RFC3339, "2019-10-09T09:11:20.861Z")
			updatedAt1, _ := time.Parse(time.RFC3339, "2019-10-09T09:12:20.861Z")
			pipelineRun1 := models.PipelineRun{
				PipelineID:  pipelineID,
				ExternalID:  713437220,
				Ref:         "main",
				Status:      "success",
				EventSource: "push",
				CreatedAt:   createdAt1,
				UpdatedAt:   updatedAt1,
				URI:         "https://gitlab.com/foobar/foobar/-/pipelines/114884002",
			}

			createdAt2, _ := time.Parse(time.RFC3339, "2019-10-11T09:11:20.861Z")
			updatedAt2, _ := time.Parse(time.RFC3339, "2019-10-11T09:12:20.861Z")
			pipelineRun2 := models.PipelineRun{
				PipelineID:  pipelineID,
				ExternalID:  713437221,
				Ref:         "main",
				Status:      "success",
				EventSource: "push",
				CreatedAt:   createdAt2,
				UpdatedAt:   updatedAt2,
				URI:         "https://gitlab.com/foobar/foobar/-/pipelines/114884003",
			}

			createdAt3, _ := time.Parse(time.RFC3339, "2019-10-11T09:13:20.861Z")
			updatedAt3, _ := time.Parse(time.RFC3339, "2019-10-11T09:14:20.861Z")
			pipelineRun3 := models.PipelineRun{
				PipelineID:  pipelineID,
				ExternalID:  713437222,
				Ref:         "main",
				Status:      "success",
				EventSource: "push",
				CreatedAt:   createdAt3,
				UpdatedAt:   updatedAt3,
				URI:         "https://gitlab.com/foobar/foobar/-/pipelines/114884004",
			}

			pipelineRuns := []models.PipelineRun{pipelineRun1, pipelineRun2, pipelineRun3}
			pipelineRunsPerDay, err := trigger.CalculatePipelineRunsPerDays(ctx, &pipelineRuns)
			Expect(err).To(BeNil())
			Expect(len(*pipelineRunsPerDay)).To(Equal(2))
			Expect((*pipelineRunsPerDay)[0].TotalPipelineRuns).To(Equal(1))
			Expect((*pipelineRunsPerDay)[1].TotalPipelineRuns).To(Equal(2))
		})
	})

	var _ = When("CreatePipelineRunsPerDays", func() {
		It("calculates and creates the pipeline runs for each day.", func() {
			// create dataflow
			repositoryIntegrationID := primitive.NewObjectID()
			pipelineIntegrationID := primitive.NewObjectID()
			deploymentIntegrationID := primitive.NewObjectID()
			repository := models.Repository{
				IntegrationID:  repositoryIntegrationID,
				ExternalID:     externalID,
				NamespacedName: "foobar/foobar",
				DefaultBranch:  "main",
				URI:            "https://gitlab.com/foobar/foobar",
			}
			pipeline := models.Pipeline{
				IntegrationID:  pipelineIntegrationID,
				ExternalID:     externalID,
				NamespacedName: "foobar/foobar",
				DefaultBranch:  "main",
				URI:            "https://gitlab.com/foobar/foobar/-/pipelines",
			}
			deployment := models.Deployment{
				IntegrationID: deploymentIntegrationID,
				TargetURI:     "https://localhost:9090",
			}
			dataflow := models.Dataflow{
				Repository: repository,
				Pipeline:   pipeline,
				Deployment: deployment,
			}
			err := daos.CreateDataflow(ctx, &dataflow)
			Expect(err).To(BeNil())

			// create pipeline runs
			createdAt1, _ := time.Parse(time.RFC3339, "2019-10-09T09:11:20.861Z")
			updatedAt1, _ := time.Parse(time.RFC3339, "2019-10-09T09:12:20.861Z")
			pipelineRun1 := models.PipelineRun{
				PipelineID:  dataflow.Pipeline.ID,
				ExternalID:  713437220,
				Ref:         "main",
				Status:      "success",
				EventSource: "push",
				CreatedAt:   createdAt1,
				UpdatedAt:   updatedAt1,
				URI:         "https://gitlab.com/foobar/foobar/-/pipelines/114884002",
			}

			createdAt2, _ := time.Parse(time.RFC3339, "2019-10-11T09:11:20.861Z")
			updatedAt2, _ := time.Parse(time.RFC3339, "2019-10-11T09:12:20.861Z")
			pipelineRun2 := models.PipelineRun{
				PipelineID:  dataflow.Pipeline.ID,
				ExternalID:  713437221,
				Ref:         "main",
				Status:      "success",
				EventSource: "push",
				CreatedAt:   createdAt2,
				UpdatedAt:   updatedAt2,
				URI:         "https://gitlab.com/foobar/foobar/-/pipelines/114884003",
			}

			createdAt3, _ := time.Parse(time.RFC3339, "2019-10-11T09:13:20.861Z")
			updatedAt3, _ := time.Parse(time.RFC3339, "2019-10-11T09:14:20.861Z")
			pipelineRun3 := models.PipelineRun{
				PipelineID:  dataflow.Pipeline.ID,
				ExternalID:  713437222,
				Ref:         "main",
				Status:      "success",
				EventSource: "push",
				CreatedAt:   createdAt3,
				UpdatedAt:   updatedAt3,
				URI:         "https://gitlab.com/foobar/foobar/-/pipelines/114884004",
			}

			pipelineRuns := []models.PipelineRun{pipelineRun1, pipelineRun2, pipelineRun3}
			err = daos.CreatePipelineRuns(ctx, dataflow.Pipeline.ID, &pipelineRuns)
			Expect(err).To(BeNil())

			err = trigger.CreatePipelineRunsPerDays(ctx, dataflow.Pipeline.ID)
			Expect(err).To(BeNil())

			var pipelineRunsPerDays []models.PipelineRunsPerDay
			err = daos.ListPipelineRunsPerDays(ctx, dataflow.Pipeline.ID, &pipelineRunsPerDays)
			Expect(err).To(BeNil())
			Expect(pipelineRunsPerDays).To(HaveLen(2))
		})
	})
})
