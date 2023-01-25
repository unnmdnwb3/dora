package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/unnmdnwb3/dora/internal/api/handler"
)

// httpRequestsTotal is a prometheus counter for all HTTP requests
var httpRequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of requests.",
	},
	[]string{"method", "path", "code"},
)

// registerMetrics registers all metrics to be served
func registerMetrics() {
	prometheus.MustRegister(httpRequestsTotal)
	return
}

// prometheusMiddleware is a middleware to set up all custom prometheus metrics
func prometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		httpRequestsTotal.WithLabelValues(c.Request.Method, c.Request.URL.Path, strconv.Itoa(c.Writer.Status())).Inc()
	}
}

// SetupRouter initializes the router and all routes to be served
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// register prometheus metrics
	registerMetrics()

	// route for prometheus metrics
	router.GET("/healthz", prometheusMiddleware(), handler.Healthz)
	router.GET("/metrics", prometheusMiddleware(), gin.WrapH(promhttp.Handler()))

	// routes for repositories
	router.GET("/api/v1/repositories", prometheusMiddleware(), handler.GetRepositories)

	// routes for integrations
	router.POST("/api/v1/integrations", prometheusMiddleware(), handler.CreateIntegration)
	router.GET("/api/v1/integrations", prometheusMiddleware(), handler.ListIntegrations)
	router.GET("/api/v1/integrations/:id", prometheusMiddleware(), handler.GetIntegration)
	router.PUT("/api/v1/integrations/:id", prometheusMiddleware(), handler.UpdateIntegration)
	router.DELETE("/api/v1/integrations/:id", prometheusMiddleware(), handler.DeleteIntegration)

	// routes for dataflows
	router.POST("/api/v1/dataflows", prometheusMiddleware(), handler.CreateDataflow)
	router.GET("/api/v1/dataflows", prometheusMiddleware(), handler.ListDataflows)
	router.GET("/api/v1/dataflows/:id", prometheusMiddleware(), handler.GetDataflow)
	router.PUT("/api/v1/dataflows/:id", prometheusMiddleware(), handler.UpdateDataflow)
	router.DELETE("/api/v1/dataflows/:id", prometheusMiddleware(), handler.DeleteDataflow)

	// routes for dataflow metrics
	router.POST("/api/v1/metrics/deployment-frequency", prometheusMiddleware(), handler.DeploymentFrequency)
	router.POST("/api/v1/metrics/lead-time-for-changes", prometheusMiddleware(), handler.LeadTimeForChanges)
	router.POST("/api/v1/metrics/mean-time-to-restore", prometheusMiddleware(), handler.MeanTimeToRestore)
	router.POST("/api/v1/metrics/change-failure-rate", prometheusMiddleware(), handler.ChangeFailureRate)

	// routes for general dataflow metrics
	router.POST("/api/v1/metrics/general/deployment-frequency", prometheusMiddleware(), handler.GeneralDeploymentFrequency)
	router.POST("/api/v1/metrics/general/lead-time-for-changes", prometheusMiddleware(), handler.GeneralLeadTimeForChanges)
	router.POST("/api/v1/metrics/general/mean-time-to-restore", prometheusMiddleware(), handler.GeneralMeanTimeToRestore)
	router.POST("/api/v1/metrics/general/change-failure-rate", prometheusMiddleware(), handler.GeneralChangeFailureRate)

	return router
}
