package mongodb_test

import (
	"context"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/unnmdnwb3/dora/internal/database/mongodb"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "mongodb.Client Suite")
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

var _ = Describe("mongodb.Client", func() {
	var _ = When("ConnectionString", func() {
		It("can build the correct connection string", func() {
			expectedConn := "mongodb://user:password@127.0.0.1:27017"
			Expect(mongodb.ConnectionString()).To(Equal(expectedConn))
		})
	})

	var _ = When("Init", func() {
		It("creates a new MongoDB client with connection to a database", func() {
			ctx := context.Background()
			err := mongodb.Init(&ctx)
			defer mongodb.Client.Disconnect(ctx)

			Expect(err).To(BeNil())
			Expect(mongodb.Client).To(Not(BeNil()))
			Expect(mongodb.DB).To(Not(BeNil()))
			Expect(mongodb.Client.Ping(ctx, readpref.Primary())).To(BeNil())
		})
	})
})
