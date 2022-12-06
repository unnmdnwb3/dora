package daos_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/unnmdnwb3/dora/internal/database/mongodb/daos"
	"github.com/unnmdnwb3/dora/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

var _ = Describe("daos.integration", func() {
	ctx := context.Background()

	var _ = When("CreateIntegration", func() {
		It("creates creates a new Integration.", func() {
			integration := models.Integration{
				Type:        "sc",
				Provider:    "gitlab",
				BearerToken: "bearertoken",
				URI:         "https://gitlab.com",
			}
			err := daos.CreateIntegration(ctx, &integration)
			Expect(err).To(BeNil())
			Expect(integration.ID).To(Not(BeEmpty()))
		})
	})

	var _ = When("GetIntegration", func() {
		It("retrieves an Integration.", func() {
			integration := models.Integration{
				Type:        "sc",
				Provider:    "gitlab",
				BearerToken: "bearertoken",
				URI:         "https://gitlab.com",
			}
			err := daos.CreateIntegration(ctx, &integration)
			Expect(err).To(BeNil())
			Expect(integration.ID).To(Not(BeEmpty()))

			var findIntegration models.Integration
			err = daos.GetIntegration(ctx, integration.ID, &findIntegration)
			Expect(err).To(BeNil())
			Expect(findIntegration.ID).To(Equal(integration.ID))
		})
	})

	var _ = When("ListIntegrations", func() {
		It("retrieves many Integrations.", func() {
			integration1 := models.Integration{
				Type:        "sc",
				Provider:    "gitlab",
				BearerToken: "bearertoken",
				URI:         "https://gitlab.onprem.com",
			}
			integration2 := models.Integration{
				Type:        "sc",
				Provider:    "cicd",
				BearerToken: "bearertoken",
				URI:         "https://gitlab.com",
			}
			integration3 := models.Integration{
				Type:        "sc",
				Provider:    "gihub",
				BearerToken: "bearertoken",
				URI:         "https://github.com",
			}
			_ = daos.CreateIntegration(ctx, &integration1)
			_ = daos.CreateIntegration(ctx, &integration2)
			_ = daos.CreateIntegration(ctx, &integration3)
			Expect(integration1.ID).To(Not(BeNil()))
			Expect(integration2.ID).To(Not(BeNil()))
			Expect(integration3.ID).To(Not(BeNil()))

			var findIntegrations []models.Integration
			err := daos.ListIntegrations(ctx, &findIntegrations)
			Expect(err).To(BeNil())
			Expect(findIntegrations).To(HaveLen(3))
		})
	})

	var _ = When("ListIntegrationsByFilter", func() {
		It("retrieves many Integrations conforming to a filter.", func() {
			integration1 := models.Integration{
				Type:        "sc",
				Provider:    "gitlab",
				BearerToken: "bearertoken",
				URI:         "https://gitlab.onprem.com",
			}
			integration2 := models.Integration{
				Type:        "cicd",
				Provider:    "gitlab",
				BearerToken: "bearertoken",
				URI:         "https://gitlab.com",
			}
			integration3 := models.Integration{
				Type:        "sc",
				Provider:    "github",
				BearerToken: "bearertoken",
				URI:         "https://github.com",
			}
			_ = daos.CreateIntegration(ctx, &integration1)
			_ = daos.CreateIntegration(ctx, &integration2)
			_ = daos.CreateIntegration(ctx, &integration3)
			Expect(integration1.ID).To(Not(BeNil()))
			Expect(integration2.ID).To(Not(BeNil()))
			Expect(integration3.ID).To(Not(BeNil()))

			var findIntegrations []models.Integration
			filter := bson.M{"type": "sc"}
			err := daos.ListIntegrationsByFilter(ctx, filter, &findIntegrations)
			Expect(err).To(BeNil())
			Expect(findIntegrations).To(HaveLen(2))
		})
	})

	var _ = When("UpdateIntegration", func() {
		It("updates an Integration.", func() {
			integration := models.Integration{
				Type:        "sc",
				Provider:    "gitlab",
				BearerToken: "bearertoken",
				URI:         "https://gitlab.com",
			}
			err := daos.CreateIntegration(ctx, &integration)
			Expect(err).To(BeNil())
			Expect(integration.ID).To(Not(BeEmpty()))

			updateIntegration := models.Integration{
				Type:        "sc",
				Provider:    "gitlab",
				BearerToken: "newbearertoken",
				URI:         "https://gitlab.com",
			}
			err = daos.UpdateIntegration(ctx, integration.ID, &updateIntegration)
			Expect(err).To(BeNil())
			Expect(updateIntegration.BearerToken).To(Equal("newbearertoken"))
		})
	})

	var _ = When("DeleteOne", func() {
		It("deletes a document with ID in a collection", func() {
			integration := models.Integration{
				Type:        "sc",
				Provider:    "gitlab",
				BearerToken: "bearertoken",
				URI:         "https://gitlab.com",
			}
			err := daos.CreateIntegration(ctx, &integration)
			Expect(err).To(BeNil())
			Expect(integration.ID).To(Not(BeEmpty()))

			err = daos.DeleteIntegration(ctx, integration.ID)
			Expect(err).To(BeNil())

			var findIntegration models.Integration
			err = daos.GetIntegration(ctx, integration.ID, &findIntegration)
			Expect(err).To(Not(BeNil()))
		})
	})
})
