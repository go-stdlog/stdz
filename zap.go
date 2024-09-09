package stdz

import (
	"fmt"
	"github.com/go-stdlog/stdlog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var lmap = map[stdlog.Level]zapcore.Level{
	stdlog.LevelDebug:   zapcore.DebugLevel,
	stdlog.LevelInfo:    zapcore.InfoLevel,
	stdlog.LevelWarning: zapcore.WarnLevel,
	stdlog.LevelError:   zapcore.ErrorLevel,
}

type Z struct {
	*zap.Logger
	cfg zap.Config
}

func New(cfg zap.Config) (*Z, error) {
	l, err := cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}
	return &Z{l, cfg}, nil
}

func (z *Z) Named(name string) stdlog.Logger {
	n := new(Z)
	*n = *z
	n.Logger = n.Logger.Named(name)
	return n
}

func (z *zLogger) SetLevel(level stdlog.Level) {
	l := z.Logger.Level()
	if err := l.Set(level.String()); err != nil {
		panic(err)
	}
}

func (z *Z) Debug(msg string, fields ...any) {
	z.Logger.Debug(msg, handleFields(stdlog.LevelDebug, fields)...)
}

func (z *Z) Info(msg string, fields ...any) {
	z.Logger.Info(msg, handleFields(stdlog.LevelInfo, fields)...)
}

func (z *Z) Warning(msg string, fields ...any) {
	z.Logger.Warn(msg, handleFields(stdlog.LevelWarning, fields)...)
}

func (z *Z) Error(err error, msg string, fields ...any) {
	z.Logger.Error(msg, handleFields(stdlog.LevelError, append(fields, zap.Error(err)))...)
}

func (z *Z) Fatal(msg string, fields ...any) {
	z.Logger.Fatal(msg, handleFields(stdlog.LevelFatal, fields)...)
}

func (z *Z) FatalError(err error, msg string, fields ...any) {
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
