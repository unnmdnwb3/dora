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
			repositoryID := primitive.NewObjectID()
			commit := models.Commit{
				RepositoryID: repositoryID,
				CreatedAt:    time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
				Sha:          "1db209656ad1ab0e14aaa4e2fe79b6caf8b2a9e7",
				ParentIds: []string{
					"487d6aedb92ab76bdc03957aceece75db906796e",
					"3d95fe3bf954501d3832e50fdd803c5f9eae3f94",
				},
			}
			err := daos.CreateCommit(ctx, repositoryID, &commit)
			Expect(err).To(BeNil())
			Expect(commit.ID).To(Not(BeEmpty()))
		})
	})

	var _ = When("CreateIncidents", func() {
		It("creates creates many new Incidents.", func() {
			repositoryID := primitive.NewObjectID()
			commits := []models.Commit{
				{
					RepositoryID: repositoryID,
					CreatedAt:    time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
					Sha:          "1db209656ad1ab0e14aaa4e2fe79b6caf8b2a9e7",
					ParentIds: []string{
						"487d6aedb92ab76bdc03957aceece75db906796e",
						"3d95fe3bf954501d3832e50fdd803c5f9eae3f94",
					},
				},
				{
					RepositoryID: repositoryID,
					CreatedAt:    time.Date(2022, 12, 27, 13, 21, 42, 0, time.UTC),
					Sha:          "487d6aedb92ab76bdc03957aceece75db906796e",
					ParentIds: []string{
						"398dc0ca313035ea4eb7ab3f29a5500631660fb7",
					},
				},
			}
			err := daos.CreateCommits(ctx, repositoryID, &commits)
			Expect(err).To(BeNil())
			Expect(commits[0].ID).To(Not(BeEmpty()))
			Expect(commits[1].ID).To(Not(BeEmpty()))
		})
	})

	var _ = When("GetCommit", func() {
		It("retrieves an Commit.", func() {
			repositoryID := primitive.NewObjectID()
			commit := models.Commit{
				RepositoryID: repositoryID,
				CreatedAt:    time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
				Sha:          "1db209656ad1ab0e14aaa4e2fe79b6caf8b2a9e7",
				ParentIds: []string{
					"487d6aedb92ab76bdc03957aceece75db906796e",
					"3d95fe3bf954501d3832e50fdd803c5f9eae3f94",
				},
			}
			err := daos.CreateCommit(ctx, repositoryID, &commit)
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
			repositoryID := primitive.NewObjectID()
			commits := []models.Commit{
				{
					RepositoryID: repositoryID,
					CreatedAt:    time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
					Sha:          "1db209656ad1ab0e14aaa4e2fe79b6caf8b2a9e7",
					ParentIds: []string{
						"487d6aedb92ab76bdc03957aceece75db906796e",
						"3d95fe3bf954501d3832e50fdd803c5f9eae3f94",
					},
				},
				{
					RepositoryID: repositoryID,
					CreatedAt:    time.Date(2022, 12, 27, 13, 21, 42, 0, time.UTC),
					Sha:          "487d6aedb92ab76bdc03957aceece75db906796e",
					ParentIds: []string{
						"398dc0ca313035ea4eb7ab3f29a5500631660fb7",
					},
				},
			}
			err := daos.CreateCommits(ctx, repositoryID, &commits)
			Expect(err).To(BeNil())
			Expect(commits[0].ID).To(Not(BeNil()))
			Expect(commits[1].ID).To(Not(BeNil()))

			var findCommits []models.Commit
			err = daos.ListCommits(ctx, repositoryID, &findCommits)
			Expect(err).To(BeNil())
			Expect(len(findCommits)).To(Equal(2))
		})
	})

	var _ = When("ListCommitsByFilter", func() {
		It("retrieves many Commits conforming to a filter.", func() {
			repositoryID := primitive.NewObjectID()
			differentRepositoryID := primitive.NewObjectID()
			commits := []models.Commit{
				{
					RepositoryID: repositoryID,
					CreatedAt:    time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
					Sha:          "1db209656ad1ab0e14aaa4e2fe79b6caf8b2a9e7",
					ParentIds: []string{
						"487d6aedb92ab76bdc03957aceece75db906796e",
						"3d95fe3bf954501d3832e50fdd803c5f9eae3f94",
					},
				},
				{
					RepositoryID: differentRepositoryID,
					CreatedAt:    time.Date(2022, 12, 27, 13, 21, 42, 0, time.UTC),
					Sha:          "487d6aedb92ab76bdc03957aceece75db906796e",
					ParentIds: []string{
						"398dc0ca313035ea4eb7ab3f29a5500631660fb7",
					},
				},
			}

			err := daos.CreateCommits(ctx, repositoryID, &commits)
			Expect(err).To(BeNil())
			Expect(commits[0].ID).To(Not(BeNil()))
			Expect(commits[1].ID).To(Not(BeNil()))

			var findCommits []models.Commit
			filter := bson.M{"repository_id": repositoryID}
			err = daos.ListCommitsByFilter(ctx, filter, &findCommits)
			Expect(err).To(BeNil())
			Expect(findCommits).To(HaveLen(2))
		})
	})

	var _ = When("UpdateCommit", func() {
		It("updates an Commit.", func() {
			repositoryID := primitive.NewObjectID()
			commit := models.Commit{
				RepositoryID: repositoryID,
				CreatedAt:    time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
				Sha:          "1db209656ad1ab0e14aaa4e2fe79b6caf8b2a9e7",
				ParentIds: []string{
					"487d6aedb92ab76bdc03957aceece75db906796e",
					"3d95fe3bf954501d3832e50fdd803c5f9eae3f94",
				},
			}
			err := daos.CreateCommit(ctx, repositoryID, &commit)
			Expect(err).To(BeNil())
			Expect(commit.ID).To(Not(BeEmpty()))

			updateCommit := models.Commit{
				RepositoryID: primitive.NewObjectID(),
				CreatedAt:    time.Date(2022, 12, 27, 13, 26, 42, 0, time.UTC),
				Sha:          "1db209656ad1ab0e14aaa4e2fe79b6caf8b2a9e7",
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
			repositoryID := primitive.NewObjectID()
			commit := models.Commit{
				RepositoryID: repositoryID,
				CreatedAt:    time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
				Sha:          "1db209656ad1ab0e14aaa4e2fe79b6caf8b2a9e7",
				ParentIds: []string{
					"487d6aedb92ab76bdc03957aceece75db906796e",
					"3d95fe3bf954501d3832e50fdd803c5f9eae3f94",
				},
			}
			err := daos.CreateCommit(ctx, repositoryID, &commit)
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
