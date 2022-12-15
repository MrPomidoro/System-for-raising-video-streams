package logger

import (
	"os"

	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logLevelSeverity = map[zapcore.Level]string{
	zapcore.DebugLevel:  "DEBUG",
	zapcore.InfoLevel:   "INFO",
	zapcore.WarnLevel:   "WARNING",
	zapcore.ErrorLevel:  "ERROR",
	zapcore.DPanicLevel: "CRITICAL",
	zapcore.PanicLevel:  "PANIC",
	zapcore.FatalLevel:  "FATAL",
}

func NewLogger(cfg *config.Config) *zap.Logger {
	l := logger{
		LogLevel:      cfg.LogLevel,
		LogFileEnable: cfg.LogFileEnable,
		LogFile:       cfg.LogFile,
		MaxSize:       cfg.MaxSize,
		MaxAge:        cfg.MaxAge,
		MaxBackups:    cfg.MaxBackups,
	}
	return l.initLogger(cfg)
}

func (l *logger) initLogger(cfg *config.Config) *zap.Logger {

	li := Logger(l)
	conf := li.newProductionEncoderConfig()
	conf.EncodeTime = zapcore.ISO8601TimeEncoder
	conf.EncodeLevel = li.customEncoderLevel
	conf.MessageKey = "message"
	conf.CallerKey = "caller"
	conf.TimeKey = "time"

	jsonEncoder := li.newJSONEncode(conf)
	textEncoder := li.newConsoleEncoder(conf)

	fileLogger := li.newCore(
		jsonEncoder,
		li.addSync(&lumberjack.Logger{
			Filename:   l.LogFile,
			MaxSize:    l.MaxSize,
			MaxAge:     l.MaxAge,
			MaxBackups: l.MaxBackups,
		}),
		li.logLevel(cfg),
	)

	consoleLogger := li.newCore(
		textEncoder,
		zapcore.AddSync(os.Stdout),
		li.logLevel(cfg),
	)

	return li.new(li.newTee(li.loggers(cfg, consoleLogger, fileLogger)...), li.zapOpts()...)
}
