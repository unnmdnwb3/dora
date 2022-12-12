package services_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/unnmdnwb3/dora/internal/models"
	"github.com/unnmdnwb3/dora/internal/services"
	"github.com/unnmdnwb3/dora/test"
)

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "services Suite")
}

var _ = Describe("mongodb.Service", func() {
	ctx := context.Background()
	var gitlabMock *httptest.Server

	var _ = BeforeEach(func() {
		_ = godotenv.Load("./../../test/.env")

		var pipelineRuns []models.PipelineRun
		_ = test.FromTestData("./../../test/data/gitlab/pipeline_runs.json", &pipelineRuns)

		gitlabMock = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet && r.URL.Path == "/api/v4/projects/40649465/pipelines" {
				w.WriteHeader(http.StatusOK)
				json, _ := json.Marshal(pipelineRuns)
				w.Write(json)
			}
		}))

		os.Setenv("GITLAB_URI", gitlabMock.URL)
	})

	var _ = AfterEach(func() {
		defer gitlabMock.Close()

		os.Remove("GITLAB_BEARER")
		os.Remove("GITLAB_URI")
		os.Remove("MONGODB_URI")
		os.Remove("MONGODB_PORT")
		os.Remove("MONGODB_USER")
		os.Remove("MONGODB_PASSWORD")
	})

	var _ = When("CreateDataflow", func() {
		It("creates creates a new Dataflow.", func() {
			repository := models.Repository{
				IntegrationID:  "638e00b85edd5bef25e5e9e1",
				ExternalID:     "40649465",
				NamespacedName: "foobar/foobar",
				DefaultBranch:  "main",
				URI:            "https://gitlab.com/foobar/foobar",
			}
			pipeline := models.Pipeline{
				IntegrationID:  "638e00b85edd5vef25e5e9a2",
				ExternalID:     "40649465",
				NamespacedName: "foobar/foobar",
				DefaultBranch:  "main",
				URI:            "https://gitlab.com/foobar/foobar/-/pipelines",
			}
			deployment := models.Deployment{
				IntegrationID: "638e00b85edd5veff5e51b13",
				TargetURI:     "https://localhost:9090",
			}
			dataflow := &models.Dataflow{
				Repository: &repository,
				Pipeline:   &pipeline,
				Deployment: &deployment,
			}
			dataflow, err := services.CreateDataflow(ctx, dataflow)
			Expect(err).To(BeNil())
			Expect(dataflow.ID).To(Not(BeEmpty()))
		})
	})

	var _ = When("ImportPipelineRuns", func() {
		It("gets all PipelineRuns of a Pipeline and persists them.", func() {
			_ = godotenv.Load("./../../test/.env")

			var pipelineRuns []models.PipelineRun
			_ = test.FromTestData("./../../test/data/gitlab/pipeline_runs.json", &pipelineRuns)

			gitlabMock = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/projects/40649465/pipelines" {
					w.WriteHeader(http.StatusOK)
					json, _ := json.Marshal(pipelineRuns)
					w.Write(json)
				}
			}))
			defer gitlabMock.Close()

			os.Setenv("GITLAB_URI", gitlabMock.URL)

			pipeline := models.Pipeline{
				IntegrationID:  "638e00b85edd5vef25e5e9a2",
				ExternalID:     "40649465",
				NamespacedName: "foobar/foobar",
				DefaultBranch:  "main",
				URI:            "https://gitlab.com/foobar/foobar/-/pipelines",
			}
			channel := make(chan error)
			go services.ImportPipelineRuns(ctx, channel, &pipeline)
			err := <-channel
			Expect(err).To(BeNil())
		})
	})
})
