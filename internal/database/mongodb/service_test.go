package mongodb_test

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/unnmdnwb3/dora/internal/database/mongodb"
	"github.com/unnmdnwb3/dora/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "mongodb.Service Suite")
}

var _ = Describe("mongodb.Service", func() {
	ctx := context.Background()
	var service *mongodb.Service

	var _ = BeforeEach(func() {
		_ = godotenv.Load("./../../../test/.env")

		service = mongodb.NewService()
		service.Connect(ctx, os.Getenv("MONGODB_DATABASE"))
	})

	var _ = AfterEach(func() {
		service.DB.Drop(ctx)
		defer service.Disconnect(ctx)

		os.Remove("MONGODB_URI")
		os.Remove("MONGODB_PORT")
		os.Remove("MONGODB_USER")
		os.Remove("MONGODB_PASSWORD")
	})

	var _ = When("ConnectionString", func() {
		It("can build a connection string", func() {
			conn := "mongodb://user:password@127.0.0.1:27017"
			Expect(mongodb.ConnectionString()).To(Equal(conn))
		})
	})

	var _ = When("Connect", func() {
		It("establishes a new connection to a MongoDB instance ", func() {
			Expect(service.Client.Ping(ctx, readpref.Primary())).To(BeNil())
		})
	})

	var _ = When("InsertOne", func() {
		It("creates a new document in a collection", func() {
			integration := models.Integration{
				Type:        "sc",
				Provider:    "gitlab",
				BearerToken: "bearertoken",
				URI:         "https://gitlab.com",
			}
			err := service.InsertOne(ctx, "integration", &integration)
			Expect(err).To(BeNil())
			Expect(integration.ID).To(Not(BeEmpty()))
		})
	})

	var _ = When("Find", func() {
		It("finds many documents in a collection", func() {
			integration1 := models.Integration{
				Type:        "sc",
				Provider:    "gitlab",
				BearerToken: "bearertoken",
				URI:         "https://gitlab.onprem.com",
			}
			integration2 := models.Integration{
				Type:        "sc",
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
			_ = service.InsertOne(ctx, "integration", &integration1)
			_ = service.InsertOne(ctx, "integration", &integration2)
			_ = service.InsertOne(ctx, "integration", &integration3)
			Expect(integration1.ID).To(Not(BeNil()))
			Expect(integration2.ID).To(Not(BeNil()))
			Expect(integration3.ID).To(Not(BeNil()))

			var findIntegrations []models.Integration
			filter := bson.M{"type": "sc"}
			err := service.Find(ctx, "integration", filter, &findIntegrations)
			Expect(err).To(BeNil())
			Expect(findIntegrations).To(HaveLen(3))
		})
	})

	var _ = When("FindOne", func() {
		It("finds a specific document in a collection", func() {
			integration := models.Integration{
				Type:        "sc",
				Provider:    "gitlab",
				BearerToken: "bearertoken",
				URI:         "https://gitlab.com",
			}
			service.InsertOne(ctx, "integration", &integration)
			Expect(integration.ID).To(Not(BeNil()))

			var findIntegration models.Integration
			filter := bson.M{"uri": "https://gitlab.com"}
			err := service.FindOne(ctx, "integration", filter, &findIntegration)
			Expect(err).To(BeNil())
			Expect(findIntegration.ID).To(Equal(integration.ID))
		})
	})

	var _ = When("FindOneByID", func() {
		It("finds a specific document with ID in a collection", func() {
			integration := models.Integration{
				Type:        "sc",
				Provider:    "gitlab",
				BearerToken: "bearertoken",
				URI:         "https://gitlab.com",
			}
			service.InsertOne(ctx, "integration", &integration)
			Expect(integration.ID).To(Not(BeNil()))

			var findIntegration models.Integration
			err := service.FindOneByID(ctx, "integration", integration.ID, &findIntegration)
			Expect(err).To(BeNil())
			Expect(findIntegration.ID).To(Equal(integration.ID))
		})
	})

	var _ = When("UpdateOne", func() {
		It("updates a document with ID in a collection", func() {
			integration := models.Integration{
				Type:        "sc",
				Provider:    "gitlab",
				BearerToken: "bearertoken",
				URI:         "https://gitlab.com",
			}
			service.InsertOne(ctx, "integration", &integration)
			Expect(integration.ID).To(Not(BeNil()))

			updateIntegration := models.Integration{
				Type:        "sc",
				Provider:    "gitlab",
				BearerToken: "newbearertoken",
				URI:         "https://gitlab.com",
			}
			err := service.UpdateOne(ctx, "integration", integration.ID, &updateIntegration)
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
			service.InsertOne(ctx, "integration", &integration)
			Expect(integration.ID).To(Not(BeNil()))

			err := service.DeleteOne(ctx, "integration", integration.ID)
			Expect(err).To(BeNil())

			var findIntegration models.Integration
			err = service.FindOneByID(ctx, "integration", integration.ID, &findIntegration)
			Expect(err).To(Not(BeNil()))
		})
	})
})
