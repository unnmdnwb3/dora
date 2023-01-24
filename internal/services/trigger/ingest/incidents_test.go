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
				Query:         "up{app='dashboard-service'}",
				Step:          600,
				Relation:      "st",
				Threshold:     1.0,
			}

			monitoringDataPoints, err := ingest.ImportMonitoringDataPoints(ctx, &deployment)
			Expect(err).To(BeNil())
			Expect(len(*monitoringDataPoints)).To(Equal(444))
		})
	})

	var _ = When("CalculateIncidents", func() {
		It("calculates Incidents based on MonitoringDataPoints.", func() {
			deployment := models.Deployment{
				IntegrationID: primitive.NewObjectID(),
				Query:         "up{app='dashboard-service'}",
				Step:          600,
				Relation:      "st",
				Threshold:     1.0,
			}

			dataPoints := []models.MonitoringDataPoint{
				{
					Value:     1.0,
					CreatedAt: time.Unix(1674130200, 0),
				},
				{
					Value:     1.0,
					CreatedAt: time.Unix(1674130800, 0),
				},
				{
					Value:     1.0,
					CreatedAt: time.Unix(1674131400, 0),
				},
				{
					Value:     0.0,
					CreatedAt: time.Unix(1674132000, 0),
				},
				{
					Value:     0.0,
					CreatedAt: time.Unix(1674132600, 0),
				},
				{
					Value:     1.0,
					CreatedAt: time.Unix(1674133200, 0),
				},
				{
					Value:     0.0,
					CreatedAt: time.Unix(1674133800, 0),
				},
				{
					Value:     0.0,
					CreatedAt: time.Unix(1674134400, 0),
				},
				{
					Value:     0.0,
					CreatedAt: time.Unix(1674135000, 0),
				},
				{
					Value:     1.0,
					CreatedAt: time.Unix(1674135600, 0),
				},
			}

			incidents, err := ingest.CalculateIncidents(ctx, &deployment, &dataPoints)
			Expect(err).To(BeNil())
			Expect(len(*incidents)).To(Equal(2))
			Expect((*incidents)[0].StartDate).To(Equal(time.Unix(1674132000, 0)))
			Expect((*incidents)[0].EndDate).To(Equal(time.Unix(1674132600, 0)))
			Expect((*incidents)[1].StartDate).To(Equal(time.Unix(1674133800, 0)))
			Expect((*incidents)[1].EndDate).To(Equal(time.Unix(1674135000, 0)))

		})
	})

	var _ = When("CreateIncidents", func() {
		It("creates Incidents based on MonitoringDataPoints.", func() {
			deployment := models.Deployment{
				IntegrationID: primitive.NewObjectID(),
				Query:         "up{app='dashboard-service'}",
				Step:          600,
				Relation:      "st",
				Threshold:     1.0,
			}

			dataPoints := []models.MonitoringDataPoint{
				{
					Value:     1.0,
					CreatedAt: time.Unix(1674130200, 0),
				},
				{
					Value:     1.0,
					CreatedAt: time.Unix(1674130800, 0),
				},
				{
					Value:     1.0,
					CreatedAt: time.Unix(1674131400, 0),
				},
				{
					Value:     0.0,
					CreatedAt: time.Unix(1674132000, 0),
				},
				{
					Value:     0.0,
					CreatedAt: time.Unix(1674132600, 0),
				},
				{
					Value:     1.0,
					CreatedAt: time.Unix(1674133200, 0),
				},
				{
					Value:     0.0,
					CreatedAt: time.Unix(1674133800, 0),
				},
				{
					Value:     0.0,
					CreatedAt: time.Unix(1674134400, 0),
				},
				{
					Value:     0.0,
					CreatedAt: time.Unix(1674135000, 0),
				},
				{
					Value:     1.0,
					CreatedAt: time.Unix(1674135600, 0),
				},
			}

			err := ingest.CreateIncidents(ctx, &deployment, &dataPoints)
			Expect(err).To(BeNil())

			var incidents []models.Incident
			err = daos.ListIncidentsByFilter(ctx, bson.M{"deployment_id": deployment.ID}, &incidents)
			Expect(incidents).To(HaveLen(2))
			Expect(incidents[0].StartDate).To(Equal(time.Unix(1674132000, 0).UTC()))
			Expect(incidents[0].EndDate).To(Equal(time.Unix(1674132600, 0).UTC()))
			Expect(incidents[1].StartDate).To(Equal(time.Unix(1674133800, 0).UTC()))
			Expect(incidents[1].EndDate).To(Equal(time.Unix(1674135000, 0).UTC()))
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
				Query:         "up{app='dashboard-service'}",
				Step:          600,
				Relation:      "st",
				Threshold:     1.0,
			}

			channel := make(chan error)
			defer close(channel)

			go ingest.ImportIncidents(ctx, channel, &deployment)
			err = <-channel
			Expect(err).To(BeNil())

			var incidents []models.Incident
			err = daos.ListIncidents(ctx, deployment.ID, &incidents)
			Expect(err).To(BeNil())
			Expect(len(incidents)).To(Equal(5))
			Expect(incidents[0].StartDate).To(Equal(time.Unix(1674131400, 0).UTC()))
			Expect(incidents[0].EndDate).To(Equal(time.Unix(1674132600, 0).UTC()))
		})
	})
})
