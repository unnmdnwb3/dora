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
		service.Connect(ctx, "dora_test")
	})

	var _ = AfterEach(func() {
		os.Remove("MONGODB_URI")
		os.Remove("MONGODB_PORT")
		os.Remove("MONGODB_USER")
		os.Remove("MONGODB_PASSWORD")

		service.DB.Drop(ctx)
		defer service.Disconnect(ctx)
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
			application := models.Application{
				Auth: "bearertoken",
				Type: "gitlab",
				URI:  "https://gitlab.com",
			}
			err := service.InsertOne(ctx, "application", &application)
			Expect(err).To(BeNil())
			Expect(application.ID).To(Not(BeNil()))
		})
	})

	var _ = When("Find", func() {
		It("finds many documents in a collection", func() {
			application1 := models.Application{
				Auth: "bearertoken",
				Type: "gitlab",
				URI:  "https://gitlab.onprem1.com",
			}
			application2 := models.Application{
				Auth: "bearertoken",
				Type: "gitlab",
				URI:  "https://gitlab.onprem2.com",
			}
			application3 := models.Application{
				Auth: "bearertoken",
				Type: "gitlab",
				URI:  "https://gitlab.onprem.com",
			}
			service.InsertOne(ctx, "application", &application1)
			service.InsertOne(ctx, "application", &application2)
			service.InsertOne(ctx, "application", &application3)
			Expect(application1.ID).To(Not(BeNil()))
			Expect(application2.ID).To(Not(BeNil()))
			Expect(application3.ID).To(Not(BeNil()))

			var findApplications []models.Application
			filter := bson.M{"type": "gitlab"}
			err := service.Find(ctx, "application", filter, &findApplications)
			Expect(err).To(BeNil())
			Expect(findApplications).To(HaveLen(3))
		})
	})

	var _ = When("FindOne", func() {
		It("finds a specific document in a collection", func() {
			application := models.Application{
				Auth: "bearertoken",
				Type: "gitlab",
				URI:  "https://gitlab.com",
			}
			service.InsertOne(ctx, "application", &application)
			Expect(application.ID).To(Not(BeNil()))

			var findApplication models.Application
			filter := bson.M{"uri": "https://gitlab.com"}
			err := service.FindOne(ctx, "application", filter, &findApplication)
			Expect(err).To(BeNil())
			Expect(findApplication.ID).To(Equal(application.ID))
		})
	})

	var _ = When("FindOneByID", func() {
		It("finds a specific document with ID in a collection", func() {
			application := models.Application{
				Auth: "bearertoken",
				Type: "gitlab",
				URI:  "https://gitlab.com",
			}
			service.InsertOne(ctx, "application", &application)
			Expect(application.ID).To(Not(BeNil()))

			var findApplication models.Application
			err := service.FindOneByID(ctx, "application", application.ID, &findApplication)
			Expect(err).To(BeNil())
			Expect(findApplication.ID).To(Equal(application.ID))
		})
	})

	var _ = When("UpdateOne", func() {
		It("updates a document with ID in a collection", func() {
			application := models.Application{
				Auth: "bearertoken",
				Type: "gitlab",
				URI:  "https://gitlab.com",
			}
			service.InsertOne(ctx, "application", &application)
			Expect(application.ID).To(Not(BeNil()))

			updateApplication := models.Application{
				Auth: "newbearertoken",
				Type: "gitlab",
				URI:  "https://gitlab.com",
			}
			err := service.UpdateOne(ctx, "application", application.ID, &updateApplication)
			Expect(err).To(BeNil())
			Expect(updateApplication.Auth).To(Equal("newbearertoken"))
		})
	})

	var _ = When("DeleteOne", func() {
		It("deletes a document with ID in a collection", func() {
			application := models.Application{
				Auth: "bearertoken",
				Type: "gitlab",
				URI:  "https://gitlab.com",
			}
			service.InsertOne(ctx, "application", &application)
			Expect(application.ID).To(Not(BeNil()))

			err := service.DeleteOne(ctx, "application", application.ID)
			Expect(err).To(BeNil())

			var findApplication models.Application
			err = service.FindOneByID(ctx, "application", application.ID, &findApplication)
			Expect(err).To(Not(BeNil()))
		})
	})
})
