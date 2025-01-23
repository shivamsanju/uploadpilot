package db

import (
	"context"

	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepo struct {
	collectionName string
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		collectionName: "users",
	}
}

func (u *UserRepo) Create(ctx context.Context, user *models.User) (*models.User, error) {
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

func (u *UserRepo) GetByUserID(ctx context.Context, userID string) (*models.User, error) {
	collection := db.Collection(u.collectionName)
	var user models.User
	err := collection.FindOne(ctx, bson.M{"userId": userID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	collection := db.Collection(u.collectionName)
	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserRepo) EmailExists(ctx context.Context, email string) (bool, error) {
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

func (u *UserRepo) GetProvider(ctx context.Context, email string) (string, error) {
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
