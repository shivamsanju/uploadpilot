package repo

import (
	"context"
	"time"

	"github.com/shivamsanju/uploader/internal/db/models"
	g "github.com/shivamsanju/uploader/pkg/globals"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DataStoreRepo interface {
	CreateDataStore(ctx context.Context, dataStore *models.DataStore) (primitive.ObjectID, error)
	GetDataStore(ctx context.Context, id string) (*models.DataStore, error)
	GetDataStores(ctx context.Context) ([]models.DataStore, error)
	DeleteDataStore(ctx context.Context, id string) error
}

type dataStoreRepo struct {
	collectionName string
}

func NewDataStoreRepo() DataStoreRepo {
	return &dataStoreRepo{
		collectionName: "datastores",
	}
}

func (ds *dataStoreRepo) GetDataStores(ctx context.Context) ([]models.DataStore, error) {
	collection := g.Db.Database(g.DbName).Collection(ds.collectionName)
	var DataStore []models.DataStore
	opts := options.Find().SetSort(bson.D{{Key: "updatedat", Value: -1}})
	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		g.Log.Errorf("no data stores found: %v", err.Error())
		return nil, err
	}
	defer cursor.Close(ctx)
	cursor.All(ctx, &DataStore)
	g.Log.Infof("found %d data stores", len(DataStore))
	return DataStore, nil
}

func (ds *dataStoreRepo) GetDataStore(ctx context.Context, id string) (*models.DataStore, error) {
	collection := g.Db.Database(g.DbName).Collection(ds.collectionName)
	var DataStore models.DataStore
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&DataStore)
	if err != nil {
		g.Log.Errorf("failed to find data store: %v", err.Error())
		return nil, err
	}
	return &DataStore, nil
}

func (ds *dataStoreRepo) CreateDataStore(ctx context.Context, dataStore *models.DataStore) (primitive.ObjectID, error) {
	g.Log.Infof("adding datatore: %+v", dataStore)
	dataStore.ID = primitive.NewObjectID()
	dataStore.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	dataStore.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	collection := g.Db.Database(g.DbName).Collection(ds.collectionName)
	r, err := collection.InsertOne(ctx, &dataStore)
	if err != nil {
		g.Log.Errorf("failed to create data store: %v", err.Error())
		return primitive.ObjectID{}, err
	}
	return (r.InsertedID).(primitive.ObjectID), nil
}

func (ds *dataStoreRepo) DeleteDataStore(ctx context.Context, id string) error {
	collection := g.Db.Database(g.DbName).Collection(ds.collectionName)
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
