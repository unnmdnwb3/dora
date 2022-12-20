package trigger_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/unnmdnwb3/dora/internal/daos"
	"github.com/unnmdnwb3/dora/internal/database/mongodb"
	"github.com/unnmdnwb3/dora/internal/models"
	"github.com/unnmdnwb3/dora/internal/services/trigger"
	"github.com/unnmdnwb3/dora/test"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ = Describe("services.trigger.import", func() {
	var (
		ctx        = context.Background()
		gitlabMock *httptest.Server

		externalID = 15392086
	)

	var _ = BeforeEach(func() {
		_ = godotenv.Load("./../../../test/.env")

		var pipelineRuns []models.PipelineRun
		_ = test.UnmarshalFixture("./../../../test/data/gitlab/pipeline_runs.json", &pipelineRuns)

		gitlabMock = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			json, _ := json.Marshal(pipelineRuns)
			w.Write(json)
		}))

		os.Setenv("GITLAB_URI", gitlabMock.URL)
	})

	var _ = AfterEach(func() {
		ctx := context.Background()

		service := mongodb.NewService()
		service.Connect(ctx, os.Getenv("MONGODB_DATABASE"))
		service.DB.Drop(ctx)
		defer service.Disconnect(ctx)

		defer gitlabMock.Close()

		os.Remove("GITLAB_BEARER")
		os.Remove("GITLAB_URI")
		os.Remove("MONGODB_URI")
		os.Remove("MONGODB_PORT")
		os.Remove("MONGODB_USER")
		os.Remove("MONGODB_PASSWORD")
	})

	var _ = When("ImportPipelineRuns", func() {
		It("gets all PipelineRuns of a Pipeline and persists them.", func() {
			pipeline := models.Pipeline{
				ID:             primitive.NewObjectID(),
				IntegrationID:  primitive.NewObjectID(),
				ExternalID:     externalID,
				NamespacedName: "foobar/foobar",
				DefaultBranch:  "main",
				URI:            "https://gitlab.com/foobar/foobar/-/pipelines",
			}
			channel := make(chan error)
			defer close(channel)
			go trigger.ImportPipelineRuns(ctx, channel, &pipeline)
			err := <-channel
			Expect(err).To(BeNil())

			var pipelineRuns []models.PipelineRun
			err = daos.ListPipelineRuns(ctx, pipeline.ID, &pipelineRuns)
			Expect(len(pipelineRuns)).To(Equal(4))
			Expect(err).To(BeNil())
		})
	})

	var _ = When("ImportData", func() {
		It("parallelizes the data import of all defined sources.", func() {
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
				TargetURI:     "https://localhost:9090",
			}
			dataflow := models.Dataflow{
				Repository: repository,
				Pipeline:   pipeline,
				Deployment: deployment,
			}
			err := daos.CreateDataflow(ctx, &dataflow)
			Expect(err).To(BeNil())

			err = trigger.ImportData(ctx, &dataflow)
			Expect(err).To(BeNil())

			var pipelineRuns []models.PipelineRun
			err = daos.ListPipelineRuns(ctx, dataflow.Pipeline.ID, &pipelineRuns)
			Expect(err).To(BeNil())
			Expect(len(pipelineRuns)).To(Equal(4))
		})
	})
})
