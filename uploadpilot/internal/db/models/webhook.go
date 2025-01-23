package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WebhookEvent string

const (
	WebhookEventFileUploadFailed WebhookEvent = "File upload failed"
	WebhookEventFileUploaded     WebhookEvent = "File uploaded"
	WebhookEventFileDeleted      WebhookEvent = "File deleted"
)

type Webhook struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	URL           string             `bson:"url" json:"url" validate:"required"`
	Event         WebhookEvent       `bson:"event" json:"event" validate:"required"`
	Method        string             `bson:"method" json:"method" validate:"required"`
	SigningSecret string             `bson:"signingSecret" json:"signingSecret"`
	WorkspaceID   primitive.ObjectID `bson:"workspaceId" json:"workspaceId" validate:"required"`
	Enabled       bool               `bson:"enabled" json:"enabled" validate:"required"`
	CreatedAt     primitive.DateTime `bson:"createdAt" json:"createdAt"`
	CreatedBy     string             `bson:"createdBy" json:"createdBy"`
	UpdatedAt     primitive.DateTime `bson:"updatedAt" json:"updatedAt"`
	UpdatedBy     string             `bson:"updatedBy" json:"updatedBy"`
}
