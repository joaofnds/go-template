package logger

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

var Module = fx.Options(
	fx.Provide(func(config Config) (*zap.Logger, error) {
		zapConfig := zap.NewProductionConfig()

		logLevel, err := zap.ParseAtomicLevel(config.Level)
		if err != nil {
			return nil, err
		}
		zapConfig.Level = logLevel

		return zapConfig.Build()
	}),
	fx.Provide(func(logger *zap.Logger) *zap.SugaredLogger {
		return logger.Sugar()
	}),
	fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
		return &fxevent.ZapLogger{Logger: logger}
	}),
)
