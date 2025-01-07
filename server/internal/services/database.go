package services

import (
	"context"
	"time"

	"github.com/shivamsanju/uploader/internal/config"
	g "github.com/shivamsanju/uploader/pkg/globals"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func initMongoDB(config *config.Config) error {

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(config.MongoURI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(opts)
	if err != nil {
		g.Log.Error("failed to connect to mongodb!")
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		g.Log.Error("failed to connect to mongodb!")
		return err
	}
	g.Db = client
	g.Log.Info("successfully connected to mongodb!")
	return nil
}
