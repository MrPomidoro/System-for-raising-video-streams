package config

import (
	"time"
)

// структура конфига
type Config struct {
	Logger             `yaml:"logger"`
	Database           `yaml:"database"`
	PathDir            `yaml:"pathDir"`
	Rtsp_simple_server `yaml:"rtsp_simple_server"`
}

type PathDir struct {
	ConfigPath string `yaml:"configPath"`
}

// параметры логгера
type Logger struct {
	LogLevel string `yaml:"loglevel"`
}

// параметры базы данных
type Database struct {
	Port                      string        `yaml:"port"`
	Host                      string        `yaml:"host"`
	DbName                    string        `yaml:"dbName"`
	User                      string        `yaml:"user"`
	Password                  string        `yaml:"password"`
	Driver                    string        `yaml:"driver"`
	DBConnectionTimeoutSecond time.Duration `yaml:"dbConnectionTimeoutSecond"`
}

type Rtsp_simple_server struct {
	Run          string        `yaml:"run"`
	Url          string        `yaml:"url"`
	Refresh_Time time.Duration `yaml:"refresh_time"`
}
