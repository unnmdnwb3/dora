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

var _ = Describe("services.trigger.aggregate.changes_per_days", func() {
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

	var _ = When("CalculateChangesPerDays", func() {
		It("calculates ChangesPerDays based on Changes.", func() {
			repositoryID := primitive.NewObjectID()
			pipelineID := primitive.NewObjectID()
			changes := []models.Change{
				{
					RepositoryID:    repositoryID,
					PipelineID:      pipelineID,
					FirstCommitDate: time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
					DeploymentDate:  time.Date(2022, 12, 27, 14, 16, 42, 0, time.UTC),
					LeadTime:        3600,
				},
				{
					RepositoryID:    repositoryID,
					PipelineID:      pipelineID,
					FirstCommitDate: time.Date(2022, 12, 27, 17, 39, 21, 0, time.UTC),
					DeploymentDate:  time.Date(2022, 12, 27, 17, 54, 21, 0, time.UTC),
					LeadTime:        900,
				},
				{
					RepositoryID:    repositoryID,
					PipelineID:      pipelineID,
					FirstCommitDate: time.Date(2022, 12, 29, 02, 21, 42, 0, time.UTC),
					DeploymentDate:  time.Date(2022, 12, 29, 02, 41, 42, 0, time.UTC),
					LeadTime:        1200,
				},
			}

			changesPerDays, err := aggregate.CalculateChangesPerDays(ctx, &changes)
			Expect(err).To(BeNil())
			Expect(len(*changesPerDays)).To(Equal(2))
			Expect((*changesPerDays)[0].TotalChanges).To(Equal(2))
			Expect((*changesPerDays)[0].TotalLeadTime).To(Equal(float64(4500)))
			Expect((*changesPerDays)[1].TotalChanges).To(Equal(1))
			Expect((*changesPerDays)[1].TotalLeadTime).To(Equal(float64(1200)))
		})
	})

	var _ = When("CreateChangesPerDays", func() {
		It("creates ChangesPerDays based on Changes.", func() {
			repositoryID := primitive.NewObjectID()
			pipelineID := primitive.NewObjectID()
			change1 := models.Change{
				RepositoryID:    repositoryID,
				PipelineID:      pipelineID,
				FirstCommitDate: time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
				DeploymentDate:  time.Date(2022, 12, 27, 14, 16, 42, 0, time.UTC),
			}

			change2 := models.Change{
				RepositoryID:    repositoryID,
				PipelineID:      pipelineID,
				FirstCommitDate: time.Date(2022, 12, 27, 17, 39, 21, 0, time.UTC),
				DeploymentDate:  time.Date(2022, 12, 27, 17, 54, 21, 0, time.UTC),
			}

			change3 := models.Change{
				RepositoryID:    repositoryID,
				PipelineID:      pipelineID,
				FirstCommitDate: time.Date(2022, 12, 29, 02, 21, 42, 0, time.UTC),
				DeploymentDate:  time.Date(2022, 12, 29, 02, 41, 42, 0, time.UTC),
			}

			changes := []models.Change{change1, change2, change3}
			err := daos.CreateChanges(ctx, repositoryID, &changes)
			Expect(err).To(BeNil())

			channel := make(chan error)
			defer close(channel)

			go aggregate.CreateChangesPerDays(ctx, channel, repositoryID, pipelineID)
			err = <-channel
			Expect(err).To(BeNil())

			var changesPerDays []models.ChangesPerDay
			err = daos.ListChangesPerDays(ctx, repositoryID, pipelineID, &changesPerDays)
			Expect(err).To(BeNil())
			Expect(len(changesPerDays)).To(Equal(2))
			Expect(changesPerDays[0].TotalChanges).To(Equal(2))
			Expect(changesPerDays[0].TotalLeadTime).To(Equal(float64(4500)))
			Expect(changesPerDays[1].TotalChanges).To(Equal(1))
			Expect(changesPerDays[1].TotalLeadTime).To(Equal(float64(1200)))
		})
	})
})
