package daos_test

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/unnmdnwb3/dora/internal/database/mongodb"
)

func TestDAOS(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "daos Suite")
}

var _ = BeforeEach(func() {
	_ = godotenv.Load("./../../../../test/.env")
})

var _ = AfterEach(func() {
	ctx := context.Background()
	service := mongodb.NewService()
	service.Connect(ctx, os.Getenv("MONGODB_DATABASE"))
	service.DB.Drop(ctx)
	defer service.Disconnect(ctx)

	os.Remove("MONGODB_URI")
	os.Remove("MONGODB_PORT")
	os.Remove("MONGODB_USER")
	os.Remove("MONGODB_PASSWORD")
	os.Remove("MONGODB_DATABASE")
})
