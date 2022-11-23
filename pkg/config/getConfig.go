package config

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/logger"
	"github.com/spf13/viper"
)

// Инициализация и заполнение конфига
func GetConfig() *Config {

	var configPath string
	args := os.Args
	for _, arg := range args {
		if strings.Split(arg, "=")[0][1:] == "configPath" {
			configPath = strings.Split(arg, "=")[1]
		}
	}
	if configPath == "" {
		configPath = "./"
	}

	var v = viper.New()
	var cfg Config

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(configPath)

	log := logger.NewLog(cfg.LogLevel)

	// Попытка чтения конфига
	if err := v.ReadInConfig(); err != nil {
		logger.LogError(log, fmt.Sprintf("cannot read config: %v", err))
	} else {
		logger.LogDebug(log, "Success read config file")
	}

	// Попытка заполнение структуры Config полученными данными
	if err := v.Unmarshal(&cfg); err != nil {
		logger.LogError(log, fmt.Sprintf("cannot read config: %v", err))
	} else {
		logger.LogDebug(log, "Success read config file")
	}

	// Проверка наличия параметров в командной строке
	readFlags(&cfg)

	fmt.Println(cfg)

	return &cfg
}

// Реализация возможности передачи параметров конфигурационного файла
// при запуске из командной строки
func readFlags(cfg *Config) {
	var A string
	flag.StringVar(&cfg.LogLevel, "loglevel", cfg.LogLevel, "The level of logging parameter")

	flag.DurationVar(&cfg.ReadTimeout, "readtimeout", cfg.ReadTimeout, "The readtimeout parameter")
	flag.DurationVar(&cfg.WriteTimeout, "writetimeout", cfg.WriteTimeout, "The writetimeout parameter")
	flag.DurationVar(&cfg.IdleTimeout, "idletimeout", cfg.IdleTimeout, "The idletimeout parameter")

	flag.StringVar(&cfg.Port, "port", cfg.Port, "The port parameter")
	flag.StringVar(&cfg.Host, "host", cfg.Host, "The host parameter")
	flag.StringVar(&cfg.Db_Name, "db_name", cfg.Db_Name, "The db_name parameter")
	flag.StringVar(&cfg.User, "user", cfg.User, "The user parameter")
	flag.StringVar(&cfg.Password, "password", cfg.Password, "The password parameter")
	flag.StringVar(&cfg.Driver, "driver", cfg.Driver, "The driver parameter")
	flag.BoolVar(&cfg.Database_Connect, "database_connect", cfg.Database_Connect, "The permission to connect")
	flag.DurationVar(&cfg.Db_Connection_Timeout_Second, "db_connection_timeout_second", cfg.Db_Connection_Timeout_Second, "The db_connection_timeout_second parameter")

	flag.StringVar(&cfg.Run, "Run", cfg.Run, "The run parameter")
	flag.StringVar(&cfg.Url, "url", cfg.Url, "The url parameter")
	flag.DurationVar(&cfg.Refresh_Time, "refresh_time", cfg.Refresh_Time, "The refresh_time parameter")

	flag.StringVar(&A, "configPath", `./`, "configPath")

	flag.Parse()
}
