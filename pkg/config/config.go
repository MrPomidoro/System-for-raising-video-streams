package config

import (
	"time"
)

// структура конфига
type Config struct {
	Logger             `yaml:"logger"`
	Server             `yaml:"server"`
	Database           `yaml:"database"`
	Rtsp_simple_server `yaml:"rtsp_simple_server"`
}

// параметры логгера
type Logger struct {
	LogLevel string `yaml:"loglevel"`
}

// параметры сервера
type Server struct {
	Addr         string        `yaml:"addr"`
	ReadTimeout  time.Duration `yaml:"readtimeout"`
	WriteTimeout time.Duration `yaml:"writetimeout"`
	IdleTimeout  time.Duration `yaml:"idletimeout"`
}

// параметры базы данных
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

type Rtsp_simple_server struct {
	Run          string        `yaml:"run"`
	Url          string        `yaml:"url"`
	Refresh_Time time.Duration `yaml:"refresh_time"`
}
