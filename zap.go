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
	cfg       zap.Config
	stackSkip uint
}

func New(cfg zap.Config) (*Z, error) {
	l, err := cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}
	return &Z{l, cfg, 0}, nil
}

func (z *Z) Named(name string) stdlog.Logger {
	n := new(Z)
	*n = *z
	n.Logger = n.Logger.Named(name)
	return n
}

func (z *Z) Skipping(count uint) stdlog.Logger {
	n := new(Z)
	*n = *z
	n.stackSkip = count
	l, err := z.cfg.Build(zap.AddCallerSkip(1 + int(z.stackSkip)))
	if err != nil {
		z.FatalError(err, "Failed to rebuild config with new stack skip")
	}

	if z.Logger.Name() != "" {
		n.Logger = l.Named(z.Logger.Name())
	} else {
		n.Logger = l
	}
	return n

}

func (z *Z) SetLevel(level stdlog.Level) {
	z.cfg.Level = zap.NewAtomicLevelAt(lmap[level])
	l, err := z.cfg.Build(zap.AddCallerSkip(1 + int(z.stackSkip)))
	if err != nil {
		z.FatalError(err, "Failed to rebuild config with new level", "level", level.String())
	}

	if z.Logger.Name() != "" {
		z.Logger = l.Named(z.Logger.Name())
		return
	}
	z.Logger = l
}

func (z *Z) Leveled(level stdlog.Level) stdlog.Logger {
	cp := new(Z)
	*cp = *z
	cp.SetLevel(level)
	return cp
}

func (z *Z) WithFields(keysAndValues ...any) stdlog.Logger {
	n := new(Z)
	*n = *z
	n.Logger = n.Logger.With(handleFields("WithFields", keysAndValues)...)
	return n
}

func (z *Z) Debug(msg string, fields ...any) {
	z.Logger.Debug(msg, handleFields(stdlog.LevelDebug.String(), fields)...)
}

func (z *Z) Info(msg string, fields ...any) {
	z.Logger.Info(msg, handleFields(stdlog.LevelInfo.String(), fields)...)
}

func (z *Z) Warning(msg string, fields ...any) {
	z.Logger.Warn(msg, handleFields(stdlog.LevelWarning.String(), fields)...)
}

func (z *Z) Error(err error, msg string, fields ...any) {
	z.Logger.Error(msg, handleFields(stdlog.LevelError.String(), fields, zap.Error(err))...)
}

func (z *Z) Fatal(msg string, fields ...any) {
	z.Logger.Fatal(msg, handleFields(stdlog.LevelFatal.String(), fields)...)
}

func (z *Z) FatalError(err error, msg string, fields ...any) {
	z.Logger.Fatal(msg, handleFields(stdlog.LevelFatal.String(), fields, zap.Error(err))...)
}

func handleFields(method string, kvs []any, extra ...zap.Field) []zap.Field {
	if len(kvs)%2 != 0 {
		panic(fmt.Errorf("uneven keys and values passed to %s", method))
	}

	fields := make([]zap.Field, 0, len(kvs)/2)

	for i := 0; i < len(kvs)-1; i += 2 {
		k := fmt.Sprintf("%v", kvs[i])
		fields = append(fields, zap.Any(k, kvs[i+1]))
	}

	return append(fields, extra...)
}
