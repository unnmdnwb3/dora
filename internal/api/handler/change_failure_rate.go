package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unnmdnwb3/dora/internal/daos"
	"github.com/unnmdnwb3/dora/internal/models"
	"github.com/unnmdnwb3/dora/internal/services/metrics"
)

// ChangeFailureRate retrieves the change failure rate of a Dataflow.
func ChangeFailureRate(c *gin.Context) {
	ctx := c.Request.Context()

	var request models.MetricsRequest
	err := c.ShouldBind(&request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	var dataflow models.Dataflow
	err = daos.GetDataflow(ctx, request.DataflowID, &dataflow)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	changeFailureRate, err := metrics.ChangeFailureRate(ctx, dataflow.ID, request.StartDate, request.EndDate, request.Window)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	c.JSON(http.StatusOK, changeFailureRate)
}

// GeneralChangeFailureRate retrieves the change failure rate of a Dataflow.
func GeneralChangeFailureRate(c *gin.Context) {
	ctx := c.Request.Context()

	var request models.GeneralMetricsRequest
	err := c.ShouldBind(&request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	changeFailureRate, err := metrics.GeneralChangeFailureRate(ctx, request.StartDate, request.EndDate, request.Window)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	c.JSON(http.StatusOK, changeFailureRate)
}
