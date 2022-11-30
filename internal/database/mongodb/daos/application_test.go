package daos_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/unnmdnwb3/dora/internal/database/mongodb/daos"
	"github.com/unnmdnwb3/dora/internal/models"
)

var _ = Describe("daos.Application", func() {
	ctx := context.Background()

	var _ = When("Create", func() {
		It("creates an application", func() {
			DAO, err := daos.NewApplication(&ctx)
			Expect(err).To((BeNil()))

			createRequest := models.Application{
				Auth: "bearertoken",
				Type: "gitlab",
				URI:  "https://gitlab.com",
			}

			createResponse, err := DAO.Create(&createRequest)
			Expect(err).To(BeNil())
			Expect(createResponse.Auth).To(Equal(createRequest.Auth))
			Expect(createResponse.Type).To(Equal(createRequest.Type))
			Expect(createResponse.URI).To(Equal(createRequest.URI))
		})
	})

	var _ = When("Read", func() {
		It("reads an application", func() {
			DAO, err := daos.NewApplication(&ctx)
			Expect(err).To((BeNil()))

			createRequest := models.Application{
				Auth: "bearertoken",
				Type: "gitlab",
				URI:  "https://gitlab.com",
			}

			createResponse, err := DAO.Create(&createRequest)
			Expect(err).To((BeNil()))

			readResponse, err := DAO.Read(createResponse.ID)
			Expect(err).To(BeNil())
			Expect(readResponse.ID).To(Equal(createResponse.ID))
			Expect(readResponse.Auth).To(Equal(createResponse.Auth))
			Expect(readResponse.Type).To(Equal(createResponse.Type))
			Expect(readResponse.URI).To(Equal(createResponse.URI))
		})
	})

	var _ = When("ReadAll", func() {
		It("reads all applications", func() {
			DAO, err := daos.NewApplication(&ctx)
			Expect(err).To((BeNil()))

			firstCreateRequest := models.Application{
				Auth: "bearertoken",
				Type: "gitlab",
				URI:  "https://gitlab.com",
			}
			secondCreateRequest := models.Application{
				Auth: "bearertoken",
				Type: "github",
				URI:  "https://github.com",
			}

			_, err = DAO.Create(&firstCreateRequest)
			Expect(err).To(BeNil())
			_, err = DAO.Create(&secondCreateRequest)
			Expect(err).To(BeNil())

			readAllResponse, err := DAO.ReadAll()
			Expect(err).To((BeNil()))
			Expect(len(*readAllResponse)).To(BeNumerically(">", 0))
		})
	})

	var _ = When("Update", func() {
		It("updates an application", func() {
			DAO, err := daos.NewApplication(&ctx)
			Expect(err).To((BeNil()))

			createRequest := models.Application{
				Auth: "bearertoken",
				Type: "gitlab",
				URI:  "https://gitlab.com",
			}

			createResponse, err := DAO.Create(&createRequest)
			Expect(err).To((BeNil()))

			authUpdate := "newbearertoken"
			updateRequest := models.Application{
				ID:   createResponse.ID,
				Auth: authUpdate,
				Type: createResponse.Type,
				URI:  createResponse.URI,
			}

			updateResponse, err := DAO.Update(&updateRequest)
			Expect(err).To(BeNil())
			Expect(updateResponse.Auth).To(Equal(updateRequest.Auth))
			Expect(updateResponse.Type).To(Equal(updateRequest.Type))
			Expect(updateResponse.URI).To(Equal(updateRequest.URI))
		})
	})

	var _ = When("Delete", func() {
		It("deletes an application", func() {
			DAO, err := daos.NewApplication(&ctx)
			Expect(err).To((BeNil()))

			createRequest := models.Application{
				Auth: "bearertoken",
				Type: "gitlab",
				URI:  "https://gitlab.com",
			}

			createResponse, err := DAO.Create(&createRequest)
			Expect(err).To((BeNil()))

			deleteRequest := createResponse.ID
			deleteResponse, err := DAO.Delete(deleteRequest)
			Expect(err).To(BeNil())
			Expect(deleteResponse.ID).To(Equal(deleteRequest))
			Expect(deleteResponse.Auth).To(Equal(createResponse.Auth))
			Expect(deleteResponse.Type).To(Equal(createResponse.Type))
			Expect(deleteResponse.URI).To(Equal(createResponse.URI))
		})
	})
})
