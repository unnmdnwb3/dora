package models

// TODO redo the `binding:"required"` part, split up in request and response types
type Application struct {
	Id string `json:"id" bson:"_id,omitempty" form:"id" uri:"id"`
	Auth string `json:"auth" bson:"auth" form:"auth"`
	Type string `json:"type" bson:"type" form:"type"`
	Uri string `json:"uri" bson:"uri" form:"uri"`
}