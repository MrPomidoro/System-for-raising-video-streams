package logger

type logger struct {
	LogLevel        string
	LogFileEnable   bool
	LogStdoutEnable bool
	LogFile         string
	MaxSize         int
	MaxAge          int
	MaxBackups      int
}
