package services

import (
	g "github.com/shivamsanju/uploader/pkg/globals"
	"go.uber.org/zap"
)

func initLogger() {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	g.Log = sugar
	g.Log.Info("logger initialized")
}
