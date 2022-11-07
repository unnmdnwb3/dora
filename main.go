package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "github.com/unnmdnwb3/dora/internal/api/v1"
)

func main() {
    router := gin.Default()
    router.GET("/api/v1/projects", v1.GetProjects)

    log.Println("The server is running and listening on localhost")
    err := http.ListenAndServe(":8080", router)
    if err != nil {
        log.Fatalln("The server encountered a fatal error: ", err)
    }
}