package svc

import (
	"context"

	"github.com/uploadpilot/uploadpilot/common/pkg/db/repo"
	"github.com/uploadpilot/uploadpilot/common/pkg/models"
	"github.com/uploadpilot/uploadpilot/manager/internal/dto"
)

type UserService struct {
	userRepo *repo.UserRepo
}

func NewUserService(userRepo *repo.UserRepo) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetUserDetails(ctx context.Context, userID string) (*models.User, error) {
	return s.userRepo.GetByUserID(ctx, userID)
}

func (s *UserService) GetUserDetailsFromContext(ctx context.Context) (*models.User, error) {
	userID := ctx.Value(dto.UserIDContextKey).(string)
	return s.userRepo.GetByUserID(ctx, userID)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	return s.userRepo.Create(ctx, user)
}
