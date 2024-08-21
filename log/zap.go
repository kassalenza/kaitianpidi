package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapLogger(fileName string, maxSize, maxBackups, maxAge int, compress bool) *zap.Logger {
	lumberJackLogger := NewLumberJack(fileName, maxSize, maxBackups, maxAge, true)

	writeSyncer := zapcore.AddSync(lumberJackLogger)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel)

	return zap.New(core, zap.AddCaller())
}
