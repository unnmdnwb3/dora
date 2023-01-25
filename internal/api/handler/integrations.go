package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unnmdnwb3/dora/internal/daos"
	"github.com/unnmdnwb3/dora/internal/models"
	"github.com/unnmdnwb3/dora/internal/utils/types"
)

// CreateIntegration creates a new Integration.
func CreateIntegration(c *gin.Context) {
	ctx := c.Request.Context()

	var integration models.Integration
	err := c.ShouldBind(&integration)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = daos.CreateIntegration(ctx, &integration)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, integration)
	return
}

// GetIntegration retrieves a Integration.
func GetIntegration(c *gin.Context) {
	ctx := c.Request.Context()

	var params models.Params
	err := c.BindUri(&params)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var integration models.Integration
	integrationID, err := types.StringToObjectID(params.ID)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = daos.GetIntegration(ctx, integrationID, &integration)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, integration)
	return
}

// ListIntegrations retrieves many Integrations.
func ListIntegrations(c *gin.Context) {
	ctx := c.Request.Context()

	var integrations []models.Integration
	err := daos.ListIntegrations(ctx, &integrations)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, integrations)
	return
}

// UpdateIntegration update a Integration.
func UpdateIntegration(c *gin.Context) {
	ctx := c.Request.Context()

	var params models.Params
	err := c.BindUri(&params)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var integration models.Integration
	err = c.ShouldBind(&integration)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	integrationID, err := types.StringToObjectID(params.ID)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = daos.UpdateIntegration(ctx, integrationID, &integration)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	integration.ID = integrationID
	c.JSON(http.StatusOK, integration)
}

// DeleteIntegration deletes a Integration.
func DeleteIntegration(c *gin.Context) {
	ctx := c.Request.Context()

	var params models.Params
	err := c.BindUri(&params)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	integrationID, err := types.StringToObjectID(params.ID)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = daos.DeleteIntegration(ctx, integrationID)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, params)
	return
}
