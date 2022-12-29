package prometheus_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/unnmdnwb3/dora/internal/connectors/prometheus"
	"github.com/unnmdnwb3/dora/test"
)

func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "prometheus.Client Suite")
}

var _ = Describe("prometheus.Client", func() {

	var (
		mock               *httptest.Server
		queryRangeResponse prometheus.QueryRangeResponse
		client             *prometheus.Client
	)

	var _ = BeforeEach(func() {
		_ = test.UnmarshalFixture("./../../../test/data/prometheus/query_range.json", &queryRangeResponse)
		mock = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			json, _ := json.Marshal(queryRangeResponse)
			w.Write(json)
		}))

		query := "job:http_total_requests:internal_server_error_percentage"
		client = prometheus.NewClient(mock.URL, "", query, time.Unix(1671556422, 0), time.Unix(1672161144, 0), "1m")
	})

	var _ = AfterEach(func() {
		defer mock.Close()
	})

	var _ = When("CreateMonitoringDataPoints", func() {
		It("creates monitoring data points from a query response", func() {
			monitoringDataPoints, err := client.CreateMonitoringDataPoints(queryRangeResponse)
			Expect(err).To(BeNil())
			Expect(len(*monitoringDataPoints)).To(Equal(17))
			Expect((*monitoringDataPoints)[16].Value).To(Equal(0.281))
			Expect((*monitoringDataPoints)[16].CreatedAt).To(Equal(time.Date(2022, 12, 28, 14, 43, 12, 0, time.UTC)))
		})
	})

	var _ = When("GetMonitoringDataPoints", func() {
		It("creates monitoring data points from a query response", func() {
			monitoringDataPoints, err := client.GetMonitoringDataPoints()
			Expect(err).To(BeNil())
			Expect(len(*monitoringDataPoints)).To(Equal(17))
			Expect((*monitoringDataPoints)[16].Value).To(Equal(0.281))
			Expect((*monitoringDataPoints)[16].CreatedAt).To(Equal(time.Date(2022, 12, 28, 14, 43, 12, 0, time.UTC)))
		})
	})
})
