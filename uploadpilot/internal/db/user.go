package db

import (
	"context"

	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	CheckUserExists(ctx context.Context, email string) (bool, error)
	GetUserProvider(ctx context.Context, email string) (string, error)
}

type userRepo struct {
	collectionName string
}

func NewUserRepo() UserRepo {
	return &userRepo{collectionName: "users"}
}

// CreateUser inserts a new user into the database, generating a unique ObjectID for them.
// It also initializes the user's Workspaces field as an empty slice.
// Returns the created user with its ID populated, or an error if the insertion fails.
func (u *userRepo) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	collection := db.Collection(u.collectionName)
	user.ID = primitive.NewObjectID()
	user.Workspaces = make([]models.UserWorkspace, 0)
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = result.InsertedID.(primitive.ObjectID)
	return user, nil
}

// GetUserByID retrieves a user document from the database by the given user ID.
// It returns the user data as a models.User object, or an error if the retrieval fails.
// The user ID is the unique identifier associated with the user, not the same as the user's email.
func (u *userRepo) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	collection := db.Collection(u.collectionName)
	var user models.User
	err := collection.FindOne(ctx, bson.M{"userId": userID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail retrieves a user document from the database by the given email address.
// It returns the user data as a models.User object, or an error if the retrieval fails.
func (u *userRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	collection := db.Collection(u.collectionName)
	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CheckUserExists checks if a user exists in the database with the given email address.
// It returns true if a user is found, false otherwise. If an error occurs during the
// database query, it returns false and the error.
func (u *userRepo) CheckUserExists(ctx context.Context, email string) (bool, error) {
	collection := db.Collection(u.collectionName)
	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetUserProvider retrieves the provider associated with a given email from the database.
// It returns the provider as a string, or an empty string if no document is found with the given email.
// If an error occurs during the database query, it returns an empty string and the error.

func (u *userRepo) GetUserProvider(ctx context.Context, email string) (string, error) {
	collection := db.Collection(u.collectionName)
	var user struct {
		Provider string `bson:"provider"`
	}
	err := collection.FindOne(ctx, bson.M{"email": email}, options.FindOne().SetProjection(bson.M{"provider": 1})).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", nil
		}
		return "", err
	}
	return user.Provider, nil
}
