package config

import (
	"time"

	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
)

// Config - структура конфига
type Config struct {
	Logger   `mapstructure:"logger"`
	Server   `mapstructure:"server"`
	Database `mapstructure:"database"`
	Rtsp     `mapstructure:"rtsp"`

	err ce.IError
}

// Logger содержит параметры логгера
type Logger struct {
	LogLevel        string `mapstructure:"logLevel"`
	LogFileEnable   bool   `mapstructure:"logFileEnable"`
	LogStdoutEnable bool   `mapstructure:"logStdoutEnable"`
	LogFile         string `mapstructure:"logpath"`
	MaxSize         int    `mapstructure:"maxSize"`
	MaxAge          int    `mapstructure:"maxAge"`
	MaxBackups      int    `mapstructure:"maxBackups"`
	RewriteLog      bool   `mapstructure:"rewriteLog"`
}

// Server содержит параметры сервера
type Server struct {
	ReadTimeout  time.Duration `mapstructure:"readTimeout"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout"`
	IdleTimeout  time.Duration `mapstructure:"idleTimeout"`
}

// Database содержит параметры базы данных
type Database struct {
	Port              int           `mapstructure:"port"`
	Host              string        `mapstructure:"host"`
	DbName            string        `mapstructure:"dbName"`
	User              string        `mapstructure:"user"`
	Password          string        `mapstructure:"password"`
	Driver            string        `mapstructure:"driver"`
	Connect           bool          `mapstructure:"connect"`
	ConnectionTimeout time.Duration `mapstructure:"connectionTimeout"`
}

// Rtsp RtspSimpleServer содержит параметры rtsp_simple_server
type Rtsp struct {
	Run         string        `mapstructure:"run"`
	Url         string        `mapstructure:"url"`
	RefreshTime time.Duration `mapstructure:"refreshTime"`
	Api         api           `mapstructure:"api"`
}

type api struct {
	UrlGet    string `mapstructure:"urlGet"`
	UrlAdd    string `mapstructure:"urlAdd"`
	UrlRemove string `mapstructure:"urlRemove"`
	UrlEdit   string `mapstructure:"urlEdit"`
}
