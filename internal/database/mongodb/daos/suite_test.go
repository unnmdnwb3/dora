package daos_test

import (
	"context"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/unnmdnwb3/dora/internal/database/mongodb"
)

func TestDeployRuns(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "daos Suite")
}

var _ = BeforeEach(func() {
	// TODO remove config from code
	os.Setenv("MONGODB_URI", "127.0.0.1")
	os.Setenv("MONGODB_PORT", "27017")
	os.Setenv("MONGODB_USER", "user")
	os.Setenv("MONGODB_PASSWORD", "password")

	ctx := context.Background()
	mongodb.Init(&ctx)
})

var _ = AfterEach(func() {
	// TODO remove config from code
	os.Remove("MONGODB_URI")
	os.Remove("MONGODB_PORT")
	os.Remove("MONGODB_USER")
	os.Remove("MONGODB_PASSWORD")

	ctx := context.Background()
	defer mongodb.Client.Disconnect(ctx)
})
