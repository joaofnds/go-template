package postgres

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
	gormlogger "gorm.io/gorm/logger"
)

type ZapLoggerAdapter struct {
	logger *zap.Logger
	config gormlogger.Config
}

func NewZapLoggerAdapter(logger *zap.Logger, config gormlogger.Config) gormlogger.Interface {
	return &ZapLoggerAdapter{logger: logger, config: config}
}

func (adapter *ZapLoggerAdapter) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	adapter.config.LogLevel = level
	return adapter
}

func (adapter *ZapLoggerAdapter) Info(_ context.Context, msg string, args ...interface{}) {
	if adapter.config.LogLevel >= gormlogger.Info {
		adapter.logger.Sugar().Infof(msg, args...)
	}
}

func (adapter *ZapLoggerAdapter) Warn(_ context.Context, msg string, args ...interface{}) {
	if adapter.config.LogLevel >= gormlogger.Warn {
		adapter.logger.Sugar().Warnf(msg, args...)
	}
}

func (adapter *ZapLoggerAdapter) Error(_ context.Context, msg string, args ...interface{}) {
	if adapter.config.LogLevel >= gormlogger.Error {
		adapter.logger.Sugar().Errorf(msg, args...)
	}
}

func (adapter *ZapLoggerAdapter) Trace(_ context.Context, begin time.Time, f func() (string, int64), err error) {
	if adapter.config.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)

	sql, rows := f()
	fields := []zap.Field{
		zap.Error(err),
		zap.String("location", fileWithLineNum()),
		zap.String("elapsed", elapsed.String()),
		zap.String("sql", sql),
		zap.Int64("rows", rows),
	}

	switch {
	case err != nil && adapter.config.LogLevel >= gormlogger.Error:
		adapter.logger.Error("", fields...)
	case elapsed >= adapter.config.SlowThreshold && adapter.config.LogLevel >= gormlogger.Warn:
		adapter.logger.Warn("", fields...)
	case adapter.config.LogLevel >= gormlogger.Info:
		adapter.logger.Info("", fields...)
	}
}

func fileWithLineNum() string {
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok && !strings.Contains(file, "gorm.io") && !strings.HasSuffix(file, "_test.go") {
			return fmt.Sprintf("%s:%d", file, line)
		}
	}
	return ""
}
