package types

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// StringToObjectID converts a string to a primitive.ObjectID
func StringToObjectID(s string) (primitive.ObjectID, error) {
	objectID, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		return primitive.ObjectID{}, fmt.Errorf("%s: %s", err.Error(), s)
	}
	return objectID, nil
}
