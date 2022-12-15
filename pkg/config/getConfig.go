package config

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func NewConfig() *Config {
	return &Config{}
}

// GetConfig инициализирует и заполняет структуру конфигурационного файла
func (cfg *Config) GetConfig() (*Config, error) {

	// Чтение пути до конфигурационного файла
	configPath := readConfigPath()
	// Если путь не был указан, выставляется по умолчанию ./
	if configPath == "" {
		configPath = "./"
	}

	var v = viper.New()
	// var cfg Config

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(configPath)

	err := readParametersFromConfig(v, cfg)
	if err != nil {
		return cfg, err
	}

	// Проверка наличия параметров в командной строке
	readFlags(cfg)

	return cfg, nil
}

func readParametersFromConfig(v *viper.Viper, cfg *Config) error {
	// Попытка чтения конфига
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("cannot read Config: %v", err)
	}
	// Попытка заполнение структуры Config полученными данными
	if err := v.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("cannot read Config: %v", err)
	}
	return nil
}

// checkConfigPath проверяет, есть ли среди флагов путь до конфигурационного файла
func readConfigPath() string {
	var configPath string
	args := os.Args
	for _, arg := range args {
		if strings.Split(arg, "=")[0][1:] == "configPath" {
			configPath = strings.Split(arg, "=")[1]
		}
	}
	return configPath
}

// readFlags реализует возможность передачи параметров
// конфигурационного файла при запуске из командной строки
func readFlags(cfg *Config) {
	var stub string

	flag.StringVar(&cfg.LogLevel, "logLevel", cfg.LogLevel, "The level of logging parameter")
	flag.BoolVar(&cfg.LogFileEnable, "logFileEnable", cfg.LogFileEnable, "The statement whether to log to a file")
	flag.StringVar(&cfg.LogFile, "logpath", cfg.LogFile, "The path to file of logging out")
	flag.IntVar(&cfg.MaxSize, "maxSize", cfg.MaxSize, "The path to file of logging out")
	flag.IntVar(&cfg.MaxAge, "maxAge", cfg.MaxAge, "The path to file of logging out")
	flag.IntVar(&cfg.MaxBackups, "maxBackups", cfg.MaxBackups, "The path to file of logging out")

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

	flag.StringVar(&stub, "configPath", `./`, "The path to file of configuration")

	flag.Parse()
}
