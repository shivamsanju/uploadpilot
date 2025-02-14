package repo

import (
	"context"

	"github.com/uploadpilot/uploadpilot/common/pkg/db"
	dbutils "github.com/uploadpilot/uploadpilot/common/pkg/db/utils"
	"github.com/uploadpilot/uploadpilot/common/pkg/models"
)

type UserRepo struct {
	db *db.DB
}

func NewUserRepo(db *db.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) Create(ctx context.Context, user *models.User) error {
	if err := r.db.Orm.WithContext(ctx).Create(user).Error; err != nil {
		return dbutils.DBError(err)
	}
	return nil
}

func (r *UserRepo) GetByUserID(ctx context.Context, userID string) (*models.User, error) {
	var user models.User
	if err := r.db.Orm.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, dbutils.DBError(err)
	}
	return &user, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.Orm.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, dbutils.DBError(err)
	}
	return &user, nil
}
