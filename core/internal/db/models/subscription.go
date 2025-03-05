package models

import (
	"time"
)

type Subscription struct {
	ID              string                 `gorm:"column:id;primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	TenantID        string                 `gorm:"column:tenant_id;not null;type:uuid" json:"tenantId"`
	PlanID          string                 `gorm:"column:plan_id;not null;type:uuid" json:"planId"`       // Reference to a subscription plan (if applicable)
	Status          string                 `gorm:"column:status;not null;type:varchar(20)" json:"status"` // e.g., "active", "cancelled", "trial", "expired"
	StartDate       time.Time              `gorm:"column:start_date;not null" json:"startDate"`
	EndDate         *time.Time             `gorm:"column:end_date" json:"endDate,omitempty"`            // When the subscription naturally expires
	TrialEndDate    *time.Time             `gorm:"column:trial_end_date" json:"trialEndDate,omitempty"` // For trial subscriptions
	NextBillingDate *time.Time             `gorm:"column:next_billing_date" json:"nextBillingDate,omitempty"`
	CancelledAt     *time.Time             `gorm:"column:cancelled_at" json:"cancelledAt,omitempty"`
	Metadata        map[string]interface{} `gorm:"column:metadata;type:jsonb" json:"metadata,omitempty"`
	CreatedAtColumn
	UpdatedAtColumn
}

// TableName overrides the table name used by GORM.
func (Subscription) TableName() string {
	return "subscriptions"
}
