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

var _ = Describe("daos.ChangesPerDays", func() {
	ctx := context.Background()

	var _ = When("CreateChangesPerDay", func() {
		It("creates creates a new pipelineRunsPerDay.", func() {
			repositoryID := primitive.NewObjectID()
			changesPerDay := models.ChangesPerDay{
				RepositoryID:  repositoryID,
				Date:          time.Date(2022, 12, 27, 0, 0, 0, 0, time.UTC),
				TotalChanges:  1,
				TotalLeadTime: 300,
			}
			err := daos.CreateChangesPerDay(ctx, repositoryID, &changesPerDay)
			Expect(err).To(BeNil())
			Expect(changesPerDay.ID).To(Not(BeEmpty()))
		})
	})

	var _ = When("CreateChangesPerDays", func() {
		It("creates creates many new ChangesPerDays.", func() {
			repositoryID := primitive.NewObjectID()
			changesPerDays := []models.ChangesPerDay{
				{
					RepositoryID:  repositoryID,
					Date:          time.Date(2022, 12, 27, 0, 0, 0, 0, time.UTC),
					TotalChanges:  1,
					TotalLeadTime: 300,
				},

				{
					RepositoryID:  repositoryID,
					Date:          time.Date(2022, 12, 28, 0, 0, 0, 0, time.UTC),
					TotalChanges:  2,
					TotalLeadTime: 450,
				},
			}
			err := daos.CreateChangesPerDays(ctx, repositoryID, &changesPerDays)
			Expect(err).To(BeNil())
			Expect(changesPerDays[0].ID).To(Not(BeEmpty()))
			Expect(changesPerDays[1].ID).To(Not(BeEmpty()))
		})
	})

	var _ = When("GetChangesPerDay", func() {
		It("retrieves an ChangesPerDay.", func() {
			repositoryID := primitive.NewObjectID()
			changesPerDay := models.ChangesPerDay{
				RepositoryID:  repositoryID,
				Date:          time.Date(2022, 12, 27, 0, 0, 0, 0, time.UTC),
				TotalChanges:  1,
				TotalLeadTime: 300,
			}
			err := daos.CreateChangesPerDay(ctx, repositoryID, &changesPerDay)
			Expect(err).To(BeNil())
			Expect(changesPerDay.ID).To(Not(BeEmpty()))

			var findChangesPerDay models.ChangesPerDay
			err = daos.GetChangesPerDay(ctx, changesPerDay.ID, &findChangesPerDay)
			Expect(err).To(BeNil())
			Expect(findChangesPerDay.ID).To(Equal(changesPerDay.ID))
		})
	})

	var _ = When("ListChangesPerDays", func() {
		It("retrieves many ChangesPerDays.", func() {
			repositoryID := primitive.NewObjectID()
			changesPerDays := []models.ChangesPerDay{
				{
					RepositoryID:  repositoryID,
					Date:          time.Date(2022, 12, 27, 0, 0, 0, 0, time.UTC),
					TotalChanges:  1,
					TotalLeadTime: 300,
				},

				{
					RepositoryID:  repositoryID,
					Date:          time.Date(2022, 12, 28, 0, 0, 0, 0, time.UTC),
					TotalChanges:  2,
					TotalLeadTime: 450,
				},
			}
			err := daos.CreateChangesPerDays(ctx, repositoryID, &changesPerDays)
			Expect(err).To(BeNil())

			var findChangesPerDays []models.ChangesPerDay
			err = daos.ListChangesPerDays(ctx, repositoryID, &findChangesPerDays)
			Expect(err).To(BeNil())
			Expect(len(findChangesPerDays)).To(Equal(2))
		})
	})

	var _ = When("ListChangesPerDaysByFilter", func() {
		It("retrieves many ChangesPerDays conforming to a filter.", func() {
			repositoryID := primitive.NewObjectID()
			changesPerDays := []models.ChangesPerDay{
				{
					RepositoryID:  repositoryID,
					Date:          time.Date(2022, 12, 27, 0, 0, 0, 0, time.UTC),
					TotalChanges:  1,
					TotalLeadTime: 300,
				},

				{
					RepositoryID:  repositoryID,
					Date:          time.Date(2022, 12, 28, 0, 0, 0, 0, time.UTC),
					TotalChanges:  2,
					TotalLeadTime: 450,
				},
			}
			err := daos.CreateChangesPerDays(ctx, repositoryID, &changesPerDays)
			Expect(err).To(BeNil())

			var findChangesPerDays []models.ChangesPerDay
			date := time.Date(2022, 12, 27, 0, 0, 1, 0, time.UTC)
			filter := bson.M{"date": bson.M{"$gte": date}}
			err = daos.ListChangesPerDaysByFilter(ctx, filter, &findChangesPerDays)
			Expect(err).To(BeNil())
			Expect(findChangesPerDays).To(HaveLen(1))
		})
	})

	var _ = When("UpdateChangesPerDay", func() {
		It("updates an ChangesPerDay.", func() {
			repositoryID := primitive.NewObjectID()
			changesPerDay := models.ChangesPerDay{
				RepositoryID:  repositoryID,
				Date:          time.Date(2022, 12, 27, 0, 0, 0, 0, time.UTC),
				TotalChanges:  1,
				TotalLeadTime: 300,
			}
			err := daos.CreateChangesPerDay(ctx, repositoryID, &changesPerDay)
			Expect(err).To(BeNil())
			Expect(changesPerDay.ID).To(Not(BeEmpty()))

			updateChangesPerDay := models.ChangesPerDay{
				RepositoryID:  repositoryID,
				Date:          time.Date(2022, 12, 27, 0, 0, 0, 0, time.UTC),
				TotalChanges:  2,
				TotalLeadTime: 600,
			}
			err = daos.UpdateChangesPerDay(ctx, changesPerDay.ID, &updateChangesPerDay)
			Expect(err).To(BeNil())
			Expect(updateChangesPerDay.TotalChanges).To(Equal(2))
		})
	})

	var _ = When("DeleteOne", func() {
		It("deletes a document with ID in a collection", func() {
			repositoryID := primitive.NewObjectID()
			changesPerDay := models.ChangesPerDay{
				RepositoryID:  repositoryID,
				Date:          time.Date(2022, 12, 27, 0, 0, 0, 0, time.UTC),
				TotalChanges:  1,
				TotalLeadTime: 300,
			}
			err := daos.CreateChangesPerDay(ctx, repositoryID, &changesPerDay)
			Expect(err).To(BeNil())
			Expect(changesPerDay.ID).To(Not(BeEmpty()))

			err = daos.DeleteChangesPerDay(ctx, changesPerDay.ID)
			Expect(err).To(BeNil())

			var findChangesPerDay models.ChangesPerDay
			err = daos.GetChangesPerDay(ctx, changesPerDay.ID, &findChangesPerDay)
			Expect(err).To(Not(BeNil()))
		})
	})
})
