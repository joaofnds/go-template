package logger

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

var Module = fx.Options(
	fx.Provide(func() (*zap.Logger, error) {
		zapConfig := zap.NewProductionConfig()
		zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		return zapConfig.Build()
	}),
	fx.Provide(func(logger *zap.Logger) *zap.SugaredLogger {
		return logger.Sugar()
	}),
	fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
		return &fxevent.ZapLogger{Logger: logger}
	}),
)
