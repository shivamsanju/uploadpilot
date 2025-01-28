package db

import (
	"context"
	"fmt"
	"time"

	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/dto"
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"github.com/uploadpilot/uploadpilot/internal/msg"
	"github.com/uploadpilot/uploadpilot/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WorkspaceRepo struct {
	collectionName     string
	userColelctionName string
}

// NewWorkspaceRepo initializes a new instance of workspaceRepo with predefined
// collection names for workspaces and users
func NewWorkspaceRepo() *WorkspaceRepo {
	return &WorkspaceRepo{
		collectionName:     "workspaces",
		userColelctionName: "users",
	}
}

func (wr *WorkspaceRepo) GetAll(ctx context.Context, userId string) ([]dto.WorkspaceNameID, error) {
	collection := db.Collection(wr.collectionName)
	var workspaces []dto.WorkspaceNameID

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

func (wr *WorkspaceRepo) Get(ctx context.Context, workspaceID string) (*models.Workspace, error) {
	id, err := primitive.ObjectIDFromHex(workspaceID)
	if err != nil {
		return nil, fmt.Errorf(msg.InvalidObjectID, workspaceID)
	}
	collection := db.Collection(wr.collectionName)
	var workspace models.Workspace
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&workspace)
	if err != nil {
		return nil, err
	}
	return &workspace, nil
}

func (wr *WorkspaceRepo) Create(ctx context.Context, workspace *models.Workspace) (*models.Workspace, error) {
	user, err := utils.GetUserDetailsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	workspace.ID = primitive.NewObjectID()
	workspace.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	workspace.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	workspace.CreatedBy = user.Email
	workspace.UpdatedBy = user.Email

	workspace.Users = []models.WorkspaceUser{{UserID: user.UserID, Role: models.UserRoleOwner}}

	session, err := db.Client().StartSession()
	if err != nil {
		return nil, err
	}

	defer session.EndSession(ctx)
	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		collection := db.Collection(wr.collectionName)
		// validate before insert
		if err := infra.Validator.ValidateBody(workspace); err != nil {
			return nil, err
		}
		result, err := collection.InsertOne(ctx, workspace)
		if err != nil {
			return nil, err
		}

		// Update user collection
		userCollection := db.Collection(wr.userColelctionName)
		userFilter := bson.M{"userId": user.UserID}
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

func (wr *WorkspaceRepo) Delete(ctx context.Context, workspaceID primitive.ObjectID) error {
	collection := db.Collection(wr.collectionName)
	_, err := collection.DeleteOne(ctx, bson.M{"_id": workspaceID})
	return err
}

func (wr *WorkspaceRepo) Exists(ctx context.Context, workspaceID string) (bool, error) {
	id, err := primitive.ObjectIDFromHex(workspaceID)
	if err != nil {
		return false, fmt.Errorf(msg.InvalidObjectID, workspaceID)
	}
	collection := db.Collection(wr.collectionName)
	count, err := collection.CountDocuments(ctx, bson.M{"_id": id})
	return count > 0, err
}

func (wr *WorkspaceRepo) GetUsersInWorkspace(ctx context.Context, workspaceID string) ([]models.WorkspaceUserWithDetails, error) {
	id, err := primitive.ObjectIDFromHex(workspaceID)
	if err != nil {
		return nil, fmt.Errorf(msg.InvalidObjectID, workspaceID)
	}
	collection := db.Collection(wr.collectionName)
	var users []models.WorkspaceUserWithDetails

	pipeline := []bson.M{
		{
			"$match": bson.M{"_id": id},
		},
		{
			"$unwind": "$users",
		},
		{
			"$lookup": bson.M{
				"from":         wr.userColelctionName,
				"localField":   "users.userId",
				"foreignField": "userId",
				"as":           "userDetails",
			},
		},
		{
			"$unwind": "$userDetails",
		},
		{
			"$project": bson.M{
				"userId": "$users.userId",
				"role":   "$users.role",
				"name":   "$userDetails.name",
				"email":  "$userDetails.email",
			},
		},
		{
			"$sort": bson.M{"name": 1},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (wr *WorkspaceRepo) CheckIfUserExistsInWorkspace(ctx context.Context, workspaceID string, userID string) (bool, error) {
	id, err := primitive.ObjectIDFromHex(workspaceID)
	if err != nil {
		return false, fmt.Errorf(msg.InvalidObjectID, workspaceID)
	}
	collection := db.Collection(wr.collectionName)
	count, err := collection.CountDocuments(ctx, bson.M{"_id": id, "users.userId": userID})
	return count > 0, err
}

func (wr *WorkspaceRepo) AddUserToWorkspace(ctx context.Context, workspaceID string, user *models.WorkspaceUser) error {
	id, err := primitive.ObjectIDFromHex(workspaceID)
	if err != nil {
		return fmt.Errorf(msg.InvalidObjectID, workspaceID)
	}
	session, err := db.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	// Transactional function
	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Update workspace collection
		workspaceCollection := db.Collection(wr.collectionName)
		workspaceFilter := bson.M{"_id": id}
		workspaceUpdate := bson.M{"$addToSet": bson.M{"users": user}}
		if _, err := workspaceCollection.UpdateOne(sessCtx, workspaceFilter, workspaceUpdate); err != nil {
			return nil, err
		}

		// Update user collection
		userCollection := db.Collection(wr.userColelctionName)
		userFilter := bson.M{"userId": user.UserID}
		userUpdate := bson.M{"$addToSet": bson.M{"workspaces": &models.UserWorkspace{WorkspaceID: id, Role: user.Role}}}
		if _, err := userCollection.UpdateOne(sessCtx, userFilter, userUpdate); err != nil {
			return nil, err
		}

		return nil, nil
	}

	// Execute the transaction
	_, err = session.WithTransaction(ctx, callback, options.Transaction())
	return err
}

func (wr *WorkspaceRepo) RemoveUserFromWorkspace(ctx context.Context, workspaceID string, userID string) error {
	id, err := primitive.ObjectIDFromHex(workspaceID)
	if err != nil {
		return fmt.Errorf(msg.InvalidObjectID, workspaceID)
	}
	session, err := db.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	userRole, err := wr.GetUserRoleInWorkspace(ctx, id, userID)
	if err != nil || userRole == models.UserRoleNotFound {
		return fmt.Errorf("user %s not found in workspace %s", userID, workspaceID)
	}

	// Transactional function
	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Update workspace collection
		workspaceCollection := db.Collection(wr.collectionName)
		workspaceFilter := bson.M{"_id": id}
		workspaceUpdate := bson.M{"$pull": bson.M{"users": &models.WorkspaceUser{UserID: userID, Role: userRole}}}
		if _, err := workspaceCollection.UpdateOne(sessCtx, workspaceFilter, workspaceUpdate); err != nil {
			return nil, err
		}

		// Update user collection
		userCollection := db.Collection(wr.userColelctionName)
		userFilter := bson.M{"userId": userID}
		userUpdate := bson.M{"$pull": bson.M{"workspaces": &models.UserWorkspace{WorkspaceID: id, Role: userRole}}}
		if _, err := userCollection.UpdateOne(sessCtx, userFilter, userUpdate); err != nil {
			return nil, err
		}

		return nil, nil
	}

	// Execute the transaction
	_, err = session.WithTransaction(ctx, callback, options.Transaction())
	return err
}

func (wr *WorkspaceRepo) GetUserRoleInWorkspace(ctx context.Context, workspaceID primitive.ObjectID, userID string) (models.UserRole, error) {
	collection := db.Collection(wr.collectionName)
	var workspace models.Workspace
	err := collection.FindOne(ctx, bson.M{"_id": workspaceID, "users.userId": userID}).Decode(&workspace)
	if err != nil {
		return "", err
	}
	for _, user := range workspace.Users {
		if user.UserID == userID {
			return user.Role, nil
		}
	}
	return models.UserRoleNotFound, nil
}

func (wr *WorkspaceRepo) ChangeUserRoleInWorkspace(ctx context.Context, workspaceID string, userID string, role models.UserRole) error {
	id, err := primitive.ObjectIDFromHex(workspaceID)
	if err != nil {
		return fmt.Errorf(msg.InvalidObjectID, workspaceID)
	}
	collection := db.Collection(wr.collectionName)
	_, err = collection.UpdateOne(ctx, bson.M{"_id": id, "users.userId": userID}, bson.M{"$set": bson.M{"users.$.role": role}})
	return err
}

func (wr *WorkspaceRepo) IsUserLastOwner(ctx context.Context, workspaceID string, userID string) (bool, error) {
	isLastOwner := true
	id, err := primitive.ObjectIDFromHex(workspaceID)
	if err != nil {
		return isLastOwner, fmt.Errorf(msg.InvalidObjectID, workspaceID)
	}
	collection := db.Collection(wr.collectionName)
	var workspace models.Workspace
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&workspace)
	if err != nil {
		return isLastOwner, err
	}

	for _, user := range workspace.Users {
		if user.UserID != userID && user.Role == models.UserRoleOwner {
			isLastOwner = false
			break
		}
	}

	return isLastOwner, nil
}

func (wr *WorkspaceRepo) GetUploaderConfig(ctx context.Context, workspaceID string) (*models.UploaderConfig, error) {
	id, err := primitive.ObjectIDFromHex(workspaceID)
	if err != nil {
		return nil, fmt.Errorf(msg.InvalidObjectID, workspaceID)
	}
	collection := db.Collection(wr.collectionName)
	var workspace models.Workspace
	err = collection.FindOne(ctx, bson.M{"_id": id}, options.FindOne().SetProjection(bson.M{"uploaderConfig": 1})).Decode(&workspace)
	if err != nil {
		return nil, err
	}
	return workspace.UploaderConfig, nil
}

func (wr *WorkspaceRepo) SetUploaderConfig(ctx context.Context, workspaceID string, config *models.UploaderConfig) error {
	id, err := primitive.ObjectIDFromHex(workspaceID)
	if err != nil {
		return fmt.Errorf(msg.InvalidObjectID, workspaceID)
	}

	user, err := utils.GetUserDetailsFromContext(ctx)
	if err != nil {
		return err
	}
	updatedBy := user.Email

	updateFields := utils.FilterNonEmptyBsonFields(bson.M{
		"uploaderConfig.maxFileSize":      config.MaxFileSize,
		"uploaderConfig.minFileSize":      config.MinFileSize,
		"uploaderConfig.maxNumberOfFiles": config.MaxNumberOfFiles,
		"uploaderConfig.minNumberOfFiles": config.MinNumberOfFiles,
		"uploaderConfig.maxTotalFileSize": config.MaxTotalFileSize,
	})

	updateFields["uploaderConfig.allowedFileTypes"] = config.AllowedFileTypes
	updateFields["uploaderConfig.allowedSources"] = config.AllowedSources
	updateFields["uploaderConfig.requiredMetadataFields"] = config.RequiredMetadataFields
	updateFields["uploaderConfig.allowPauseAndResume"] = config.AllowPauseAndResume
	updateFields["uploaderConfig.enableImageEditing"] = config.EnableImageEditing
	updateFields["uploaderConfig.useCompression"] = config.UseCompression
	updateFields["uploaderConfig.useFaultTolerantMode"] = config.UseFaultTolerantMode
	updateFields["uploaderConfig.authEndpoint"] = config.AuthEndpoint
	updateFields["updatedBy"] = updatedBy
	updateFields["updatedAt"] = primitive.NewDateTimeFromTime(time.Now())

	update := bson.M{"$set": updateFields}

	collection := db.Collection(wr.collectionName)
	_, err = collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}
