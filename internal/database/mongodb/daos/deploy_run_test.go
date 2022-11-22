package daos_test

import (
	"context"
	"os"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/unnmdnwb3/dora/internal/database/mongodb"
	"github.com/unnmdnwb3/dora/internal/database/mongodb/daos"
	"github.com/unnmdnwb3/dora/internal/models"
)

func TestDeployRuns(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "daos.DeployRun Suite")
}

var _ = BeforeEach(func() {
	// TODO remove config from code
	os.Setenv("MONGODB_URI", "127.0.0.1")
	os.Setenv("MONGODB_PORT", "27017")
	os.Setenv("MONGODB_USER", "user")
	os.Setenv("MONGODB_PASSWORD", "password")

	ctx := context.Background()
	mongodb.Init(&ctx)
})

var _ = AfterEach(func() {
	// TODO remove config from code
	os.Remove("MONGODB_URI")
	os.Remove("MONGODB_PORT")
	os.Remove("MONGODB_USER")
	os.Remove("MONGODB_PASSWORD")

	ctx := context.Background()
	defer mongodb.Client.Disconnect(ctx)
})

var _ = Describe("daos.DeployRun", func() {
	ctx := context.Background()

	var _ = When("Create", func() {
		It("creates a deploy run", func() {
			DAO, err := daos.NewDeployRun(&ctx)
			Expect(err).To((BeNil()))

			createdAt, _ := time.Parse(time.RFC3339, "2020-02-04T14:29:50.092Z")
			updatedAt, _ := time.Parse(time.RFC3339, "2020-02-04T14:35:51.459Z")
			createRequest := models.DeployRun{
				ProjectID: "15392086",
				Ref:       "main",
				Status:    "success",
				Source:    "push",
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
				URI:       "https://gitlab.com/foobar/foobar/-/pipelines/114883218",
			}

			ID, err := DAO.Create(&createRequest)
			Expect(err).To(BeNil())
			Expect(ID).To(Not(BeEmpty()))
		})
	})

	var _ = When("CreateMany", func() {
		It("creates many deploy runs", func() {
			DAO, err := daos.NewDeployRun(&ctx)
			Expect(err).To((BeNil()))

			firstCreatedAt, _ := time.Parse(time.RFC3339, "2020-02-04T14:29:50.092Z")
			firstUpdatedAt, _ := time.Parse(time.RFC3339, "2020-02-04T14:35:51.459Z")
			firstCreateRequest := models.DeployRun{
				ProjectID: "15392086",
				Ref:       "main",
				Status:    "success",
				Source:    "push",
				CreatedAt: firstCreatedAt,
				UpdatedAt: firstUpdatedAt,
				URI:       "https://gitlab.com/foobar/foobar/-/pipelines/114883218",
			}

			secondCreatedAt, _ := time.Parse(time.RFC3339, "2020-02-04T14:31:51.235Z")
			secondUpdatedAt, _ := time.Parse(time.RFC3339, "2020-02-04T14:38:18.252Z")
			secondCreateRequest := models.DeployRun{
				ProjectID: "15392086",
				Ref:       "main",
				Status:    "success",
				Source:    "push",
				CreatedAt: secondCreatedAt,
				UpdatedAt: secondUpdatedAt,
				URI:       "https://gitlab.com/foobar/foobar/-/pipelines/114884002",
			}

			createRequest := []models.DeployRun{
				firstCreateRequest,
				secondCreateRequest,
			}

			IDs, err := DAO.CreateMany(createRequest)
			Expect(err).To(BeNil())
			Expect(len(*IDs)).To(Equal(2))
		})
	})

	var _ = When("Read", func() {
		It("reads a deploy run", func() {
			DAO, err := daos.NewDeployRun(&ctx)
			Expect(err).To((BeNil()))

			createdAt, _ := time.Parse(time.RFC3339, "2020-02-04T14:29:50.092Z")
			updatedAt, _ := time.Parse(time.RFC3339, "2020-02-04T14:35:51.459Z")
			createRequest := models.DeployRun{
				ProjectID: "15392086",
				Ref:       "main",
				Status:    "success",
				Source:    "push",
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
				URI:       "https://gitlab.com/foobar/foobar/-/pipelines/114883218",
			}

			ID, err := DAO.Create(&createRequest)
			Expect(err).To(BeNil())
			Expect(ID).To(Not(BeEmpty()))

			deployRun, err := DAO.Read(ID)
			Expect(err).To(BeNil())
			Expect(deployRun.ID).To(Equal(ID))
			Expect(deployRun.ProjectID).To(Equal(createRequest.ProjectID))
		})
	})

	var _ = When("ReadAll", func() {
		It("reads all applications", func() {
			DAO, err := daos.NewDeployRun(&ctx)
			Expect(err).To((BeNil()))

			firstCreatedAt, _ := time.Parse(time.RFC3339, "2020-02-04T14:29:50.092Z")
			firstUpdatedAt, _ := time.Parse(time.RFC3339, "2020-02-04T14:35:51.459Z")
			firstCreateRequest := models.DeployRun{
				ProjectID: "15392086",
				Ref:       "main",
				Status:    "success",
				Source:    "push",
				CreatedAt: firstCreatedAt,
				UpdatedAt: firstUpdatedAt,
				URI:       "https://gitlab.com/foobar/foobar/-/pipelines/114883218",
			}

			secondCreatedAt, _ := time.Parse(time.RFC3339, "2020-02-04T14:31:51.235Z")
			secondUpdatedAt, _ := time.Parse(time.RFC3339, "2020-02-04T14:38:18.252Z")
			secondCreateRequest := models.DeployRun{
				ProjectID: "15392086",
				Ref:       "main",
				Status:    "success",
				Source:    "push",
				CreatedAt: secondCreatedAt,
				UpdatedAt: secondUpdatedAt,
				URI:       "https://gitlab.com/foobar/foobar/-/pipelines/114884002",
			}

			createRequest := []models.DeployRun{
				firstCreateRequest,
				secondCreateRequest,
			}

			IDs, err := DAO.CreateMany(createRequest)
			Expect(err).To(BeNil())
			Expect(len(*IDs)).To(Equal(2))

			readAllResponse, err := DAO.ReadAll()
			Expect(err).To((BeNil()))
			Expect(len(*readAllResponse)).To(BeNumerically(">", 0))
		})
	})

	var _ = When("Update", func() {
		It("updates a deploy run", func() {
			DAO, err := daos.NewDeployRun(&ctx)
			Expect(err).To((BeNil()))

			createdAt, _ := time.Parse(time.RFC3339, "2020-02-04T14:29:50.092Z")
			updatedAt, _ := time.Parse(time.RFC3339, "2020-02-04T14:35:51.459Z")
			createRequest := models.DeployRun{
				Ref:       "master",
				Status:    "success",
				Source:    "push",
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
				URI:       "https://gitlab.com/foobar/foobar/-/pipelines/114883218",
			}

			ID, err := DAO.Create(&createRequest)
			Expect(err).To((BeNil()))

			refUpdate := "main"
			createdAt, _ = time.Parse(time.RFC3339, "2020-02-04T14:29:50.092Z")
			updatedAt, _ = time.Parse(time.RFC3339, "2020-02-04T14:35:51.459Z")
			updateRequest := models.DeployRun{
				ID:        ID,
				ProjectID: createRequest.ProjectID,
				Ref:       refUpdate,
				Status:    createRequest.Status,
				Source:    createRequest.Source,
				CreatedAt: createRequest.CreatedAt,
				UpdatedAt: createRequest.UpdatedAt,
				URI:       "https://gitlab.com/foobar/foobar/-/pipelines/114883218",
			}

			updateResponse, err := DAO.Update(&updateRequest)
			Expect(err).To(BeNil())
			Expect(updateResponse.Ref).To(Equal(refUpdate))
		})
	})

	var _ = When("Delete", func() {
		It("deletes a deploy run", func() {
			DAO, err := daos.NewDeployRun(&ctx)
			Expect(err).To((BeNil()))

			createdAt, _ := time.Parse(time.RFC3339, "2020-02-04T14:29:50.092Z")
			updatedAt, _ := time.Parse(time.RFC3339, "2020-02-04T14:35:51.459Z")
			createRequest := models.DeployRun{
				ProjectID: "15392086",
				Ref:       "master",
				Status:    "success",
				Source:    "push",
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
				URI:       "https://gitlab.com/foobar/foobar/-/pipelines/114883218",
			}

			ID, err := DAO.Create(&createRequest)
			Expect(err).To((BeNil()))

			deleteResponse, err := DAO.Delete(ID)
			Expect(err).To(BeNil())
			Expect(deleteResponse.ID).To(Equal(ID))
		})
	})
})
