package repo

import (
	"context"

	"github.com/uploadpilot/manager/internal/db/driver"
	"github.com/uploadpilot/manager/internal/db/models"
	dbutils "github.com/uploadpilot/manager/internal/db/utils"
)

type SubscriptionRepo struct {
	db *driver.Driver
}

func NewSubscriptionRepo(db *driver.Driver) *SubscriptionRepo {
	return &SubscriptionRepo{
		db: db,
	}
}

func (r *SubscriptionRepo) Get(ctx context.Context, subscriptionID string) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.db.Orm.WithContext(ctx).First(&subscription, "id = ?", subscriptionID).Error

	if err != nil {
		return nil, dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}

	return &subscription, nil
}

func (r *SubscriptionRepo) GetSubscriptionByTenantID(ctx context.Context, tenantID string) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.db.Orm.WithContext(ctx).First(&subscription, "tenant_id = ?", tenantID).Error

	if err != nil {
		return nil, dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}

	return &subscription, nil
}

func (r *SubscriptionRepo) Create(ctx context.Context, subscription *models.Subscription) error {
	err := r.db.Orm.WithContext(ctx).Create(subscription).Error
	if err != nil {
		return dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}
	return nil
}

func (r *SubscriptionRepo) Delete(ctx context.Context, subscriptionID string) error {
	if err := r.db.Orm.WithContext(ctx).Delete(&models.Subscription{}, "id = ?", subscriptionID).Error; err != nil {
		return dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}
	return nil
}

func (r *SubscriptionRepo) Update(ctx context.Context, subscriptionID string, subscription *models.Subscription) error {
	if err := r.db.Orm.WithContext(ctx).Save(subscription).Error; err != nil {
		return dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}
	return nil
}
