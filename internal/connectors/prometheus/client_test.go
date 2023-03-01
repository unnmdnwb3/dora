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
		queryResponse prometheus.QueryResponse

		client *prometheus.Client
	)

	var _ = BeforeEach(func() {
		_ = test.UnmarshalFixture("./../../../test/data/prometheus/query.json", &queryResponse)
		mock = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			json, _ := json.Marshal(queryResponse)
			w.Write(json)
		}))

		query := "ALERTS{alertname='TargetDown', job='ak-core/log-processor-service'}[6w]"

		client = prometheus.NewClient(mock.URL, "", query)
	})

	var _ = AfterEach(func() {
		defer mock.Close()
	})

	var _ = When("CreateAlerts", func() {
		It("creates alerts from a query response", func() {
			alerts, err := client.CreateAlerts(queryResponse)
			Expect(err).To(BeNil())
			Expect(len(*alerts)).To(Equal(62))
			Expect((*alerts)[0].CreatedAt).To(Equal(time.Unix(1674486526, 0)))
		})
	})
})
