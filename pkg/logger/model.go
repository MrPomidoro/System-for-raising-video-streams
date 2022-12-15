package logger

type logger struct {
	LogLevel      string
	LogFileEnable bool
	LogFile       string
	MaxSize       int
	MaxAge        int

	MaxBackups int
}
