package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unnmdnwb3/dora/internal/daos"
	"github.com/unnmdnwb3/dora/internal/models"
	"github.com/unnmdnwb3/dora/internal/services/metrics"
)

// LeadTimeForChanges retrieves the lead time for changes of a Dataflow.
func LeadTimeForChanges(c *gin.Context) {
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

	leadTimeForChanges, err := metrics.LeadTimeForChanges(ctx, dataflow.ID, request.StartDate, request.EndDate, request.Window)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, leadTimeForChanges)
	return
}

// GeneralLeadTimeForChanges retrieves the lead time for changes of a Dataflow.
func GeneralLeadTimeForChanges(c *gin.Context) {
	ctx := c.Request.Context()

	var request models.GeneralMetricsRequest
	err := c.ShouldBind(&request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	leadTimeForChanges, err := metrics.GeneralLeadTimeForChanges(ctx, request.StartDate, request.EndDate, request.Window)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, leadTimeForChanges)
	return
}
