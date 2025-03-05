package repo

import (
	"context"
	"errors"

	"github.com/uploadpilot/core/internal/db/driver"
	"github.com/uploadpilot/core/internal/db/models"
	dbutils "github.com/uploadpilot/core/internal/db/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type APIKeyRepo struct {
	db *driver.Driver
}

func NewAPIKeyRepo(db *driver.Driver) *APIKeyRepo {
	return &APIKeyRepo{
		db: db,
	}
}

func (r *APIKeyRepo) CreateApiKey(ctx context.Context, apiKey *models.APIKey) error {
	if err := r.db.Orm.WithContext(ctx).Create(apiKey).Error; err != nil {
		return dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}
	return nil
}

func (r *APIKeyRepo) GetApiKeyDetailsByID(ctx context.Context, id string) (*models.APIKey, error) {
	var apiKey models.APIKey
	if err := r.db.Orm.WithContext(ctx).First(&apiKey, "id = ?", id).Error; err != nil {
		return nil, dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}
	return &apiKey, nil
}

func (r *APIKeyRepo) GetApiKeyLimitedDetailsByID(ctx context.Context, id string) (*models.APIKey, error) {
	var apiKey models.APIKey
	if err := r.db.Orm.WithContext(ctx).
		Preload("Permissions").
		Select("id", "name", "user_id", "created_at", "expires_at", "revoked_at", "revoked_by").
		First(&apiKey, "id = ?", id).Error; err != nil {
		return nil, dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}
	return &apiKey, nil
}

func (r *APIKeyRepo) GetApiKeyDetails(ctx context.Context, hash string) (*models.APIKey, error) {
	var apiKey models.APIKey
	if err := r.db.Orm.WithContext(ctx).
		Preload("Permissions").
		First(&apiKey, "api_key_hash = ?", hash).Error; err != nil {
		return nil, dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}
	return &apiKey, nil
}

func (r *APIKeyRepo) GetAllApiKeys(ctx context.Context, userID string) ([]models.APIKey, error) {
	var apiKeys []models.APIKey
	if err := r.db.Orm.WithContext(ctx).
		Select("id", "name", "user_id", "created_at", "expires_at", "revoked_at", "revoked_by").
		Order(clause.OrderBy{Columns: []clause.OrderByColumn{
			{Column: clause.Column{Name: "revoked_at"}, Desc: true},
			{Column: clause.Column{Name: "created_at"}, Desc: true},
		}}).
		Find(&apiKeys, "user_id = ?", userID).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return make([]models.APIKey, 0), nil
		}
		return nil, dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}
	return apiKeys, nil
}

func (r *APIKeyRepo) Update(ctx context.Context, apiKey *models.APIKey) error {
	if err := r.db.Orm.WithContext(ctx).Save(apiKey).Error; err != nil {
		return dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}
	return nil
}
