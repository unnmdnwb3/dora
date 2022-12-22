package api

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/unnmdnwb3/dora/internal/api/handler"
)

// SetupRouter initializes the router and all routes to be served
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// route for prometheus metrics
	router.GET("/healthz", handler.Healthz)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// routes for repositories
	router.GET("/api/v1/repositories", handler.GetRepositories)

	// routes for integrations
	router.POST("/api/v1/integrations", handler.CreateIntegration)
	router.GET("/api/v1/integrations", handler.ListIntegrations)
	router.GET("/api/v1/integrations/:id", handler.GetIntegration)
	router.PUT("/api/v1/integrations/:id", handler.UpdateIntegration)
	router.DELETE("/api/v1/integrations/:id", handler.DeleteIntegration)

	// routes for dataflows
	router.POST("/api/v1/dataflows", handler.CreateDataflow)
	router.GET("/api/v1/dataflows", handler.ListDataflows)
	router.GET("/api/v1/dataflows/:id", handler.GetDataflow)
	router.PUT("/api/v1/dataflows/:id", handler.UpdateDataflow)
	router.DELETE("/api/v1/dataflows/:id", handler.DeleteDataflow)

	// routes for metrics
	router.POST("/api/v1/metrics/deployment-frequency", handler.DeploymentFrequency)

	return router
}
