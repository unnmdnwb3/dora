package ingest_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/unnmdnwb3/dora/internal/daos"
	"github.com/unnmdnwb3/dora/internal/database/mongodb"
	"github.com/unnmdnwb3/dora/internal/models"
	"github.com/unnmdnwb3/dora/internal/services/trigger/ingest"
	"github.com/unnmdnwb3/dora/test"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ = Describe("services.ingest.commits", func() {
	var (
		ctx                  = context.Background()
		gitlabRepositoryMock *httptest.Server
	)

	var _ = BeforeEach(func() {
		_ = godotenv.Load("./../../../../test/.env")

		var commits []models.Commit
		_ = test.UnmarshalFixture("./../../../../test/data/gitlab/commits.json", &commits)
		gitlabRepositoryMock = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			json, _ := json.Marshal(commits)
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
				ExternalID:     15392086,
				NamespacedName: "foobar/foobar",
				DefaultBranch:  "main",
			}

			channel := make(chan error)
			defer close(channel)

			go ingest.ImportCommits(ctx, channel, &repository)
			err = <-channel
			Expect(err).To(BeNil())

			var commits []models.Commit
			err = daos.ListCommits(ctx, repository.ID, &commits)
			Expect(len(commits)).To(Equal(10))
			Expect(err).To(BeNil())
		})
	})
})
