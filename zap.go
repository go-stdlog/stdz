package stdz

import (
	"fmt"
	"github.com/go-stdlog/stdlog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zLogger struct {
var lmap = map[stdlog.Level]zapcore.Level{
	stdlog.LevelDebug:   zapcore.DebugLevel,
	stdlog.LevelInfo:    zapcore.InfoLevel,
	stdlog.LevelWarning: zapcore.WarnLevel,
	stdlog.LevelError:   zapcore.ErrorLevel,
}
	*zap.Logger
}

func New(logger *zap.Logger) stdlog.Logger {
	return &zLogger{logger}
}

func (z *zLogger) Named(name string) stdlog.Logger {
	return New(z.Logger.Named(name))
}

func (z *zLogger) SetLevel(level stdlog.Level) {
	l := z.Logger.Level()
	if err := l.Set(level.String()); err != nil {
		panic(err)
	}
}

func (z *zLogger) Debug(msg string, fields ...any) {
	z.Logger.Debug(msg, handleFields(stdlog.LevelDebug, fields)...)
}

func (z *zLogger) Info(msg string, fields ...any) {
	z.Logger.Info(msg, handleFields(stdlog.LevelInfo, fields)...)
}

func (z *zLogger) Warning(msg string, fields ...any) {
	z.Logger.Warn(msg, handleFields(stdlog.LevelWarning, fields)...)
}

func (z *zLogger) Error(err error, msg string, fields ...any) {
	z.Logger.Error(msg, handleFields(stdlog.LevelError, append(fields, zap.Error(err)))...)
}

func (z *zLogger) Fatal(msg string, fields ...any) {
	z.Logger.Fatal(msg, handleFields(stdlog.LevelFatal, fields)...)
}

func (z *zLogger) FatalError(err error, msg string, fields ...any) {
	z.Logger.Fatal(msg, handleFields(stdlog.LevelFatal, append(fields, zap.Error(err)))...)
}

func handleFields(level stdlog.Level, kvs []any) []zap.Field {
	if len(kvs)%2 != 0 {
		panic(fmt.Errorf("uneven keys and values passed to %s", level.String()))
	}

	fields := make([]zap.Field, 0, len(kvs)/2)

	for i := 0; i < len(kvs)-1; i += 2 {
		k := fmt.Sprintf("%v", kvs[i])
		v := fmt.Sprintf("%v", kvs[i+1])
		fields = append(fields, zap.Any(k, v))
	}

	return fields
}
