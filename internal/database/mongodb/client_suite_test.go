package mongodb_test

import (
	"context"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/unnmdnwb3/dora/internal/database/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var ctx context.Context
var client *mongo.Client


func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MongoDB Client Suite")
}

var _ = BeforeSuite(func() {
	os.Setenv("MONGODB_URI", "127.0.0.1")
	os.Setenv("MONGODB_PORT", "27017")
	os.Setenv("MONGODB_USER", "user")
	os.Setenv("MONGODB_PASSWORD", "password")
})
  
var _ = AfterSuite(func() {
	os.Remove("MONGODB_URI")
	os.Remove("MONGODB_PORT")
	os.Remove("MONGODB_USER")
	os.Remove("MONGODB_PASSWORD")
})

var _ = Describe("MongoDB", func() {
	It("can build the correct connection string", func() {
		expectedConn := "mongodb://user:password@127.0.0.1:27017"
		Expect(mongodb.ConnectionString()).To(Equal(expectedConn))
	})

	It("can create a new MongoDB client with connection", func() {
		client, err := mongodb.NewClient(&ctx)
		defer client.Disconnect(ctx)
		
		Expect(err).To(BeNil())
		Expect(client.Ping(ctx, readpref.Primary())).To(BeNil())
	})
})
