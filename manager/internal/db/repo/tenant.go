package repo

import (
	"context"

	"github.com/uploadpilot/manager/internal/db/driver"
	"github.com/uploadpilot/manager/internal/db/models"
	dbutils "github.com/uploadpilot/manager/internal/db/utils"
)

type TenantRepo struct {
	db *driver.Driver
}

func NewTenantRepo(db *driver.Driver) *TenantRepo {
	return &TenantRepo{
		db: db,
	}
}

func (r *TenantRepo) Get(ctx context.Context, tenantID string) (*models.Tenant, error) {
	var tenant models.Tenant
	err := r.db.Orm.WithContext(ctx).First(&tenant, "id = ?", tenantID).Error

	if err != nil {
		return nil, dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}

	return &tenant, nil
}

func (r *TenantRepo) GetAll(ctx context.Context) ([]models.Tenant, error) {
	var tenants []models.Tenant
	err := r.db.Orm.WithContext(ctx).Find(&tenants).Error

	if err != nil {
		return nil, dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}

	return tenants, nil
}

func (r *TenantRepo) Create(ctx context.Context, tenant *models.Tenant, tenantUpdateCallback func(*models.Tenant) error) error {
	tx := r.db.Orm.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create the tenant within the transaction
	if err := tx.Create(tenant).Error; err != nil {
		tx.Rollback()
		return dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}

	// Update user metadata
	if err := tenantUpdateCallback(tenant); err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction if everything is successful
	if err := tx.Commit().Error; err != nil {
		return dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}
	return nil
}

func (r *TenantRepo) Delete(ctx context.Context, tenantID string) error {
	if err := r.db.Orm.WithContext(ctx).Delete(&models.Tenant{}, "id = ?", tenantID).Error; err != nil {
		return dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}
	return nil
}
