package gitlab

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/unnmdnwb3/dora/test"

	"github.com/unnmdnwb3/dora/internal/models"
)

func TestGitlabClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "gitlab.Client Suite")
}

var _ = Describe("gitlab.Client", func() {
	var repositories []models.Repository
	_ = test.FromTestData("./../../../test/data/gitlab/repositories.json", &repositories)

	var _ = When("GetRepositories", func() {
		It("get all repositories", func() {
			mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json, _ := json.Marshal(repositories)
				w.Write(json)

			}))
			defer mock.Close()

			client := Client{
				Auth: "token",
				URI:  mock.URL,
			}

			repositories, err := client.GetRepositories()
			Expect(err).To(BeNil())
			Expect(len(*repositories)).To(Equal(1))
		})
	})

	var organisations []models.Organisation
	_ = test.FromTestData("./../../../test/data/gitlab/organisations.json", &organisations)

	var _ = When("GetOrganisations", func() {
		It("get all organisations", func() {
			mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json, _ := json.Marshal(organisations)
				w.Write(json)

			}))
			defer mock.Close()

			client := Client{
				Auth: "token",
				URI:  mock.URL,
			}

			organisations, err := client.GetOrganisations()
			Expect(err).To(BeNil())
			Expect(len(*organisations)).To(Equal(2))
		})
	})

	var deployRuns []models.DeployRun
	_ = test.FromTestData("./../../../test/data/gitlab/deploy_runs.json", &deployRuns)

	var _ = When("GetDeployRuns", func() {
		It("get all deploy runs", func() {
			mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json, _ := json.Marshal(organisations)
				w.Write(json)

			}))
			defer mock.Close()

			client := Client{
				Auth: "token",
				URI:  mock.URL,
			}

			projectID := "15392086"
			ref := "main"
			deployRuns, err := client.GetDeployRuns(projectID, ref)
			Expect(err).To(BeNil())
			Expect(len(*deployRuns)).To(Equal(2))
		})
	})
})
