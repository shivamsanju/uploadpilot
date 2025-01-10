package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ImportPolicy struct {
	ID               primitive.ObjectID `bson:"_id" json:"id" validate:"required"`
	Name             string             `json:"name" validate:"required,min=2,max=100"`
	AllowedMimeTypes []string           `json:"allowedMimeTypes" validate:"required"`
	AllowedSources   []string           `json:"allowedSources" validate:"required"`
	MaxFileSizeKb    int64              `json:"maxFileSizeKb" validate:"required,min=1"`
	MaxFileCount     int64              `json:"maxFileCount" validate:"required,min=1"`
	CreatedAt        primitive.DateTime `json:"createdAt"`
	CreatedBy        string             `json:"createdBy"`
	UpdatedAt        primitive.DateTime `json:"updatedAt"`
	UpdatedBy        string             `json:"updatedBy"`
}
