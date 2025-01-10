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

type WorkflowRepo interface {
	GetWorkflows(ctx context.Context) ([]models.Workflow, error)
	GetWorkflow(ctx context.Context, id string) (*models.Workflow, error)
	CreateWorkflow(ctx context.Context, cb *models.Workflow) (primitive.ObjectID, error)
	DeleteWorkflow(ctx context.Context, id string) error
	UpdateWorkflow(ctx context.Context, id string, metadata map[string]interface{}) error
}

type workflowRepo struct {
	collectionName string
}

func NewWorkflowRepo() WorkflowRepo {
	return &workflowRepo{
		collectionName: "workflows",
	}
}

func (wr *workflowRepo) GetWorkflows(ctx context.Context) ([]models.Workflow, error) {
	collection := g.Db.Database(g.DbName).Collection(wr.collectionName)
	var cb []models.Workflow
	opts := options.Find().SetSort(bson.D{{Key: "updatedat", Value: -1}})
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

func (wr *workflowRepo) GetWorkflow(ctx context.Context, id string) (*models.Workflow, error) {
	collection := g.Db.Database(g.DbName).Collection(wr.collectionName)
	var cb models.Workflow
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&cb)
	if err != nil {
		g.Log.Errorf("failed to find workflow: %s", err.Error())
		return nil, err
	}
	return &cb, nil
}

func (wr *workflowRepo) CreateWorkflow(ctx context.Context, cb *models.Workflow) (primitive.ObjectID, error) {
	cb.ID = primitive.NewObjectID()
	cb.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	cb.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	collection := g.Db.Database(g.DbName).Collection(wr.collectionName)
	r, err := collection.InsertOne(ctx, &cb)
	if err != nil {
		g.Log.Errorf("failed to add workflow: %v", err.Error())
		return primitive.ObjectID{}, err
	}
	return (r.InsertedID).(primitive.ObjectID), nil
}

func (wr *workflowRepo) DeleteWorkflow(ctx context.Context, id string) error {
	collection := g.Db.Database(g.DbName).Collection(wr.collectionName)
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (wr *workflowRepo) UpdateWorkflow(ctx context.Context, id string, metadata map[string]interface{}) error {
	metadata["updatedat"] = primitive.NewDateTimeFromTime(time.Now())

	collection := g.Db.Database(g.DbName).Collection(wr.collectionName)
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": metadata})
	return err
}
