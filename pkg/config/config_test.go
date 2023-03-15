package config

import (
	"errors"
	"flag"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
)

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
			path := "./config.yaml"
			file, err := os.Create(path)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			defer os.Remove(path)

			if tt.yaml != "" {
				file.WriteString(tt.yaml)
			}
			defer file.Close()

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

func TestReadFlags(t *testing.T) {

	cfg := &Config{}
	readFlags(cfg)

	tests := []struct {
		name string
		args []string
		want Config
	}{
		{
			name: "empty args",
			args: []string{},
			want: Config{},
		},
		{
			name: "all flags",
			args: []string{
				"--logLevel=debug",
				"--logFileEnable=true",
				"--logStdoutEnable=false",
				"--logpath=/var/log/myapp.log",
				"--maxSize=1024",
				"--maxAge=7",
				"--maxBackups=3",
				"--rewriteLog=true",
				"--readTimeout=5s",
				"--wriTetimeout=10s",
				"--idleTimeout=60s",
				"--port=8080",
				"--host=localhost",
				"--dbName=mydb",
				"--tableName=mytable",
				"--user=myuser",
				"--password=mypassword",
				"--driver=mysql",
				"--connect=true",
				"--connectionTimeout=30s",
				"--run=server",
				"--url=http://localhost:8080",
				"--refreshTime=1m",
				"--urlGet=http://localhost:8000/get",
				"--urlAdd=http://localhost:8000/add",
				"--urlRemove=http://localhost:8000/remove",
				"--urlEdit=http://localhost:8000/edit",
				"--configPath=./",
			},
			want: Config{
				Logger: Logger{
					LogLevel:        "debug",
					LogFileEnable:   true,
					LogStdoutEnable: false,
					LogFile:         "/var/log/myapp.log",
					MaxSize:         1024,
					MaxAge:          7,
					MaxBackups:      3,
					RewriteLog:      true,
				},
				Server: Server{
					ReadTimeout:  5 * time.Second,
					WriteTimeout: 10 * time.Second,
					IdleTimeout:  60 * time.Second,
				},
				Database: Database{
					Port:              8080,
					Host:              "localhost",
					DbName:            "mydb",
					TableName:         "mytable",
					User:              "myuser",
					Password:          "mypassword",
					Driver:            "mysql",
					Connect:           true,
					ConnectionTimeout: 30 * time.Second,
				},
				Rtsp: Rtsp{
					Run:         "server",
					Url:         "http://localhost:8080",
					RefreshTime: 1 * time.Minute,
					Api: api{
						UrlGet:    "http://localhost:8000/get",
						UrlAdd:    "http://localhost:8000/add",
						UrlRemove: "http://localhost:8000/remove",
						UrlEdit:   "http://localhost:8000/edit",
					},
				},
			},
		},
	}

	// Run tests
	for _, tt := range tests {

		// Set command-line arguments
		err := flag.CommandLine.Parse(tt.args)
		if err != nil {
			t.Fatalf("test %s: %v", tt.name, err)
		}
		// Compare the resulting Config struct with the expected one
		if *cfg != tt.want {
			t.Errorf("test %s readFlags(%v)\n got: %v\nwant: %v", tt.name, tt.args, *cfg, tt.want)
			continue
		}
		t.Log("Good test", tt.name, *cfg)
	}
}

// func TestIsFileEmpty(t *testing.T) {

// 	tests := []struct {
// 		name           string
// 		fileContents   string
// 		expectedResult bool
// 	}{
// 		{
// 			name:           "Empty file",
// 			fileContents:   "",
// 			expectedResult: true,
// 		},
// 		{
// 			name:           "Non-empty file",
// 			fileContents:   "Hello, World!",
// 			expectedResult: false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			f, err := ioutil.TempFile("", "test")
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			defer os.Remove(f.Name())

// 			_, err = f.WriteString(tt.fileContents)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			err = f.Close()
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			if got := isFileEmpty(f.Name()); got != tt.expectedResult {
// 				t.Errorf("isFileEmpty() = %v, want %v", got, tt.expectedResult)
// 			}
// 		})
// 	}
// }
