package db

import (
	"context"
	"time"

	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ImportRepo interface {
	FindAll(ctx context.Context, skip int64, limit int64, search string) ([]models.Import, int64, error)
	FindAllImportsByUploaderId(ctx context.Context, uploaderId string, skip int64, limit int64, search string) ([]models.Import, int64, error)
	Get(ctx context.Context, id string) (*models.Import, error)
	Create(ctx context.Context, imp *models.Import) (*models.Import, error)
	Update(ctx context.Context, id *primitive.ObjectID, imp *models.Import) (*models.Import, error)
	Delete(ctx context.Context, id string) error
}

type importRepo struct {
	collectionName string
}

func NewImportRepo() ImportRepo {
	return &importRepo{
		collectionName: "imports",
	}
}

func (i *importRepo) FindAll(ctx context.Context, skip int64, limit int64, search string) ([]models.Import, int64, error) {
	// Build filter with optional search
	filter := bson.M{}
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
	var cb []models.Import
	if err := cursor.All(ctx, &cb); err != nil {
		infra.Log.Errorf("failed to decode imports: %s", err.Error())
		return nil, 0, err
	}

	infra.Log.Infof("found %d imports out of %d", len(cb), totalRecords)
	return cb, totalRecords, nil
}

func (i *importRepo) FindAllImportsByUploaderId(ctx context.Context, id string, skip int64, limit int64, search string) ([]models.Import, int64, error) {
	uploaderId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		infra.Log.Errorf("not a valid id: %v", err.Error())
		return nil, 0, err
	}

	// Build filter with uploaderId and optional search
	filter := bson.M{"uploaderId": uploaderId}
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
	var cb []models.Import
	if err := cursor.All(ctx, &cb); err != nil {
		infra.Log.Errorf("failed to decode imports: %s", err.Error())
		return nil, 0, err
	}

	return cb, totalRecords, nil
}

func (i *importRepo) Get(ctx context.Context, id string) (*models.Import, error) {
	importId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		infra.Log.Errorf("not a valid id: %v", err.Error())
		return nil, err
	}
	var cb models.Import
	collection := db.Collection(i.collectionName)
	err = collection.FindOne(ctx, bson.M{"_id": importId}).Decode(&cb)
	if err != nil {
		infra.Log.Errorf("failed to find import: %s", err.Error())
		return nil, err
	}
	return &cb, nil
}

func (i *importRepo) Create(ctx context.Context, imp *models.Import) (*models.Import, error) {
	imp.FinishedAt = primitive.NewDateTimeFromTime(time.Now())
	collection := db.Collection(i.collectionName)
	_, err := collection.InsertOne(ctx, imp)
	if err != nil {
		infra.Log.Errorf("failed to create import: %s", err.Error())
		return nil, err
	}
	return imp, nil
}

func (i *importRepo) Update(ctx context.Context, id *primitive.ObjectID, imp *models.Import) (*models.Import, error) {
	collection := db.Collection(i.collectionName)
	_, err := collection.UpdateOne(ctx, bson.M{"_id": &id}, bson.M{"$set": imp})
	if err != nil {
		infra.Log.Errorf("failed to update import: %s", err.Error())
		return nil, err
	}
	return imp, nil
}

func (i *importRepo) Delete(ctx context.Context, id string) error {
	importId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		infra.Log.Errorf("not a valid id: %v", err.Error())
		return err
	}
	collection := db.Collection(i.collectionName)
	_, err = collection.DeleteOne(ctx, bson.M{"_id": importId})
	if err != nil {
		infra.Log.Errorf("failed to delete import: %s", err.Error())
		return err
	}
	return nil
}
