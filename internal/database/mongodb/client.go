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

const database = "dora"

func NewClient() (*mongo.Client, error) {
	user := os.Getenv("MONGODB_USER")
	if user == "" {
		return nil, errors.New(fmt.Sprintf("Could not find env:  %s", user))
	}

	password := os.Getenv("MONGODB_PASSWORD")
	if password == "" {
		return nil, errors.New(fmt.Sprintf("Could not find env:  %s", password))
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return nil, errors.New(fmt.Sprintf("Could not find env:  %s", uri))
	}

	port := os.Getenv("MONGODB_PORT")
	if port == "" {
		return nil, errors.New(fmt.Sprintf("Could not find env:  %s", port))
	}

	conn := fmt.Sprintf("mongodb://%s:%s@%s:%s", user, password, uri, port)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conn))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not establish connection to database"))
	}

	return client, nil
}
