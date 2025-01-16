package db

import (
	"context"
	"time"

	"github.com/uploadpilot/uploadpilot/pkg/db/models"
	g "github.com/uploadpilot/uploadpilot/pkg/globals"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ImportRepo interface {
	Get(ctx context.Context, id string) (*models.Import, error)
	GetAll(ctx context.Context) ([]models.Import, error)
	Create(ctx context.Context, imp *models.Import) (*models.Import, error)
	Update(ctx context.Context, id *primitive.ObjectID, imp *models.Import) (*models.Import, error)
	Delete(ctx context.Context, id string) error
	FindImportsByUploaderId(ctx context.Context, uploaderId string) ([]models.Import, error)
}

type importRepo struct {
	collectionName string
}

func NewImportRepo() ImportRepo {
	return &importRepo{
		collectionName: "imports",
	}
}

func (i *importRepo) Get(ctx context.Context, id string) (*models.Import, error) {
	collection := g.Db.Database(g.DbName).Collection(i.collectionName)
	importId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		g.Log.Errorf("not a valid id: %v", err.Error())
		return nil, err
	}
	var cb models.Import
	err = collection.FindOne(ctx, bson.M{"_id": importId}).Decode(&cb)
	if err != nil {
		g.Log.Errorf("failed to find import: %s", err.Error())
		return nil, err
	}
	return &cb, nil
}

func (i *importRepo) GetAll(ctx context.Context) ([]models.Import, error) {
	collection := g.Db.Database(g.DbName).Collection(i.collectionName)
	var cb []models.Import
	opts := options.Find().SetSort(bson.D{{Key: "finishedAt", Value: -1}})
	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		g.Log.Errorf("no imports found: %s", err.Error())
		return nil, err
	}
	defer cursor.Close(ctx)
	cursor.All(ctx, &cb)
	g.Log.Infof("found %d imports", len(cb))
	return cb, nil
}

func (i *importRepo) Create(ctx context.Context, imp *models.Import) (*models.Import, error) {
	collection := g.Db.Database(g.DbName).Collection(i.collectionName)
	imp.FinishedAt = primitive.NewDateTimeFromTime(time.Now())
	_, err := collection.InsertOne(ctx, imp)
	if err != nil {
		g.Log.Errorf("failed to create import: %s", err.Error())
		return nil, err
	}
	return imp, nil
}

func (i *importRepo) Update(ctx context.Context, id *primitive.ObjectID, imp *models.Import) (*models.Import, error) {
	collection := g.Db.Database(g.DbName).Collection(i.collectionName)
	_, err := collection.UpdateOne(ctx, bson.M{"_id": &id}, bson.M{"$set": imp})
	if err != nil {
		g.Log.Errorf("failed to update import: %s", err.Error())
		return nil, err
	}
	return imp, nil
}

func (i *importRepo) Delete(ctx context.Context, id string) error {
	collection := g.Db.Database(g.DbName).Collection(i.collectionName)
	importId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		g.Log.Errorf("not a valid id: %v", err.Error())
		return err
	}
	_, err = collection.DeleteOne(ctx, bson.M{"_id": importId})
	if err != nil {
		g.Log.Errorf("failed to delete import: %s", err.Error())
		return err
	}
	return nil
}

func (i *importRepo) FindImportsByUploaderId(ctx context.Context, id string) ([]models.Import, error) {
	collection := g.Db.Database(g.DbName).Collection(i.collectionName)
	uploaderId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		g.Log.Errorf("not a valid id: %v", err.Error())
		return nil, err
	}
	var cb []models.Import
	opts := options.Find().SetSort(bson.D{{Key: "finishedAt", Value: -1}})
	cursor, err := collection.Find(ctx, bson.M{"uploaderId": uploaderId}, opts)
	if err != nil {
		g.Log.Errorf("no imports found: %s", err.Error())
		return nil, err
	}
	defer cursor.Close(ctx)
	cursor.All(ctx, &cb)
	g.Log.Infof("found %d imports", len(cb))
	return cb, nil
}
