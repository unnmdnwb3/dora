package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unnmdnwb3/dora/internal/daos"
	"github.com/unnmdnwb3/dora/internal/models"
)

func CreateApplication(c *gin.Context) {
	ctx := c.Request.Context()
	var integration models.Application
	err := c.ShouldBind(&integration)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	integrationDAO, err := daos.NewApplication(&ctx)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	response, err := integrationDAO.Create(&integration)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	c.JSON(http.StatusOK, response)
}

func GetApplications(c *gin.Context) {
	ctx := c.Request.Context()
	integrationDAO, err := daos.NewApplication(&ctx)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	response, err := integrationDAO.ReadAll()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	
	c.JSON(http.StatusOK, response)
}

func GetApplication(c *gin.Context) {
	ctx := c.Request.Context()
	var integration models.Application
	err := c.BindUri(&integration)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	integrationDAO, err := daos.NewApplication(&ctx)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	response, err := integrationDAO.Read(integration.Id)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	
	c.JSON(http.StatusOK, response)
}

// TODO get id from uri instead of form
func UpdateApplication(c *gin.Context) {
	ctx := c.Request.Context()
	var integration models.Application
	err := c.ShouldBind(&integration)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	integrationDAO, err := daos.NewApplication(&ctx)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	
	response, err := integrationDAO.Update(&integration)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	
	c.JSON(http.StatusOK, response)
}

func DeleteApplication(c *gin.Context) {
	ctx := c.Request.Context()
	var integration models.Application
	err := c.BindUri(&integration)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	integrationDAO, err := daos.NewApplication(&ctx)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	response, err := integrationDAO.Delete(integration.Id)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	
	c.JSON(http.StatusOK, response)
}