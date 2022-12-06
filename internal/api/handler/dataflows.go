package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unnmdnwb3/dora/internal/database/mongodb/daos"
	"github.com/unnmdnwb3/dora/internal/models"
)

// CreateDataflow creates a new Dataflow.
func CreateDataflow(c *gin.Context) {
	ctx := c.Request.Context()

	var dataflow models.Dataflow
	err := c.ShouldBind(&dataflow)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	err = daos.CreateDataflow(ctx, &dataflow)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	c.JSON(http.StatusOK, dataflow)
}

// GetDataflow retrieves a Dataflow.
func GetDataflow(c *gin.Context) {
	ctx := c.Request.Context()

	var params models.Params
	err := c.BindUri(&params)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	var dataflow models.Dataflow
	err = daos.GetDataflow(ctx, params.ID, &dataflow)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	c.JSON(http.StatusOK, dataflow)
}

// ListDataflows retrieves many Dataflows.
func ListDataflows(c *gin.Context) {
	ctx := c.Request.Context()

	var dataflows []models.Dataflow
	err := daos.ListDataflows(ctx, &dataflows)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	c.JSON(http.StatusOK, dataflows)
}

// UpdateDataflow update a Dataflow.
func UpdateDataflow(c *gin.Context) {
	ctx := c.Request.Context()

	var params models.Params
	err := c.BindUri(&params)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	var dataflow models.Dataflow
	err = c.ShouldBind(&dataflow)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	err = daos.UpdateDataflow(ctx, params.ID, &dataflow)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	dataflow.ID = params.ID
	c.JSON(http.StatusOK, dataflow)
}

// DeleteDataflow deletes a Dataflow.
func DeleteDataflow(c *gin.Context) {
	ctx := c.Request.Context()

	var params models.Params
	err := c.BindUri(&params)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	err = daos.DeleteDataflow(ctx, params.ID)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	c.JSON(http.StatusOK, params)
}
