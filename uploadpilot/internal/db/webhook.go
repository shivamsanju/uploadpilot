package db

import (
	"context"
	"fmt"
	"time"

	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/msg"
	"github.com/uploadpilot/uploadpilot/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WebhookRepo struct {
	collectionName string
}

func NewWebhookRepo() *WebhookRepo {
	return &WebhookRepo{collectionName: "webhooks"}
}

func (wr *WebhookRepo) GetAll(ctx context.Context, workspaceID string) ([]models.Webhook, error) {
	id, err := primitive.ObjectIDFromHex(workspaceID)
	if err != nil {
		return nil, fmt.Errorf(msg.InvalidObjectID, workspaceID)
	}

	collection := db.Collection(wr.collectionName)
	var webhooks []models.Webhook

	opts := options.Find().SetProjection(bson.M{"signingSecret": 0}).SetSort(bson.D{{Key: "updatedAt", Value: -1}})
	cursor, err := collection.Find(ctx, bson.M{"workspaceId": id}, opts)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &webhooks); err != nil {
		return nil, err
	}

	return webhooks, nil
}

func (wr *WebhookRepo) Get(ctx context.Context, workspaceID, webhookID string) (*models.Webhook, error) {
	wsID, err := primitive.ObjectIDFromHex(workspaceID)
	if err != nil {
		return nil, fmt.Errorf(msg.InvalidObjectID, workspaceID)
	}

	id, err := primitive.ObjectIDFromHex(webhookID)
	if err != nil {
		return nil, fmt.Errorf(msg.InvalidObjectID, webhookID)
	}

	collection := db.Collection(wr.collectionName)
	var webhook models.Webhook

	opts := options.FindOne().SetProjection(bson.M{"signingSecret": 0})
	err = collection.FindOne(ctx, bson.M{"workspaceId": wsID, "_id": id}, opts).Decode(&webhook)
	if err != nil {
		return nil, err
	}

	return &webhook, nil
}

func (wr *WebhookRepo) Create(ctx context.Context, workspaceID string, webhook *models.Webhook) error {
	wsID, err := primitive.ObjectIDFromHex(workspaceID)
	if err != nil {
		return fmt.Errorf(msg.InvalidObjectID, workspaceID)
	}
	user, err := utils.GetUserDetailsFromContext(ctx)
	if err != nil {
		return err
	}

	webhook.CreatedBy = user.Email
	webhook.UpdatedBy = user.Email
	webhook.Enabled = true
	webhook.WorkspaceID = wsID
	webhook.ID = primitive.NewObjectID()
	webhook.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	webhook.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	collection := db.Collection(wr.collectionName)
	_, err = collection.InsertOne(ctx, webhook)
	if err != nil {
		return err
	}

	return nil
}

func (wr *WebhookRepo) Delete(ctx context.Context, workspaceID, webhookID string) error {
	wsID, err := primitive.ObjectIDFromHex(workspaceID)
	if err != nil {
		return fmt.Errorf(msg.InvalidObjectID, workspaceID)
	}

	id, err := primitive.ObjectIDFromHex(webhookID)
	if err != nil {
		return fmt.Errorf(msg.InvalidObjectID, webhookID)
	}

	collection := db.Collection(wr.collectionName)
	if _, err := collection.DeleteOne(ctx, bson.M{"workspaceId": wsID, "_id": id}); err != nil {
		return err
	}

	return nil
}

func (wr *WebhookRepo) Update(ctx context.Context, workspaceID string, webhookID string, webhook *models.Webhook) error {
	wsID, err := primitive.ObjectIDFromHex(workspaceID)
	if err != nil {
		return fmt.Errorf(msg.InvalidObjectID, workspaceID)
	}

	id, err := primitive.ObjectIDFromHex(webhookID)
	if err != nil {
		return fmt.Errorf(msg.InvalidObjectID, webhookID)
	}

	user, err := utils.GetUserDetailsFromContext(ctx)
	if err != nil {
		return err
	}

	collection := db.Collection(wr.collectionName)
	webhook.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	webhook.UpdatedBy = user.Email

	if _, err := collection.UpdateOne(ctx, bson.M{"workspaceId": wsID, "_id": id}, bson.M{"$set": webhook}); err != nil {
		return err
	}

	return nil
}

func (wr *WebhookRepo) Patch(ctx context.Context, workspaceID, webhookID string, patch bson.M) error {
	wsID, err := primitive.ObjectIDFromHex(workspaceID)
	if err != nil {
		return fmt.Errorf(msg.InvalidObjectID, workspaceID)
	}

	id, err := primitive.ObjectIDFromHex(webhookID)
	if err != nil {
		return fmt.Errorf(msg.InvalidObjectID, webhookID)
	}

	user, err := utils.GetUserDetailsFromContext(ctx)
	if err != nil {
		return err
	}

	patch["updatedBy"] = user.Email
	patch["updatedAt"] = primitive.NewDateTimeFromTime(time.Now())

	collection := db.Collection(wr.collectionName)
	_, err = collection.UpdateOne(ctx, bson.M{"workspaceId": wsID, "_id": id}, bson.M{"$set": patch})
	return err
}

func (wr *WebhookRepo) GetEnabledWebhooksWithSecret(ctx context.Context, workspaceID string) ([]models.Webhook, error) {
	id, err := primitive.ObjectIDFromHex(workspaceID)
	if err != nil {
		return nil, fmt.Errorf(msg.InvalidObjectID, workspaceID)
	}

	collection := db.Collection(wr.collectionName)
	var webhooks []models.Webhook

	cursor, err := collection.Find(ctx, bson.M{"workspaceId": id, "enabled": true}, nil)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &webhooks); err != nil {
		return nil, err
	}

	return webhooks, nil
}
