// copied from https://github.com/garsue/watermillzap/blob/master/adapter.go
package watermill

import (
	"github.com/ThreeDotsLabs/watermill"
	"go.uber.org/zap"
)

type Logger struct {
	backend *zap.Logger
	fields  watermill.LogFields
}

func newLogger(z *zap.Logger) watermill.LoggerAdapter {
	return &Logger{backend: z}
}

func (l *Logger) Error(msg string, err error, fields watermill.LogFields) {
	fields = l.fields.Add(fields)
	fs := make([]zap.Field, 0, len(fields)+1)
	fs = append(fs, zap.Error(err))
	for k, v := range fields {
		fs = append(fs, zap.Any(k, v))
	}
	l.backend.Error(msg, fs...)
}

func (l *Logger) Info(msg string, fields watermill.LogFields) {
	fields = l.fields.Add(fields)
	fs := make([]zap.Field, 0, len(fields)+1)
	for k, v := range fields {
		fs = append(fs, zap.Any(k, v))
	}
	l.backend.Info(msg, fs...)
}

func (l *Logger) Debug(msg string, fields watermill.LogFields) {
	fields = l.fields.Add(fields)
	fs := make([]zap.Field, 0, len(fields)+1)
	for k, v := range fields {
		fs = append(fs, zap.Any(k, v))
	}
	l.backend.Debug(msg, fs...)
}

func (l *Logger) Trace(msg string, fields watermill.LogFields) {
	fields = l.fields.Add(fields)
	fs := make([]zap.Field, 0, len(fields)+1)
	for k, v := range fields {
		fs = append(fs, zap.Any(k, v))
	}
	l.backend.Debug(msg, fs...)
}

func (l *Logger) With(fields watermill.LogFields) watermill.LoggerAdapter {
	return &Logger{
		backend: l.backend,
		fields:  l.fields.Add(fields),
	}
}
