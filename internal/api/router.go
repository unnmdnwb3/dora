package api

import (
	"github.com/gin-gonic/gin"
	"github.com/unnmdnwb3/dora/internal/api/handler"
)

// SetupRouter initializes the router and all routes to be served
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// routes for repositories
	router.GET("/api/v1/repositories", handler.GetRepositories)

	// routes for integrations
	router.POST("/api/v1/integrations", handler.CreateIntegration)
	router.GET("/api/v1/integrations", handler.ListIntegrations)
	router.GET("/api/v1/integrations/:id", handler.GetIntegration)
	router.PUT("/api/v1/integrations/:id", handler.UpdateIntegration)
	router.DELETE("/api/v1/integrations/:id", handler.DeleteIntegration)

	return router
}
