package daos_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/unnmdnwb3/dora/internal/daos"
	"github.com/unnmdnwb3/dora/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ = Describe("daos.pipelineRun", func() {
	ctx := context.Background()
	externalID := 713437228

	var _ = When("CreatePipelineRun", func() {
		It("creates creates a new PipelineRun.", func() {
			pipelineID := primitive.NewObjectID()
			createdAt, _ := time.Parse(time.RFC3339, "2020-02-04T14:29:50.092Z")
			updatedAt, _ := time.Parse(time.RFC3339, "2020-02-04T14:35:51.459Z")
			pipelineRun := models.PipelineRun{
				PipelineID:  pipelineID,
				ExternalID:  externalID,
				Ref:         "main",
				Status:      "success",
				EventSource: "push",
				CreatedAt:   createdAt,
				UpdatedAt:   updatedAt,
				URI:         "https://gitlab.com/foobar/foobar/-/pipelines/114884002",
			}
			err := daos.CreatePipelineRun(ctx, pipelineID, &pipelineRun)
			Expect(err).To(BeNil())
			Expect(pipelineRun.ID).To(Not(BeEmpty()))
		})
	})

	var _ = When("CreatePipelineRuns", func() {
		It("creates creates many new PipelineRuns.", func() {
			pipelineID := primitive.NewObjectID()
			createdAt1, _ := time.Parse(time.RFC3339, "2020-02-04T14:29:50.092Z")
			updatedAt1, _ := time.Parse(time.RFC3339, "2020-02-04T14:35:51.459Z")
			pipelineRun1 := models.PipelineRun{
				PipelineID:  pipelineID,
				ExternalID:  externalID,
				Ref:         "main",
				Status:      "success",
				EventSource: "push",
				CreatedAt:   createdAt1,
				UpdatedAt:   updatedAt1,
				URI:         "https://gitlab.com/foobar/foobar/-/pipelines/114884002",
			}

			createdAt2, _ := time.Parse(time.RFC3339, "2020-02-04T14:39:50.092Z")
			updatedAt2, _ := time.Parse(time.RFC3339, "2020-02-04T14:45:51.459Z")
			pipelineRun2 := models.PipelineRun{
				PipelineID:  pipelineID,
				ExternalID:  externalID,
				Ref:         "main",
				Status:      "success",
				EventSource: "push",
				CreatedAt:   createdAt2,
				UpdatedAt:   updatedAt2,
				URI:         "https://gitlab.com/foobar/foobar/-/pipelines/114884003",
			}
			pipelineRuns := []models.PipelineRun{pipelineRun1, pipelineRun2}
			err := daos.CreatePipelineRuns(ctx, pipelineID, &pipelineRuns)
			Expect(err).To(BeNil())
			Expect(pipelineRuns[0].ID).To(Not(BeEmpty()))
			Expect(pipelineRuns[1].ID).To(Not(BeEmpty()))
		})
	})

	var _ = When("GetPipelineRun", func() {
		It("retrieves an PipelineRun.", func() {
			pipelineID := primitive.NewObjectID()
			createdAt, _ := time.Parse(time.RFC3339, "2020-02-04T14:29:50.092Z")
			updatedAt, _ := time.Parse(time.RFC3339, "2020-02-04T14:35:51.459Z")
			pipelineRun := models.PipelineRun{
				PipelineID:  pipelineID,
				ExternalID:  externalID,
				Ref:         "main",
				Status:      "success",
				EventSource: "push",
				CreatedAt:   createdAt,
				UpdatedAt:   updatedAt,
				URI:         "https://gitlab.com/foobar/foobar/-/pipelines/114884002",
			}
			err := daos.CreatePipelineRun(ctx, pipelineID, &pipelineRun)
			Expect(err).To(BeNil())
			Expect(pipelineRun.ID).To(Not(BeEmpty()))

			var findPipelineRun models.PipelineRun
			err = daos.GetPipelineRun(ctx, pipelineRun.ID, &findPipelineRun)
			Expect(err).To(BeNil())
			Expect(findPipelineRun.ID).To(Equal(pipelineRun.ID))
		})
	})

	var _ = When("ListPipelineRuns", func() {
		It("retrieves many PipelineRuns.", func() {
			pipelineID := primitive.NewObjectID()
			createdAt1, _ := time.Parse(time.RFC3339, "2020-02-04T14:29:50.092Z")
			updatedAt1, _ := time.Parse(time.RFC3339, "2020-02-04T14:35:51.459Z")
			pipelineRun1 := models.PipelineRun{
				PipelineID:  pipelineID,
				ExternalID:  externalID,
				Ref:         "main",
				Status:      "success",
				EventSource: "push",
				CreatedAt:   createdAt1,
				UpdatedAt:   updatedAt1,
				URI:         "https://gitlab.com/foobar/foobar/-/pipelines/114884002",
			}

			createdAt2, _ := time.Parse(time.RFC3339, "2020-02-04T14:39:50.092Z")
			updatedAt2, _ := time.Parse(time.RFC3339, "2020-02-04T14:45:51.459Z")
			pipelineRun2 := models.PipelineRun{
				PipelineID:  pipelineID,
				ExternalID:  externalID,
				Ref:         "main",
				Status:      "success",
				EventSource: "push",
				CreatedAt:   createdAt2,
				UpdatedAt:   updatedAt2,
				URI:         "https://gitlab.com/foobar/foobar/-/pipelines/114884003",
			}
			pipelineRuns := []models.PipelineRun{pipelineRun1, pipelineRun2}
			err := daos.CreatePipelineRuns(ctx, pipelineID, &pipelineRuns)
			Expect(err).To(BeNil())

			var findPipelineRuns []models.PipelineRun
			err = daos.ListPipelineRuns(ctx, pipelineID, &findPipelineRuns)
			Expect(err).To(BeNil())
			Expect(len(findPipelineRuns)).To(Equal(2))
		})
	})

	var _ = When("ListPipelineRunsByFilter", func() {
		It("retrieves many PipelineRuns conforming to a filter.", func() {
			pipelineID := primitive.NewObjectID()
			createdAt1, _ := time.Parse(time.RFC3339, "2020-02-04T14:29:50.092Z")
			updatedAt1, _ := time.Parse(time.RFC3339, "2020-02-04T14:35:51.459Z")
			pipelineRun1 := models.PipelineRun{
				PipelineID:  pipelineID,
				ExternalID:  externalID,
				Ref:         "main",
				Status:      "success",
				EventSource: "push",
				CreatedAt:   createdAt1,
				UpdatedAt:   updatedAt1,
				URI:         "https://gitlab.com/foobar/foobar/-/pipelines/114884002",
			}

			createdAt2, _ := time.Parse(time.RFC3339, "2020-02-04T14:39:50.092Z")
			updatedAt2, _ := time.Parse(time.RFC3339, "2020-02-04T14:45:51.459Z")
			pipelineRun2 := models.PipelineRun{
				PipelineID:  pipelineID,
				ExternalID:  externalID,
				Ref:         "main",
				Status:      "success",
				EventSource: "push",
				CreatedAt:   createdAt2,
				UpdatedAt:   updatedAt2,
				URI:         "https://gitlab.com/foobar/foobar/-/pipelines/114884003",
			}
			pipelineRuns := []models.PipelineRun{pipelineRun1, pipelineRun2}
			err := daos.CreatePipelineRuns(ctx, pipelineID, &pipelineRuns)
			Expect(err).To(BeNil())

			var findPipelineRuns []models.PipelineRun
			filter := bson.M{"uri": pipelineRun1.URI}
			err = daos.ListPipelineRunsByFilter(ctx, filter, &findPipelineRuns)
			Expect(err).To(BeNil())
			Expect(findPipelineRuns).To(HaveLen(1))
		})
	})

	var _ = When("UpdatePipelineRun", func() {
		It("updates an PipelineRun.", func() {
			pipelineID := primitive.NewObjectID()
			createdAt, _ := time.Parse(time.RFC3339, "2020-02-04T14:29:50.092Z")
			updatedAt, _ := time.Parse(time.RFC3339, "2020-02-04T14:35:51.459Z")
			pipelineRun := models.PipelineRun{
				PipelineID:  pipelineID,
				ExternalID:  externalID,
				Ref:         "main",
				Status:      "success",
				EventSource: "push",
				CreatedAt:   createdAt,
				UpdatedAt:   updatedAt,
				URI:         "https://gitlab.com/foobar/foobar/-/pipelines/114884002",
			}
			err := daos.CreatePipelineRun(ctx, pipelineID, &pipelineRun)
			Expect(err).To(BeNil())
			Expect(pipelineRun.ID).To(Not(BeEmpty()))

			newUpdatedAt, _ := time.Parse(time.RFC3339, "2020-02-04T14:45:51.459Z")
			updatePipelineRun := models.PipelineRun{
				PipelineID:  pipelineID,
				ExternalID:  externalID,
				Ref:         "main",
				Status:      "success",
				EventSource: "push",
				CreatedAt:   createdAt,
				UpdatedAt:   newUpdatedAt,
				URI:         "https://gitlab.com/foobar/foobar/-/pipelines/114884002",
			}
			err = daos.UpdatePipelineRun(ctx, pipelineRun.ID, &updatePipelineRun)
			Expect(err).To(BeNil())
			Expect(updatePipelineRun.UpdatedAt).To(Equal(newUpdatedAt))
		})
	})

	var _ = When("DeleteOne", func() {
		It("deletes a document with ID in a collection", func() {
			pipelineID := primitive.NewObjectID()
			createdAt, _ := time.Parse(time.RFC3339, "2020-02-04T14:29:50.092Z")
			updatedAt, _ := time.Parse(time.RFC3339, "2020-02-04T14:35:51.459Z")
			pipelineRun := models.PipelineRun{
				PipelineID:  pipelineID,
				ExternalID:  externalID,
				Ref:         "main",
				Status:      "success",
				EventSource: "push",
				CreatedAt:   createdAt,
				UpdatedAt:   updatedAt,
				URI:         "https://gitlab.com/foobar/foobar/-/pipelines/114884002",
			}
			err := daos.CreatePipelineRun(ctx, pipelineID, &pipelineRun)
			Expect(err).To(BeNil())
			Expect(pipelineRun.ID).To(Not(BeEmpty()))

			err = daos.DeletePipelineRun(ctx, pipelineRun.ID)
			Expect(err).To(BeNil())

			var findPipelineRun models.PipelineRun
			err = daos.GetPipelineRun(ctx, pipelineRun.ID, &findPipelineRun)
			Expect(err).To(Not(BeNil()))
		})
	})
})
