package daos_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/unnmdnwb3/dora/internal/daos"
	"github.com/unnmdnwb3/dora/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ = Describe("daos.Dataflow", func() {
	ctx := context.Background()
	externalID := 40649465

	var _ = When("CreateDataflow", func() {
		It("creates creates a new Dataflow.", func() {
			repository := models.Repository{
				IntegrationID:  primitive.NewObjectID(),
				ExternalID:     externalID,
				NamespacedName: "foobar/foobar",
				DefaultBranch:  "main",
				URI:            "https://gitlab.com/foobar/foobar",
			}
			pipeline := models.Pipeline{
				IntegrationID:  primitive.NewObjectID(),
				ExternalID:     externalID,
				NamespacedName: "foobar/foobar",
				DefaultBranch:  "main",
				URI:            "https://gitlab.com/foobar/foobar/-/pipelines",
			}
			deployment := models.Deployment{
				IntegrationID: primitive.NewObjectID(),
			}
			dataflow := models.Dataflow{
				Repository: repository,
				Pipeline:   pipeline,
				Deployment: deployment,
			}
			err := daos.CreateDataflow(ctx, &dataflow)
			Expect(err).To(BeNil())
			Expect(dataflow.ID).To(Not(BeEmpty()))
			Expect(dataflow.Repository.ID).To(Not(BeEmpty()))
			Expect(dataflow.Pipeline.ID).To(Not(BeEmpty()))
			Expect(dataflow.Deployment.ID).To(Not(BeEmpty()))
		})
	})

	var _ = When("GetDataflow", func() {
		It("retrieves an Dataflow.", func() {
			repository := models.Repository{
				IntegrationID:  primitive.NewObjectID(),
				ExternalID:     externalID,
				NamespacedName: "foobar/foobar",
				DefaultBranch:  "main",
				URI:            "https://gitlab.com/foobar/foobar",
			}
			pipeline := models.Pipeline{
				IntegrationID:  primitive.NewObjectID(),
				ExternalID:     externalID,
				NamespacedName: "foobar/foobar",
				DefaultBranch:  "main",
				URI:            "https://gitlab.com/foobar/foobar/-/pipelines",
			}
			deployment := models.Deployment{
				IntegrationID: primitive.NewObjectID(),
			}
			dataflow := models.Dataflow{
				Repository: repository,
				Pipeline:   pipeline,
				Deployment: deployment,
			}
			err := daos.CreateDataflow(ctx, &dataflow)
			Expect(err).To(BeNil())
			Expect(dataflow.ID).To(Not(BeEmpty()))

			var findDataflow models.Dataflow
			err = daos.GetDataflow(ctx, dataflow.ID, &findDataflow)
			Expect(err).To(BeNil())
			Expect(findDataflow.ID).To(Equal(dataflow.ID))
		})
	})

	var _ = When("ListDataflows", func() {
		It("retrieves many Dataflows.", func() {
			repository1 := models.Repository{
				IntegrationID:  primitive.NewObjectID(),
				ExternalID:     externalID,
				NamespacedName: "foobar/foobar",
				DefaultBranch:  "main",
				URI:            "https://gitlab.com/foobar/foobar",
			}
			pipeline1 := models.Pipeline{
				IntegrationID:  primitive.NewObjectID(),
				ExternalID:     externalID,
				NamespacedName: "foobar/foobar",
				DefaultBranch:  "main",
				URI:            "https://gitlab.com/foobar/foobar/-/pipelines",
			}
			deployment1 := models.Deployment{
				IntegrationID: primitive.NewObjectID(),
			}
			dataflow1 := models.Dataflow{
				Repository: repository1,
				Pipeline:   pipeline1,
				Deployment: deployment1,
			}
			_ = daos.CreateDataflow(ctx, &dataflow1)
			Expect(dataflow1.ID).To(Not(BeNil()))

			repository2 := models.Repository{
				IntegrationID:  primitive.NewObjectID(),
				ExternalID:     externalID,
				NamespacedName: "fizz/buzz",
				DefaultBranch:  "main",
				URI:            "https://gitlab.com/foobar/foobar",
			}
			pipeline2 := models.Pipeline{
				IntegrationID:  primitive.NewObjectID(),
				ExternalID:     externalID,
				NamespacedName: "fizz/buzz",
				DefaultBranch:  "main",
				URI:            "https://gitlab.com/foobar/foobar/-/pipelines",
			}
			deployment2 := models.Deployment{
				IntegrationID: primitive.NewObjectID(),
			}
			dataflow2 := models.Dataflow{
				Repository: repository2,
				Pipeline:   pipeline2,
				Deployment: deployment2,
			}
			_ = daos.CreateDataflow(ctx, &dataflow2)
			Expect(dataflow2.ID).To(Not(BeNil()))

			var findDataflows []models.Dataflow
			err := daos.ListDataflows(ctx, &findDataflows)
			Expect(err).To(BeNil())
			Expect(findDataflows).To(HaveLen(2))
		})
	})

	var _ = When("ListDataflowsByFilter", func() {
		It("retrieves many Dataflows conforming to a filter.", func() {
			repository1 := models.Repository{
				IntegrationID:  primitive.NewObjectID(),
				ExternalID:     externalID,
				NamespacedName: "foobar/foobar",
				DefaultBranch:  "main",
				URI:            "https://gitlab.com/foobar/foobar",
			}
			pipeline1 := models.Pipeline{
				IntegrationID:  primitive.NewObjectID(),
				ExternalID:     externalID,
				NamespacedName: "foobar/foobar",
				DefaultBranch:  "main",
				URI:            "https://gitlab.com/foobar/foobar/-/pipelines",
			}
			deployment1 := models.Deployment{
				IntegrationID: primitive.NewObjectID(),
			}
			dataflow1 := models.Dataflow{
				Repository: repository1,
				Pipeline:   pipeline1,
				Deployment: deployment1,
			}
			_ = daos.CreateDataflow(ctx, &dataflow1)
			Expect(dataflow1.ID).To(Not(BeNil()))

			repository2 := models.Repository{
				IntegrationID:  primitive.NewObjectID(),
				ExternalID:     externalID,
				NamespacedName: "fizz/buzz",
				DefaultBranch:  "main",
				URI:            "https://gitlab.com/foobar/foobar",
			}
			pipeline2 := models.Pipeline{
				IntegrationID:  primitive.NewObjectID(),
				ExternalID:     externalID,
				NamespacedName: "fizz/buzz",
				DefaultBranch:  "main",
				URI:            "https://gitlab.com/foobar/foobar/-/pipelines",
			}
			deployment2 := models.Deployment{
				IntegrationID: primitive.NewObjectID(),
			}
			dataflow2 := models.Dataflow{
				Repository: repository2,
				Pipeline:   pipeline2,
				Deployment: deployment2,
			}
			_ = daos.CreateDataflow(ctx, &dataflow2)
			Expect(dataflow2.ID).To(Not(BeNil()))

			var findDataflows []models.Dataflow
			err := daos.ListDataflows(ctx, &findDataflows)
			Expect(err).To(BeNil())
			Expect(findDataflows).To(HaveLen(2))
		})
	})

	var _ = When("UpdateDataflow", func() {
		It("updates an Dataflow.", func() {
			repository := models.Repository{
				IntegrationID:  primitive.NewObjectID(),
				ExternalID:     externalID,
				NamespacedName: "foobar/foobar",
				DefaultBranch:  "main",
				URI:            "https://gitlab.com/foobar/foobar",
			}
			pipeline := models.Pipeline{
				IntegrationID:  primitive.NewObjectID(),
				ExternalID:     externalID,
				NamespacedName: "foobar/foobar",
				DefaultBranch:  "main",
				URI:            "https://gitlab.com/foobar/foobar/-/pipelines",
			}
			deployment := models.Deployment{
				IntegrationID: primitive.NewObjectID(),
			}
			dataflow := models.Dataflow{
				Repository: repository,
				Pipeline:   pipeline,
				Deployment: deployment,
			}
			err := daos.CreateDataflow(ctx, &dataflow)
			Expect(err).To(BeNil())
			Expect(dataflow.ID).To(Not(BeEmpty()))

			newDeployment := models.Deployment{
				IntegrationID: primitive.NewObjectID(),
			}
			updateDataflow := models.Dataflow{
				Repository: repository,
				Pipeline:   pipeline,
				Deployment: newDeployment,
			}

			err = daos.UpdateDataflow(ctx, dataflow.ID, &updateDataflow)
			Expect(err).To(BeNil())
			Expect(updateDataflow.Deployment.IntegrationID).To(Equal(newDeployment.IntegrationID))
		})
	})

	var _ = When("DeleteOne", func() {
		It("deletes a document with ID in a collection", func() {
			repository := models.Repository{
				IntegrationID:  primitive.NewObjectID(),
				ExternalID:     externalID,
				NamespacedName: "foobar/foobar",
				DefaultBranch:  "main",
				URI:            "https://gitlab.com/foobar/foobar",
			}
			pipeline := models.Pipeline{
				IntegrationID:  primitive.NewObjectID(),
				ExternalID:     externalID,
				NamespacedName: "foobar/foobar",
				DefaultBranch:  "main",
				URI:            "https://gitlab.com/foobar/foobar/-/pipelines",
			}
			deployment := models.Deployment{
				IntegrationID: primitive.NewObjectID(),
			}
			dataflow := models.Dataflow{
				Repository: repository,
				Pipeline:   pipeline,
				Deployment: deployment,
			}
			err := daos.CreateDataflow(ctx, &dataflow)
			Expect(err).To(BeNil())
			Expect(dataflow.ID).To(Not(BeEmpty()))

			err = daos.DeleteDataflow(ctx, dataflow.ID)
			Expect(err).To(BeNil())

			var findDataflow models.Dataflow
			err = daos.GetDataflow(ctx, dataflow.ID, &findDataflow)
			Expect(err).To(Not(BeNil()))
		})
	})
})
