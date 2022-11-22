package mongodb

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// Client provides all functionality for MongoDB
	// and gets instatiated in the main function
	Client *mongo.Client

	// DB points to the standard database
	DB *mongo.Database
)

const (
	// Timeout holds the general TimeOut
	Timeout = 5 * time.Second

	// StandardDatabase defines the standard database MongoDB needs to connect to
	StandardDatabase = "dora"
)

// ConnectionString creates the URI needed to connect to a mongodb instance
func ConnectionString() (string, error) {
	// necessary
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return "", errors.New("could not find env: MONGODB_URI")
	}

	// optional
	auth := ""
	user := os.Getenv("MONGODB_USER")
	password := os.Getenv("MONGODB_PASSWORD")
	if user != "" && password != "" {
		auth = fmt.Sprintf("%s:%s@", user, password)
	}
	port := os.Getenv("MONGODB_PORT")
	if port != "" {
		port = fmt.Sprintf(":%s", port)
	}

	return fmt.Sprintf("mongodb://%s%s%s", auth, uri, port), nil
}

// Init creates a new mongodb client and pointer to the database
func Init(ctx *context.Context) error {
	conn, err := ConnectionString()
	if err != nil {
		return err
	}

	Client, err = mongo.Connect(*ctx, options.Client().ApplyURI(conn))
	if err != nil {
		return fmt.Errorf("could not establish connection to database: %s", err.Error())
	}

	DB = Client.Database(StandardDatabase)

	return nil
}
