package aggregate_test

import (
	"context"
	"os"
	"time"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/unnmdnwb3/dora/internal/daos"
	"github.com/unnmdnwb3/dora/internal/database/mongodb"
	"github.com/unnmdnwb3/dora/internal/models"
	"github.com/unnmdnwb3/dora/internal/services/trigger/aggregate"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ = Describe("services.trigger.aggregate", func() {
	ctx := context.Background()

	var _ = BeforeEach(func() {
		_ = godotenv.Load("./../../../../test/.env")
	})

	var _ = AfterEach(func() {
		service := mongodb.NewService()
		service.Connect(ctx, os.Getenv("MONGODB_DATABASE"))
		service.DB.Drop(ctx)
		defer service.Disconnect(ctx)

		os.Remove("MONGODB_URI")
		os.Remove("MONGODB_PORT")
		os.Remove("MONGODB_USER")
		os.Remove("MONGODB_PASSWORD")
	})

	var _ = When("CalculatePipelineRunsPerDays", func() {
		It("calculates the pipeline runs per day.", func() {
			pipelineID := primitive.NewObjectID()
			pipelineRuns := []models.PipelineRun{
				{
					PipelineID:  pipelineID,
					ExternalID:  713437220,
					Sha:         "1cfffa2ae16528e36115ece8b1f2601bcf74414e",
					Ref:         "main",
					Status:      "success",
					EventSource: "push",
					CreatedAt:   time.Date(2019, 10, 9, 9, 11, 20, 0, time.UTC),
					UpdatedAt:   time.Date(2019, 10, 9, 9, 12, 20, 0, time.UTC),
					URI:         "https://gitlab.com/foobar/foobar/-/pipelines/114884002",
				},
				{
					PipelineID:  pipelineID,
					ExternalID:  713437221,
					Sha:         "345207c839e94a939aebdc86835ae2e2a6c85acb",
					Ref:         "main",
					Status:      "success",
					EventSource: "push",
					CreatedAt:   time.Date(2019, 10, 11, 9, 11, 20, 0, time.UTC),
					UpdatedAt:   time.Date(2019, 10, 11, 9, 12, 20, 0, time.UTC),
					URI:         "https://gitlab.com/foobar/foobar/-/pipelines/114884003",
				},
				{
					PipelineID:  pipelineID,
					ExternalID:  713437222,
					Sha:         "dcc7ef44dc6a376854c5f2cc42b0b24aa3a9ed10",
					Ref:         "main",
					Status:      "success",
					EventSource: "push",
					CreatedAt:   time.Date(2019, 10, 11, 9, 13, 20, 0, time.UTC),
					UpdatedAt:   time.Date(2019, 10, 11, 9, 14, 20, 0, time.UTC),
					URI:         "https://gitlab.com/foobar/foobar/-/pipelines/114884004",
				},
			}

			pipelineRunsPerDay, err := aggregate.CalculatePipelineRunsPerDays(ctx, &pipelineRuns)
			Expect(err).To(BeNil())
			Expect(len(*pipelineRunsPerDay)).To(Equal(2))
			Expect((*pipelineRunsPerDay)[0].TotalPipelineRuns).To(Equal(1))
			Expect((*pipelineRunsPerDay)[1].TotalPipelineRuns).To(Equal(2))
		})
	})

	var _ = When("CreatePipelineRunsPerDays", func() {
		It("calculates and creates the pipeline runs for each day.", func() {
			pipelineID := primitive.NewObjectID()
			pipelineRuns := []models.PipelineRun{
				{
					PipelineID:  pipelineID,
					ExternalID:  713437220,
					Sha:         "1cfffa2ae16528e36115ece8b1f2601bcf74414e",
					Ref:         "main",
					Status:      "success",
					EventSource: "push",
					CreatedAt:   time.Date(2019, 10, 9, 9, 11, 20, 0, time.UTC),
					UpdatedAt:   time.Date(2019, 10, 9, 9, 12, 20, 0, time.UTC),
					URI:         "https://gitlab.com/foobar/foobar/-/pipelines/114884002",
				},
				{
					PipelineID:  pipelineID,
					ExternalID:  713437221,
					Sha:         "345207c839e94a939aebdc86835ae2e2a6c85acb",
					Ref:         "main",
					Status:      "success",
					EventSource: "push",
					CreatedAt:   time.Date(2019, 10, 11, 9, 11, 20, 0, time.UTC),
					UpdatedAt:   time.Date(2019, 10, 11, 9, 12, 20, 0, time.UTC),
					URI:         "https://gitlab.com/foobar/foobar/-/pipelines/114884003",
				},
				{
					PipelineID:  pipelineID,
					ExternalID:  713437222,
					Sha:         "dcc7ef44dc6a376854c5f2cc42b0b24aa3a9ed10",
					Ref:         "main",
					Status:      "success",
					EventSource: "push",
					CreatedAt:   time.Date(2019, 10, 11, 9, 13, 20, 0, time.UTC),
					UpdatedAt:   time.Date(2019, 10, 11, 9, 14, 20, 0, time.UTC),
					URI:         "https://gitlab.com/foobar/foobar/-/pipelines/114884004",
				},
			}

			err := daos.CreatePipelineRuns(ctx, pipelineID, &pipelineRuns)
			Expect(err).To(BeNil())

			channel := make(chan error)
			defer close(channel)

			go aggregate.CreatePipelineRunsPerDays(ctx, channel, pipelineID)
			err = <-channel
			Expect(err).To(BeNil())

			var pipelineRunsPerDays []models.PipelineRunsPerDay
			err = daos.ListPipelineRunsPerDays(ctx, pipelineID, &pipelineRunsPerDays)
			Expect(err).To(BeNil())
			Expect(pipelineRunsPerDays).To(HaveLen(2))
		})
	})
})
