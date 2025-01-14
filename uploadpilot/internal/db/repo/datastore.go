package repo

import (
	"context"
	"time"

	"github.com/shivamsanju/uploader/internal/db/models"
	g "github.com/shivamsanju/uploader/pkg/globals"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DataStoreRepo interface {
	GetAll(ctx context.Context) ([]bson.M, error)
	Get(ctx context.Context, id string) (*bson.M, error)
	Create(ctx context.Context, dataStore *models.DataStore) (primitive.ObjectID, error)
	Delete(ctx context.Context, id string) error
}

type dataStoreRepo struct {
	collectionName string
}

func NewDataStoreRepo() DataStoreRepo {
	return &dataStoreRepo{
		collectionName: "datastores",
	}
}

func (ds *dataStoreRepo) GetAll(ctx context.Context) ([]bson.M, error) {
	collection := g.Db.Database(g.DbName).Collection(ds.collectionName)
	var cb []bson.M
	pipeline := []bson.M{
		{"$lookup": bson.M{
			"from":         "storageconnectors",
			"localField":   "connectorId",
			"foreignField": "_id",
			"as":           "connectorDetails",
		}},
		{"$unwind": bson.M{"path": "$connectorDetails"}},
		{"$addFields": bson.M{
			"connectorName": "$connectorDetails.name",
			"connectorType": "$connectorDetails.type",
			"connectorId":   "$connectorDetails._id",
		}},
		{"$unset": "connectorDetails"},
		{"$sort": bson.M{"updatedAt": -1}},
	}
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		g.Log.Errorf("no data stores found: %v", err.Error())
		return nil, err
	}
	defer cursor.Close(ctx)
	cursor.All(ctx, &cb)
	g.Log.Infof("found %d data stores", len(cb))
	return cb, nil
}

func (ds *dataStoreRepo) Get(ctx context.Context, id string) (*bson.M, error) {
	collection := g.Db.Database(g.DbName).Collection(ds.collectionName)
	connectorId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		g.Log.Errorf("not a valid id: %v", err.Error())
		return nil, err
	}
	var cb []bson.M
	pipeline := []bson.M{
		{"$match": bson.M{
			"_id": connectorId,
		}},
		{"$lookup": bson.M{
			"from":         "storageconnectors",
			"localField":   "connectorId",
			"foreignField": "_id",
			"as":           "connectorDetails",
		}},
		{"$unwind": bson.M{"path": "$connectorDetails"}},
	}
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		g.Log.Errorf("failed to find data store: %v", err.Error())
		return nil, err
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, &cb)
	if err != nil {
		g.Log.Errorf("failed to find data store: %v", err.Error())
		return nil, err
	}
	if len(cb) == 0 {
		return nil, nil
	}
	return &cb[0], nil
}

func (ds *dataStoreRepo) Create(ctx context.Context, dataStore *models.DataStore) (primitive.ObjectID, error) {
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

func (ds *dataStoreRepo) Delete(ctx context.Context, id string) error {
	collection := g.Db.Database(g.DbName).Collection(ds.collectionName)
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
