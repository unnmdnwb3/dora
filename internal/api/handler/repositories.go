package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unnmdnwb3/dora/internal/connectors/gitlab"
	"github.com/unnmdnwb3/dora/internal/daos"
	"github.com/unnmdnwb3/dora/internal/models"
)

// GetRepositories gets all repositories
func GetRepositories(c *gin.Context) {
	ctx := c.Request.Context()

	var integrations []models.Integration
	err := daos.ListIntegrations(ctx, &integrations)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	allRepositories := []models.Repository{}
	for _, integration := range integrations {
		// TODO: currently only gitlab is supported
		client := gitlab.NewClient(integration.URI, integration.BearerToken)

		repositories, err := client.GetRepositories()
		if err != nil {
			log.Fatalln(err.Error())
		}

		allRepositories = append(allRepositories, *repositories...)
	}

	c.IndentedJSON(http.StatusOK, allRepositories)
}
