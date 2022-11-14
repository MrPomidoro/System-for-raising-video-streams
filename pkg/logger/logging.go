package logger

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// Выводит сообщение msg на уровне "Fatal"
// и завершает работу программы
func LogFatal(log *logrus.Logger, msg interface{}) {
	log.Fatalf(fmt.Sprintf("%v", msg))
}

// Выводит сообщение msg на уровне "Error"
func LogError(log *logrus.Logger, msg interface{}) {
	log.Errorf(fmt.Sprintf("%v", msg))
}

// Выводит сообщение msg на уровне "Warn"
func LogWarn(log *logrus.Logger, msg interface{}) {
	log.Warnf(fmt.Sprintf("%v", msg))
}

// Выводит сообщение msg на уровне "Info"
func LogInfo(log *logrus.Logger, msg interface{}) {
	log.Infof(fmt.Sprintf("%v", msg))
}

// Выводит сообщение msg на уровне "Debug"
func LogDebug(log *logrus.Logger, msg interface{}) {
	log.Debugf(fmt.Sprintf("%v", msg))
}

// Выводит сообщение msg на уровне "Trace"
func LogTrace(log *logrus.Logger, msg interface{}) {
	log.Tracef(fmt.Sprintf("%v", msg))
}

//

// Выводит сообщение msg на уровне "Info" с указанием статус кода и вызванного метода
func LogInfoStatusCode(log *logrus.Logger, msg interface{}, method, status string) {
	log.WithField("meth", method).WithField("status", status).Printf(fmt.Sprintf("%v", msg))
}

// Выводит сообщение msg на уровне "Error" с указанием статус кода и вызванного метода
func LogErrorStatusCode(log *logrus.Logger, msg interface{}, method, status string) {
	log.WithField("meth", method).WithField("status", status).Errorf(fmt.Sprintf("%v", msg))
}

// func LogSMTH(log *logrus.Logger, msg interface{}, method, status string)
