package db

import (
	"context"
	"time"

	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WorkspaceRepo interface {
	Create(ctx context.Context, workspace *models.Workspace, userId string) (*models.Workspace, error)
	GetWorkspacesForUser(ctx context.Context, userId string) ([]models.Workspace, error)
	AddUserToWorkspace(ctx context.Context, workspaceID primitive.ObjectID, user *models.WorkspaceUser) error
}

type workspaceRepo struct {
	collectionName     string
	userColelctionName string
}

func NewWorkspaceRepo() WorkspaceRepo {
	return &workspaceRepo{
		collectionName:     "workspaces",
		userColelctionName: "users",
	}
}

func (u *workspaceRepo) GetWorkspacesForUser(ctx context.Context, userId string) ([]models.Workspace, error) {
	collection := db.Collection(u.collectionName)
	var workspaces []models.Workspace

	findOptions := options.Find().SetProjection(bson.M{
		"id":   1,
		"name": 1,
	})
	findOptions.SetSort(bson.M{"updatedAt": -1})

	cursor, err := collection.Find(ctx, bson.M{"users.userId": userId}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &workspaces); err != nil {
		return nil, err
	}

	return workspaces, nil
}

func (u *workspaceRepo) Create(ctx context.Context, workspace *models.Workspace, userId string) (*models.Workspace, error) {
	workspace.ID = primitive.NewObjectID()
	workspace.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	workspace.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	workspace.Users = []models.WorkspaceUser{{UserID: userId, Role: models.UserRoleOwner}}

	session, err := db.Client().StartSession()
	if err != nil {
		return nil, err
	}

	defer session.EndSession(ctx)
	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		collection := db.Collection(u.collectionName)
		result, err := collection.InsertOne(ctx, workspace)
		if err != nil {
			return nil, err
		}

		// Update user collection
		userCollection := db.Collection(u.userColelctionName)
		userFilter := bson.M{"userId": userId}
		userUpdate := bson.M{"$addToSet": bson.M{"workspaces": &models.UserWorkspace{WorkspaceID: workspace.ID, Role: models.UserRoleOwner}}}
		if _, err := userCollection.UpdateOne(sessCtx, userFilter, userUpdate); err != nil {
			return nil, err
		}
		return result, nil
	}

	// Execute the transaction
	_, err = session.WithTransaction(ctx, callback, options.Transaction())
	if err != nil {
		return nil, err
	}
	return workspace, nil

}

func (u *workspaceRepo) AddUserToWorkspace(ctx context.Context, workspaceID primitive.ObjectID, user *models.WorkspaceUser) error {
	session, err := db.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	// Transactional function
	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Update workspace collection
		workspaceCollection := db.Collection(u.collectionName)
		workspaceFilter := bson.M{"_id": workspaceID}
		workspaceUpdate := bson.M{"$addToSet": bson.M{"users": user}}
		if _, err := workspaceCollection.UpdateOne(sessCtx, workspaceFilter, workspaceUpdate); err != nil {
			return nil, err
		}

		// Update user collection
		userCollection := db.Collection(u.userColelctionName)
		userFilter := bson.M{"userId": user.UserID}
		userUpdate := bson.M{"$addToSet": bson.M{"workspaces": &models.UserWorkspace{WorkspaceID: workspaceID, Role: user.Role}}}
		if _, err := userCollection.UpdateOne(sessCtx, userFilter, userUpdate); err != nil {
			return nil, err
		}

		return nil, nil
	}

	// Execute the transaction
	_, err = session.WithTransaction(ctx, callback, options.Transaction())
	return err
}
