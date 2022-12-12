package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

// NewLog - функция для инициализации логгера с настраиваемым уровнем логирования
// на основе данных из конфига
func NewLog(level, fileName string) *logrus.Logger {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Errorf("cannot open file for logging: %v", err)
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

// initLogLevel реализует выбор уровня логирования на основе переданной строковой переменной
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
