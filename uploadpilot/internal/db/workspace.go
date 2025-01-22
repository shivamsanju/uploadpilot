package db

import (
	"context"
	"fmt"
	"time"

	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WorkspaceRepo interface {
	GetWorkspace(ctx context.Context, workspaceID primitive.ObjectID) (*models.Workspace, error)
	Create(ctx context.Context, workspace *models.Workspace) (*models.Workspace, error)
	Delete(ctx context.Context, workspaceID primitive.ObjectID) error
	CheckWorkspaceExists(ctx context.Context, workspaceID primitive.ObjectID) (bool, error)

	// users related
	GetUsersInWorkspace(ctx context.Context, workspaceID primitive.ObjectID) ([]models.WorkspaceUserWithDetails, error)
	GetWorkspacesForUser(ctx context.Context, userId string) ([]models.Workspace, error)
	CheckIfUserExistsInWorkspace(ctx context.Context, workspaceID primitive.ObjectID, userID string) (bool, error)
	AddUserToWorkspace(ctx context.Context, workspaceID primitive.ObjectID, user *models.WorkspaceUser) error
	RemoveUserFromWorkspace(ctx context.Context, workspaceID primitive.ObjectID, userID string) error
	GetUserRoleInWorkspace(ctx context.Context, workspaceID primitive.ObjectID, userID string) (models.UserRole, error)
	ChangeUserRoleInWorkspace(ctx context.Context, workspaceID primitive.ObjectID, userID string, role models.UserRole) error

	// uploader config related
	GetUploaderConfig(ctx context.Context, workspaceID primitive.ObjectID) (*models.UploaderConfig, error)
	UpdateUploaderConfig(ctx context.Context, workspaceID primitive.ObjectID, cb *models.UploaderConfig) error
}

type workspaceRepo struct {
	collectionName     string
	userColelctionName string
}

// NewWorkspaceRepo initializes a new instance of workspaceRepo with predefined
// collection names for workspaces and users. It returns an interface that
// defines the contract for workspace-related operations.
func NewWorkspaceRepo() WorkspaceRepo {
	return &workspaceRepo{
		collectionName:     "workspaces",
		userColelctionName: "users",
	}
}

// GetWorkspace retrieves the workspace with the given ID from the workspaces collection.
func (wr *workspaceRepo) GetWorkspace(ctx context.Context, workspaceID primitive.ObjectID) (*models.Workspace, error) {
	collection := db.Collection(wr.collectionName)
	var workspace models.Workspace
	err := collection.FindOne(ctx, bson.M{"_id": workspaceID}).Decode(&workspace)
	if err != nil {
		return nil, err
	}
	return &workspace, nil
}

// Create creates a new workspace and adds the given user as owner to it.
//
// This method uses a transaction to ensure that either both the workspace and the user are updated,
// or neither of them is.
//
// Note that the ID of the workspace is set by this method.
func (wr *workspaceRepo) Create(ctx context.Context, workspace *models.Workspace) (*models.Workspace, error) {
	creatorEmail := ctx.Value("email").(string)
	creatorUserID := ctx.Value("userId").(string)

	workspace.ID = primitive.NewObjectID()
	workspace.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	workspace.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	workspace.CreatedBy = creatorEmail
	workspace.UpdatedBy = creatorEmail

	workspace.Users = []models.WorkspaceUser{{UserID: creatorUserID, Role: models.UserRoleOwner}}

	session, err := db.Client().StartSession()
	if err != nil {
		return nil, err
	}

	defer session.EndSession(ctx)
	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		collection := db.Collection(wr.collectionName)
		result, err := collection.InsertOne(ctx, workspace)
		if err != nil {
			return nil, err
		}

		// Update user collection
		userCollection := db.Collection(wr.userColelctionName)
		userFilter := bson.M{"userId": creatorUserID}
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

// Delete deletes the workspace with the given ID.
//
// The deletion is performed by finding and removing the document with the given ID from the workspaces collection.
//
// It is the caller's responsibility to ensure that the workspace is not in use by any other operations.
func (wr *workspaceRepo) Delete(ctx context.Context, workspaceID primitive.ObjectID) error {
	collection := db.Collection(wr.collectionName)
	_, err := collection.DeleteOne(ctx, bson.M{"_id": workspaceID})
	return err
}

// CheckWorkspaceExists checks if a workspace with the given ID exists in the workspaces collection.
//
// It returns true if the workspace exists, false otherwise.
func (wr *workspaceRepo) CheckWorkspaceExists(ctx context.Context, workspaceID primitive.ObjectID) (bool, error) {
	collection := db.Collection(wr.collectionName)
	count, err := collection.CountDocuments(ctx, bson.M{"_id": workspaceID})
	return count > 0, err
}

// GetUsersInWorkspace retrieves a list of users associated with the given workspace ID.
func (wr *workspaceRepo) GetUsersInWorkspace(ctx context.Context, workspaceID primitive.ObjectID) ([]models.WorkspaceUserWithDetails, error) {
	collection := db.Collection(wr.collectionName)
	var users []models.WorkspaceUserWithDetails

	pipeline := []bson.M{
		{
			"$match": bson.M{"_id": workspaceID},
		},
		{
			"$unwind": "$users", // Unwind the users array to process each user individually
		},
		{
			"$lookup": bson.M{
				"from":         wr.userColelctionName, // Collection containing user details
				"localField":   "users.userId",        // Match workspace userId
				"foreignField": "userId",              // With user collection's userId
				"as":           "userDetails",         // Output array of matched user details
			},
		},
		{
			"$unwind": "$userDetails", // Unwind userDetails array to access individual user details
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

// GetWorkspacesForUser retrieves a list of workspaces associated with the specified user ID.
// It returns a slice of models.Workspace containing only the ID and name fields, sorted by the
// updatedAt field in descending order. If an error occurs during the database query, it returns
// an error.
func (wr *workspaceRepo) GetWorkspacesForUser(ctx context.Context, userId string) ([]models.Workspace, error) {
	collection := db.Collection(wr.collectionName)
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

func (wr *workspaceRepo) CheckIfUserExistsInWorkspace(ctx context.Context, workspaceID primitive.ObjectID, userID string) (bool, error) {
	collection := db.Collection(wr.collectionName)
	count, err := collection.CountDocuments(ctx, bson.M{"_id": workspaceID, "users.userId": userID})
	return count > 0, err
}

// AddUserToWorkspace adds the given user to the given workspace using a transaction.
//
// The transaction ensures that either both the workspace and the user are updated,
// or neither of them is.
func (wr *workspaceRepo) AddUserToWorkspace(ctx context.Context, workspaceID primitive.ObjectID, user *models.WorkspaceUser) error {
	session, err := db.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	// Transactional function
	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Update workspace collection
		workspaceCollection := db.Collection(wr.collectionName)
		workspaceFilter := bson.M{"_id": workspaceID}
		workspaceUpdate := bson.M{"$addToSet": bson.M{"users": user}}
		if _, err := workspaceCollection.UpdateOne(sessCtx, workspaceFilter, workspaceUpdate); err != nil {
			return nil, err
		}

		// Update user collection
		userCollection := db.Collection(wr.userColelctionName)
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

// RemoveUserFromWorkspace removes the given user from the given workspace using a transaction.
//
// The transaction ensures that either both the workspace and the user are updated,
// or neither of them is.
func (wr *workspaceRepo) RemoveUserFromWorkspace(ctx context.Context, workspaceID primitive.ObjectID, userID string) error {
	session, err := db.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	userRole, err := wr.GetUserRoleInWorkspace(ctx, workspaceID, userID)
	if err != nil || userRole == models.UserRoleNotFound {
		return fmt.Errorf("user %s not found in workspace %s", userID, workspaceID)
	}

	// Transactional function
	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Update workspace collection
		workspaceCollection := db.Collection(wr.collectionName)
		workspaceFilter := bson.M{"_id": workspaceID}
		workspaceUpdate := bson.M{"$pull": bson.M{"users": &models.WorkspaceUser{UserID: userID, Role: userRole}}}
		if _, err := workspaceCollection.UpdateOne(sessCtx, workspaceFilter, workspaceUpdate); err != nil {
			return nil, err
		}

		// Update user collection
		userCollection := db.Collection(wr.userColelctionName)
		userFilter := bson.M{"userId": userID}
		userUpdate := bson.M{"$pull": bson.M{"workspaces": &models.UserWorkspace{WorkspaceID: workspaceID, Role: userRole}}}
		if _, err := userCollection.UpdateOne(sessCtx, userFilter, userUpdate); err != nil {
			return nil, err
		}

		return nil, nil
	}

	// Execute the transaction
	_, err = session.WithTransaction(ctx, callback, options.Transaction())
	return err
}

// GetUserRoleInWorkspace retrieves the role of the user in the given workspace.
func (wr *workspaceRepo) GetUserRoleInWorkspace(ctx context.Context, workspaceID primitive.ObjectID, userID string) (models.UserRole, error) {
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

func (wr *workspaceRepo) ChangeUserRoleInWorkspace(ctx context.Context, workspaceID primitive.ObjectID, userID string, role models.UserRole) error {
	collection := db.Collection(wr.collectionName)
	_, err := collection.UpdateOne(ctx, bson.M{"_id": workspaceID, "users.userId": userID}, bson.M{"$set": bson.M{"users.$.role": role}})
	return err
}

// GetUploaderConfig retrieves the uploader configuration associated with the given workspace ID.
//
// The retrieval is performed by finding the workspace document with the given ID and extracting the
// uploader configuration from it.
//
// If the ID is not a valid ObjectID, an error is returned.
//
// It is the caller's responsibility to ensure that the workspace is not in use by any other operations.
func (wr *workspaceRepo) GetUploaderConfig(ctx context.Context, workspaceID primitive.ObjectID) (*models.UploaderConfig, error) {
	collection := db.Collection(wr.collectionName)
	var workspace models.Workspace
	err := collection.FindOne(ctx, bson.M{"_id": workspaceID}, options.FindOne().SetProjection(bson.M{"uploaderConfig": 1})).Decode(&workspace)
	if err != nil {
		return nil, err
	}
	return workspace.UploaderConfig, nil
}

// UpdateUploaderConfig updates the uploader configuration for a specific workspace.
//
// This function takes the workspace ID, updated uploader configuration data, and the identifier
// of the user making the update. It prepares a BSON update object containing the non-empty fields of the
// provided configuration. The function then updates the workspace document in the database by
// setting the specified fields in the uploader configuration, and records the user who made the
// update along with the timestamp of the update.
func (wr *workspaceRepo) UpdateUploaderConfig(ctx context.Context, workspaceID primitive.ObjectID, updatedData *models.UploaderConfig) error {
	updatedBy := ctx.Value("email").(string)

	updateFields := utils.FilterNonEmptyBsonFields(bson.M{
		"uploaderConfig.maxFileSize":            updatedData.MaxFileSize,
		"uploaderConfig.minFileSize":            updatedData.MinFileSize,
		"uploaderConfig.maxNumberOfFiles":       updatedData.MaxNumberOfFiles,
		"uploaderConfig.minNumberOfFiles":       updatedData.MinNumberOfFiles,
		"uploaderConfig.maxTotalFileSize":       updatedData.MaxTotalFileSize,
		"uploaderConfig.allowedFileTypes":       updatedData.AllowedFileTypes,
		"uploaderConfig.allowedSources":         updatedData.AllowedSources,
		"uploaderConfig.requiredMetadataFields": updatedData.RequiredMetadataFields,
		"uploaderConfig.allowPauseAndResume":    updatedData.AllowPauseAndResume,
		"uploaderConfig.enableImageEditing":     updatedData.EnableImageEditing,
		"uploaderConfig.useCompression":         updatedData.UseCompression,
		"uploaderConfig.useFaultTolerantMode":   updatedData.UseFaultTolerantMode,
		"updatedBy":                             updatedBy,
		"updatedAt":                             primitive.NewDateTimeFromTime(time.Now()),
	})

	update := bson.M{"$set": updateFields}

	collection := db.Collection(wr.collectionName)
	_, err := collection.UpdateOne(ctx, bson.M{"_id": workspaceID}, update)
	return err
}
