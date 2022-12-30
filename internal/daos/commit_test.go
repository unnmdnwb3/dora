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

var _ = Describe("daos.Commit", func() {
	ctx := context.Background()

	var _ = When("CreateCommit", func() {
		It("creates a new Commit.", func() {
			commit := models.Commit{
				PipelineID: primitive.NewObjectID(),
				CreatedAt:  time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
				Sha:        "1db209656ad1ab0e14aaa4e2fe79b6caf8b2a9e7",
				ParentIds: []string{
					"487d6aedb92ab76bdc03957aceece75db906796e",
					"3d95fe3bf954501d3832e50fdd803c5f9eae3f94",
				},
			}
			err := daos.CreateCommit(ctx, &commit)
			Expect(err).To(BeNil())
			Expect(commit.ID).To(Not(BeEmpty()))
		})
	})

	var _ = When("CreateIncidents", func() {
		It("creates creates many new Incidents.", func() {
			pipelineID := primitive.NewObjectID()
			commits := []models.Commit{
				{
					PipelineID: pipelineID,
					CreatedAt:  time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
					Sha:        "1db209656ad1ab0e14aaa4e2fe79b6caf8b2a9e7",
					ParentIds: []string{
						"487d6aedb92ab76bdc03957aceece75db906796e",
						"3d95fe3bf954501d3832e50fdd803c5f9eae3f94",
					},
				},
				{
					PipelineID: pipelineID,
					CreatedAt:  time.Date(2022, 12, 27, 13, 21, 42, 0, time.UTC),
					Sha:        "487d6aedb92ab76bdc03957aceece75db906796e",
					ParentIds: []string{
						"398dc0ca313035ea4eb7ab3f29a5500631660fb7",
					},
				},
			}
			err := daos.CreateCommits(ctx, &commits)
			Expect(err).To(BeNil())
			Expect(commits[0].ID).To(Not(BeEmpty()))
			Expect(commits[1].ID).To(Not(BeEmpty()))
		})
	})

	var _ = When("GetCommit", func() {
		It("retrieves an Commit.", func() {
			commit := models.Commit{
				PipelineID: primitive.NewObjectID(),
				CreatedAt:  time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
				Sha:        "1db209656ad1ab0e14aaa4e2fe79b6caf8b2a9e7",
				ParentIds: []string{
					"487d6aedb92ab76bdc03957aceece75db906796e",
					"3d95fe3bf954501d3832e50fdd803c5f9eae3f94",
				},
			}
			err := daos.CreateCommit(ctx, &commit)
			Expect(err).To(BeNil())
			Expect(commit.ID).To(Not(BeEmpty()))

			var findCommit models.Commit
			err = daos.GetCommit(ctx, commit.ID, &findCommit)
			Expect(err).To(BeNil())
			Expect(findCommit.ID).To(Equal(commit.ID))
		})
	})

	var _ = When("ListCommits", func() {
		It("retrieves many Commits.", func() {
			pipelineID := primitive.NewObjectID()
			commits := []models.Commit{
				{
					PipelineID: pipelineID,
					CreatedAt:  time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
					Sha:        "1db209656ad1ab0e14aaa4e2fe79b6caf8b2a9e7",
					ParentIds: []string{
						"487d6aedb92ab76bdc03957aceece75db906796e",
						"3d95fe3bf954501d3832e50fdd803c5f9eae3f94",
					},
				},
				{
					PipelineID: pipelineID,
					CreatedAt:  time.Date(2022, 12, 27, 13, 21, 42, 0, time.UTC),
					Sha:        "487d6aedb92ab76bdc03957aceece75db906796e",
					ParentIds: []string{
						"398dc0ca313035ea4eb7ab3f29a5500631660fb7",
					},
				},
			}
			err := daos.CreateCommits(ctx, &commits)
			Expect(err).To(BeNil())
			Expect(commits[0].ID).To(Not(BeNil()))
			Expect(commits[1].ID).To(Not(BeNil()))

			var findCommits []models.Commit
			err = daos.ListCommits(ctx, pipelineID, &findCommits)
			Expect(err).To(BeNil())
			Expect(findCommits).To(HaveLen(2))
		})
	})

	var _ = When("ListCommitsByFilter", func() {
		It("retrieves many Commits conforming to a filter.", func() {
			pipelineID := primitive.NewObjectID()
			differentPipelineID := primitive.NewObjectID()
			commits := []models.Commit{
				{
					PipelineID: pipelineID,
					CreatedAt:  time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
					Sha:        "1db209656ad1ab0e14aaa4e2fe79b6caf8b2a9e7",
					ParentIds: []string{
						"487d6aedb92ab76bdc03957aceece75db906796e",
						"3d95fe3bf954501d3832e50fdd803c5f9eae3f94",
					},
				},
				{
					PipelineID: differentPipelineID,
					CreatedAt:  time.Date(2022, 12, 27, 13, 21, 42, 0, time.UTC),
					Sha:        "487d6aedb92ab76bdc03957aceece75db906796e",
					ParentIds: []string{
						"398dc0ca313035ea4eb7ab3f29a5500631660fb7",
					},
				},
			}

			err := daos.CreateCommits(ctx, &commits)
			Expect(err).To(BeNil())
			Expect(commits[0].ID).To(Not(BeNil()))
			Expect(commits[1].ID).To(Not(BeNil()))

			var findCommits []models.Commit
			filter := bson.M{"pipeline_id": pipelineID}
			err = daos.ListCommitsByFilter(ctx, filter, &findCommits)
			Expect(err).To(BeNil())
			Expect(findCommits).To(HaveLen(1))
		})
	})

	var _ = When("UpdateCommit", func() {
		It("updates an Commit.", func() {
			commit := models.Commit{
				PipelineID: primitive.NewObjectID(),
				CreatedAt:  time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
				Sha:        "1db209656ad1ab0e14aaa4e2fe79b6caf8b2a9e7",
				ParentIds: []string{
					"487d6aedb92ab76bdc03957aceece75db906796e",
					"3d95fe3bf954501d3832e50fdd803c5f9eae3f94",
				},
			}
			err := daos.CreateCommit(ctx, &commit)
			Expect(err).To(BeNil())
			Expect(commit.ID).To(Not(BeEmpty()))

			updateCommit := models.Commit{
				PipelineID: primitive.NewObjectID(),
				CreatedAt:  time.Date(2022, 12, 27, 13, 26, 42, 0, time.UTC),
				Sha:        "1db209656ad1ab0e14aaa4e2fe79b6caf8b2a9e7",
				ParentIds: []string{
					"487d6aedb92ab76bdc03957aceece75db906796e",
					"3d95fe3bf954501d3832e50fdd803c5f9eae3f94",
				},
			}
			err = daos.UpdateCommit(ctx, commit.ID, &updateCommit)
			Expect(err).To(BeNil())
			Expect(updateCommit.CreatedAt).To(Equal(time.Date(2022, 12, 27, 13, 26, 42, 0, time.UTC)))
		})
	})

	var _ = When("DeleteOne", func() {
		It("deletes a document with ID in a collection", func() {
			commit := models.Commit{
				PipelineID: primitive.NewObjectID(),
				CreatedAt:  time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
				Sha:        "1db209656ad1ab0e14aaa4e2fe79b6caf8b2a9e7",
				ParentIds: []string{
					"487d6aedb92ab76bdc03957aceece75db906796e",
					"3d95fe3bf954501d3832e50fdd803c5f9eae3f94",
				},
			}
			err := daos.CreateCommit(ctx, &commit)
			Expect(err).To(BeNil())
			Expect(commit.ID).To(Not(BeEmpty()))

			err = daos.DeleteCommit(ctx, commit.ID)
			Expect(err).To(BeNil())

			var findCommit models.Commit
			err = daos.GetCommit(ctx, commit.ID, &findCommit)
			Expect(err).To(Not(BeNil()))
		})
	})
})
