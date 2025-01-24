package db

import (
	"context"
	"fmt"

	"github.com/uploadpilot/uploadpilot/internal/msg"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UploadLogsRepo struct {
	collectionName string
}

func NewUploadLogsRepo() *UploadLogsRepo {
	return &UploadLogsRepo{
		collectionName: "uploadlogs",
	}
}

func (u *UploadLogsRepo) GetLogs(ctx context.Context, uploadID string) ([]bson.M, error) {
	id, err := primitive.ObjectIDFromHex(uploadID)
	if err != nil {
		return nil, fmt.Errorf(msg.InvalidObjectID, uploadID)
	}

	collection := db.Collection(u.collectionName)
	var logs []bson.M

	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}}).SetProjection(bson.M{"uploadId": 0, "workspaceId": 0})
	cursor, err := collection.Find(ctx, bson.M{"uploadId": id}, opts)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &logs); err != nil {
		return nil, err
	}

	return logs, nil
}

func (u *UploadLogsRepo) BatchAddLogs(ctx context.Context, logs []interface{}) error {
	collection := db.Collection(u.collectionName)
	_, err := collection.InsertMany(context.Background(), logs)
	if err != nil {
		return err
	}
	return nil
}
