package logger

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Options(fx.Provide(NewLogger), fx.Provide(NewSugarLogger))

func NewLogger() (*zap.Logger, error) {
	return zap.NewProduction()
}

func NewSugarLogger(logger *zap.Logger) *zap.SugaredLogger {
	return logger.Sugar()
}
