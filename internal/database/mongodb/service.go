package mongodb

import (
	"context"
	"errors"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Service provides the functionality to use
type Service struct {
	Client *mongo.Client
	DB     *mongo.Database
}

// NewService creates a new service providing a Mongos.db connection
func NewService() *Service {
	return &Service{}
}

// Connect establishes a connection to a MongoDB instance
func (s *Service) Connect(ctx context.Context, databaseName string) error {
	conn, err := ConnectionString()
	if err != nil {
		return err
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conn))
	if err != nil {
		return fmt.Errorf("could not establish connection to database: %s", err.Error())
	}
	db := client.Database(databaseName)

	s.Client = client
	s.DB = db

	return nil
}

// Disconnect removes a connection to a MongoDB instance
func (s *Service) Disconnect(ctx context.Context) error {
	err := s.Client.Disconnect(ctx)
	return err
}

// ConnectionString creates the connection string with the full URI for a Mongos.db instance.
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

// InsertOne inserts a document into a collection.
func (s *Service) InsertOne(ctx context.Context, collection string, v any) error {
	coll := s.DB.Collection(collection)

	insertOneResult, err := coll.InsertOne(ctx, v)
	if err != nil {
		return err
	}

	err = s.FindOneByID(ctx, collection, insertOneResult.InsertedID.(primitive.ObjectID).Hex(), v)

	return err
}

// Find finds many documents in a collection.
func (s *Service) Find(ctx context.Context, collection string, filter bson.M, vs any) error {
	coll := s.DB.Collection(collection)

	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return err
	}

	err = cursor.All(ctx, vs)

	return err
}

// FindOne finds a documents in a collection.
func (s *Service) FindOne(ctx context.Context, collection string, filter bson.M, v any) error {
	coll := s.DB.Collection(collection)

	findOneResult := coll.FindOne(ctx, filter)
	if findOneResult.Err() != nil {
		return findOneResult.Err()
	}

	err := findOneResult.Decode(v)

	return err
}

// FindOneByID finds a documents with a specific ID in a collection.
func (s *Service) FindOneByID(ctx context.Context, collection string, ID string, v any) error {
	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	err = s.FindOne(ctx, collection, filter, v)

	return err
}

// UpdateOne updates a document in a collection.
func (s *Service) UpdateOne(ctx context.Context, collection string, ID string, v any) error {
	coll := s.DB.Collection(collection)

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
func (s *Service) DeleteOne(ctx context.Context, collection string, ID string) error {
	coll := s.DB.Collection(collection)

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
