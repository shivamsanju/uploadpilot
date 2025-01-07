package db

import (
	"context"

	"github.com/shivamsanju/uploader/internal/db/models"
	g "github.com/shivamsanju/uploader/pkg/globals"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func GetWorkflows(ctx context.Context) ([]models.Workflow, error) {
	collection := g.Db.Database("uploader").Collection("workflows")
	var cb []models.Workflow
	opts := options.Find().SetSort(bson.M{"updatedAt": 1})
	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		g.Log.Errorf("no workflows found: %s", err.Error())
		return nil, err
	}
	defer cursor.Close(ctx)
	cursor.All(ctx, &cb)
	g.Log.Infof("found %d workflows", len(cb))
	return cb, nil
}

func AddWorkflow(ctx context.Context, cb *models.Workflow) bson.ObjectID {
	g.Log.Infof("adding workflow: %+v", cb)
	collection := g.Db.Database("uploader").Collection("workflows")
	r, err := collection.InsertOne(ctx, &cb)
	if err != nil {
		g.Log.Errorf("failed to add workflow: %v", err.Error())
	}
	return (r.InsertedID).(bson.ObjectID)
}

func GetWorkflow(ctx context.Context, id bson.ObjectID) *models.Workflow {
	collection := g.Db.Database("uploader").Collection("workflows")
	var cb models.Workflow
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&cb)
	if err != nil {
		g.Log.Errorf("failed to find workflow: %s", err.Error())
	}
	return &cb
}

func UpdateWorkflow(ctx context.Context, id bson.ObjectID, metadata map[string]interface{}) error {
	collection := g.Db.Database("uploader").Collection("workflows")
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": metadata})
	return err
}

func DeleteWorkflow(ctx context.Context, id bson.ObjectID) error {
	collection := g.Db.Database("uploader").Collection("workflows")
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
