package test

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var NopLogger = fx.Options(
	fx.NopLogger,
	fx.Provide(func() *zap.Logger {
		return zap.NewNop()
	}),
	fx.Provide(func(logger *zap.Logger) *zap.SugaredLogger {
		return logger.Sugar()
	}),
)
