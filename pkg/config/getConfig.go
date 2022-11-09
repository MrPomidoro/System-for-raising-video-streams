package config

import (
	"flag"
	"fmt"

	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/logger"
	"github.com/spf13/viper"
)

// инициализация и заполнение конфига
func GetConfig() *Config {
	var v = viper.New()
	var cfg Config

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./")

	log := logger.NewLog(cfg.LogLevel)

	if err := v.ReadInConfig(); err != nil {
		logger.LogError(log, fmt.Sprintf("cannot read config: %v", err))
	} else {
		logger.LogDebug(log, "Success read config file")
	}

	if err := v.Unmarshal(&cfg); err != nil {
		logger.LogError(log, fmt.Sprintf("cannot read config: %v", err))
	} else {
		logger.LogDebug(log, "Success read config file")
	}

	readFlags(&cfg)

	return &cfg
}

func readFlags(cfg *Config) {
	flag.StringVar(&cfg.LogLevel, "loglevel", cfg.LogLevel, "The loglevel parameter")
	flag.StringVar(&cfg.Port, "port", cfg.Port, "The port parameter")
	flag.StringVar(&cfg.Host, "host", cfg.Host, "The host parameter")
	flag.StringVar(&cfg.DbName, "dbName", cfg.DbName, "The dbName parameter")
	flag.StringVar(&cfg.User, "user", cfg.User, "The user parameter")
	flag.StringVar(&cfg.Password, "password", cfg.Password, "The password parameter")
	flag.StringVar(&cfg.Driver, "driver", cfg.Driver, "The driver parameter")
	flag.DurationVar(&cfg.DBConnectionTimeoutSecond, "dbConnectionTimeoutSecond", cfg.DBConnectionTimeoutSecond, "The dbConnectionTimeoutSecond parameter")
	flag.StringVar(&cfg.ConfigPath, "configPath", cfg.ConfigPath, "The configPath parameter")
	flag.Parse()
}
