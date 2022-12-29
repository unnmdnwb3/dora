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

var _ = Describe("daos.PipelineRunsPerDays", func() {
	ctx := context.Background()

	var _ = When("CreatePipelineRunsPerDay", func() {
		It("creates creates a new PipelineRunsPerDay.", func() {
			pipelineID := primitive.NewObjectID()
			date, _ := time.Parse(time.RFC3339, "2020-02-04T14:29:50.092Z")
			pipelineRunsPerDay := models.PipelineRunsPerDay{
				PipelineID:        pipelineID,
				Date:              date,
				TotalPipelineRuns: 1,
			}
			err := daos.CreatePipelineRunsPerDay(ctx, pipelineID, &pipelineRunsPerDay)
			Expect(err).To(BeNil())
			Expect(pipelineRunsPerDay.ID).To(Not(BeEmpty()))
		})
	})

	var _ = When("CreatePipelineRunsPerDays", func() {
		It("creates creates many new PipelineRunsPerDays.", func() {
			pipelineID := primitive.NewObjectID()
			date1, _ := time.Parse(time.RFC3339, "2020-02-04T00:00:00.000Z")
			pipelineRunsPerDay1 := models.PipelineRunsPerDay{
				PipelineID:        pipelineID,
				Date:              date1,
				TotalPipelineRuns: 1,
			}

			date2, _ := time.Parse(time.RFC3339, "2020-02-05T00:00:00.000Z")
			pipelineRunsPerDay2 := models.PipelineRunsPerDay{
				PipelineID:        pipelineID,
				Date:              date2,
				TotalPipelineRuns: 2,
			}
			pipelineRunsPerDays := []models.PipelineRunsPerDay{pipelineRunsPerDay1, pipelineRunsPerDay2}
			err := daos.CreatePipelineRunsPerDays(ctx, pipelineID, &pipelineRunsPerDays)
			Expect(err).To(BeNil())
			Expect(pipelineRunsPerDays[0].ID).To(Not(BeEmpty()))
			Expect(pipelineRunsPerDays[1].ID).To(Not(BeEmpty()))
		})
	})

	var _ = When("GetPipelineRunsPerDay", func() {
		It("retrieves an PipelineRunsPerDay.", func() {
			pipelineID := primitive.NewObjectID()
			date, _ := time.Parse(time.RFC3339, "2020-02-04T14:29:50.092Z")
			pipelineRunsPerDay := models.PipelineRunsPerDay{
				PipelineID:        pipelineID,
				Date:              date,
				TotalPipelineRuns: 1,
			}
			err := daos.CreatePipelineRunsPerDay(ctx, pipelineID, &pipelineRunsPerDay)
			Expect(err).To(BeNil())
			Expect(pipelineRunsPerDay.ID).To(Not(BeEmpty()))

			var findPipelineRunsPerDay models.PipelineRunsPerDay
			err = daos.GetPipelineRunsPerDay(ctx, pipelineRunsPerDay.ID, &findPipelineRunsPerDay)
			Expect(err).To(BeNil())
			Expect(findPipelineRunsPerDay.ID).To(Equal(pipelineRunsPerDay.ID))
		})
	})

	var _ = When("ListPipelineRunsPerDays", func() {
		It("retrieves many PipelineRunsPerDays.", func() {
			pipelineID := primitive.NewObjectID()
			date1, _ := time.Parse(time.RFC3339, "2020-02-04T00:00:00.000Z")
			pipelineRunsPerDay1 := models.PipelineRunsPerDay{
				PipelineID:        pipelineID,
				Date:              date1,
				TotalPipelineRuns: 1,
			}

			date2, _ := time.Parse(time.RFC3339, "2020-02-05T00:00:00.000Z")
			pipelineRunsPerDay2 := models.PipelineRunsPerDay{
				PipelineID:        pipelineID,
				Date:              date2,
				TotalPipelineRuns: 2,
			}
			pipelineRunsPerDays := []models.PipelineRunsPerDay{pipelineRunsPerDay1, pipelineRunsPerDay2}
			err := daos.CreatePipelineRunsPerDays(ctx, pipelineID, &pipelineRunsPerDays)
			Expect(err).To(BeNil())

			var findPipelineRunsPerDays []models.PipelineRunsPerDay
			err = daos.ListPipelineRunsPerDays(ctx, pipelineID, &findPipelineRunsPerDays)
			Expect(err).To(BeNil())
			Expect(len(findPipelineRunsPerDays)).To(Equal(2))
		})
	})

	var _ = When("ListPipelineRunsPerDaysByFilter", func() {
		It("retrieves many PipelineRunsPerDays conforming to a filter.", func() {
			pipelineID := primitive.NewObjectID()
			date1, _ := time.Parse(time.RFC3339, "2020-02-04T00:00:00.000Z")
			pipelineRunsPerDay1 := models.PipelineRunsPerDay{
				PipelineID:        pipelineID,
				Date:              date1,
				TotalPipelineRuns: 1,
			}

			date2, _ := time.Parse(time.RFC3339, "2020-02-05T00:00:00.000Z")
			pipelineRunsPerDay2 := models.PipelineRunsPerDay{
				PipelineID:        pipelineID,
				Date:              date2,
				TotalPipelineRuns: 2,
			}
			pipelineRunsPerDays := []models.PipelineRunsPerDay{pipelineRunsPerDay1, pipelineRunsPerDay2}
			err := daos.CreatePipelineRunsPerDays(ctx, pipelineID, &pipelineRunsPerDays)
			Expect(err).To(BeNil())

			var findPipelineRunsPerDays []models.PipelineRunsPerDay
			filter := bson.M{"date": bson.M{"$gte": date2}}
			err = daos.ListPipelineRunsPerDaysByFilter(ctx, filter, &findPipelineRunsPerDays)
			Expect(err).To(BeNil())
			Expect(findPipelineRunsPerDays).To(HaveLen(1))
		})
	})

	var _ = When("UpdatePipelineRunsPerDay", func() {
		It("updates an PipelineRunsPerDay.", func() {
			pipelineID := primitive.NewObjectID()
			date, _ := time.Parse(time.RFC3339, "2020-02-04T14:29:50.092Z")
			pipelineRunsPerDay := models.PipelineRunsPerDay{
				PipelineID:        pipelineID,
				Date:              date,
				TotalPipelineRuns: 1,
			}
			err := daos.CreatePipelineRunsPerDay(ctx, pipelineID, &pipelineRunsPerDay)
			Expect(err).To(BeNil())
			Expect(pipelineRunsPerDay.ID).To(Not(BeEmpty()))

			updatePipelineRunsPerDay := models.PipelineRunsPerDay{
				PipelineID:        pipelineID,
				Date:              date,
				TotalPipelineRuns: 2,
			}
			err = daos.UpdatePipelineRunsPerDay(ctx, pipelineRunsPerDay.ID, &updatePipelineRunsPerDay)
			Expect(err).To(BeNil())
			Expect(updatePipelineRunsPerDay.TotalPipelineRuns).To(Equal(2))
		})
	})

	var _ = When("DeleteOne", func() {
		It("deletes a document with ID in a collection", func() {
			pipelineID := primitive.NewObjectID()
			date, _ := time.Parse(time.RFC3339, "2020-02-04T14:29:50.092Z")
			pipelineRunsPerDay := models.PipelineRunsPerDay{
				PipelineID:        pipelineID,
				Date:              date,
				TotalPipelineRuns: 1,
			}
			err := daos.CreatePipelineRunsPerDay(ctx, pipelineID, &pipelineRunsPerDay)
			Expect(err).To(BeNil())
			Expect(pipelineRunsPerDay.ID).To(Not(BeEmpty()))

			err = daos.DeletePipelineRunsPerDay(ctx, pipelineRunsPerDay.ID)
			Expect(err).To(BeNil())

			var findPipelineRunsPerDay models.PipelineRunsPerDay
			err = daos.GetPipelineRunsPerDay(ctx, pipelineRunsPerDay.ID, &findPipelineRunsPerDay)
			Expect(err).To(Not(BeNil()))
		})
	})
})
