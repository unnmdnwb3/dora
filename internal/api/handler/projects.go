package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unnmdnwb3/dora/internal/integrations/gitlab"
)

func ReadAllProjects(c *gin.Context) {
	client, err := gitlab.NewClient()
	if err != nil {
		log.Fatalln(err.Error())
	}

	projects, err := client.GetProjects()
	if err != nil {
		log.Fatalln(err.Error())
	}
	
	c.IndentedJSON(http.StatusOK, projects)
}