package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unnmdnwb3/dora/internal/applications/gitlab"
)

func GetRepositories(c *gin.Context) {
	client, err := gitlab.NewClient()
	if err != nil {
		log.Fatalln(err.Error())
	}

	repositories, err := client.GetRepositories()
	if err != nil {
		log.Fatalln(err.Error())
	}
	
	c.IndentedJSON(http.StatusOK, repositories)
}