package mongodb

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// Client provides the driver functionality for MongoDB.
	// It gets initialized in the main function as a singleton.
	Client *mongo.Client

	// DB provides the connectivity to a MongoDB instance.
	// It gets initialized in the main function as a singleton.
	DB *mongo.Database
)

const (
	// Timeout defines the timeout for database.
	Timeout = 5 * time.Second

	// DefaultDatabase defines the default database name.
	DefaultDatabase = "dora"
)

// ConnectionString creates the connection string with the full URI for a MongoDB instance.
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

// Init a Client and DB connection to MongoDB as singletons.
func Init(ctx *context.Context) error {
	conn, err := ConnectionString()
	if err != nil {
		return err
	}

	Client, err = mongo.Connect(*ctx, options.Client().ApplyURI(conn))
	if err != nil {
		return fmt.Errorf("could not establish connection to database: %s", err.Error())
	}

	DB = Client.Database(DefaultDatabase)

	return nil
}

// InsertOne inserts a document into a collection.
func InsertOne(ctx context.Context, collection string, v any) error {
	coll := DB.Collection(collection)
	insertOneResult, err := coll.InsertOne(ctx, v)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": insertOneResult.InsertedID}
	applicationResult := coll.FindOne(ctx, filter)
	if applicationResult.Err() != nil {
		return err
	}

	err = applicationResult.Decode(v)
	return err
}

// Find finds many documents in a collection.
func Find(ctx context.Context, collection string, filter bson.M, vs any) error {
	coll := DB.Collection(collection)

	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return err
	}

	err = cursor.All(ctx, vs)
	return err
}

// FindOne finds a documents in a collection.
func FindOne(ctx context.Context, collection string, filter bson.M, v any) error {
	coll := DB.Collection(collection)

	findOneResult := coll.FindOne(ctx, filter)
	if findOneResult.Err() != nil {
		return findOneResult.Err()
	}

	err := findOneResult.Decode(v)
	return err
}

// FindOneByID finds a documents with a specific ID in a collection.
func FindOneByID(ctx context.Context, collection string, ID string, v any) error {
	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	err = FindOne(ctx, collection, filter, v)

	return err
}

// UpdateOne updates a document in a collection.
func UpdateOne(ctx context.Context, collection string, ID string, v any) error {
	coll := DB.Collection(collection)

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": v}
	updateOneResult, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if updateOneResult.MatchedCount == 0 {
		return fmt.Errorf("id for update not found: %s", ID)
	}

	return nil
}

// DeleteOne deletes a document in a collection.
func DeleteOne(ctx context.Context, collection string, ID string) error {
	coll := DB.Collection(collection)

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	deleteResult, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if deleteResult.DeletedCount == 0 {
		return fmt.Errorf("id for delete not found: %s", ID)
	}

	return nil
}
