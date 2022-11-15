package daos_test

import (
	"context"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/unnmdnwb3/dora/internal/daos"
	"github.com/unnmdnwb3/dora/internal/models"
)

var ctx context.Context
var integrationDao daos.Integration

func TestIntegrations(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Daos Suite")
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

var _ = Describe("Integration Daos", func() {
	It("creates an integration", func() {
		integrationDao, err := daos.NewIntegration(&ctx)
		Expect(err).To((BeNil()))

		createRequest := models.Integration{
			Auth: "bearertoken",
			Type: "gitlab",
			Uri: "https://gitlab.com",
		}

		createResponse, err := integrationDao.Create(&createRequest)
		Expect(err).To((BeNil()))
		Expect(createResponse.Auth).To((Equal(createRequest.Auth)))
		Expect(createResponse.Type).To((Equal(createRequest.Type)))
		Expect(createResponse.Uri).To((Equal(createRequest.Uri)))
	})

	It("reads all integrations", func() {
		integrationDao, err := daos.NewIntegration(&ctx)
		Expect(err).To((BeNil()))

		readAllResponse, err := integrationDao.ReadAll()
		Expect(err).To((BeNil()))
		Expect(len(*readAllResponse)).To(BeNumerically(">", 0))
	})

	It("updates an integration", func() {
		integrationDao, err := daos.NewIntegration(&ctx)
		Expect(err).To((BeNil()))

		updateRequest := models.Integration{
			Auth: "newbearertoken",
			Type: "gitlab",
			Uri: "https://gitlab.com",
		}

		updateResponse, err := integrationDao.Update(&updateRequest)
		Expect(err).To(BeNil())
		Expect(updateResponse.Auth).To(Equal(updateRequest.Auth))
		Expect(updateResponse.Type).To(Equal(updateRequest.Type))
		Expect(updateResponse.Uri).To(Equal(updateRequest.Uri))
	})

	It("deletes an integration", func() {
		integrationDao, err := daos.NewIntegration(&ctx)
		Expect(err).To((BeNil()))

		createRequest := models.Integration{
			Auth: "bearertoken",
			Type: "gitlab",
			Uri: "https://gitlab.com",
		}

		createResponse, err := integrationDao.Create(&createRequest)
		Expect(err).To((BeNil()))

		deleteRequest := createResponse.Id
		deleteResponse, err := integrationDao.Delete(deleteRequest)
		Expect(err).To(BeNil())
		Expect(deleteResponse.Id).To(Equal(deleteRequest))
		Expect(deleteResponse.Auth).To(Equal(createResponse.Auth))
		Expect(deleteResponse.Type).To(Equal(createResponse.Type))
		Expect(deleteResponse.Uri).To(Equal(createResponse.Uri))
		Expect(err).To(BeNil())
	})
})