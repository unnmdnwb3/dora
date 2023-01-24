package ingest_test

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
	"github.com/unnmdnwb3/dora/internal/services/trigger/ingest"
	"github.com/unnmdnwb3/dora/test"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ = Describe("services.ingest.pipeline_runs", func() {
	var (
		ctx                = context.Background()
		gitlabPipelineMock *httptest.Server
	)

	var _ = BeforeEach(func() {
		_ = godotenv.Load("./../../../../test/.env")

		var pipelineRuns []models.PipelineRun
		_ = test.UnmarshalFixture("./../../../../test/data/gitlab/pipeline_runs.json", &pipelineRuns)
		gitlabPipelineMock = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			json, _ := json.Marshal(pipelineRuns)
			w.Write(json)
		}))
	})

	var _ = AfterEach(func() {
		ctx := context.Background()

		service := mongodb.NewService()
		service.Connect(ctx, os.Getenv("MONGODB_DATABASE"))
		service.DB.Drop(ctx)
		defer service.Disconnect(ctx)

		defer gitlabPipelineMock.Close()

		os.Remove("MONGODB_URI")
		os.Remove("MONGODB_PORT")
		os.Remove("MONGODB_USER")
		os.Remove("MONGODB_PASSWORD")
	})

	var _ = When("ImportPipelineRuns", func() {
		It("gets all PipelineRuns of a Pipeline and persists them.", func() {
			integration := models.Integration{
				ID:          primitive.NewObjectID(),
				Provider:    "gitlab",
				Type:        "cicd",
				URI:         gitlabPipelineMock.URL,
				BearerToken: "bearertoken",
			}
			err := daos.CreateIntegration(ctx, &integration)
			Expect(err).To(BeNil())

			pipeline := models.Pipeline{
				ID:             primitive.NewObjectID(),
				IntegrationID:  integration.ID,
				ExternalID:     15392086,
				NamespacedName: "foobar/foobar",
				DefaultBranch:  "main",
			}

			channel := make(chan error)
			defer close(channel)

			go ingest.ImportPipelineRuns(ctx, channel, &pipeline)
			err = <-channel
			Expect(err).To(BeNil())

			var pipelineRuns []models.PipelineRun
			err = daos.ListPipelineRuns(ctx, pipeline.ID, &pipelineRuns)
			Expect(len(pipelineRuns)).To(Equal(4))
			Expect(err).To(BeNil())
		})
	})
})
