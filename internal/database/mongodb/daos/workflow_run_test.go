package daos_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/unnmdnwb3/dora/internal/database/mongodb/daos"
	"github.com/unnmdnwb3/dora/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

var _ = Describe("daos.WorkflowRun", func() {
	ctx := context.Background()

	var _ = When("Create", func() {
		It("creates a workflow run", func() {
			DAO, err := daos.NewWorkflowRun(&ctx)
			Expect(err).To((BeNil()))

			createdAt, _ := time.Parse(time.RFC3339, "2020-02-04T14:29:50.092Z")
			updatedAt, _ := time.Parse(time.RFC3339, "2020-02-04T14:35:51.459Z")
			createRequest := models.WorkflowRun{
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
		It("creates many workflow runs", func() {
			DAO, err := daos.NewWorkflowRun(&ctx)
			Expect(err).To((BeNil()))

			firstCreatedAt, _ := time.Parse(time.RFC3339, "2020-02-04T14:29:50.092Z")
			firstUpdatedAt, _ := time.Parse(time.RFC3339, "2020-02-04T14:35:51.459Z")
			firstCreateRequest := models.WorkflowRun{
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
			secondCreateRequest := models.WorkflowRun{
				ProjectID: "15392086",
				Ref:       "main",
				Status:    "success",
				Source:    "push",
				CreatedAt: secondCreatedAt,
				UpdatedAt: secondUpdatedAt,
				URI:       "https://gitlab.com/foobar/foobar/-/pipelines/114884002",
			}

			createRequest := []models.WorkflowRun{
				firstCreateRequest,
				secondCreateRequest,
			}

			IDs, err := DAO.CreateMany(createRequest)
			Expect(err).To(BeNil())
			Expect(len(*IDs)).To(Equal(2))
		})
	})

	var _ = When("Read", func() {
		It("reads a workflow run", func() {
			DAO, err := daos.NewWorkflowRun(&ctx)
			Expect(err).To((BeNil()))

			createdAt, _ := time.Parse(time.RFC3339, "2020-02-04T14:29:50.092Z")
			updatedAt, _ := time.Parse(time.RFC3339, "2020-02-04T14:35:51.459Z")
			createRequest := models.WorkflowRun{
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

			workflowRun, err := DAO.Read(ID)
			Expect(err).To(BeNil())
			Expect(workflowRun.ID).To(Equal(ID))
			Expect(workflowRun.ProjectID).To(Equal(createRequest.ProjectID))
		})
	})

	var _ = When("ReadAll", func() {
		It("reads all workflow runs", func() {
			DAO, err := daos.NewWorkflowRun(&ctx)
			Expect(err).To((BeNil()))

			firstCreatedAt, _ := time.Parse(time.RFC3339, "2020-02-04T14:29:50.092Z")
			firstUpdatedAt, _ := time.Parse(time.RFC3339, "2020-02-04T14:35:51.459Z")
			firstCreateRequest := models.WorkflowRun{
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
			secondCreateRequest := models.WorkflowRun{
				ProjectID: "15392086",
				Ref:       "main",
				Status:    "success",
				Source:    "push",
				CreatedAt: secondCreatedAt,
				UpdatedAt: secondUpdatedAt,
				URI:       "https://gitlab.com/foobar/foobar/-/pipelines/114884002",
			}

			createRequest := []models.WorkflowRun{
				firstCreateRequest,
				secondCreateRequest,
			}

			IDs, err := DAO.CreateMany(createRequest)
			Expect(err).To(BeNil())
			Expect(len(*IDs)).To(Equal(2))

			filter := bson.M{}
			readAllResponse, err := DAO.ReadAll(filter)
			Expect(err).To((BeNil()))
			Expect(len(*readAllResponse)).To(BeNumerically(">", 0))
		})
	})

	var _ = When("Update", func() {
		It("updates a workflow run", func() {
			DAO, err := daos.NewWorkflowRun(&ctx)
			Expect(err).To((BeNil()))

			createdAt, _ := time.Parse(time.RFC3339, "2020-02-04T14:29:50.092Z")
			updatedAt, _ := time.Parse(time.RFC3339, "2020-02-04T14:35:51.459Z")
			createRequest := models.WorkflowRun{
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
			updateRequest := models.WorkflowRun{
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
		It("deletes a workflow run", func() {
			DAO, err := daos.NewWorkflowRun(&ctx)
			Expect(err).To((BeNil()))

			createdAt, _ := time.Parse(time.RFC3339, "2020-02-04T14:29:50.092Z")
			updatedAt, _ := time.Parse(time.RFC3339, "2020-02-04T14:35:51.459Z")
			createRequest := models.WorkflowRun{
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