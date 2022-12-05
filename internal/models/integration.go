package models

// Integration represents an third-party integration
type Integration struct {
	ID          string `bson:"_id,omitempty" json:"id"`
	Type        string `bson:"type" json:"type"` // TODO change this to an enum
	Provider    string `bson:"provider" json:"provider"`
	URI         string `bson:"uri" json:"uri"`
	BearerToken string `bson:"bearer_token" json:"bearer_token"`
}
