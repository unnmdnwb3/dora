package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unnmdnwb3/dora/internal/daos"
	"github.com/unnmdnwb3/dora/internal/models"
	"github.com/unnmdnwb3/dora/internal/services/metrics"
)

// MeanTimeToRestore retrieves the mean time to restore of a Dataflow.
func MeanTimeToRestore(c *gin.Context) {
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

	meanTimeToRestore, err := metrics.MeanTimeToRestore(ctx, dataflow.ID, request.StartDate, request.EndDate, request.Window)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, meanTimeToRestore)
}

// GeneralMeanTimeToRestore retrieves the mean time to restore of a Dataflow.
func GeneralMeanTimeToRestore(c *gin.Context) {
	ctx := c.Request.Context()

	var request models.GeneralMetricsRequest
	err := c.ShouldBind(&request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	meanTimeToRestore, err := metrics.GeneralMeanTimeToRestore(ctx, request.StartDate, request.EndDate, request.Window)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, meanTimeToRestore)
}
