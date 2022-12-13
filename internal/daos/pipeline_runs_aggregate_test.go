package daos_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/unnmdnwb3/dora/internal/daos"
	"github.com/unnmdnwb3/dora/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

var _ = Describe("daos.pipelineRunsAggregate", func() {
	ctx := context.Background()

	var _ = When("CreatePipelineRunsAggregate", func() {
		It("creates creates a new PipelineRunsAggregate.", func() {
			date, _ := time.Parse(time.RFC3339, "2020-02-04T14:29:50.092Z")
			pipelineRunsAggregate := models.PipelineRunsAggregate{
				PipelineID:        "15392086",
				Date:              date,
				TotalPipelineRuns: 1,
			}
			err := daos.CreatePipelineRunsAggregate(ctx, &pipelineRunsAggregate)
			Expect(err).To(BeNil())
			Expect(pipelineRunsAggregate.ID).To(Not(BeEmpty()))
		})
	})

	var _ = When("CreatePipelineRunsAggregates", func() {
		It("creates creates many new PipelineRunsAggregates.", func() {
			date1, _ := time.Parse(time.RFC3339, "2020-02-04T00:00:00.000Z")
			pipelineRunsAggregate1 := models.PipelineRunsAggregate{
				PipelineID:        "15392086",
				Date:              date1,
				TotalPipelineRuns: 1,
			}

			date2, _ := time.Parse(time.RFC3339, "2020-02-05T00:00:00.000Z")
			pipelineRunsAggregate2 := models.PipelineRunsAggregate{
				PipelineID:        "15392086",
				Date:              date2,
				TotalPipelineRuns: 2,
			}
			pipelineRunsAggregates := []models.PipelineRunsAggregate{pipelineRunsAggregate1, pipelineRunsAggregate2}
			err := daos.CreatePipelineRunsAggregates(ctx, &pipelineRunsAggregates)
			Expect(err).To(BeNil())
			Expect(pipelineRunsAggregates[0].ID).To(Not(BeEmpty()))
			Expect(pipelineRunsAggregates[1].ID).To(Not(BeEmpty()))
		})
	})

	var _ = When("GetPipelineRunsAggregate", func() {
		It("retrieves an PipelineRunsAggregate.", func() {
			date, _ := time.Parse(time.RFC3339, "2020-02-04T14:29:50.092Z")
			pipelineRunsAggregate := models.PipelineRunsAggregate{
				PipelineID:        "15392086",
				Date:              date,
				TotalPipelineRuns: 1,
			}
			err := daos.CreatePipelineRunsAggregate(ctx, &pipelineRunsAggregate)
			Expect(err).To(BeNil())
			Expect(pipelineRunsAggregate.ID).To(Not(BeEmpty()))

			var findPipelineRunsAggregate models.PipelineRunsAggregate
			err = daos.GetPipelineRunsAggregate(ctx, pipelineRunsAggregate.ID, &findPipelineRunsAggregate)
			Expect(err).To(BeNil())
			Expect(findPipelineRunsAggregate.ID).To(Equal(pipelineRunsAggregate.ID))
		})
	})

	var _ = When("ListPipelineRunsAggregates", func() {
		It("retrieves many PipelineRunsAggregates.", func() {
			date1, _ := time.Parse(time.RFC3339, "2020-02-04T00:00:00.000Z")
			pipelineRunsAggregate1 := models.PipelineRunsAggregate{
				PipelineID:        "15392086",
				Date:              date1,
				TotalPipelineRuns: 1,
			}

			date2, _ := time.Parse(time.RFC3339, "2020-02-05T00:00:00.000Z")
			pipelineRunsAggregate2 := models.PipelineRunsAggregate{
				PipelineID:        "15392086",
				Date:              date2,
				TotalPipelineRuns: 2,
			}
			_ = daos.CreatePipelineRunsAggregate(ctx, &pipelineRunsAggregate1)
			_ = daos.CreatePipelineRunsAggregate(ctx, &pipelineRunsAggregate2)
			Expect(pipelineRunsAggregate1.ID).To(Not(BeNil()))
			Expect(pipelineRunsAggregate2.ID).To(Not(BeNil()))

			var findPipelineRunsAggregates []models.PipelineRunsAggregate
			err := daos.ListPipelineRunsAggregates(ctx, pipelineRunsAggregate1.PipelineID, &findPipelineRunsAggregates)
			Expect(err).To(BeNil())
			Expect(findPipelineRunsAggregates).To(HaveLen(2))
		})
	})

	var _ = When("ListPipelineRunsAggregatesByFilter", func() {
		It("retrieves many PipelineRunsAggregates conforming to a filter.", func() {
			date1, _ := time.Parse(time.RFC3339, "2020-02-04T00:00:00.000Z")
			pipelineRunsAggregate1 := models.PipelineRunsAggregate{
				PipelineID:        "15392086",
				Date:              date1,
				TotalPipelineRuns: 1,
			}

			date2, _ := time.Parse(time.RFC3339, "2020-02-05T00:00:00.000Z")
			pipelineRunsAggregate2 := models.PipelineRunsAggregate{
				PipelineID:        "15392086",
				Date:              date2,
				TotalPipelineRuns: 2,
			}
			_ = daos.CreatePipelineRunsAggregate(ctx, &pipelineRunsAggregate1)
			_ = daos.CreatePipelineRunsAggregate(ctx, &pipelineRunsAggregate2)
			Expect(pipelineRunsAggregate1.ID).To(Not(BeNil()))
			Expect(pipelineRunsAggregate2.ID).To(Not(BeNil()))

			var findPipelineRunsAggregates []models.PipelineRunsAggregate
			filter := bson.M{"date": bson.M{"$gte": date2}}
			err := daos.ListPipelineRunsAggregatesByFilter(ctx, filter, &findPipelineRunsAggregates)
			Expect(err).To(BeNil())
			Expect(findPipelineRunsAggregates).To(HaveLen(1))
		})
	})

	var _ = When("UpdatePipelineRunsAggregate", func() {
		It("updates an PipelineRunsAggregate.", func() {
			date, _ := time.Parse(time.RFC3339, "2020-02-04T14:29:50.092Z")
			pipelineRunsAggregate := models.PipelineRunsAggregate{
				PipelineID:        "15392086",
				Date:              date,
				TotalPipelineRuns: 1,
			}
			err := daos.CreatePipelineRunsAggregate(ctx, &pipelineRunsAggregate)
			Expect(err).To(BeNil())
			Expect(pipelineRunsAggregate.ID).To(Not(BeEmpty()))

			updatePipelineRunsAggregate := models.PipelineRunsAggregate{
				PipelineID:        "15392086",
				Date:              date,
				TotalPipelineRuns: 2,
			}
			err = daos.UpdatePipelineRunsAggregate(ctx, pipelineRunsAggregate.ID, &updatePipelineRunsAggregate)
			Expect(err).To(BeNil())
			Expect(updatePipelineRunsAggregate.TotalPipelineRuns).To(Equal(2))
		})
	})

	var _ = When("DeleteOne", func() {
		It("deletes a document with ID in a collection", func() {
			date, _ := time.Parse(time.RFC3339, "2020-02-04T14:29:50.092Z")
			pipelineRunsAggregate := models.PipelineRunsAggregate{
				PipelineID:        "15392086",
				Date:              date,
				TotalPipelineRuns: 1,
			}
			err := daos.CreatePipelineRunsAggregate(ctx, &pipelineRunsAggregate)
			Expect(err).To(BeNil())
			Expect(pipelineRunsAggregate.ID).To(Not(BeEmpty()))

			err = daos.DeletePipelineRunsAggregate(ctx, pipelineRunsAggregate.ID)
			Expect(err).To(BeNil())

			var findPipelineRunsAggregate models.PipelineRunsAggregate
			err = daos.GetPipelineRunsAggregate(ctx, pipelineRunsAggregate.ID, &findPipelineRunsAggregate)
			Expect(err).To(Not(BeNil()))
		})
	})
})
