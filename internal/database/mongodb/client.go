package mongodb

import (
	"context"
	"errors"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const Database = "dora"

func NewClient(ctx *context.Context) (*mongo.Client, error) {
	conn, err := ConnectionString()
	if err != nil {
		return nil, err
	}

	client, err := mongo.Connect(*ctx, options.Client().ApplyURI(conn))
	if err != nil {
		return nil, fmt.Errorf("could not establish connection to database: %s", err.Error())
	}

	return client, nil
}

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