package ingest_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/unnmdnwb3/dora/internal/connectors/prometheus"
	"github.com/unnmdnwb3/dora/internal/daos"
	"github.com/unnmdnwb3/dora/internal/database/mongodb"
	"github.com/unnmdnwb3/dora/internal/models"
	"github.com/unnmdnwb3/dora/internal/services/trigger/ingest"
	"github.com/unnmdnwb3/dora/test"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ = Describe("services.trigger.import.incidents", func() {
	var (
		ctx            = context.Background()
		prometheusMock *httptest.Server
	)

	var _ = BeforeEach(func() {
		_ = godotenv.Load("./../../../../test/.env")

		var queryRangeResponse prometheus.QueryRangeResponse
		_ = test.UnmarshalFixture("./../../../../test/data/prometheus/query_range.json", &queryRangeResponse)
		prometheusMock = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			json, _ := json.Marshal(queryRangeResponse)
			w.Write(json)
		}))
	})

	var _ = AfterEach(func() {
		service := mongodb.NewService()
		service.Connect(ctx, os.Getenv("MONGODB_DATABASE"))
		service.DB.Drop(ctx)
		defer service.Disconnect(ctx)

		defer prometheusMock.Close()

		os.Remove("MONGODB_URI")
		os.Remove("MONGODB_PORT")
		os.Remove("MONGODB_USER")
		os.Remove("MONGODB_PASSWORD")
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

			monitoringDataPoints, err := ingest.ImportMonitoringDataPoints(ctx, &deployment)
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

			go ingest.ImportIncidents(ctx, channel, &deployment)
			err = <-channel
			Expect(err).To(BeNil())

			var incidents []models.Incident
			err = daos.ListIncidents(ctx, deployment.ID, &incidents)
			Expect(err).To(BeNil())
			Expect(len(incidents)).To(Equal(3))

		})
	})

	var _ = When("CalculateIncidents", func() {
		It("calculates Incidents based on MonitoringDataPoints.", func() {
			deployment := models.Deployment{
				IntegrationID: primitive.NewObjectID(),
				Query:         "job:http_total_requests:internal_server_error_percentage",
				Step:          "5m",
				Relation:      "gt",
				Threshold:     0.2,
			}

			dataPoints := []models.MonitoringDataPoint{
				{
					Value:     0.1,
					CreatedAt: time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
				},
				{
					Value:     0.3,
					CreatedAt: time.Date(2022, 12, 27, 13, 21, 43, 0, time.UTC),
				},
				{
					Value:     0.35,
					CreatedAt: time.Date(2022, 12, 27, 13, 26, 42, 0, time.UTC),
				},
				{
					Value:     0.4,
					CreatedAt: time.Date(2022, 12, 27, 13, 36, 44, 0, time.UTC),
				},
				{
					Value:     0.3,
					CreatedAt: time.Date(2022, 12, 27, 13, 41, 42, 0, time.UTC),
				},
				{
					Value:     0.1,
					CreatedAt: time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
				},
			}

			incidents, err := ingest.CalculateIncidents(ctx, &deployment, &dataPoints)
			Expect(err).To(BeNil())
			Expect(len(*incidents)).To(Equal(2))
			Expect((*incidents)[0].StartDate).To(Equal(time.Date(2022, 12, 27, 13, 21, 43, 0, time.UTC)))
			Expect((*incidents)[0].EndDate).To(Equal(time.Date(2022, 12, 27, 13, 26, 42, 0, time.UTC)))
			Expect((*incidents)[1].StartDate).To(Equal(time.Date(2022, 12, 27, 13, 36, 44, 0, time.UTC)))
			Expect((*incidents)[1].EndDate).To(Equal(time.Date(2022, 12, 27, 13, 41, 42, 0, time.UTC)))
		})
	})

	var _ = When("CreateIncidents", func() {
		It("creates Incidents based on MonitoringDataPoints.", func() {
			deployment := models.Deployment{
				IntegrationID: primitive.NewObjectID(),
				Query:         "job:http_total_requests:internal_server_error_percentage",
				Step:          "5m",
				Relation:      "gt",
				Threshold:     0.2,
			}

			dataPoints := []models.MonitoringDataPoint{
				{
					Value:     0.1,
					CreatedAt: time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
				},
				{
					Value:     0.3,
					CreatedAt: time.Date(2022, 12, 27, 13, 21, 43, 0, time.UTC),
				},
				{
					Value:     0.35,
					CreatedAt: time.Date(2022, 12, 27, 13, 26, 42, 0, time.UTC),
				},
				{
					Value:     0.4,
					CreatedAt: time.Date(2022, 12, 27, 13, 36, 44, 0, time.UTC),
				},
				{
					Value:     0.3,
					CreatedAt: time.Date(2022, 12, 27, 13, 41, 42, 0, time.UTC),
				},
				{
					Value:     0.1,
					CreatedAt: time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
				},
			}

			err := ingest.CreateIncidents(ctx, &deployment, &dataPoints)
			Expect(err).To(BeNil())

			var incidents []models.Incident
			err = daos.ListIncidentsByFilter(ctx, bson.M{"deployment_id": deployment.ID}, &incidents)
			Expect(incidents).To(HaveLen(2))
			Expect(incidents[0].StartDate).To(Equal(time.Date(2022, 12, 27, 13, 21, 43, 0, time.UTC)))
			Expect(incidents[0].EndDate).To(Equal(time.Date(2022, 12, 27, 13, 26, 42, 0, time.UTC)))
			Expect(incidents[1].StartDate).To(Equal(time.Date(2022, 12, 27, 13, 36, 44, 0, time.UTC)))
			Expect(incidents[1].EndDate).To(Equal(time.Date(2022, 12, 27, 13, 41, 42, 0, time.UTC)))
		})
	})
})
