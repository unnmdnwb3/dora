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

	// connect to database
	err := mongodb.Init(&ctx)
	if err != nil {
		log.Fatalln(err)
	}

	err = mongodb.Client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Successfully connected to database.")
	defer mongodb.Client.Disconnect(ctx)

	router := gin.Default()

	// routes for repositories
	router.GET("/api/v1/repositories", handler.GetRepositories)

	// routes for applications
	router.POST("/api/applications", handler.CreateApplication)
	router.GET("/api/applications", handler.GetApplications)
	router.GET("/api/applications/:id", handler.GetApplication)
	router.PUT("/api/applications", handler.UpdateApplication)
	router.DELETE("/api/applications/:id", handler.DeleteApplication)

	log.Println("\nThe server is running and listening on localhost! ")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalln("The server encountered a fatal error: 🚀", err)
	}
}
