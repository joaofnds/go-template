package logger

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

var Module = fx.Options(
	fx.Provide(NewLogger),
	fx.Provide(NewSugarLogger),
	fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
		return &fxevent.ZapLogger{Logger: logger}
	}),
)

func NewLogger() (*zap.Logger, error) {
	return zap.NewProduction()
}

func NewSugarLogger(logger *zap.Logger) *zap.SugaredLogger {
	return logger.Sugar()
}
