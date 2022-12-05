package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unnmdnwb3/dora/internal/database/mongodb"
	"github.com/unnmdnwb3/dora/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

const collection = "integration"

// CreateIntegration creates a new Integration.
func CreateIntegration(c *gin.Context) {
	ctx := c.Request.Context()
	service := mongodb.NewService()
	err := service.Connect(ctx, mongodb.DefaultDatabase)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	defer service.Disconnect(ctx)

	var integration models.Integration
	err = c.ShouldBind(&integration)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	err = service.InsertOne(ctx, collection, &integration)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	c.JSON(http.StatusOK, integration)
}

// GetIntegration retrieves a Integration.
func GetIntegration(c *gin.Context) {
	ctx := c.Request.Context()
	service := mongodb.NewService()
	err := service.Connect(ctx, mongodb.DefaultDatabase)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	defer service.Disconnect(ctx)

	var params models.Params
	err = c.BindUri(&params)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	var integration models.Integration
	err = service.FindOneByID(ctx, collection, params.ID, &integration)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	c.JSON(http.StatusOK, integration)
}

// ListIntegrations retrieves many Integrations.
func ListIntegrations(c *gin.Context) {
	ctx := c.Request.Context()
	service := mongodb.NewService()
	err := service.Connect(ctx, mongodb.DefaultDatabase)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	defer service.Disconnect(ctx)

	var integrations []models.Integration
	err = service.Find(ctx, collection, bson.M{}, &integrations)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	c.JSON(http.StatusOK, integrations)
}

// UpdateIntegration update a Integration.
func UpdateIntegration(c *gin.Context) {
	ctx := c.Request.Context()
	service := mongodb.NewService()
	err := service.Connect(ctx, mongodb.DefaultDatabase)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	defer service.Disconnect(ctx)

	var params models.Params
	err = c.BindUri(&params)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	var integration models.Integration
	err = c.ShouldBind(&integration)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	err = service.UpdateOne(ctx, collection, params.ID, &integration)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	integration.ID = params.ID
	c.JSON(http.StatusOK, integration)
}

// DeleteIntegration deletes a Integration.
func DeleteIntegration(c *gin.Context) {
	ctx := c.Request.Context()
	service := mongodb.NewService()
	err := service.Connect(ctx, mongodb.DefaultDatabase)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	defer service.Disconnect(ctx)

	var params models.Params
	err = c.BindUri(&params)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	err = service.DeleteOne(ctx, collection, params.ID)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	c.JSON(http.StatusOK, params)
}
