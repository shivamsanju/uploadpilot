package models

import (
	"time"
)

type User struct {
	ID            string    `gorm:"column:id;not null;type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Email         string    `gorm:"column:email;not null;size:255;unique" json:"email,omitempty"`
	EmailVerified bool      `gorm:"column:email_verified;default:false" json:"emailVerified,omitempty"`
	FirstName     *string   `gorm:"column:first_name;size:100" json:"firstName,omitempty"`
	LastName      *string   `gorm:"column:last_name;size:100" json:"lastName,omitempty"`
	Name          *string   `gorm:"column:name" json:"name,omitempty"`
	NickName      *string   `gorm:"column:nick_name;size:100" json:"nickName,omitempty"`
	Phone         *string   `gorm:"column:phone;size:25" json:"phone,omitempty"`
	Provider      *string   `gorm:"column:provider;size:100" json:"provider,omitempty"`
	Description   *string   `gorm:"column:description;size:500" json:"description,omitempty"`
	AvatarURL     *string   `gorm:"column:avatar_url;size:255" json:"avatarUrl,omitempty"`
	Location      *string   `gorm:"column:location;size:255" json:"location,omitempty"`
	IsUserBanned  bool      `gorm:"column:is_user_banned;default:false" json:"isUserBanned,omitempty"`
	BanReason     *string   `gorm:"column:ban_reason;size:255" json:"banReason,omitempty"`
	TrialStartsAt time.Time `gorm:"column:trial_starts_at;not null" json:"trialStartsAt,omitempty"`
	TrialEndsAt   time.Time `gorm:"column:trial_ends_at;not null" json:"trialEndsAt,omitempty"`
	CreatedAtColumn
	UpdatedAtColumn
}

func (User) TableName() string {
	return "users"
}

type UserRole string

const (
	UserRoleOwner       UserRole = "Owner"
	UserRoleContributor UserRole = "Contributor"
	UserRoleViewer      UserRole = "Viewer"
	UserRoleNotFound    UserRole = "Not Found"
)

type UserInfo string

const (
	UserID UserInfo = "ID"
	Email  UserInfo = "Email"
	Name   UserInfo = "Name"
)
