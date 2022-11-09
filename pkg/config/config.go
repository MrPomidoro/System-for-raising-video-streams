package config

import (
	"time"
)

// структура конфига
type Config struct {
	Logger      `yaml:"logger"`
	Database    `yaml:"database"`
	PathDir     `yaml:"pathDir"`
	MqttConnect `yaml:"mqttConnect"`
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

// параметры mqtt коннекта
type MqttConnect struct {
	MqttLogin      string `yaml:"mqttLogin"`
	MqttPassword   string `yaml:"mqttPassword"`
	MqttHost       string `yaml:"mqttHost"`
	MqttDomainName string `yaml:"mqttDomainName"`
}
