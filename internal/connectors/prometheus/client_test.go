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

		client *prometheus.Client
	)

	var _ = BeforeEach(func() {
		_ = test.UnmarshalFixture("./../../../test/data/prometheus/query_range.json", &queryRangeResponse)
		mock = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			json, _ := json.Marshal(queryRangeResponse)
			w.Write(json)
		}))

		query := "up{app='dashboard-service'}"
		start := time.Unix(1673186400, 0)
		end := time.Unix(1674396000, 0)
		step := 600
		client = prometheus.NewClient(mock.URL, "", query, start, end, step)
	})

	var _ = AfterEach(func() {
		defer mock.Close()
	})

	var _ = When("CreateMonitoringDataPoints", func() {
		It("creates monitoring data points from a query response", func() {
			monitoringDataPoints, err := client.CreateMonitoringDataPoints(queryRangeResponse)
			Expect(err).To(BeNil())
			Expect(len(*monitoringDataPoints)).To(Equal(444))
			Expect((*monitoringDataPoints)[0].Value).To(Equal(1.0))
			Expect((*monitoringDataPoints)[0].CreatedAt).To(Equal(time.Unix(1674130200, 0)))
		})
	})

	var _ = When("GetMonitoringDataPoints", func() {
		It("creates monitoring data points from a query response", func() {
			monitoringDataPoints, err := client.GetMonitoringDataPoints()
			Expect(err).To(BeNil())
			Expect(len(*monitoringDataPoints)).To(Equal(444))
			Expect((*monitoringDataPoints)[0].Value).To(Equal(1.0))
			Expect((*monitoringDataPoints)[0].CreatedAt).To(Equal(time.Unix(1674130200, 0)))
		})
	})
})
