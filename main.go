package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/unnmdnwb3/dora/internal/api"
	"github.com/unnmdnwb3/dora/internal/database/mongodb"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	ctx := context.Background()

	service := mongodb.NewService()
	database := os.Getenv("MONGODB_DATABASE")
	err := service.Connect(ctx, database)
	if err != nil {
		log.Fatalln("Could not connect to database: ", err.Error())
	}

	err = service.Client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Fatalln("Could not ping database: ", err.Error())
	}

	fmt.Println("Successfully connected to database.")
	defer service.Disconnect(ctx)

	router := api.SetupRouter()

	log.Println("\nThe server is running and listening on localhost! 🚀")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalln("The server encountered a fatal error: ", err.Error())
	}
}
