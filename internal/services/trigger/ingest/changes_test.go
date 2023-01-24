package ingest_test

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
	"github.com/unnmdnwb3/dora/internal/services/trigger/ingest"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ = Describe("services.trigger.ingest.changes", func() {
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

	var _ = When("GetFirstCommits", func() {
		It("gets the first commits of a change.", func() {
			pipelineID := primitive.NewObjectID()
			pipelineRuns := []models.PipelineRun{
				{
					PipelineID:  pipelineID,
					ExternalID:  713437228,
					Sha:         "3d95fe3bf954501d3832e50fdd803c5f9eae3f94",
					Ref:         "main",
					Status:      "success",
					EventSource: "push",
					CreatedAt:   time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
					UpdatedAt:   time.Date(2022, 12, 27, 13, 21, 42, 0, time.UTC),
					URI:         "https://gitlab.com/foobar/foobar/-/pipelines/713437228",
				},
				{
					PipelineID:  pipelineID,
					ExternalID:  713437229,
					Sha:         "1db209656ad1ab0e14aaa4e2fe79b6caf8b2a9e7",
					Ref:         "main",
					Status:      "success",
					EventSource: "push",
					CreatedAt:   time.Date(2022, 12, 28, 15, 37, 28, 0, time.UTC),
					UpdatedAt:   time.Date(2022, 12, 28, 15, 43, 17, 0, time.UTC),
					URI:         "https://gitlab.com/foobar/foobar/-/pipelines/713437229",
				},
			}

			err := daos.CreatePipelineRuns(ctx, pipelineID, &pipelineRuns)
			Expect(err).To(BeNil())

			repositoryID := primitive.NewObjectID()
			// commits := []models.Commit{
			// 	{
			// 		RepositoryID: repositoryID,
			// 		CreatedAt:    time.Date(2022, 12, 28, 13, 01, 11, 0, time.UTC),
			// 		Sha:          "1db209656ad1ab0e14aaa4e2fe79b6caf8b2a9e7",
			// 		ParentShas: []string{
			// 			"3d95fe3bf954501d3832e50fdd803c5f9eae3f94",
			// 			"487d6aedb92ab76bdc03957aceece75db906796e",
			// 		},
			// 	},
			// 	{
			// 		RepositoryID: repositoryID,
			// 		CreatedAt:    time.Date(2022, 12, 28, 12, 46, 21, 0, time.UTC),
			// 		Sha:          "487d6aedb92ab76bdc03957aceece75db906796e",
			// 		ParentShas: []string{
			// 			"3d95fe3bf954501d3832e50fdd803c5f9eae3f94",
			// 		},
			// 	},
			// 	{
			// 		RepositoryID: repositoryID,
			// 		CreatedAt:    time.Date(2022, 12, 28, 12, 21, 5, 0, time.UTC),
			// 		Sha:          "3d95fe3bf954501d3832e50fdd803c5f9eae3f94",
			// 		ParentShas: []string{
			// 			"5da8e92e9f9243f7ee937170474531393a2cf48f",
			// 			"b9b48bcf26ab79c77e4aa4dcf28ca466bdc3b9fa",
			// 		},
			// 	},
			// 	{
			// 		RepositoryID: repositoryID,
			// 		CreatedAt:    time.Date(2022, 12, 27, 15, 55, 34, 0, time.UTC),
			// 		Sha:          "b9b48bcf26ab79c77e4aa4dcf28ca466bdc3b9fa",
			// 		ParentShas: []string{
			// 			"5da8e92e9f9243f7ee937170474531393a2cf48f",
			// 		},
			// 	},
			// 	{
			// 		RepositoryID: repositoryID,
			// 		CreatedAt:    time.Date(2022, 12, 27, 14, 00, 2, 0, time.UTC),
			// 		Sha:          "5da8e92e9f9243f7ee937170474531393a2cf48f",
			// 		ParentShas: []string{
			// 			"0c9e7c4b194a4a5c7066301a8c4f0c6c061ce9bc",
			// 		},
			// 	},
			// }

			commits := []models.Commit{
				{
					RepositoryID: repositoryID,
					CreatedAt:    time.Date(2022, 12, 27, 14, 00, 2, 0, time.UTC),
					Sha:          "5da8e92e9f9243f7ee937170474531393a2cf48f",
					ParentShas: []string{
						"0c9e7c4b194a4a5c7066301a8c4f0c6c061ce9bc",
					},
				},
				{
					RepositoryID: repositoryID,
					CreatedAt:    time.Date(2022, 12, 27, 15, 55, 34, 0, time.UTC),
					Sha:          "b9b48bcf26ab79c77e4aa4dcf28ca466bdc3b9fa",
					ParentShas: []string{
						"5da8e92e9f9243f7ee937170474531393a2cf48f",
					},
				},
				{
					RepositoryID: repositoryID,
					CreatedAt:    time.Date(2022, 12, 28, 12, 21, 5, 0, time.UTC),
					Sha:          "3d95fe3bf954501d3832e50fdd803c5f9eae3f94",
					ParentShas: []string{
						"5da8e92e9f9243f7ee937170474531393a2cf48f",
						"b9b48bcf26ab79c77e4aa4dcf28ca466bdc3b9fa",
					},
				},
				{
					RepositoryID: repositoryID,
					CreatedAt:    time.Date(2022, 12, 28, 12, 46, 21, 0, time.UTC),
					Sha:          "487d6aedb92ab76bdc03957aceece75db906796e",
					ParentShas: []string{
						"3d95fe3bf954501d3832e50fdd803c5f9eae3f94",
					},
				},
				{
					RepositoryID: repositoryID,
					CreatedAt:    time.Date(2022, 12, 28, 13, 01, 11, 0, time.UTC),
					Sha:          "1db209656ad1ab0e14aaa4e2fe79b6caf8b2a9e7",
					ParentShas: []string{
						"3d95fe3bf954501d3832e50fdd803c5f9eae3f94",
						"487d6aedb92ab76bdc03957aceece75db906796e",
					},
				},
			}

			err = daos.CreateCommits(ctx, repositoryID, &commits)
			Expect(err).To(BeNil())

			firstCommits, err := ingest.GetFirstCommits(ctx, repositoryID, &pipelineRuns)
			Expect(err).To(BeNil())
			Expect(len(*firstCommits)).To(Equal(2))
			Expect((*firstCommits)[0].Sha).To(Equal("b9b48bcf26ab79c77e4aa4dcf28ca466bdc3b9fa"))
			Expect((*firstCommits)[1].Sha).To(Equal("487d6aedb92ab76bdc03957aceece75db906796e"))
		})
	})

	var _ = When("CalculateChanges", func() {
		It("calculates changes from pipeline runs and commits.", func() {
			pipelineID := primitive.NewObjectID()
			pipelineRuns := []models.PipelineRun{
				{
					PipelineID:  pipelineID,
					ExternalID:  713437228,
					Sha:         "3d95fe3bf954501d3832e50fdd803c5f9eae3f94",
					Ref:         "main",
					Status:      "success",
					EventSource: "push",
					CreatedAt:   time.Date(2022, 12, 27, 15, 57, 42, 0, time.UTC),
					UpdatedAt:   time.Date(2022, 12, 27, 15, 59, 12, 0, time.UTC),
					URI:         "https://gitlab.com/foobar/foobar/-/pipelines/713437228",
				},
				{
					PipelineID:  pipelineID,
					ExternalID:  713437229,
					Sha:         "1db209656ad1ab0e14aaa4e2fe79b6caf8b2a9e7",
					Ref:         "main",
					Status:      "success",
					EventSource: "push",
					CreatedAt:   time.Date(2022, 12, 28, 12, 47, 16, 0, time.UTC),
					UpdatedAt:   time.Date(2022, 12, 28, 12, 50, 2, 0, time.UTC),
					URI:         "https://gitlab.com/foobar/foobar/-/pipelines/713437229",
				},
			}

			repositoryID := primitive.NewObjectID()
			commits := []models.Commit{
				{
					RepositoryID: repositoryID,
					CreatedAt:    time.Date(2022, 12, 27, 14, 00, 2, 0, time.UTC),
					Sha:          "5da8e92e9f9243f7ee937170474531393a2cf48f",
					ParentShas: []string{
						"0c9e7c4b194a4a5c7066301a8c4f0c6c061ce9bc",
					},
				},
				{
					RepositoryID: repositoryID,
					CreatedAt:    time.Date(2022, 12, 27, 15, 55, 34, 0, time.UTC),
					Sha:          "b9b48bcf26ab79c77e4aa4dcf28ca466bdc3b9fa",
					ParentShas: []string{
						"5da8e92e9f9243f7ee937170474531393a2cf48f",
					},
				},
				{
					RepositoryID: repositoryID,
					CreatedAt:    time.Date(2022, 12, 28, 12, 21, 5, 0, time.UTC),
					Sha:          "3d95fe3bf954501d3832e50fdd803c5f9eae3f94",
					ParentShas: []string{
						"5da8e92e9f9243f7ee937170474531393a2cf48f",
						"b9b48bcf26ab79c77e4aa4dcf28ca466bdc3b9fa",
					},
				},
				{
					RepositoryID: repositoryID,
					CreatedAt:    time.Date(2022, 12, 28, 12, 46, 21, 0, time.UTC),
					Sha:          "487d6aedb92ab76bdc03957aceece75db906796e",
					ParentShas: []string{
						"3d95fe3bf954501d3832e50fdd803c5f9eae3f94",
					},
				},
				{
					RepositoryID: repositoryID,
					CreatedAt:    time.Date(2022, 12, 28, 13, 01, 11, 0, time.UTC),
					Sha:          "1db209656ad1ab0e14aaa4e2fe79b6caf8b2a9e7",
					ParentShas: []string{
						"3d95fe3bf954501d3832e50fdd803c5f9eae3f94",
						"487d6aedb92ab76bdc03957aceece75db906796e",
					},
				},
			}

			err := daos.CreateCommits(ctx, repositoryID, &commits)
			Expect(err).To(BeNil())

			firstCommits, err := ingest.GetFirstCommits(ctx, repositoryID, &pipelineRuns)

			changes, err := ingest.CalculateChanges(ctx, firstCommits, &pipelineRuns)
			Expect(err).To(BeNil())
			Expect(len(*changes)).To(Equal(2))
			Expect((*changes)[0].LeadTime).To(Equal(float64(218)))
			Expect((*changes)[1].LeadTime).To(Equal(float64(221)))
		})
	})

	var _ = When("CreateChanges", func() {
		It("creates changes from pipeline runs and commits.", func() {
			pipelineID := primitive.NewObjectID()
			pipelineRuns := []models.PipelineRun{
				{
					PipelineID:  pipelineID,
					ExternalID:  713437228,
					Sha:         "3d95fe3bf954501d3832e50fdd803c5f9eae3f94",
					Ref:         "main",
					Status:      "success",
					EventSource: "push",
					CreatedAt:   time.Date(2022, 12, 27, 15, 57, 42, 0, time.UTC),
					UpdatedAt:   time.Date(2022, 12, 27, 15, 59, 12, 0, time.UTC),
					URI:         "https://gitlab.com/foobar/foobar/-/pipelines/713437228",
				},
				{
					PipelineID:  pipelineID,
					ExternalID:  713437229,
					Sha:         "1db209656ad1ab0e14aaa4e2fe79b6caf8b2a9e7",
					Ref:         "main",
					Status:      "success",
					EventSource: "push",
					CreatedAt:   time.Date(2022, 12, 28, 12, 47, 16, 0, time.UTC),
					UpdatedAt:   time.Date(2022, 12, 28, 12, 50, 2, 0, time.UTC),
					URI:         "https://gitlab.com/foobar/foobar/-/pipelines/713437229",
				},
			}

			err := daos.CreatePipelineRuns(ctx, pipelineID, &pipelineRuns)
			Expect(err).To(BeNil())

			repositoryID := primitive.NewObjectID()
			commits := []models.Commit{
				{
					RepositoryID: repositoryID,
					CreatedAt:    time.Date(2022, 12, 27, 14, 00, 2, 0, time.UTC),
					Sha:          "5da8e92e9f9243f7ee937170474531393a2cf48f",
					ParentShas: []string{
						"0c9e7c4b194a4a5c7066301a8c4f0c6c061ce9bc",
					},
				},
				{
					RepositoryID: repositoryID,
					CreatedAt:    time.Date(2022, 12, 27, 15, 55, 34, 0, time.UTC),
					Sha:          "b9b48bcf26ab79c77e4aa4dcf28ca466bdc3b9fa",
					ParentShas: []string{
						"5da8e92e9f9243f7ee937170474531393a2cf48f",
					},
				},
				{
					RepositoryID: repositoryID,
					CreatedAt:    time.Date(2022, 12, 28, 12, 21, 5, 0, time.UTC),
					Sha:          "3d95fe3bf954501d3832e50fdd803c5f9eae3f94",
					ParentShas: []string{
						"5da8e92e9f9243f7ee937170474531393a2cf48f",
						"b9b48bcf26ab79c77e4aa4dcf28ca466bdc3b9fa",
					},
				},
				{
					RepositoryID: repositoryID,
					CreatedAt:    time.Date(2022, 12, 28, 12, 46, 21, 0, time.UTC),
					Sha:          "487d6aedb92ab76bdc03957aceece75db906796e",
					ParentShas: []string{
						"3d95fe3bf954501d3832e50fdd803c5f9eae3f94",
					},
				},
				{
					RepositoryID: repositoryID,
					CreatedAt:    time.Date(2022, 12, 28, 13, 01, 11, 0, time.UTC),
					Sha:          "1db209656ad1ab0e14aaa4e2fe79b6caf8b2a9e7",
					ParentShas: []string{
						"3d95fe3bf954501d3832e50fdd803c5f9eae3f94",
						"487d6aedb92ab76bdc03957aceece75db906796e",
					},
				},
			}

			err = daos.CreateCommits(ctx, repositoryID, &commits)
			Expect(err).To(BeNil())

			err = ingest.CreateChanges(ctx, repositoryID, pipelineID)
			Expect(err).To(BeNil())

			var changes []models.Change
			err = daos.ListChanges(ctx, repositoryID, &changes)

			Expect(len(changes)).To(Equal(2))
			Expect(changes[0].LeadTime).To(Equal(float64(218)))
			Expect(changes[1].LeadTime).To(Equal(float64(221)))
		})
	})
})
