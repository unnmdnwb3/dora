package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unnmdnwb3/dora/internal/daos"
	"github.com/unnmdnwb3/dora/internal/models"
)

func CreateIntegration(c *gin.Context) {
	ctx := c.Request.Context()
	var integration models.Integration
	err := c.ShouldBind(&integration)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	integrationDAO, err := daos.NewIntegration(&ctx)
	response, err := integrationDAO.Create(&integration)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	c.JSON(http.StatusOK, response)
}

func ReadAllIntegrations(c *gin.Context) {
	ctx := c.Request.Context()
	integrationDAO, err := daos.NewIntegration(&ctx)
	response, err := integrationDAO.ReadAll()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	
	c.JSON(http.StatusOK, response)
}

func ReadIntegration(c *gin.Context) {
	ctx := c.Request.Context()
	var integration models.Integration
	err := c.BindUri(&integration)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	integrationDAO, err := daos.NewIntegration(&ctx)
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
func UpdateIntegration(c *gin.Context) {
	ctx := c.Request.Context()
	var integration models.Integration
	err := c.ShouldBind(&integration)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	integrationDAO, err := daos.NewIntegration(&ctx)
	response, err := integrationDAO.Update(&integration)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	
	c.JSON(http.StatusOK, response)
}

func DeleteIntegration(c *gin.Context) {
	ctx := c.Request.Context()
	var integration models.Integration
	err := c.BindUri(&integration)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	integrationDAO, err := daos.NewIntegration(&ctx)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	response, err := integrationDAO.Delete(integration.Id)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	
	c.JSON(http.StatusOK, response)
}