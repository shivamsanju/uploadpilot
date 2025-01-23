package infra

import (
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

var (
	Log      *zap.SugaredLogger
	Validate = validator.New()
)
