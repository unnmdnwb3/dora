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

func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "mongodb.Client Suite")
}

var _ = BeforeEach(func() {
	_ = godotenv.Load("./../../../test/.env")

	ctx := context.Background()
	mongodb.Init(&ctx)
})

var _ = AfterEach(func() {
	os.Remove("MONGODB_URI")
	os.Remove("MONGODB_PORT")
	os.Remove("MONGODB_USER")
	os.Remove("MONGODB_PASSWORD")

	ctx := context.Background()
	mongodb.DB.Drop(ctx)
	defer mongodb.Client.Disconnect(ctx)
})

var _ = Describe("mongodb.Client", func() {
	ctx := context.Background()

	var _ = When("ConnectionString", func() {
		It("can build a connection string", func() {
			conn := "mongodb://user:password@127.0.0.1:27017"
			Expect(mongodb.ConnectionString()).To(Equal(conn))
		})
	})

	var _ = When("Init", func() {
		It("creates a new client and connection to a MongoDB instance ", func() {
			ctx := context.Background()

			Expect(mongodb.Client).To(Not(BeNil()))
			Expect(mongodb.DB).To(Not(BeNil()))
			Expect(mongodb.Client.Ping(ctx, readpref.Primary())).To(BeNil())
		})
	})

	var _ = When("InsertOne", func() {
		It("creates a new document in a collection", func() {
			application := models.Application{
				Auth: "bearertoken",
				Type: "gitlab",
				URI:  "https://gitlab.com",
			}
			err := mongodb.InsertOne(ctx, "application", &application)
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
			mongodb.InsertOne(ctx, "application", &application1)
			mongodb.InsertOne(ctx, "application", &application2)
			mongodb.InsertOne(ctx, "application", &application3)
			Expect(application1.ID).To(Not(BeNil()))
			Expect(application2.ID).To(Not(BeNil()))
			Expect(application3.ID).To(Not(BeNil()))

			var findApplications []models.Application
			filter := bson.M{"type": "gitlab"}
			err := mongodb.Find(ctx, "application", filter, &findApplications)
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
			mongodb.InsertOne(ctx, "application", &application)
			Expect(application.ID).To(Not(BeNil()))

			var findApplication models.Application
			filter := bson.M{"uri": "https://gitlab.com"}
			err := mongodb.FindOne(ctx, "application", filter, &findApplication)
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
			mongodb.InsertOne(ctx, "application", &application)
			Expect(application.ID).To(Not(BeNil()))

			var findApplication models.Application
			err := mongodb.FindOneByID(ctx, "application", application.ID, &findApplication)
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
			mongodb.InsertOne(ctx, "application", &application)
			Expect(application.ID).To(Not(BeNil()))

			updateApplication := models.Application{
				Auth: "newbearertoken",
				Type: "gitlab",
				URI:  "https://gitlab.com",
			}
			err := mongodb.UpdateOne(ctx, "application", application.ID, &updateApplication)
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
			mongodb.InsertOne(ctx, "application", &application)
			Expect(application.ID).To(Not(BeNil()))

			err := mongodb.DeleteOne(ctx, "application", application.ID)
			Expect(err).To(BeNil())

			var findApplication models.Application
			err = mongodb.FindOneByID(ctx, "application", application.ID, &findApplication)
			Expect(err).To(Not(BeNil()))
		})
	})
})
