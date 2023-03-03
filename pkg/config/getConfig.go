package config

import (
	"flag"
	"os"
	"strings"

	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
	"github.com/spf13/viper"
)

// GetConfig инициализирует и заполняет структуру конфигурационного файла
func GetConfig() (*Config, ce.IError) {
	var cfg Config
	cfg.err = ce.ErrorConfig

	// Чтение пути до конфигурационного файла
	configPath := readConfigPath()
	// Если путь не был указан, выставляется по умолчанию ./
	if configPath == "" {
		configPath = "./"
	}

	var v = viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(configPath)

	err := readParametersFromConfig(v, &cfg)
	if err != nil {
		return &cfg, cfg.err.SetError(err)
	}

	// Проверка наличия параметров в командной строке
	readFlags(&cfg)
	return &cfg, nil
}

func readParametersFromConfig(v *viper.Viper, cfg *Config) *ce.Error {
	// Попытка чтения конфига
	if err := v.ReadInConfig(); err != nil {
		return cfg.err.SetError(err)
	}
	// Попытка заполнение структуры Config полученными данными
	if err := v.Unmarshal(&cfg); err != nil {
		return cfg.err.SetError(err)
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
	flag.BoolVar(&cfg.LogStdoutEnable, "logStdoutEnable", cfg.LogStdoutEnable, "The statement whether to log to console")
	flag.StringVar(&cfg.LogFile, "logpath", cfg.LogFile, "The path to file of logging out")
	flag.IntVar(&cfg.MaxSize, "maxSize", cfg.MaxSize, "The path to file of logging out")
	flag.IntVar(&cfg.MaxAge, "maxAge", cfg.MaxAge, "The path to file of logging out")
	flag.IntVar(&cfg.MaxBackups, "maxBackups", cfg.MaxBackups, "The path to file of logging out")
	flag.BoolVar(&cfg.RewriteLog, "rewriteLog", cfg.RewriteLog, "Is rewrite log file")

	flag.DurationVar(&cfg.ReadTimeout, "readTimeout", cfg.ReadTimeout, "The readtimeout parameter")
	flag.DurationVar(&cfg.WriteTimeout, "wriTetimeout", cfg.WriteTimeout, "The writetimeout parameter")
	flag.DurationVar(&cfg.IdleTimeout, "idleTimeout", cfg.IdleTimeout, "The idletimeout parameter")

	flag.IntVar(&cfg.Port, "port", cfg.Port, "The host parameter")
	flag.StringVar(&cfg.Host, "host", cfg.Host, "The host parameter")
	flag.StringVar(&cfg.DbName, "dbName", cfg.DbName, "The db_name parameter")
	flag.StringVar(&cfg.User, "user", cfg.User, "The user parameter")
	flag.StringVar(&cfg.Password, "password", cfg.Password, "The password parameter")
	flag.StringVar(&cfg.Driver, "driver", cfg.Driver, "The driver parameter")
	flag.BoolVar(&cfg.Connect, "connect", cfg.Connect, "The permission to connect")
	flag.DurationVar(&cfg.ConnectionTimeout, "connectionTimeout", cfg.ConnectionTimeout, "The db_connection_timeout_second parameter")

	flag.StringVar(&cfg.Run, "run", cfg.Run, "The run parameter")
	flag.StringVar(&cfg.Url, "url", cfg.Url, "The url parameter")
	flag.DurationVar(&cfg.RefreshTime, "refreshTime", cfg.RefreshTime, "The refresh time parameter")

	flag.StringVar(&cfg.Api.UrlGet, "urlGet", cfg.Api.UrlGet, "The url for get from rtsp")
	flag.StringVar(&cfg.Api.UrlAdd, "urlAdd", cfg.Api.UrlAdd, "The url for add from rtsp")
	flag.StringVar(&cfg.Api.UrlRemove, "urlRemove", cfg.Api.UrlRemove, "The url for remove into rtsp")
	flag.StringVar(&cfg.Api.UrlEdit, "urlEdit", cfg.Api.UrlEdit, "The url for edit into rtsp")

	flag.StringVar(&stub, "configPath", `./`, "The path to file of configuration")

	flag.Parse()
}
