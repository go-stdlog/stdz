package stdz

import (
	"github.com/go-stdlog/stdlog"
	"go.uber.org/zap"
)

type zLogger struct {
	*zap.Logger
}

func New(logger *zap.Logger) stdlog.Logger[zap.Field] {
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

func (z *zLogger) Debug(msg string, fields ...zap.Field) {
	z.Logger.Debug(msg, fields...)
}

func (z *zLogger) Info(msg string, fields ...zap.Field) {
	z.Logger.Info(msg, fields...)
}

func (z *zLogger) Warning(msg string, fields ...zap.Field) {
	z.Logger.Warn(msg, fields...)
}

func (z *zLogger) Error(err error, msg string, fields ...zap.Field) {
	z.Logger.Error(msg, append(fields, zap.Error(err))...)
}

func (z *zLogger) Fatal(msg string, fields ...zap.Field) {
	z.Logger.Fatal(msg, fields...)
}

func (z *zLogger) FatalError(err error, msg string, fields ...zap.Field) {
	z.Logger.Fatal(msg, append(fields, zap.Error(err))...)
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
