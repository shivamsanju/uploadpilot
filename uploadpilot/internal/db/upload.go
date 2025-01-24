package db

import (
	"context"
	"fmt"
	"time"

	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"github.com/uploadpilot/uploadpilot/internal/messages"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UploadRepo struct {
	collectionName string
}

func NewUploadRepo() *UploadRepo {
	return &UploadRepo{
		collectionName: "uploads",
	}
}

func (i *UploadRepo) GetAll(ctx context.Context, workspaceID string, skip int64, limit int64, search string) ([]models.Upload, int64, error) {
	id, err := primitive.ObjectIDFromHex(workspaceID)
	if err != nil {
		return nil, 0, fmt.Errorf(messages.InvalidObjectID, workspaceID)
	}

	filter := bson.M{"workspaceId": id}
	if search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": search, "$options": "i"}},
			{"status": bson.M{"$regex": search, "$options": "i"}},
			{"storedFileName": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	// Count total records matching the filter
	collection := db.Collection(i.collectionName)
	totalRecords, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		infra.Log.Errorf("failed to count imports: %s", err.Error())
		return nil, 0, err
	}

	// Apply pagination and sorting
	opts := options.Find().
		SetSort(bson.D{{Key: "finishedAt", Value: -1}}).
		SetSkip(skip).
		SetLimit(limit)

	// Fetch documents
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		infra.Log.Errorf("failed to find imports: %s", err.Error())
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// Decode documents into the result slice
	var cb []models.Upload
	if err := cursor.All(ctx, &cb); err != nil {
		infra.Log.Errorf("failed to decode imports: %s", err.Error())
		return nil, 0, err
	}

	return cb, totalRecords, nil
}

func (i *UploadRepo) GetAllFilterByMetadata(ctx context.Context, workspaceID string, skip, limit int64, search map[string]string) ([]models.Upload, int64, error) {
	id, err := primitive.ObjectIDFromHex(workspaceID)
	if err != nil {
		return nil, 0, fmt.Errorf(messages.InvalidObjectID, workspaceID)
	}

	filter := bson.M{"workspaceId": id}
	for key, value := range search {
		if key != "" && value != "" {
			filter["metadata."+key] = value
		}
	}

	// Count total records matching the filter
	collection := db.Collection(i.collectionName)
	totalRecords, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		infra.Log.Errorf("failed to count imports: %s", err.Error())
		return nil, 0, err
	}

	// Apply pagination and sorting
	opts := options.Find().
		SetSort(bson.D{{Key: "finishedAt", Value: -1}}).
		SetSkip(skip).
		SetLimit(limit)

	// Fetch documents
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		infra.Log.Errorf("failed to find imports: %s", err.Error())
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// Decode documents into the result slice
	var cb []models.Upload
	if err := cursor.All(ctx, &cb); err != nil {
		infra.Log.Errorf("failed to decode imports: %s", err.Error())
		return nil, 0, err
	}

	return cb, totalRecords, nil
}

func (i *UploadRepo) Get(ctx context.Context, workspaceID, uploadID string) (*models.Upload, error) {
	id, err := primitive.ObjectIDFromHex(uploadID)
	if err != nil {
		return nil, fmt.Errorf(messages.InvalidObjectID, uploadID)
	}
	var cb models.Upload
	collection := db.Collection(i.collectionName)
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&cb)
	if err != nil {
		infra.Log.Errorf("failed to find import: %s", err.Error())
		return nil, err
	}
	return &cb, nil
}

func (i *UploadRepo) Create(ctx context.Context, workspaceID string, upload *models.Upload) error {
	wsID, err := primitive.ObjectIDFromHex(workspaceID)
	if err != nil {
		return fmt.Errorf(messages.InvalidObjectID, workspaceID)
	}

	upload.ID = primitive.NewObjectID()
	upload.WorkspaceID = wsID
	upload.StartedAt = primitive.NewDateTimeFromTime(time.Now())
	upload.Metadata["uploadId"] = upload.ID

	infra.Log.Infof("creating upload: %+v", upload)
	collection := db.Collection(i.collectionName)
	_, err = collection.InsertOne(ctx, upload)
	if err != nil {
		infra.Log.Errorf("failed to create upload: %s", err.Error())
		return err
	}

	return nil
}

func (i *UploadRepo) Update(ctx context.Context, workspaceID, uploadID string, upload *models.Upload) error {
	id, err := primitive.ObjectIDFromHex(uploadID)
	if err != nil {
		return fmt.Errorf(messages.InvalidObjectID, uploadID)
	}

	wsID, err := primitive.ObjectIDFromHex(workspaceID)
	if err != nil {
		return fmt.Errorf(messages.InvalidObjectID, uploadID)
	}

	collection := db.Collection(i.collectionName)
	_, err = collection.UpdateOne(ctx, bson.M{"_id": id, "workspaceId": wsID}, bson.M{"$set": upload})
	if err != nil {
		infra.Log.Errorf("failed to update upload: %s", err.Error())
		return err
	}
	return nil
}

func (i *UploadRepo) Delete(ctx context.Context, workspaceID, uploadID string) error {
	id, err := primitive.ObjectIDFromHex(uploadID)
	if err != nil {
		return fmt.Errorf(messages.InvalidObjectID, uploadID)
	}
	collection := db.Collection(i.collectionName)
	_, err = collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		infra.Log.Errorf("failed to delete upload: %s", err.Error())
		return err
	}
	return nil
}

func (i *UploadRepo) AddLogs(ctx context.Context, workspaceID, uploadID string, logs []models.Log) error {
	id, err := primitive.ObjectIDFromHex(uploadID)
	if err != nil {
		return fmt.Errorf(messages.InvalidObjectID, uploadID)
	}

	wsID, err := primitive.ObjectIDFromHex(workspaceID)
	if err != nil {
		return fmt.Errorf(messages.InvalidObjectID, workspaceID)
	}

	collection := db.Collection(i.collectionName)
	_, err = collection.UpdateOne(ctx, bson.M{"_id": id, "workspaceId": wsID}, bson.M{"$push": bson.M{"logs": bson.M{"$each": logs}}})
	if err != nil {
		infra.Log.Errorf("failed to add logs to upload: %s", err.Error())
		return err
	}
	return nil
}

func (i *UploadRepo) SetStatus(ctx context.Context, workspaceID, uploadID string, status models.UploadStatus) error {
	id, err := primitive.ObjectIDFromHex(uploadID)
	if err != nil {
		return fmt.Errorf(messages.InvalidObjectID, uploadID)
	}

	wsID, err := primitive.ObjectIDFromHex(workspaceID)
	if err != nil {
		return fmt.Errorf(messages.InvalidObjectID, workspaceID)
	}

	collection := db.Collection(i.collectionName)

	_, err = collection.UpdateOne(ctx, bson.M{"_id": id, "workspaceId": wsID}, bson.M{"$set": bson.M{"status": status}})
	if err != nil {
		infra.Log.Errorf("failed to update upload status: %s", err.Error())
		return err
	}
	return nil
}
