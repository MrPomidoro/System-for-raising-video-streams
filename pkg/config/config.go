package config

import (
	"time"
)

// Config - структура конфига
type Config struct {
	Logger             `yaml:"logger"`
	Server             `yaml:"server"`
	Database           `yaml:"database"`
	Rtsp_simple_server `yaml:"rtsp_simple_server"`
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
	ReadTimeout  time.Duration `yaml:"readtimeout"`
	WriteTimeout time.Duration `yaml:"writetimeout"`
	IdleTimeout  time.Duration `yaml:"idletimeout"`
}

// Database содержит параметры базы данных
type Database struct {
	Port                      string        `yaml:"port"`
	Host                      string        `yaml:"host"`
	DbName                    string        `yaml:"db_name"`
	User                      string        `yaml:"user"`
	Password                  string        `yaml:"password"`
	Driver                    string        `yaml:"driver"`
	DatabaseConnect           bool          `yaml:"database_connect"`
	DbConnectionTimeoutSecond time.Duration `yaml:"db_connection_timeout_second"`
}

// Rtsp_simple_server содержит параметры rtsp_simple_server
type Rtsp_simple_server struct {
	Run         string        `yaml:"run"`
	Url         string        `yaml:"url"`
	RefreshTime time.Duration `yaml:"refresh_time"`
}
