package gitlab_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/unnmdnwb3/dora/test"

	"github.com/unnmdnwb3/dora/internal/connectors/gitlab"
	"github.com/unnmdnwb3/dora/internal/models"
)

func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "gitlab.Client Suite")
}

var _ = Describe("gitlab.Client", func() {
	projectID := "15392086"
	referenceBranch := "main"

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

			client := gitlab.Client{
				Auth: "token",
				URI:  mock.URL,
			}

			organisations, err := client.GetOrganisations()
			Expect(err).To(BeNil())
			Expect(len(*organisations)).To(Equal(2))
		})
	})

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

			client := gitlab.Client{
				Auth: "token",
				URI:  mock.URL,
			}

			repositories, err := client.GetRepositories()
			Expect(err).To(BeNil())
			Expect(len(*repositories)).To(Equal(1))
		})
	})

	var commits []models.Commit
	_ = test.FromTestData("./../../../test/data/gitlab/commits.json", &commits)

	var _ = When("GetCommits", func() {
		It("get all commits of a repository", func() {
			mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json, _ := json.Marshal(commits)
				w.Write(json)

			}))
			defer mock.Close()

			client := gitlab.Client{
				Auth: "token",
				URI:  mock.URL,
			}

			commits, err := client.GetCommits(projectID, referenceBranch)
			Expect(err).To(BeNil())
			Expect(len(*commits)).To(Equal(10))
		})
	})

	var pullRequests []models.PullRequest
	_ = test.FromTestData("./../../../test/data/gitlab/pull_requests.json", &pullRequests)

	var _ = When("GetPullRequests", func() {
		It("get all pull requests of a repository", func() {
			mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json, _ := json.Marshal(pullRequests)
				w.Write(json)

			}))
			defer mock.Close()

			client := gitlab.Client{
				Auth: "token",
				URI:  mock.URL,
			}

			pullRequests, err := client.GetPullRequests(projectID, referenceBranch)
			Expect(err).To(BeNil())
			Expect(len(*pullRequests)).To(Equal(3))
		})
	})

	var deployRuns []models.WorkflowRun
	_ = test.FromTestData("./../../../test/data/gitlab/workflow_runs.json", &deployRuns)

	var _ = When("WorkflowRuns", func() {
		It("get all deploy runs", func() {
			mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json, _ := json.Marshal(organisations)
				w.Write(json)

			}))
			defer mock.Close()

			client := gitlab.Client{
				Auth: "token",
				URI:  mock.URL,
			}

			deployRuns, err := client.GetWorkflowRuns(projectID, referenceBranch)
			Expect(err).To(BeNil())
			Expect(len(*deployRuns)).To(Equal(2))
		})
	})
})
