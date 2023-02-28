package logger

import (
	"io"

	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (l *logger) logLevel(cfg *config.Config) zapcore.Level {
	// Если переданный уровень возможно распарсить - это выполняется,
	// иначе - передаётся уровень "Info" по дефолту
	level, err := zapcore.ParseLevel(cfg.LogLevel)
	if err != nil {
		return zapcore.InfoLevel
	}
	return level
}

func (l *logger) loggers(cfg *config.Config, consoleLogger, fileLogger zapcore.Core) []zapcore.Core {
	cores := make([]zapcore.Core, 0)
	// Логирование в консоль - всегда
	cores = append(cores, consoleLogger)

	// Если нужно логировать в файл, файл добавляется к слайсу
	if cfg.LogFileEnable {
		cores = append(cores, fileLogger)
	}
	return cores
}

func (l *logger) customEncoderLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + logLevelSeverity[level] + "]")
}

func (l *logger) zapOpts() []zap.Option {
	return []zap.Option{zap.AddCaller()}
}

// zapcore transcode
func (l *logger) newJSONEncode(cfg zapcore.EncoderConfig) zapcore.Encoder {
	return zapcore.NewJSONEncoder(cfg)
}

func (l *logger) newConsoleEncoder(cfg zapcore.EncoderConfig) zapcore.Encoder {
	return zapcore.NewConsoleEncoder(cfg)
}

func (l *logger) newCore(enc zapcore.Encoder, ws zapcore.WriteSyncer, enab zapcore.LevelEnabler) zapcore.Core {
	return zapcore.NewCore(enc, ws, enab)
}

func (l *logger) newTee(cores ...zapcore.Core) zapcore.Core {
	return zapcore.NewTee(cores...)
}

func (l *logger) addSync(w io.Writer) zapcore.WriteSyncer {
	return zapcore.AddSync(w)
}

// zap transcode
func (l *logger) newProductionEncoderConfig() zapcore.EncoderConfig {
	return zap.NewProductionEncoderConfig()
}

func (l *logger) new(core zapcore.Core, options ...zap.Option) *zap.Logger {
	return zap.New(core, options...)
}
