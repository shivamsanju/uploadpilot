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

type ImportPolicyRepo interface {
	GetImportPolicies(ctx context.Context) ([]models.ImportPolicy, error)
	GetImportPolicy(ctx context.Context, id string) (*models.ImportPolicy, error)
	CreateImportPolicy(ctx context.Context, cb *models.ImportPolicy) (primitive.ObjectID, error)
	DeleteImportPolicy(ctx context.Context, id string) error
	UpdateImportPolicy(ctx context.Context, id string, updateMap map[string]interface{}) error
}

type importPolicyRepo struct {
	collectionName string
}

func NewImportPolicyRepo() ImportPolicyRepo {
	return &importPolicyRepo{collectionName: "importpolicies"}
}

func (ip *importPolicyRepo) GetImportPolicies(ctx context.Context) ([]models.ImportPolicy, error) {
	collection := g.Db.Database(g.DbName).Collection(ip.collectionName)
	var cb []models.ImportPolicy
	opts := options.Find().SetSort(bson.D{{Key: "updatedat", Value: -1}}).SetLimit(100)
	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	cursor.All(ctx, &cb)
	return cb, nil
}

func (ip *importPolicyRepo) GetImportPolicy(ctx context.Context, id string) (*models.ImportPolicy, error) {
	collection := g.Db.Database(g.DbName).Collection(ip.collectionName)
	var cb models.ImportPolicy
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&cb)
	if err != nil {
		return nil, err
	}
	return &cb, nil
}

func (ip *importPolicyRepo) CreateImportPolicy(ctx context.Context, cb *models.ImportPolicy) (primitive.ObjectID, error) {
	cb.ID = primitive.NewObjectID()
	cb.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	cb.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	collection := g.Db.Database(g.DbName).Collection(ip.collectionName)
	r, err := collection.InsertOne(ctx, &cb)
	if err != nil {
		g.Log.Errorf("failed to add import policy: %v", err.Error())
		return primitive.ObjectID{}, err
	}
	return (r.InsertedID).(primitive.ObjectID), nil
}

func (ip *importPolicyRepo) DeleteImportPolicy(ctx context.Context, id string) error {
	collection := g.Db.Database(g.DbName).Collection(ip.collectionName)
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}

func (ip *importPolicyRepo) UpdateImportPolicy(ctx context.Context, id string, updateMap map[string]interface{}) error {
	updateMap["updatedAt"] = primitive.NewDateTimeFromTime(time.Now())

	collection := g.Db.Database(g.DbName).Collection(ip.collectionName)
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updateMap})
	if err != nil {
		return err
	}
	return nil
}
