package db

import (
	"context"
	"fmt"
	"time"

	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"github.com/uploadpilot/uploadpilot/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UploaderRepo interface {
	FindAll(ctx context.Context, skip int64, limit int64, search string) ([]models.Uploader, int64, error)
	Get(ctx context.Context, id string) (*models.Uploader, error)
	Create(ctx context.Context, cb *models.Uploader) (primitive.ObjectID, error)
	Delete(ctx context.Context, id string) error
	GetDataStoreCreds(ctx context.Context, id string) (map[string]interface{}, error)
	UpdateConfig(ctx context.Context, id string, cb *models.UploaderConfig, updatedBy string) error
}

type uploaderRepo struct {
	collectionName string
}

func NewUploaderRepo() UploaderRepo {
	return &uploaderRepo{
		collectionName: "uploaders",
	}
}

func (ur *uploaderRepo) FindAll(ctx context.Context, skip int64, limit int64, search string) ([]models.Uploader, int64, error) {
	filter := bson.M{}
	if search != "" {
		filter["name"] = bson.M{"$regex": search, "$options": "i"}
	}

	collection := db.Collection(ur.collectionName)
	totalRecords, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		infra.Log.Errorf("failed to count documents: %s", err.Error())
		return nil, 0, err
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "updatedAt", Value: -1}}).
		SetSkip(skip).
		SetLimit(limit)

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		infra.Log.Errorf("failed to find uploaders: %s", err.Error())
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// Decode results
	var cb []models.Uploader
	if err := cursor.All(ctx, &cb); err != nil {
		infra.Log.Errorf("failed to decode uploaders: %s", err.Error())
		return nil, 0, err
	}

	return cb, totalRecords, nil
}

func (ur *uploaderRepo) Create(ctx context.Context, cb *models.Uploader) (primitive.ObjectID, error) {
	cb.ID = primitive.NewObjectID()
	cb.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	cb.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	collection := db.Collection(ur.collectionName)
	r, err := collection.InsertOne(ctx, &cb)
	if err != nil {
		infra.Log.Errorf("failed to add Uploader: %v", err.Error())
		return primitive.ObjectID{}, err
	}
	return (r.InsertedID).(primitive.ObjectID), nil
}

func (ur *uploaderRepo) Get(ctx context.Context, id string) (*models.Uploader, error) {
	uploaderID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var cb models.Uploader
	collection := db.Collection(ur.collectionName)
	err = collection.FindOne(ctx, bson.M{"_id": uploaderID}).Decode(&cb)
	if err != nil {
		infra.Log.Errorf("failed to find storage: %s", err.Error())
		return nil, err
	}
	return &cb, nil
}

func (ur *uploaderRepo) Delete(ctx context.Context, id string) error {
	uploaderID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	collection := db.Collection(ur.collectionName)
	_, err = collection.DeleteOne(ctx, bson.M{"_id": uploaderID})
	return err
}

func (ur *uploaderRepo) UpdateConfig(ctx context.Context, id string, updatedData *models.UploaderConfig, updatedBy string) error {
	uploaderID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	updateFields := utils.FilterNonEmptyBsonFields(bson.M{
		"config.maxFileSize":            updatedData.MaxFileSize,
		"config.minFileSize":            updatedData.MinFileSize,
		"config.maxNumberOfFiles":       updatedData.MaxNumberOfFiles,
		"config.minNumberOfFiles":       updatedData.MinNumberOfFiles,
		"config.maxTotalFileSize":       updatedData.MaxTotalFileSize,
		"config.allowedFileTypes":       updatedData.AllowedFileTypes,
		"config.allowedSources":         updatedData.AllowedSources,
		"config.requiredMetadataFields": updatedData.RequiredMetadataFields,
		"config.allowPauseAndResume":    updatedData.AllowPauseAndResume,
		"config.enableImageEditing":     updatedData.EnableImageEditing,
		"config.useCompression":         updatedData.UseCompression,
		"config.useFaultTolerantMode":   updatedData.UseFaultTolerantMode,
		"updatedBy":                     updatedBy,
		"updatedAt":                     primitive.NewDateTimeFromTime(time.Now()),
	})

	update := bson.M{"$set": updateFields}

	collection := db.Collection(ur.collectionName)
	_, err = collection.UpdateOne(ctx, bson.M{"_id": uploaderID}, update)
	return err
}

func (ur *uploaderRepo) GetDataStoreCreds(ctx context.Context, id string) (map[string]interface{}, error) {
	uploaderID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	pipeline := []bson.M{
		{"$match": bson.M{"_id": uploaderID}},
		{"$lookup": bson.M{
			"from":         "storageconnectors",
			"localField":   "dataStore.connectorId",
			"foreignField": "_id",
			"as":           "connectorDetails",
		}},
		{"$unwind": bson.M{"path": "$connectorDetails"}},
		{"$addFields": bson.M{"bucket": "$dataStore.bucket"}},
	}

	// Execute the aggregation pipeline
	collection := db.Collection(ur.collectionName)
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to execute aggregation: %w", err)
	}
	defer cursor.Close(ctx)

	// Decode results into a slice of maps
	var results []map[string]interface{}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode aggregation results: %w", err)
	}

	if len(results) == 0 {
		return nil, nil
	}
	return results[0], nil
}
