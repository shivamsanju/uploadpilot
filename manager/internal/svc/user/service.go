package user

import (
	"context"

	"github.com/uploadpilot/go-core/db/pkg/models"
	"github.com/uploadpilot/go-core/db/pkg/repo"
	"github.com/uploadpilot/manager/internal/dto"
)

type Service struct {
	userRepo *repo.UserRepo
}

func NewService(userRepo *repo.UserRepo) *Service {
	return &Service{userRepo: userRepo}
}

func (s *Service) GetUserDetails(ctx context.Context, userID string) (*models.User, error) {
	return s.userRepo.GetByUserID(ctx, userID)
}

func (s *Service) GetUserDetailsFromContext(ctx context.Context) (*models.User, error) {
	userID := ctx.Value(dto.UserIDContextKey).(string)
	return s.userRepo.GetByUserID(ctx, userID)
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

func (s *Service) CreateUser(ctx context.Context, user *models.User) error {
	return s.userRepo.Create(ctx, user)
}
