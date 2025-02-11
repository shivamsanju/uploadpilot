package dbutils

import (
	"errors"

	"github.com/google/uuid"
	"github.com/uploadpilot/uploadpilot/common/pkg/infra"
	"github.com/uploadpilot/uploadpilot/common/pkg/msg"
	"gorm.io/gorm"
)

func DBError(err error) error {
	if err == nil {
		return nil
	}

	infra.Log.Errorf(msg.DBError, err.Error())

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New(msg.DBErrRecordNotFound)
	}
	return errors.New(msg.ErrUnknown)
}

func GenerateUUID() string {
	return uuid.New().String()
}
