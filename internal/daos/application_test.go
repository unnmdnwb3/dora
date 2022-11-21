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

func TestApplications(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "daos.Application Suite")
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

var _ = Describe("daos.Application", func() {
	ctx := context.Background()
	
	var _ = When("Create", func() {
		It("creates an application", func() {
			applicationDao, err := daos.NewApplication(&ctx)
			Expect(err).To((BeNil()))

			createRequest := models.Application{
				Auth: "bearertoken",
				Type: "gitlab",
				Uri: "https://gitlab.com",
			}

			createResponse, err := applicationDao.Create(&createRequest)
			Expect(err).To((BeNil()))
			Expect(createResponse.Auth).To((Equal(createRequest.Auth)))
			Expect(createResponse.Type).To((Equal(createRequest.Type)))
			Expect(createResponse.Uri).To((Equal(createRequest.Uri)))
		})
	})

	var _ = When("Read", func() {
		It("reads an application", func() {
			applicationDao, err := daos.NewApplication(&ctx)
			Expect(err).To((BeNil()))

			createRequest := models.Application{
				Auth: "bearertoken",
				Type: "gitlab",
				Uri: "https://gitlab.com",
			}

			createResponse, err := applicationDao.Create(&createRequest)
			Expect(err).To((BeNil()))
			
			readResponse, err := applicationDao.Read(createResponse.Id)
			Expect(err).To((BeNil()))
			Expect(readResponse.Id).To((Equal(createResponse.Id)))
			Expect(readResponse.Auth).To((Equal(createResponse.Auth)))
			Expect(readResponse.Type).To((Equal(createResponse.Type)))
			Expect(readResponse.Uri).To((Equal(createResponse.Uri)))
		})
	})

	var _ = When("ReadAll", func() {
		It("reads all applications", func() {
			applicationDao, err := daos.NewApplication(&ctx)
			Expect(err).To((BeNil()))

			readAllResponse, err := applicationDao.ReadAll()
			Expect(err).To((BeNil()))
			Expect(len(*readAllResponse)).To(BeNumerically(">", 0))
		})
	})

	var _ = When("Update", func() {
		It("updates an application", func() {
			applicationDao, err := daos.NewApplication(&ctx)
			Expect(err).To((BeNil()))

			createRequest := models.Application{
				Auth: "bearertoken",
				Type: "gitlab",
				Uri: "https://gitlab.com",
			}

			createResponse, err := applicationDao.Create(&createRequest)
			Expect(err).To((BeNil()))

			authUpdate := "newbearertoken"
			updateRequest := models.Application{
				Id: createResponse.Id,
				Auth: authUpdate,
				Type: createResponse.Type,
				Uri: createResponse.Uri,
			}

			updateResponse, err := applicationDao.Update(&updateRequest)
			Expect(err).To(BeNil())
			Expect(updateResponse.Auth).To(Equal(updateRequest.Auth))
			Expect(updateResponse.Type).To(Equal(updateRequest.Type))
			Expect(updateResponse.Uri).To(Equal(updateRequest.Uri))
		})
	})

	var _ = When("Delete", func() {
		It("deletes an application", func() {
			applicationDao, err := daos.NewApplication(&ctx)
			Expect(err).To((BeNil()))

			createRequest := models.Application{
				Auth: "bearertoken",
				Type: "gitlab",
				Uri: "https://gitlab.com",
			}

			createResponse, err := applicationDao.Create(&createRequest)
			Expect(err).To((BeNil()))

			deleteRequest := createResponse.Id
			deleteResponse, err := applicationDao.Delete(deleteRequest)
			Expect(err).To(BeNil())
			Expect(deleteResponse.Id).To(Equal(deleteRequest))
			Expect(deleteResponse.Auth).To(Equal(createResponse.Auth))
			Expect(deleteResponse.Type).To(Equal(createResponse.Type))
			Expect(deleteResponse.Uri).To(Equal(createResponse.Uri))
			Expect(err).To(BeNil())
		})
	})
})