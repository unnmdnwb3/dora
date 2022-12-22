package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Healthz returns the health status of the service.
func Healthz(c *gin.Context) {
	c.Status(http.StatusOK)
}
