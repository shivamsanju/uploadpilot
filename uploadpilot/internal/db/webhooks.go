package db

import (
	"context"
	"time"

	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WebhookRepo interface {
	GetWebhooks(ctx context.Context, workspaceID primitive.ObjectID) ([]models.Webhook, error)
	GetEnabledWebhooksWithSecret(ctx context.Context, workspaceID primitive.ObjectID) ([]models.Webhook, error)
	GetWebhook(ctx context.Context, workspaceID primitive.ObjectID, webhookID primitive.ObjectID) (*models.Webhook, error)
	CreateWebhook(ctx context.Context, workspaceID primitive.ObjectID, webhook *models.Webhook) (*models.Webhook, error)
	UpdateWebhook(ctx context.Context, workspaceID primitive.ObjectID, webhook *models.Webhook, updatedBy string) (*models.Webhook, error)
	DeleteWebhook(ctx context.Context, workspaceID primitive.ObjectID, webhookID primitive.ObjectID) error
	PatchWebhook(ctx context.Context, workspaceID primitive.ObjectID, webhookID primitive.ObjectID, patch *bson.M) error
}

type webhookRepo struct {
	collectionName string
}

func NewWebhookRepo() WebhookRepo {
	return &webhookRepo{collectionName: "webhooks"}
}

func (wr *webhookRepo) GetWebhooks(ctx context.Context, workspaceID primitive.ObjectID) ([]models.Webhook, error) {
	collection := db.Collection(wr.collectionName)
	var webhooks []models.Webhook
	opts := options.Find().SetProjection(bson.M{"signingSecret": 0}).SetSort(bson.D{{Key: "updatedAt", Value: -1}})
	cursor, err := collection.Find(ctx, bson.M{"workspaceId": workspaceID}, opts)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &webhooks)
	if err != nil {
		return nil, err
	}
	return webhooks, nil
}

func (wr *webhookRepo) GetEnabledWebhooksWithSecret(ctx context.Context, workspaceID primitive.ObjectID) ([]models.Webhook, error) {
	collection := db.Collection(wr.collectionName)
	var webhooks []models.Webhook
	cursor, err := collection.Find(ctx, bson.M{"workspaceId": workspaceID, "enabled": true}, nil)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &webhooks)
	if err != nil {
		return nil, err
	}
	return webhooks, nil
}

func (wr *webhookRepo) GetWebhook(ctx context.Context, workspaceID primitive.ObjectID, webhookID primitive.ObjectID) (*models.Webhook, error) {
	collection := db.Collection(wr.collectionName)
	var webhook models.Webhook
	opts := options.FindOne().SetProjection(bson.M{"signingSecret": 0})
	err := collection.FindOne(ctx, bson.M{"workspaceId": workspaceID, "_id": webhookID}, opts).Decode(&webhook)
	return &webhook, err
}

func (wr *webhookRepo) CreateWebhook(ctx context.Context, workspaceID primitive.ObjectID, webhook *models.Webhook) (*models.Webhook, error) {
	webhook.ID = primitive.NewObjectID()
	webhook.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	webhook.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	collection := db.Collection(wr.collectionName)
	_, err := collection.InsertOne(ctx, webhook)
	if err != nil {
		return nil, err
	}

	return webhook, nil
}

func (wr *webhookRepo) DeleteWebhook(ctx context.Context, workspaceID primitive.ObjectID, webhookID primitive.ObjectID) error {
	collection := db.Collection(wr.collectionName)
	_, err := collection.DeleteOne(ctx, bson.M{"workspaceId": workspaceID, "_id": webhookID})
	return err
}

func (wr *webhookRepo) UpdateWebhook(ctx context.Context, workspaceID primitive.ObjectID, webhook *models.Webhook, updatedBy string) (*models.Webhook, error) {
	collection := db.Collection(wr.collectionName)
	webhook.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	webhook.UpdatedBy = updatedBy

	_, err := collection.UpdateOne(ctx, bson.M{"workspaceId": workspaceID, "_id": webhook.ID}, bson.M{"$set": webhook})
	if err != nil {
		return nil, err
	}

	return webhook, nil
}

func (wr *webhookRepo) PatchWebhook(ctx context.Context, workspaceID primitive.ObjectID, webhookID primitive.ObjectID, patch *bson.M) error {
	updatedBy := ctx.Value("email").(string)
	patchReq := *patch
	patchReq["updatedBy"] = updatedBy
	patchReq["updatedAt"] = primitive.NewDateTimeFromTime(time.Now())

	collection := db.Collection(wr.collectionName)
	_, err := collection.UpdateOne(ctx, bson.M{"workspaceId": workspaceID, "_id": webhookID}, bson.M{
		"$set": patchReq,
	})
	return err
}
