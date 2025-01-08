package repo

import (
	"context"

	"github.com/shivamsanju/uploader/internal/db/models"
	g "github.com/shivamsanju/uploader/pkg/globals"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func GetStorageConnectors(ctx context.Context) ([]models.StorageConnector, error) {
	collection := g.Db.Database("uploader").Collection("storage")
	var cb []models.StorageConnector
	opts := options.Find().SetSort(bson.D{{Key: "updatedat", Value: -1}})
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

func AddStorageConnector(ctx context.Context, cb *models.StorageConnector) bson.ObjectID {
	g.Log.Infof("adding storage: %+v", cb)
	collection := g.Db.Database("uploader").Collection("storage")
	r, err := collection.InsertOne(ctx, &cb)
	if err != nil {
		g.Log.Errorf("failed to add storage: %v", err.Error())
	}
	return (r.InsertedID).(bson.ObjectID)
}

func DeleteStorageConnector(ctx context.Context, id bson.ObjectID) error {
	collection := g.Db.Database("uploader").Collection("storage")
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
