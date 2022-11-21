package models

// Application describes any integrated tool to gather data from
type Application struct {
	// TODO redo the `binding:"required"` part, split up in request and response types
	ID   string `json:"id" bson:"_id,omitempty" form:"id" uri:"id"`
	Auth string `json:"auth" bson:"auth" form:"auth"`
	Type string `json:"type" bson:"type" form:"type"`
	URI  string `json:"uri" bson:"uri" form:"uri"`
}
