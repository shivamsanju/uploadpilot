package db

import (
	"context"

	"github.com/uploadpilot/uploadpilot/internal/cache"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"github.com/uploadpilot/uploadpilot/internal/utils"
)

type UserRepo struct {
}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (u *UserRepo) Create(ctx context.Context, user *models.User) error {
	dbMutateFn := func(user *models.User) error {
		return sqlDB.WithContext(ctx).Create(user).Error
	}
	cl := cache.NewClient[*models.User](0)
	if err := cl.Mutate(ctx, UserIDKey(user.ID), []string{UserEmailKey(user.Email)}, user, dbMutateFn, 0); err != nil {
		return err
	}

	infra.Log.Infof("created user: %+v", user)
	return nil
}

func (u *UserRepo) GetByUserID(ctx context.Context, userID string) (*models.User, error) {
	dbFetch := func(user *models.User) error {
		if err := sqlDB.WithContext(ctx).Where("id = ?", userID).First(user).Error; err != nil {
			return utils.DBError(err)
		}
		return nil
	}

	var user models.User
	cl := cache.NewClient[*models.User](0)
	if err := cl.Query(ctx, UserIDKey(userID), &user, dbFetch); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	dbFetch := func(user *models.User) error {
		//check type of user

		if err := sqlDB.WithContext(ctx).Where("email = ?", email).First(user).Error; err != nil {
			if err.Error() == "record not found" {
				return nil
			}
			return utils.DBError(err)
		}
		return nil
	}

	var user models.User
	cl := cache.NewClient[*models.User](0)
	if err := cl.Query(ctx, UserEmailKey(email), &user, dbFetch); err != nil {
		return nil, err
	}
	return &user, nil
}

func UserIDKey(userID string) string {
	return "user:" + userID
}

func UserEmailKey(email string) string {
	return "user:" + email
}
