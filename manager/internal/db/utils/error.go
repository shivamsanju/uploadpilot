package dbutils

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/uploadpilot/manager/internal/db/errs"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func DBError(ctx context.Context, logger logger.Interface, err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errs.ErrRecordNotFound
	}

	logger.Error(ctx, "[db_error]", err)
	return errs.ErrUnknownDBError
}

func GenerateUUID() string {
	return uuid.New().String()
}
