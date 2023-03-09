package config

import (
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
)

// const (
// 	logS =
// )

func TestGetConfig(t *testing.T) {
	tests := []struct {
		name      string
		yaml      string
		expectCfg *Config
		expectErr error
	}{
		{
			name: "GetConfigOK",
			yaml: `
logger:	
  logLevel: DEBUG
  logFileEnable: true
  logStdoutEnable: true
  maxSize: 500
  maxAge: 28
  maxBackups: 7
  rewriteLog: true
server:
  readTimeout: 200ms
  writeTimeout: 200ms
  idleTimeout: 10s
database:
  port: 5432
  host: 192.168.0.32
  dbName: www
  user: sysadmin
  password: w3X{77PpCR
  driver: postgres
  connect: true
  connectionTimeout: 10s
rtsp:
  run: #usr/bin
  refreshTime: 10s
  url: http://10.100.100.228:9997
  api:  
    urlGet: /v1/paths/list
    urlAdd: /v1/config/paths/add/
    urlRemove: /v1/config/paths/remove/
    urlEdit: /v1/config/paths/edit/`,
			expectCfg: &Config{
				Logger:   Logger{LogLevel: "DEBUG", LogFileEnable: true, LogStdoutEnable: true, MaxSize: 500, MaxAge: 28, MaxBackups: 7, RewriteLog: true},
				Server:   Server{ReadTimeout: time.Duration(200 * time.Millisecond), WriteTimeout: time.Duration(200 * time.Millisecond), IdleTimeout: time.Duration(10 * time.Second)},
				Database: Database{Port: 5432, Host: "192.168.0.32", User: "sysadmin", Password: "w3X{77PpCR", DbName: "www", Driver: "postgres", Connect: true, ConnectionTimeout: time.Duration(10 * time.Second)},
				Rtsp: Rtsp{RefreshTime: time.Duration(10 * time.Second), Url: "http://10.100.100.228:9997",
					Api: api{UrlGet: "/v1/paths/list", UrlAdd: "/v1/config/paths/add/", UrlRemove: "/v1/config/paths/remove/", UrlEdit: "/v1/config/paths/edit/"}},
				err: customError.ErrorConfig,
			},
			expectErr: nil,
		},
		{
			name: "GetConfigErrorFileExtraSpace",
			yaml: `
logger:	
    logLevel: DEBUG
  logFileEnable: true
  logStdoutEnable: true
  maxSize: 500
  maxAge: 28
  maxBackups: 7
  rewriteLog: true`,
			expectCfg: &Config{
				err: customError.ErrorConfig,
			},
			expectErr: errors.New("While parsing config: yaml: line 1: did not find expected key"),
		},
		{
			name: "GetConfigErrorFileHaveTabs",
			yaml: `
database:
  port: 5432
		  host: 192.168.0.32
		  unknown: 1`,
			expectCfg: &Config{
				err: customError.ErrorConfig,
			},
			expectErr: errors.New("While parsing config: yaml: line 3: found a tab character that violates indentation"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := "/home/ksenia/go/src/github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config/config.yaml"
			file, err := os.Create(path)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			defer os.Remove(path)

			file.WriteString(tt.yaml)
			file.Close()

			cfg, errGetCfg := GetConfig()

			var expErr *customError.Error
			if tt.expectErr != nil {
				expErr = customError.NewError(2, "50.0.1", "error at the level of reading and processing the config").SetError(tt.expectErr)
			}

			if errGetCfg != nil && tt.expectErr != nil {

				errGetCfgS := strings.Split(errGetCfg.Error(), ":")
				errExpS := strings.Split(expErr.Error(), ":")

				if errGetCfgS[len(errGetCfgS)-1] != errExpS[len(errExpS)-1] {
					t.Errorf("unexpected error %v, expect %v", errGetCfg, expErr)
				}
			} else if (errGetCfg == nil && tt.expectErr != nil) || (errGetCfg != nil && tt.expectErr == nil) {
				t.Errorf("unexpected error %v, expect %v", errGetCfg, expErr)
			}

			if !reflect.DeepEqual(cfg, tt.expectCfg) {
				t.Errorf("\nexpect\n%v, \ngot \n%v", tt.expectCfg, cfg)
			}
		})
	}
}
