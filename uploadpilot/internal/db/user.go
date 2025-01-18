package db

import (
	"context"

	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	CheckUserExists(ctx context.Context, email string) (bool, error)
}

type userRepo struct {
	collectionName string
}

func NewUserRepo() UserRepo {
	return &userRepo{collectionName: "users"}
}

func (u *userRepo) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	collection := db.Collection(u.collectionName)
	user.ID = primitive.NewObjectID()
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = result.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (u *userRepo) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	collection := db.Collection(u.collectionName)
	var user models.User
	err := collection.FindOne(ctx, bson.M{"userId": userID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepo) CheckUserExists(ctx context.Context, userID string) (bool, error) {
	collection := db.Collection(u.collectionName)
	var user models.User
	err := collection.FindOne(ctx, bson.M{"userId": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
