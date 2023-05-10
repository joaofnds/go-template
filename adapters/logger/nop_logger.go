package logger

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

var NopLoggerProvider = fx.Options(
	fx.NopLogger,
	fx.Provide(func() *zap.Logger {
		return zap.NewNop()
	}),
	fx.Provide(func(logger *zap.Logger) *zap.SugaredLogger {
		return logger.Sugar()
	}),
	fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
		return &fxevent.ZapLogger{Logger: logger}
	}),
)
