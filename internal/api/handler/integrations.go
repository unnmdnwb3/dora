package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unnmdnwb3/dora/internal/daos"
	"github.com/unnmdnwb3/dora/internal/models"
)

// CreateIntegration creates a new Integration.
func CreateIntegration(c *gin.Context) {
	ctx := c.Request.Context()

	var integration models.Integration
	err := c.ShouldBind(&integration)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	err = daos.CreateIntegration(ctx, &integration)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	c.JSON(http.StatusOK, integration)
}

// GetIntegration retrieves a Integration.
func GetIntegration(c *gin.Context) {
	ctx := c.Request.Context()

	var params models.Params
	err := c.BindUri(&params)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	var integration models.Integration
	err = daos.GetIntegration(ctx, params.ID, &integration)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	c.JSON(http.StatusOK, integration)
}

// ListIntegrations retrieves many Integrations.
func ListIntegrations(c *gin.Context) {
	ctx := c.Request.Context()

	var integrations []models.Integration
	err := daos.ListIntegrations(ctx, &integrations)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	c.JSON(http.StatusOK, integrations)
}

// UpdateIntegration update a Integration.
func UpdateIntegration(c *gin.Context) {
	ctx := c.Request.Context()

	var params models.Params
	err := c.BindUri(&params)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	var integration models.Integration
	err = c.ShouldBind(&integration)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	err = daos.UpdateIntegration(ctx, params.ID, &integration)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	integration.ID = params.ID
	c.JSON(http.StatusOK, integration)
}

// DeleteIntegration deletes a Integration.
func DeleteIntegration(c *gin.Context) {
	ctx := c.Request.Context()

	var params models.Params
	err := c.BindUri(&params)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	err = daos.DeleteIntegration(ctx, params.ID)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	c.JSON(http.StatusOK, params)
}
