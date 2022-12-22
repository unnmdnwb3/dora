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

	var request models.DeploymentFrequencyRequest
	err := c.ShouldBind(&request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	var dataflow models.Dataflow
	err = daos.GetDataflow(ctx, request.DataflowID, &dataflow)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	deploymentFrequency, err := metrics.CalculateDeploymentFrequency(ctx, dataflow.ID, request.Window, request.EndDate)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	c.JSON(http.StatusOK, deploymentFrequency)
}
