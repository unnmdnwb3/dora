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

var _ = Describe("daos.Change", func() {
	ctx := context.Background()

	var _ = When("CreateChange", func() {
		It("creates a new Change.", func() {
			change := models.Change{
				PipelineID:      primitive.NewObjectID(),
				FirstCommitDate: time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
				DeploymentDate:  time.Date(2022, 12, 27, 13, 16, 21, 0, time.UTC),
			}
			err := daos.CreateChange(ctx, &change)
			Expect(err).To(BeNil())
			Expect(change.ID).To(Not(BeEmpty()))
		})
	})

	var _ = When("CreateIncidents", func() {
		It("creates creates many new Incidents.", func() {
			pipelineID := primitive.NewObjectID()
			change1 := models.Change{
				PipelineID:      pipelineID,
				FirstCommitDate: time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
				DeploymentDate:  time.Date(2022, 12, 27, 13, 21, 42, 0, time.UTC),
			}
			change2 := models.Change{
				PipelineID:      pipelineID,
				FirstCommitDate: time.Date(2022, 12, 27, 14, 51, 21, 0, time.UTC),
				DeploymentDate:  time.Date(2022, 12, 27, 14, 59, 34, 0, time.UTC),
			}
			changes := []models.Change{change1, change2}
			err := daos.CreateChanges(ctx, &changes)
			Expect(err).To(BeNil())
			Expect(changes[0].ID).To(Not(BeEmpty()))
			Expect(changes[1].ID).To(Not(BeEmpty()))
		})
	})

	var _ = When("GetChange", func() {
		It("retrieves an Change.", func() {
			change := models.Change{
				PipelineID:      primitive.NewObjectID(),
				FirstCommitDate: time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
				DeploymentDate:  time.Date(2022, 12, 27, 13, 16, 21, 0, time.UTC),
			}
			err := daos.CreateChange(ctx, &change)
			Expect(err).To(BeNil())
			Expect(change.ID).To(Not(BeEmpty()))

			var findChange models.Change
			err = daos.GetChange(ctx, change.ID, &findChange)
			Expect(err).To(BeNil())
			Expect(findChange.ID).To(Equal(change.ID))
		})
	})

	var _ = When("ListChanges", func() {
		It("retrieves many Changes.", func() {
			pipelineID := primitive.NewObjectID()
			change1 := models.Change{
				PipelineID:      pipelineID,
				FirstCommitDate: time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
				DeploymentDate:  time.Date(2022, 12, 27, 13, 21, 42, 0, time.UTC),
			}
			change2 := models.Change{
				PipelineID:      pipelineID,
				FirstCommitDate: time.Date(2022, 12, 27, 14, 51, 21, 0, time.UTC),
				DeploymentDate:  time.Date(2022, 12, 27, 14, 59, 34, 0, time.UTC),
			}
			change3 := models.Change{
				PipelineID:      pipelineID,
				FirstCommitDate: time.Date(2022, 12, 28, 21, 27, 40, 0, time.UTC),
				DeploymentDate:  time.Date(2022, 12, 28, 21, 45, 46, 0, time.UTC),
			}
			_ = daos.CreateChange(ctx, &change1)
			_ = daos.CreateChange(ctx, &change2)
			_ = daos.CreateChange(ctx, &change3)
			Expect(change1.ID).To(Not(BeNil()))
			Expect(change2.ID).To(Not(BeNil()))
			Expect(change3.ID).To(Not(BeNil()))

			var findChanges []models.Change
			err := daos.ListChanges(ctx, pipelineID, &findChanges)
			Expect(err).To(BeNil())
			Expect(findChanges).To(HaveLen(3))
		})
	})

	var _ = When("ListChangesByFilter", func() {
		It("retrieves many Changes conforming to a filter.", func() {
			pipelineID := primitive.NewObjectID()
			differentPipelineID := primitive.NewObjectID()
			changes := []models.Change{
				{
					PipelineID:      pipelineID,
					FirstCommitDate: time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
					DeploymentDate:  time.Date(2022, 12, 27, 13, 21, 42, 0, time.UTC),
				},
				{
					PipelineID:      pipelineID,
					FirstCommitDate: time.Date(2022, 12, 27, 14, 51, 21, 0, time.UTC),
					DeploymentDate:  time.Date(2022, 12, 27, 14, 59, 34, 0, time.UTC),
				},
				{
					PipelineID:      differentPipelineID,
					FirstCommitDate: time.Date(2022, 12, 28, 21, 27, 40, 0, time.UTC),
					DeploymentDate:  time.Date(2022, 12, 28, 21, 45, 46, 0, time.UTC),
				},
			}

			err := daos.CreateChanges(ctx, &changes)
			Expect(err).To(BeNil())

			Expect(changes[0].ID).To(Not(BeNil()))
			Expect(changes[1].ID).To(Not(BeNil()))
			Expect(changes[2].ID).To(Not(BeNil()))

			var findChanges []models.Change
			filter := bson.M{"pipeline_id": pipelineID}
			err = daos.ListChangesByFilter(ctx, filter, &findChanges)
			Expect(err).To(BeNil())
			Expect(findChanges).To(HaveLen(2))
		})
	})

	var _ = When("UpdateChange", func() {
		It("updates an Change.", func() {
			change := models.Change{
				PipelineID:      primitive.NewObjectID(),
				FirstCommitDate: time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
				DeploymentDate:  time.Date(2022, 12, 27, 13, 21, 42, 0, time.UTC),
			}
			err := daos.CreateChange(ctx, &change)
			Expect(err).To(BeNil())
			Expect(change.ID).To(Not(BeEmpty()))

			updateChange := models.Change{
				PipelineID:      primitive.NewObjectID(),
				FirstCommitDate: time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
				DeploymentDate:  time.Date(2022, 12, 27, 13, 26, 42, 0, time.UTC),
			}
			err = daos.UpdateChange(ctx, change.ID, &updateChange)
			Expect(err).To(BeNil())
			Expect(updateChange.DeploymentDate).To(Equal(time.Date(2022, 12, 27, 13, 26, 42, 0, time.UTC)))
		})
	})

	var _ = When("DeleteOne", func() {
		It("deletes a document with ID in a collection", func() {
			change := models.Change{
				PipelineID:      primitive.NewObjectID(),
				FirstCommitDate: time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
				DeploymentDate:  time.Date(2022, 12, 27, 13, 21, 42, 0, time.UTC),
			}
			err := daos.CreateChange(ctx, &change)
			Expect(err).To(BeNil())
			Expect(change.ID).To(Not(BeEmpty()))

			err = daos.DeleteChange(ctx, change.ID)
			Expect(err).To(BeNil())

			var findChange models.Change
			err = daos.GetChange(ctx, change.ID, &findChange)
			Expect(err).To(Not(BeNil()))
		})
	})
})
