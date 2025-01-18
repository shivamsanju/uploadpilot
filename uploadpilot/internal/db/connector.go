package db

import (
	"context"
	"fmt"
	"time"

	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StorageConnectorRepo interface {
	FindAll(ctx context.Context, skip int64, limit int64, search string) ([]models.StorageConnector, int64, error)
	Get(ctx context.Context, id string) (*models.StorageConnector, error)
	Create(ctx context.Context, cb *models.StorageConnector) (primitive.ObjectID, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, updateMap map[string]interface{}) error
}

type storageConnectorRepo struct {
	collectionName string
}

func NewStorageConnectorRepo() StorageConnectorRepo {
	return &storageConnectorRepo{
		collectionName: "storageconnectors",
	}
}
func (sr *storageConnectorRepo) FindAll(ctx context.Context, skip int64, limit int64, search string) ([]models.StorageConnector, int64, error) {
	filter := bson.M{}
	if search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": search, "$options": "i"}},
			{"type": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	collection := db.Collection(sr.collectionName)
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
		infra.Log.Errorf("failed to find storage connectors: %s", err.Error())
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// Decode results
	var cb []models.StorageConnector
	if err := cursor.All(ctx, &cb); err != nil {
		infra.Log.Errorf("failed to decode storage connectors: %s", err.Error())
		return nil, 0, err
	}

	infra.Log.Infof("found %d storage connectors (total records: %d)", len(cb), totalRecords)
	return cb, totalRecords, nil
}

func (sr *storageConnectorRepo) Get(ctx context.Context, id string) (*models.StorageConnector, error) {
	collection := db.Collection(sr.collectionName)
	var cb models.StorageConnector
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&cb)
	if err != nil {
		infra.Log.Errorf("failed to find storage: %s", err.Error())
		return nil, err
	}
	return &cb, nil
}

func (sr *storageConnectorRepo) Create(ctx context.Context, cb *models.StorageConnector) (primitive.ObjectID, error) {
	cb.ID = primitive.NewObjectID()
	cb.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	cb.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	collection := db.Collection(sr.collectionName)
	r, err := collection.InsertOne(ctx, &cb)
	if err != nil {
		infra.Log.Errorf("failed to add storage: %v", err.Error())
		return primitive.ObjectID{}, err
	}
	return (r.InsertedID).(primitive.ObjectID), nil
}

func (sr *storageConnectorRepo) Delete(ctx context.Context, id string) error {
	collection := db.Collection(sr.collectionName)
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (sr *storageConnectorRepo) Update(ctx context.Context, id string, updateMap map[string]interface{}) error {
	updateMap["updatedAt"] = primitive.NewDateTimeFromTime(time.Now())

	// TODO: Fix this method
	connectorType, ok := updateMap["type"].(string)
	if !ok {
		return fmt.Errorf("invalid connector type")
	}
	collection := db.Collection(sr.collectionName)
	if connectorType == string(models.StorageTypeAzure) {
		accountKey, ok := updateMap["accountKey"].(string)
		if !ok {
			return fmt.Errorf("invalid account key")
		}
		collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"accountKey": accountKey}})
	} else if connectorType == string(models.StorageTypeS3) {
		secretKey, ok := updateMap["secretKey"].(string)
		if !ok {
			return fmt.Errorf("invalid secret key")
		}
		collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"secretKey": secretKey}})
	} else if connectorType == string(models.StorageTypeGCS) {
		serviceAccountKey, ok := updateMap["serviceAccountKey"].(string)
		if !ok {
			return fmt.Errorf("invalid service account key")
		}
		collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"serviceAccountKey": serviceAccountKey}})
	} else {
		return fmt.Errorf("invalid connector type")
	}
	return nil
}
