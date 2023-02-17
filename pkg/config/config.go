package config

import (
	"time"

	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
)

// Config - структура конфига
type Config struct {
	Logger             `yaml:"logger"`
	Server             `yaml:"server"`
	Database           `yaml:"database"`
	Rtsp_simple_server `yaml:"rtsp_simple_server"`
	err                ce.Error
}

// Logger содержит параметры логгера
type Logger struct {
	LogLevel        string `yaml:"logLevel"`
	LogFileEnable   bool   `yaml:"logFileEnable"`
	LogStdoutEnable bool   `yaml:"logStdoutEnable"`
	LogFile         string `yaml:"logpath"`
	MaxSize         int    `yaml:"maxSize"`
	MaxAge          int    `yaml:"maxAge"`
	MaxBackups      int    `yaml:"maxBackups"`
	RewriteLog      bool   `yaml:"rewriteLog"`
}

// Server содержит параметры сервера
type Server struct {
	ReadTimeout  time.Duration `yaml:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
	IdleTimeout  time.Duration `yaml:"idleTimeout"`
}

// Database содержит параметры базы данных
type Database struct {
	Port                      string        `yaml:"port"`
	Host                      string        `yaml:"host"`
	DbName                    string        `yaml:"dbName"`
	User                      string        `yaml:"user"`
	Password                  string        `yaml:"password"`
	Driver                    string        `yaml:"driver"`
	DatabaseConnect           bool          `yaml:"databaseConnect"`
	DbConnectionTimeoutSecond time.Duration `yaml:"dbConnectionTimeoutSecond"`
}

// Rtsp_simple_server содержит параметры rtsp_simple_server
type Rtsp_simple_server struct {
	Run         string        `yaml:"run"`
	Url         string        `yaml:"url"`
	RefreshTime time.Duration `yaml:"refreshTime"`
}
