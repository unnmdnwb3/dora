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

	// routes for applications
	router.POST("/api/applications", handler.CreateApplication)
	router.GET("/api/applications", handler.GetApplications)
	router.GET("/api/applications/:id", handler.GetApplication)
	router.PUT("/api/applications", handler.UpdateApplication)
	router.DELETE("/api/applications/:id", handler.DeleteApplication)

	return router
}
