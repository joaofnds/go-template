package queue

import (
	"fmt"

	"go.uber.org/zap"
)

type AsyncZapLogger struct {
	zap *zap.Logger
}

func NewAsyncZapLogger(logger *zap.Logger) *AsyncZapLogger {
	return &AsyncZapLogger{logger}
}

func (adapter AsyncZapLogger) Debug(args ...interface{}) {
	adapter.zap.Debug(fmt.Sprint(args...))
}

func (adapter AsyncZapLogger) Info(args ...interface{}) {
	adapter.zap.Info(fmt.Sprint(args...))
}

func (adapter AsyncZapLogger) Warn(args ...interface{}) {
	adapter.zap.Warn(fmt.Sprint(args...))
}

func (adapter AsyncZapLogger) Error(args ...interface{}) {
	adapter.zap.Error(fmt.Sprint(args...))
}

func (adapter AsyncZapLogger) Fatal(args ...interface{}) {
	adapter.zap.Fatal(fmt.Sprint(args...))
}
