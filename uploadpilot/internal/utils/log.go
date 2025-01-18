package utils

import (
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"go.uber.org/zap"
)

func InitLogger() error {
	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	infra.Log = sugar
	infra.Log.Info("logger initialized")
	return nil
}
