package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unnmdnwb3/dora/internal/daos"
	"github.com/unnmdnwb3/dora/internal/models"
	"github.com/unnmdnwb3/dora/internal/services/metrics"
)

// DeploymentFrequency retrieves the deployment frequency of a Dataflow.
func DeploymentFrequency(c *gin.Context) {
	ctx := c.Request.Context()

	var request models.MetricsRequest
	err := c.ShouldBind(&request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var dataflow models.Dataflow
	err = daos.GetDataflow(ctx, request.DataflowID, &dataflow)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	deploymentFrequency, err := metrics.DeploymentFrequency(ctx, dataflow.ID, request.StartDate, request.EndDate, request.Window)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, deploymentFrequency)
	return
}

// GeneralDeploymentFrequency retrieves the general deployment frequency.
func GeneralDeploymentFrequency(c *gin.Context) {
	ctx := c.Request.Context()

	var request models.GeneralMetricsRequest
	err := c.ShouldBind(&request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	deploymentFrequency, err := metrics.GeneralDeploymentFrequency(ctx, request.StartDate, request.EndDate, request.Window)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, deploymentFrequency)
	return
}
