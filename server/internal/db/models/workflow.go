package models

import "go.mongodb.org/mongo-driver/v2/bson"

type Workflow struct {
	ID              bson.ObjectID          `bson:"_id,omitempty" json:"id"`
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	Tags            []string               `json:"tags"`
	Metadata        map[string]interface{} `json:"metadata"`
	ImportPolicyId  bson.ObjectID          `json:"importPolicyId"`
	DestStorageUnit *StorageUnit           `json:"destStorageUnit"`
	UpdatedAt       int64                  `json:"updatedAt"`
	UpdatedBy       string                 `json:"updatedBy"`
	CreatedAt       int64                  `json:"createdAt"`
	CreatedBy       string                 `json:"createdBy"`
}
