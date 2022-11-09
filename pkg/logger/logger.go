package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

func NewLog(level string) *logrus.Logger {
	file, err := os.OpenFile(FileNameConst, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf(OpenFileErrConst, err)
	}

	log := &logrus.Logger{
		Out:   io.MultiWriter(file, os.Stdout),
		Level: initLogLevel(level),
		Formatter: &easy.Formatter{
			TimestampFormat: ServTimestampFormatConst,
			LogFormat:       ServLogFormatConst,
		},
	}

	return log
}

func initLogLevel(level string) logrus.Level {
	switch level {
	case "FATAL":
		return logrus.FatalLevel
	case "ERROR":
		return logrus.ErrorLevel
	case "WARN":
		return logrus.WarnLevel
	case "INFO":
		return logrus.InfoLevel
	case "DEBUG":
		return logrus.DebugLevel
	case "TRACE":
		return logrus.TraceLevel
	default:
		return logrus.InfoLevel
	}
}
