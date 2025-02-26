package models

import (
	"time"
)

type User struct {
	ID            string    `gorm:"column:id;not null;type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Email         string    `gorm:"column:email;not null;size:255;unique" json:"email"`
	EmailVerified bool      `gorm:"column:email_verified;default:false" json:"emailVerified"`
	FirstName     *string   `gorm:"column:first_name;size:100" json:"firstName"`
	LastName      *string   `gorm:"column:last_name;size:100" json:"lastName"`
	Name          *string   `gorm:"column:name" json:"name"`
	NickName      *string   `gorm:"column:nick_name;size:100" json:"nickName"`
	Phone         *string   `gorm:"column:phone;size:25" json:"phone"`
	Provider      *string   `gorm:"column:provider;size:100" json:"provider"`
	Description   *string   `gorm:"column:description;size:500" json:"description"`
	AvatarURL     *string   `gorm:"column:avatar_url;size:255" json:"avatarUrl"`
	Location      *string   `gorm:"column:location;size:255" json:"location"`
	IsUserBanned  bool      `gorm:"column:is_user_banned;default:false" json:"isUserBanned"`
	BanReason     *string   `gorm:"column:ban_reason;size:255" json:"banReason"`
	TrialStartsAt time.Time `gorm:"column:trial_starts_at;not null" json:"trialStartsAt"`
	TrialEndsAt   time.Time `gorm:"column:trial_ends_at;not null" json:"trialEndsAt"`
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
