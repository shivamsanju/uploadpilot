package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/shivamsanju/uploader/internal/db/models"
	g "github.com/shivamsanju/uploader/pkg/globals"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StorageConnectorRepo interface {
	GetAll(ctx context.Context) ([]models.StorageConnector, error)
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
func (sr *storageConnectorRepo) GetAll(ctx context.Context) ([]models.StorageConnector, error) {
	collection := g.Db.Database(g.DbName).Collection(sr.collectionName)
	var cb []models.StorageConnector
	opts := options.Find().SetSort(bson.D{{Key: "updatedAt", Value: -1}})
	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		g.Log.Errorf("no storages found: %s", err.Error())
		return nil, err
	}
	defer cursor.Close(ctx)
	cursor.All(ctx, &cb)
	g.Log.Infof("found %d storages", len(cb))
	return cb, nil
}

func (sr *storageConnectorRepo) Get(ctx context.Context, id string) (*models.StorageConnector, error) {
	collection := g.Db.Database(g.DbName).Collection(sr.collectionName)
	var cb models.StorageConnector
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&cb)
	if err != nil {
		g.Log.Errorf("failed to find storage: %s", err.Error())
		return nil, err
	}
	return &cb, nil
}

func (sr *storageConnectorRepo) Create(ctx context.Context, cb *models.StorageConnector) (primitive.ObjectID, error) {
	cb.ID = primitive.NewObjectID()
	cb.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	cb.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	collection := g.Db.Database(g.DbName).Collection(sr.collectionName)
	r, err := collection.InsertOne(ctx, &cb)
	if err != nil {
		g.Log.Errorf("failed to add storage: %v", err.Error())
		return primitive.ObjectID{}, err
	}
	return (r.InsertedID).(primitive.ObjectID), nil
}

func (sr *storageConnectorRepo) Delete(ctx context.Context, id string) error {
	collection := g.Db.Database(g.DbName).Collection(sr.collectionName)
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (sr *storageConnectorRepo) Update(ctx context.Context, id string, updateMap map[string]interface{}) error {
	updateMap["updatedAt"] = primitive.NewDateTimeFromTime(time.Now())

	// TODO: Fix this method
	collection := g.Db.Database(g.DbName).Collection(sr.collectionName)
	connectorType, ok := updateMap["type"].(string)
	if !ok {
		return fmt.Errorf("invalid connector type")
	}
	if connectorType == string(models.StorageTypeAzure) {
		accountKey, ok := updateMap["accountKey"].(string)
		if !ok {
			return fmt.Errorf("invalid account key")
		}
		collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"accountKey": accountKey}})
	}
	return nil
}
