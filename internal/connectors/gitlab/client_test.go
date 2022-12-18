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

	var _ = When("GetOrganisations", func() {
		It("get all organisations", func() {
			var fixture []models.Organisation
			err := test.UnmarshalFixture("./../../../test/data/gitlab/organisations.json", &fixture)
			Expect(err).To(BeNil())

			mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json, _ := json.Marshal(fixture)
				w.Write(json)

			}))
			defer mock.Close()

			client := gitlab.Client{
				Auth: "token",
				URI:  mock.URL,
			}

			organisations, err := client.GetOrganisations()
			Expect(err).To(BeNil())
			Expect(len(*organisations)).To(Equal(4))
		})
	})

	var _ = When("GetRepositories", func() {
		It("get all repositories", func() {
			var fixture []models.Repository
			err := test.UnmarshalFixture("./../../../test/data/gitlab/repositories.json", &fixture)
			Expect(err).To(BeNil())

			mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json, _ := json.Marshal(fixture)
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

	var _ = When("GetCommits", func() {
		It("get all commits of a repository", func() {
			var fixture []models.Commit
			err := test.UnmarshalFixture("./../../../test/data/gitlab/commits.json", &fixture)
			Expect(err).To(BeNil())

			mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json, _ := json.Marshal(fixture)
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

	var _ = When("GetPullRequests", func() {
		It("get all pull requests of a repository", func() {
			var fixture []models.PullRequest
			err := test.UnmarshalFixture("./../../../test/data/gitlab/pull_requests.json", &fixture)
			Expect(err).To(BeNil())

			mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json, _ := json.Marshal(fixture)
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

	var _ = When("PipelineRuns", func() {
		It("get all pipeline runs", func() {
			var fixture []models.PipelineRun
			err := test.UnmarshalFixture("./../../../test/data/gitlab/pipeline_runs.json", &fixture)
			Expect(err).To(BeNil())

			mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json, _ := json.Marshal(fixture)
				w.Write(json)

			}))
			defer mock.Close()

			client := gitlab.Client{
				Auth: "token",
				URI:  mock.URL,
			}

			pipelineRuns, err := client.GetPipelineRuns(15392086, referenceBranch)
			Expect(err).To(BeNil())
			Expect(len(*pipelineRuns)).To(Equal(4))
		})
	})
})
