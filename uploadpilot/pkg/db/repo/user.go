package repo

import (
	"context"
	"time"

	"github.com/uploadpilot/uploadpilot/pkg/db/models"
	g "github.com/uploadpilot/uploadpilot/pkg/globals"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	CheckUserExists(ctx context.Context, email string) (bool, error)
}

type userRepo struct {
	collectionName string
}

func NewUserRepo() UserRepo {
	return &userRepo{collectionName: "users"}
}

func (u *userRepo) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	collection := g.Db.Database(g.DbName).Collection(u.collectionName)
	user.ID = primitive.NewObjectID()
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = result.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (u *userRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	collection := g.Db.Database(g.DbName).Collection(u.collectionName)
	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepo) UpdateUserToken(ctx context.Context, userId string, signedToken string, signedRefreshToken string) error {
	collection := g.Db.Database(g.DbName).Collection(u.collectionName)

	updateObj := bson.D{
		{"token", signedToken},
		{"refreshToken", signedRefreshToken},
		{"updatedAt", primitive.NewDateTimeFromTime(time.Now())},
	}

	upsert := true
	opt := options.UpdateOptions{Upsert: &upsert}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": userId}, bson.D{{"$set", updateObj}}, &opt)
	if err != nil {
		return err
	}

	return nil
}

func (u *userRepo) CheckUserExists(ctx context.Context, email string) (bool, error) {
	collection := g.Db.Database(g.DbName).Collection(u.collectionName)
	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return false, err
	}
	return true, nil
}
