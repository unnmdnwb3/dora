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
	"github.com/unnmdnwb3/dora/internal/connectors/prometheus"
	"github.com/unnmdnwb3/dora/internal/daos"
	"github.com/unnmdnwb3/dora/internal/database/mongodb"
	"github.com/unnmdnwb3/dora/internal/models"
	"github.com/unnmdnwb3/dora/internal/services/trigger"
	"github.com/unnmdnwb3/dora/test"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ = Describe("services.trigger.import", func() {
	var (
		gitlabRepositoryMock *httptest.Server
		gitlabPipelineMock   *httptest.Server
		prometheusMock       *httptest.Server

		ctx        = context.Background()
		externalID = 15392086
	)

	var _ = BeforeEach(func() {
		_ = godotenv.Load("./../../../test/.env")

		var commits []models.Commit
		_ = test.UnmarshalFixture("./../../../test/data/gitlab/commits.json", &commits)

		gitlabRepositoryMock = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			json, _ := json.Marshal(commits)
			w.Write(json)
		}))

		var pipelineRuns []models.PipelineRun
		_ = test.UnmarshalFixture("./../../../test/data/gitlab/pipeline_runs.json", &pipelineRuns)

		gitlabPipelineMock = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			json, _ := json.Marshal(pipelineRuns)
			w.Write(json)
		}))

		var queryRangeResponse prometheus.QueryRangeResponse
		_ = test.UnmarshalFixture("./../../../test/data/prometheus/query_range.json", &queryRangeResponse)
		prometheusMock = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			json, _ := json.Marshal(queryRangeResponse)
			w.Write(json)
		}))
	})

	var _ = AfterEach(func() {
		ctx := context.Background()

		service := mongodb.NewService()
		service.Connect(ctx, os.Getenv("MONGODB_DATABASE"))
		service.DB.Drop(ctx)
		defer service.Disconnect(ctx)

		defer gitlabRepositoryMock.Close()
		defer gitlabPipelineMock.Close()
		defer prometheusMock.Close()

		os.Remove("MONGODB_URI")
		os.Remove("MONGODB_PORT")
		os.Remove("MONGODB_USER")
		os.Remove("MONGODB_PASSWORD")
	})

	var _ = When("ImportCommits", func() {
		It("gets all Commits of a Repository and persists them.", func() {
			integration := models.Integration{
				ID:          primitive.NewObjectID(),
				Provider:    "gitlab",
				Type:        "vc",
				URI:         gitlabRepositoryMock.URL,
				BearerToken: "bearertoken",
			}
			err := daos.CreateIntegration(ctx, &integration)
			Expect(err).To(BeNil())

			repository := models.Repository{
				ID:             primitive.NewObjectID(),
				IntegrationID:  integration.ID,
				ExternalID:     externalID,
				NamespacedName: "foobar/foobar",
				DefaultBranch:  "main",
				URI:            "https://gitlab.com/foobar/foobar/-/pipelines",
			}

			channel := make(chan error)
			defer close(channel)

			go trigger.ImportCommits(ctx, channel, &repository)
			err = <-channel
			Expect(err).To(BeNil())

			var commits []models.Commit
			err = daos.ListCommits(ctx, repository.ID, &commits)
			Expect(len(commits)).To(Equal(10))
			Expect(err).To(BeNil())
		})
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
				ExternalID:     externalID,
				NamespacedName: "foobar/foobar",
				DefaultBranch:  "main",
				URI:            "https://gitlab.com/foobar/foobar/-/pipelines",
			}

			channel := make(chan error)
			defer close(channel)

			go trigger.ImportPipelineRuns(ctx, channel, &pipeline)
			err = <-channel
			Expect(err).To(BeNil())

			var pipelineRuns []models.PipelineRun
			err = daos.ListPipelineRuns(ctx, pipeline.ID, &pipelineRuns)
			Expect(len(pipelineRuns)).To(Equal(4))
			Expect(err).To(BeNil())
		})
	})

	var _ = When("ImportMonitoringDataPoints", func() {
		It("gets all MonitoringDataPoints of a Deployment.", func() {
			integration := models.Integration{
				ID:          primitive.NewObjectID(),
				Provider:    "prometheus",
				Type:        "im",
				URI:         prometheusMock.URL,
				BearerToken: "bearertoken",
			}
			err := daos.CreateIntegration(ctx, &integration)
			Expect(err).To(BeNil())

			deployment := models.Deployment{
				IntegrationID: integration.ID,
				Query:         "job:http_total_requests:internal_server_error_percentage",
				Step:          "5m",
				Relation:      "gt",
				Threshold:     0.2,
			}

			monitoringDataPoints, err := trigger.ImportMonitoringDataPoints(ctx, &deployment)
			Expect(err).To(BeNil())
			Expect(len(*monitoringDataPoints)).To(Equal(17))
		})
	})

	var _ = When("ImportIncidents", func() {
		It("gets all Incidents of a Deployment and persists them.", func() {
			integration := models.Integration{
				ID:          primitive.NewObjectID(),
				Provider:    "prometheus",
				Type:        "im",
				URI:         prometheusMock.URL,
				BearerToken: "bearertoken",
			}
			err := daos.CreateIntegration(ctx, &integration)
			Expect(err).To(BeNil())

			deployment := models.Deployment{
				IntegrationID: integration.ID,
				Query:         "job:http_total_requests:internal_server_error_percentage",
				Step:          "5m",
				Relation:      "gt",
				Threshold:     0.2,
			}

			channel := make(chan error)
			defer close(channel)

			go trigger.ImportIncidents(ctx, channel, &deployment)
			err = <-channel
			Expect(err).To(BeNil())

			var incidents []models.Incident
			err = daos.ListIncidents(ctx, deployment.ID, &incidents)
			Expect(err).To(BeNil())
			Expect(len(incidents)).To(Equal(3))

		})
	})

	var _ = When("ImportData", func() {
		It("parallelizes the data import of all defined sources.", func() {
			repositoryIntegration := models.Integration{
				ID:          primitive.NewObjectID(),
				Provider:    "gitlab",
				Type:        "vc",
				URI:         gitlabRepositoryMock.URL,
				BearerToken: "bearertoken",
			}
			err := daos.CreateIntegration(ctx, &repositoryIntegration)
			Expect(err).To(BeNil())

			pipelineIntegration := models.Integration{
				ID:          primitive.NewObjectID(),
				Provider:    "gitlab",
				Type:        "cicd",
				URI:         gitlabPipelineMock.URL,
				BearerToken: "bearertoken",
			}
			err = daos.CreateIntegration(ctx, &pipelineIntegration)
			Expect(err).To(BeNil())

			deploymentIntegration := models.Integration{
				ID:          primitive.NewObjectID(),
				Provider:    "prometheus",
				Type:        "im",
				URI:         prometheusMock.URL,
				BearerToken: "bearertoken",
			}
			err = daos.CreateIntegration(ctx, &deploymentIntegration)
			Expect(err).To(BeNil())

			repository := models.Repository{
				IntegrationID:  repositoryIntegration.ID,
				ExternalID:     externalID,
				NamespacedName: "foobar/foobar",
				DefaultBranch:  "main",
				URI:            "https://gitlab.com/foobar/foobar",
			}
			pipeline := models.Pipeline{
				IntegrationID:  pipelineIntegration.ID,
				ExternalID:     externalID,
				NamespacedName: "foobar/foobar",
				DefaultBranch:  "main",
				URI:            "https://gitlab.com/foobar/foobar/-/pipelines",
			}
			deployment := models.Deployment{
				IntegrationID: deploymentIntegration.ID,
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
			err = daos.CreateDataflow(ctx, &dataflow)
			Expect(err).To(BeNil())

			err = trigger.ImportData(ctx, &dataflow)
			Expect(err).To(BeNil())

			// check if data was imported
			var commits []models.Commit
			err = daos.ListCommits(ctx, dataflow.Repository.ID, &commits)
			Expect(err).To(BeNil())
			Expect(len(commits)).To(Equal(10))

			var pipelineRuns []models.PipelineRun
			err = daos.ListPipelineRuns(ctx, dataflow.Pipeline.ID, &pipelineRuns)
			Expect(err).To(BeNil())
			Expect(len(pipelineRuns)).To(Equal(4))

			var incidents []models.Incident
			err = daos.ListIncidents(ctx, dataflow.Deployment.ID, &incidents)
			Expect(err).To(BeNil())
			Expect(len(incidents)).To(Equal(3))
		})
	})
})
