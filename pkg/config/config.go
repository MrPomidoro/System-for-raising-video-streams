package config

import (
	"time"
)

// Структура конфига
type Config struct {
	Logger             `yaml:"logger"`
	Server             `yaml:"server"`
	Database           `yaml:"database"`
	Rtsp_simple_server `yaml:"rtsp_simple_server"`
}

// Параметры логгера
type Logger struct {
	LogLevel string `yaml:"loglevel"`
	LogFile  string `yaml:"logfile"`
}

// Параметры сервера
type Server struct {
	Server_Port  string        `yaml:"server_port"`
	Server_Host  string        `yaml:"server_host"`
	ReadTimeout  time.Duration `yaml:"readtimeout"`
	WriteTimeout time.Duration `yaml:"writetimeout"`
	IdleTimeout  time.Duration `yaml:"idletimeout"`
}

// Параметры базы данных
type Database struct {
	Port                         string        `yaml:"port"`
	Host                         string        `yaml:"host"`
	Db_Name                      string        `yaml:"db_name"`
	User                         string        `yaml:"user"`
	Password                     string        `yaml:"password"`
	Driver                       string        `yaml:"driver"`
	Database_Connect             bool          `yaml:"database_connect"`
	Db_Connection_Timeout_Second time.Duration `yaml:"db_connection_timeout_second"`
}

// Параметры rtsp_simple_server
type Rtsp_simple_server struct {
	Run string `yaml:"run"`
	// Url          string        `yaml:"url"`
	Refresh_Time time.Duration `yaml:"refresh_time"`
}
