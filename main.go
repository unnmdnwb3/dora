package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unnmdnwb3/dora/internal/api/handler"
	"github.com/unnmdnwb3/dora/internal/database/mongodb"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	ctx := context.Background()
    client, err := mongodb.NewClient(&ctx)
    if err != nil {
		log.Fatalln(err)
	}

    err = client.Ping(context.TODO(), readpref.Primary())
    if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Successfully connected to database.")
	client.Disconnect(ctx)

    router := gin.Default()

    router.GET("/api/v1/projects", handler.ReadAllProjects)

	router.POST("/api/integrations", handler.CreateIntegration)
    router.GET("/api/integrations", handler.ReadIntegration)
	router.PUT("/api/integrations", handler.UpdateIntegration)
	router.DELETE("/api/integrations", handler.DeleteIntegration)

    log.Println("The server is running and listening on localhost")
    err = http.ListenAndServe(":8080", router)
    if err != nil {
        log.Fatalln("The server encountered a fatal error: ", err)
    }
}