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

func (r *APIKeyRepo) CreateApiKey(ctx context.Context, apiKey *models.APIKey, apiKeyCreateCallback func() error) error {
	tx := r.db.Orm.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create the API key within the transaction
	if err := tx.Create(apiKey).Error; err != nil {
		tx.Rollback()
		return dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}

	// Update user metadata
	if err := apiKeyCreateCallback(); err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction if everything is successful
	if err := tx.Commit().Error; err != nil {
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
		Select("id", "name", "user_id", "created_at", "expires_at", "revoked_at", "revoked_by").
		First(&apiKey, "id = ?", id).Error; err != nil {
		return nil, dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}
	return &apiKey, nil
}

func (r *APIKeyRepo) GetApiKeyDetails(ctx context.Context, hash string) (*models.APIKey, error) {
	var apiKey models.APIKey
	if err := r.db.Orm.WithContext(ctx).
		First(&apiKey, "api_key_hash = ?", hash).Error; err != nil {
		return nil, dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}
	return &apiKey, nil
}

func (r *APIKeyRepo) GetAllApiKeysInTenant(ctx context.Context, tenantID string) ([]models.APIKey, error) {
	var apiKeys []models.APIKey
	if err := r.db.Orm.WithContext(ctx).
		Select("id", "name", "user_id", "created_at", "expires_at", "revoked_at", "revoked_by").
		Order(clause.OrderBy{Columns: []clause.OrderByColumn{
			{Column: clause.Column{Name: "revoked_at"}, Desc: true},
			{Column: clause.Column{Name: "created_at"}, Desc: true},
		}}).
		Where("tenant_id = ?", tenantID).
		Find(&apiKeys).
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
