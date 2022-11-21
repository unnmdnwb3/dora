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
	var organisations []models.Organisation
	_ = test.FromTestData("./../../../test/data/gitlab/organisations.json", &organisations)

	var _ = When("GetRepositories", func() {
		It("get all repositories", func() {
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

			projects, err := client.GetRepositories()
			Expect(err).To(BeNil())
			Expect(len(*projects)).To(Equal(2))
		})
	})
})
