package logger

import (
	"io"

	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	logLevel(cfg *config.Config) zapcore.Level
	loggers(cfg *config.Config, consoleLogger, fileLogger zapcore.Core) []zapcore.Core
	customEncoderLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder)
	zapOpts() []zap.Option

	zapcoreDTO
	zapDTO
}

type zapcoreDTO interface {
	newJSONEncode(cfg zapcore.EncoderConfig) zapcore.Encoder
	newConsoleEncoder(cfg zapcore.EncoderConfig) zapcore.Encoder
	newCore(enc zapcore.Encoder, ws zapcore.WriteSyncer, enab zapcore.LevelEnabler) zapcore.Core
	newTee(cores ...zapcore.Core) zapcore.Core
	addSync(w io.Writer) zapcore.WriteSyncer
}

type zapDTO interface {
	newProductionEncoderConfig() zapcore.EncoderConfig
	new(core zapcore.Core, options ...zap.Option) *zap.Logger
}
