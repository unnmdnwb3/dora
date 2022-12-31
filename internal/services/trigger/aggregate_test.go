package trigger_test

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
	"github.com/unnmdnwb3/dora/internal/services/trigger"
	"github.com/unnmdnwb3/dora/internal/utils/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ = Describe("services.trigger.aggregate", func() {
	ctx := context.Background()
	externalID := 40649465

	var _ = BeforeEach(func() {
		_ = godotenv.Load("./../../../test/.env")
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
			pipelineID, _ := types.StringToObjectID("638e00b85edd5bef25e5e9e1")
			createdAt1, _ := time.Parse(time.RFC3339, "2019-10-09T09:11:20.861Z")
			updatedAt1, _ := time.Parse(time.RFC3339, "2019-10-09T09:12:20.861Z")
			pipelineRun1 := models.PipelineRun{
				PipelineID:  pipelineID,
				ExternalID:  713437220,
				Sha:         "1cfffa2ae16528e36115ece8b1f2601bcf74414e",
				Ref:         "main",
				Status:      "success",
				EventSource: "push",
				CreatedAt:   createdAt1,
				UpdatedAt:   updatedAt1,
				URI:         "https://gitlab.com/foobar/foobar/-/pipelines/114884002",
			}

			createdAt2, _ := time.Parse(time.RFC3339, "2019-10-11T09:11:20.861Z")
			updatedAt2, _ := time.Parse(time.RFC3339, "2019-10-11T09:12:20.861Z")
			pipelineRun2 := models.PipelineRun{
				PipelineID:  pipelineID,
				ExternalID:  713437221,
				Sha:         "345207c839e94a939aebdc86835ae2e2a6c85acb",
				Ref:         "main",
				Status:      "success",
				EventSource: "push",
				CreatedAt:   createdAt2,
				UpdatedAt:   updatedAt2,
				URI:         "https://gitlab.com/foobar/foobar/-/pipelines/114884003",
			}

			createdAt3, _ := time.Parse(time.RFC3339, "2019-10-11T09:13:20.861Z")
			updatedAt3, _ := time.Parse(time.RFC3339, "2019-10-11T09:14:20.861Z")
			pipelineRun3 := models.PipelineRun{
				PipelineID:  pipelineID,
				ExternalID:  713437222,
				Sha:         "dcc7ef44dc6a376854c5f2cc42b0b24aa3a9ed10",
				Ref:         "main",
				Status:      "success",
				EventSource: "push",
				CreatedAt:   createdAt3,
				UpdatedAt:   updatedAt3,
				URI:         "https://gitlab.com/foobar/foobar/-/pipelines/114884004",
			}

			pipelineRuns := []models.PipelineRun{pipelineRun1, pipelineRun2, pipelineRun3}
			pipelineRunsPerDay, err := trigger.CalculatePipelineRunsPerDays(ctx, &pipelineRuns)
			Expect(err).To(BeNil())
			Expect(len(*pipelineRunsPerDay)).To(Equal(2))
			Expect((*pipelineRunsPerDay)[0].TotalPipelineRuns).To(Equal(1))
			Expect((*pipelineRunsPerDay)[1].TotalPipelineRuns).To(Equal(2))
		})
	})

	var _ = When("CreatePipelineRunsPerDays", func() {
		It("calculates and creates the pipeline runs for each day.", func() {
			// create dataflow
			repositoryIntegrationID := primitive.NewObjectID()
			pipelineIntegrationID := primitive.NewObjectID()
			deploymentIntegrationID := primitive.NewObjectID()
			repository := models.Repository{
				IntegrationID:  repositoryIntegrationID,
				ExternalID:     externalID,
				NamespacedName: "foobar/foobar",
				DefaultBranch:  "main",
				URI:            "https://gitlab.com/foobar/foobar",
			}
			pipeline := models.Pipeline{
				IntegrationID:  pipelineIntegrationID,
				ExternalID:     externalID,
				NamespacedName: "foobar/foobar",
				DefaultBranch:  "main",
				URI:            "https://gitlab.com/foobar/foobar/-/pipelines",
			}
			deployment := models.Deployment{
				IntegrationID: deploymentIntegrationID,
			}
			dataflow := models.Dataflow{
				Repository: repository,
				Pipeline:   pipeline,
				Deployment: deployment,
			}
			err := daos.CreateDataflow(ctx, &dataflow)
			Expect(err).To(BeNil())

			// create pipeline runs
			createdAt1, _ := time.Parse(time.RFC3339, "2019-10-09T09:11:20.861Z")
			updatedAt1, _ := time.Parse(time.RFC3339, "2019-10-09T09:12:20.861Z")
			pipelineID, _ := types.StringToObjectID("638e00b85edd5bef25e5e9e1")
			pipelineRun1 := models.PipelineRun{
				PipelineID:  pipelineID,
				ExternalID:  713437220,
				Sha:         "1cfffa2ae16528e36115ece8b1f2601bcf74414e",
				Ref:         "main",
				Status:      "success",
				EventSource: "push",
				CreatedAt:   createdAt1,
				UpdatedAt:   updatedAt1,
				URI:         "https://gitlab.com/foobar/foobar/-/pipelines/114884002",
			}

			createdAt2, _ := time.Parse(time.RFC3339, "2019-10-11T09:11:20.861Z")
			updatedAt2, _ := time.Parse(time.RFC3339, "2019-10-11T09:12:20.861Z")
			pipelineRun2 := models.PipelineRun{
				PipelineID:  pipelineID,
				ExternalID:  713437221,
				Sha:         "345207c839e94a939aebdc86835ae2e2a6c85acb",
				Ref:         "main",
				Status:      "success",
				EventSource: "push",
				CreatedAt:   createdAt2,
				UpdatedAt:   updatedAt2,
				URI:         "https://gitlab.com/foobar/foobar/-/pipelines/114884003",
			}

			createdAt3, _ := time.Parse(time.RFC3339, "2019-10-11T09:13:20.861Z")
			updatedAt3, _ := time.Parse(time.RFC3339, "2019-10-11T09:14:20.861Z")
			pipelineRun3 := models.PipelineRun{
				PipelineID:  pipelineID,
				ExternalID:  713437222,
				Sha:         "dcc7ef44dc6a376854c5f2cc42b0b24aa3a9ed10",
				Ref:         "main",
				Status:      "success",
				EventSource: "push",
				CreatedAt:   createdAt3,
				UpdatedAt:   updatedAt3,
				URI:         "https://gitlab.com/foobar/foobar/-/pipelines/114884004",
			}

			pipelineRuns := []models.PipelineRun{pipelineRun1, pipelineRun2, pipelineRun3}
			err = daos.CreatePipelineRuns(ctx, dataflow.Pipeline.ID, &pipelineRuns)
			Expect(err).To(BeNil())

			err = trigger.CreatePipelineRunsPerDays(ctx, dataflow.Pipeline.ID)
			Expect(err).To(BeNil())

			var pipelineRunsPerDays []models.PipelineRunsPerDay
			err = daos.ListPipelineRunsPerDays(ctx, dataflow.Pipeline.ID, &pipelineRunsPerDays)
			Expect(err).To(BeNil())
			Expect(pipelineRunsPerDays).To(HaveLen(2))
		})
	})

	var _ = When("CalculateIncidentsPerDays", func() {
		It("calculates IncidentsPerDays based on Incidents.", func() {
			deploymentID := primitive.NewObjectID()
			incident1 := models.Incident{
				DeploymentID: deploymentID,
				StartDate:    time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
				EndDate:      time.Date(2022, 12, 27, 14, 16, 42, 0, time.UTC),
			}

			incident2 := models.Incident{
				DeploymentID: deploymentID,
				StartDate:    time.Date(2022, 12, 27, 17, 39, 21, 0, time.UTC),
				EndDate:      time.Date(2022, 12, 27, 17, 54, 21, 0, time.UTC),
			}

			incident3 := models.Incident{
				DeploymentID: deploymentID,
				StartDate:    time.Date(2022, 12, 29, 02, 21, 42, 0, time.UTC),
				EndDate:      time.Date(2022, 12, 29, 02, 41, 42, 0, time.UTC),
			}

			incidents := []models.Incident{incident1, incident2, incident3}
			err := daos.CreateIncidents(ctx, &incidents)
			Expect(err).To(BeNil())

			incidentsPerDays, err := trigger.CalculateIncidentsPerDays(ctx, &incidents)
			Expect(err).To(BeNil())
			Expect(len(*incidentsPerDays)).To(Equal(2))
			Expect((*incidentsPerDays)[0].TotalIncidents).To(Equal(2))
			Expect((*incidentsPerDays)[0].TotalDuration).To(Equal(float64(4500)))
			Expect((*incidentsPerDays)[1].TotalIncidents).To(Equal(1))
			Expect((*incidentsPerDays)[1].TotalDuration).To(Equal(float64(1200)))
		})
	})

	var _ = When("CreateIncidentsPerDays", func() {
		It("creates IncidentsPerDays based on Incidents.", func() {
			deploymentID := primitive.NewObjectID()
			incident1 := models.Incident{
				DeploymentID: deploymentID,
				StartDate:    time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
				EndDate:      time.Date(2022, 12, 27, 14, 16, 42, 0, time.UTC),
			}

			incident2 := models.Incident{
				DeploymentID: deploymentID,
				StartDate:    time.Date(2022, 12, 27, 17, 39, 21, 0, time.UTC),
				EndDate:      time.Date(2022, 12, 27, 17, 54, 21, 0, time.UTC),
			}

			incident3 := models.Incident{
				DeploymentID: deploymentID,
				StartDate:    time.Date(2022, 12, 29, 02, 21, 42, 0, time.UTC),
				EndDate:      time.Date(2022, 12, 29, 02, 41, 42, 0, time.UTC),
			}

			incidents := []models.Incident{incident1, incident2, incident3}
			err := daos.CreateIncidents(ctx, &incidents)
			Expect(err).To(BeNil())

			err = trigger.CreateIncidentsPerDays(ctx, deploymentID)
			Expect(err).To(BeNil())

			var incidentsPerDays []models.IncidentsPerDay
			err = daos.ListIncidentsPerDays(ctx, deploymentID, &incidentsPerDays)
			Expect(err).To(BeNil())
			Expect(incidentsPerDays).To(HaveLen(2))
			Expect(incidentsPerDays[0].TotalIncidents).To(Equal(2))
			Expect(incidentsPerDays[0].TotalDuration).To(Equal(float64(4500)))
			Expect(incidentsPerDays[1].TotalIncidents).To(Equal(1))
			Expect(incidentsPerDays[1].TotalDuration).To(Equal(float64(1200)))
		})
	})

	var _ = When("CalculateIncidents", func() {
		It("calculates Incidents based on MonitoringDataPoints.", func() {
			deployment := models.Deployment{
				IntegrationID: primitive.NewObjectID(),
				Query:         "job:http_total_requests:internal_server_error_percentage",
				Step:          "5m",
				Relation:      "gt",
				Threshold:     0.2,
			}

			dataPoints := []models.MonitoringDataPoint{
				{
					Value:     0.1,
					CreatedAt: time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
				},
				{
					Value:     0.3,
					CreatedAt: time.Date(2022, 12, 27, 13, 21, 43, 0, time.UTC),
				},
				{
					Value:     0.35,
					CreatedAt: time.Date(2022, 12, 27, 13, 26, 42, 0, time.UTC),
				},
				{
					Value:     0.4,
					CreatedAt: time.Date(2022, 12, 27, 13, 36, 44, 0, time.UTC),
				},
				{
					Value:     0.3,
					CreatedAt: time.Date(2022, 12, 27, 13, 41, 42, 0, time.UTC),
				},
				{
					Value:     0.1,
					CreatedAt: time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
				},
			}

			incidents, err := trigger.CalculateIncidents(ctx, &deployment, &dataPoints)
			Expect(err).To(BeNil())
			Expect(len(*incidents)).To(Equal(2))
			Expect((*incidents)[0].StartDate).To(Equal(time.Date(2022, 12, 27, 13, 21, 43, 0, time.UTC)))
			Expect((*incidents)[0].EndDate).To(Equal(time.Date(2022, 12, 27, 13, 26, 42, 0, time.UTC)))
			Expect((*incidents)[1].StartDate).To(Equal(time.Date(2022, 12, 27, 13, 36, 44, 0, time.UTC)))
			Expect((*incidents)[1].EndDate).To(Equal(time.Date(2022, 12, 27, 13, 41, 42, 0, time.UTC)))
		})
	})

	var _ = When("CreateIncidents", func() {
		It("creates Incidents based on MonitoringDataPoints.", func() {
			deployment := models.Deployment{
				IntegrationID: primitive.NewObjectID(),
				Query:         "job:http_total_requests:internal_server_error_percentage",
				Step:          "5m",
				Relation:      "gt",
				Threshold:     0.2,
			}

			dataPoints := []models.MonitoringDataPoint{
				{
					Value:     0.1,
					CreatedAt: time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
				},
				{
					Value:     0.3,
					CreatedAt: time.Date(2022, 12, 27, 13, 21, 43, 0, time.UTC),
				},
				{
					Value:     0.35,
					CreatedAt: time.Date(2022, 12, 27, 13, 26, 42, 0, time.UTC),
				},
				{
					Value:     0.4,
					CreatedAt: time.Date(2022, 12, 27, 13, 36, 44, 0, time.UTC),
				},
				{
					Value:     0.3,
					CreatedAt: time.Date(2022, 12, 27, 13, 41, 42, 0, time.UTC),
				},
				{
					Value:     0.1,
					CreatedAt: time.Date(2022, 12, 27, 13, 16, 42, 0, time.UTC),
				},
			}

			err := trigger.CreateIncidents(ctx, &deployment, &dataPoints)
			Expect(err).To(BeNil())

			var incidents []models.Incident
			err = daos.ListIncidentsByFilter(ctx, bson.M{"deployment_id": deployment.ID}, &incidents)
			Expect(incidents).To(HaveLen(2))
			Expect(incidents[0].StartDate).To(Equal(time.Date(2022, 12, 27, 13, 21, 43, 0, time.UTC)))
			Expect(incidents[0].EndDate).To(Equal(time.Date(2022, 12, 27, 13, 26, 42, 0, time.UTC)))
			Expect(incidents[1].StartDate).To(Equal(time.Date(2022, 12, 27, 13, 36, 44, 0, time.UTC)))
			Expect(incidents[1].EndDate).To(Equal(time.Date(2022, 12, 27, 13, 41, 42, 0, time.UTC)))
		})
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
						"b9b48bcf26ab79c77e4aa4dcf28ca466bdc3b9fa",
						"5da8e92e9f9243f7ee937170474531393a2cf48f",
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
						"487d6aedb92ab76bdc03957aceece75db906796e",
						"3d95fe3bf954501d3832e50fdd803c5f9eae3f94",
					},
				},
			}

			err = daos.CreateCommits(ctx, repositoryID, &commits)
			Expect(err).To(BeNil())

			firstCommits, err := trigger.GetFirstCommits(ctx, repositoryID, &pipelineRuns)
			Expect(err).To(BeNil())
			Expect(len(*firstCommits)).To(Equal(2))
			Expect((*firstCommits)[0].Sha).To(Equal("b9b48bcf26ab79c77e4aa4dcf28ca466bdc3b9fa"))
			Expect((*firstCommits)[1].Sha).To(Equal("487d6aedb92ab76bdc03957aceece75db906796e"))
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
			commits := []models.Commit{
				{
					RepositoryID: repositoryID,
					CreatedAt:    time.Date(2022, 12, 27, 11, 0, 2, 0, time.UTC),
					Sha:          "5da8e92e9f9243f7ee937170474531393a2cf48f",
					ParentShas: []string{
						"0c9e7c4b194a4a5c7066301a8c4f0c6c061ce9bc",
					},
				},
				{
					RepositoryID: repositoryID,
					CreatedAt:    time.Date(2022, 12, 27, 12, 26, 34, 0, time.UTC),
					Sha:          "b9b48bcf26ab79c77e4aa4dcf28ca466bdc3b9fa",
					ParentShas: []string{
						"5da8e92e9f9243f7ee937170474531393a2cf48f",
					},
				},
				{
					RepositoryID: repositoryID,
					CreatedAt:    time.Date(2022, 12, 27, 13, 21, 41, 0, time.UTC),
					Sha:          "3d95fe3bf954501d3832e50fdd803c5f9eae3f94",
					ParentShas: []string{
						"b9b48bcf26ab79c77e4aa4dcf28ca466bdc3b9fa",
						"5da8e92e9f9243f7ee937170474531393a2cf48f",
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
					CreatedAt:    time.Date(2022, 12, 28, 15, 37, 20, 0, time.UTC),
					Sha:          "1db209656ad1ab0e14aaa4e2fe79b6caf8b2a9e7",
					ParentShas: []string{
						"487d6aedb92ab76bdc03957aceece75db906796e",
						"3d95fe3bf954501d3832e50fdd803c5f9eae3f94",
					},
				},
			}

			err = daos.CreateCommits(ctx, repositoryID, &commits)
			Expect(err).To(BeNil())

			err = trigger.CreateChanges(ctx, repositoryID, pipelineID)
			Expect(err).To(BeNil())

			var changes []models.Change
			err = daos.ListChanges(ctx, repositoryID, &changes)

			Expect(len(changes)).To(Equal(2))
			Expect(changes[0].LeadTime).To(Equal(float64(3308)))
			Expect(changes[1].LeadTime).To(Equal(float64(10616)))
		})
	})
})
