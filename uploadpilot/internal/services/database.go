package services

import (
	"context"
	"time"

	"github.com/shivamsanju/uploader/internal/config"
	"github.com/shivamsanju/uploader/internal/db/models"
	g "github.com/shivamsanju/uploader/pkg/globals"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initMongoDB(config *config.Config) error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(config.MongoURI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		g.Log.Error("failed to connect to mongodb!")
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		g.Log.Error("failed to connect to mongodb!")
		return err
	}
	g.Db = client
	g.DbName = config.DatabaseName
	g.Log.Info("successfully connected to mongodb!")

	err = seedDatabase(ctx)
	if err != nil {
		return err
	}
	return nil
}

func seedDatabase(ctx context.Context) error {
	collection := g.Db.Database(g.DbName).Collection("storageconnectors")
	// check if exists
	item := collection.FindOne(ctx, bson.M{"name": "local"})
	if item.Err() == nil {
		return nil
	}

	cb := models.StorageConnector{
		ID:   primitive.NewObjectID(),
		Name: "local",
		Type: "local",
		LocalConfig: &models.LocalConfig{
			Region:    "us-east-1",
			AccessKey: "minio",
			SecretKey: "minio",
		},
	}
	_, err := collection.InsertOne(ctx, &cb)
	if err != nil {
		g.Log.Errorf("failed to add storage: %v", err.Error())
	}
	return err
}
