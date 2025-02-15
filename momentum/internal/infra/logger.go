package infra

import "go.uber.org/zap"

func NewLogger() (*zap.SugaredLogger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	defer logger.Sync()
	zap.ReplaceGlobals(logger)
	sugar := logger.Sugar()
	return sugar, nil

}
