package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Workflow struct {
	ID             primitive.ObjectID     `bson:"_id" json:"id" validate:"required"`
	Name           string                 `json:"name" validate:"required,min=2,max=100"`
	Description    string                 `json:"description" validate:"max=500"`
	Tags           []string               `json:"tags"`
	Metadata       map[string]interface{} `json:"metadata"`
	ImportPolicyID primitive.ObjectID     `json:"importPolicyId" validate:"required"`
	DataStoreID    primitive.ObjectID     `json:"dataStoreId" validate:"required"`
	CreatedAt      primitive.DateTime     `json:"createdAt"`
	CreatedBy      string                 `json:"createdBy"`
	UpdatedAt      primitive.DateTime     `json:"updatedAt"`
	UpdatedBy      string                 `json:"updatedBy"`
}
