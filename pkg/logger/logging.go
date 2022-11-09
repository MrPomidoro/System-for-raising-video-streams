package logger

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// выводит сообщение msg на уровне "Fatal"
func LogFatal(log *logrus.Logger, msg interface{}) {
	log.Fatalf(fmt.Sprintf("%v", msg))
}

// выводит сообщение msg на уровне "Error"
func LogError(log *logrus.Logger, msg interface{}) {
	log.Errorf(fmt.Sprintf("%v", msg))
}

// выводит сообщение msg на уровне "Warn"
func LogWarn(log *logrus.Logger, msg interface{}) {
	log.Warnf(fmt.Sprintf("%v", msg))
}

// выводит сообщение msg на уровне "Info"
func LogInfo(log *logrus.Logger, msg interface{}) {
	log.Infof(fmt.Sprintf("%v", msg))
}

// выводит сообщение msg на уровне "Debug"
func LogDebug(log *logrus.Logger, msg interface{}) {
	log.Debugf(fmt.Sprintf("%v", msg))
}

// выводит сообщение msg на уровне "Trace"
func LogTrace(log *logrus.Logger, msg interface{}) {
	log.Tracef(fmt.Sprintf("%v", msg))
}
