package db

import (
	"context"
	"time"

	"github.com/uploadpilot/uploadpilot/internal/config"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	db *mongo.Database
)

func Init() error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(config.MongoURI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		infra.Log.Error("failed to connect to mongodb!")
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		infra.Log.Error("failed to connect to mongodb!")
		return err
	}
	db = client.Database(config.DatabaseName)

	infra.Log.Info("successfully connected to mongodb!")

	err = seedDatabase(ctx)
	if err != nil {
		return err
	}
	return nil
}

func seedDatabase(ctx context.Context) error {
	collection := db.Collection("storageconnectors")
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
		infra.Log.Errorf("failed to add storage: %v", err.Error())
	}
	return err
}
