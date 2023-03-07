package logger

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func TestLogLevel(t *testing.T) {
	tests := []struct {
		name string
		cfg  config.Config
	}{
		{
			name: "TestInfoLevel",
			cfg:  config.Config{Logger: config.Logger{LogLevel: "INFO"}},
		},
		{
			name: "TestInvalidLevel",
			cfg:  config.Config{Logger: config.Logger{LogLevel: "invalid"}},
		},
	}
	log := logger{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if log.logLevel(&tt.cfg) != zapcore.InfoLevel {
				t.Errorf("expect %v, got %v", log.logLevel(&tt.cfg), zapcore.InfoLevel)
			}
		})
	}
}

func TestNewLogger(t *testing.T) {
	l := &logger{}
	li := Logger(l)
	path := "/home/ksenia/go/src/github.com/Kseniya-cha/System-for-raising-video-streams/pkg/logger/logTest.out"

	conf := li.newProductionEncoderConfig()
	conf.EncodeTime = zapcore.ISO8601TimeEncoder
	conf.EncodeLevel = li.customEncoderLevel
	conf.MessageKey = "message"
	conf.CallerKey = "caller"
	conf.TimeKey = "time"
	jsonEncoder := li.newJSONEncode(conf)
	textEncoder := li.newConsoleEncoder(conf)

	tests := []struct {
		name             string
		cfg              config.Config
		logFile          zapcore.Core
		logStdout        zapcore.Core
		isNeedCreateFile bool
	}{
		{
			name: "TestFileStdout",
			cfg: config.Config{Logger: config.Logger{
				LogLevel: "INFO", RewriteLog: false, LogFile: path,
				LogFileEnable: true, LogStdoutEnable: true, MaxSize: 100,
				MaxAge: 28, MaxBackups: 7,
			}},
			logFile: li.newCore(
				jsonEncoder,
				li.addSync(&lumberjack.Logger{
					Filename:   path,
					MaxSize:    100,
					MaxAge:     28,
					MaxBackups: 7,
				}),
				zapcore.InfoLevel,
			),
			logStdout: li.newCore(
				textEncoder,
				zapcore.AddSync(os.Stdout),
				zapcore.InfoLevel,
			),
			isNeedCreateFile: false,
		},
		{
			name: "TestFile",
			cfg: config.Config{Logger: config.Logger{
				LogLevel: "INFO", RewriteLog: false, LogFile: path,
				LogFileEnable: true, LogStdoutEnable: false, MaxSize: 100,
				MaxAge: 28, MaxBackups: 7,
			}},
			logFile: li.newCore(
				jsonEncoder,
				li.addSync(&lumberjack.Logger{
					Filename:   path,
					MaxSize:    100,
					MaxAge:     28,
					MaxBackups: 7,
				}),
				zapcore.InfoLevel,
			),
			isNeedCreateFile: false,
		},
		{
			name: "TestStdout",
			cfg: config.Config{Logger: config.Logger{
				LogLevel: "INFO", RewriteLog: false, LogFile: path,
				LogFileEnable: false, LogStdoutEnable: true, MaxSize: 100,
				MaxAge: 28, MaxBackups: 7,
			}},
			logStdout: li.newCore(
				textEncoder,
				zapcore.AddSync(os.Stdout),
				zapcore.InfoLevel,
			),
			isNeedCreateFile: false,
		},
		{
			name: "TestBothFalse",
			cfg: config.Config{Logger: config.Logger{
				LogLevel: "INFO", RewriteLog: false, LogFile: path,
				LogFileEnable: false, LogStdoutEnable: false, MaxSize: 100,
				MaxAge: 28, MaxBackups: 7,
			}},
			isNeedCreateFile: false,
		},
		{
			name: "TestRewrite/FileNotExists",
			cfg: config.Config{Logger: config.Logger{
				LogLevel: "INFO", RewriteLog: true, LogFile: path,
				LogFileEnable: true, LogStdoutEnable: false, MaxSize: 100,
				MaxAge: 28, MaxBackups: 7,
			}},
			logFile: li.newCore(
				jsonEncoder,
				li.addSync(&lumberjack.Logger{
					Filename:   path,
					MaxSize:    100,
					MaxAge:     28,
					MaxBackups: 7,
				}),
				zapcore.InfoLevel,
			),
			isNeedCreateFile: false,
		},
		{
			name: "TestRewrite/FileExists",
			cfg: config.Config{Logger: config.Logger{
				LogLevel: "INFO", RewriteLog: true, LogFile: path,
				LogFileEnable: true, LogStdoutEnable: false, MaxSize: 100,
				MaxAge: 28, MaxBackups: 7,
			}},
			logFile: li.newCore(
				jsonEncoder,
				li.addSync(&lumberjack.Logger{
					Filename:   path,
					MaxSize:    100,
					MaxAge:     28,
					MaxBackups: 7,
				}),
				zapcore.InfoLevel,
			),
			isNeedCreateFile: true,
		},
	}

	isEqIdx2Log := []int{2, 3, 4, 5, 7, 8, 9}
	isEqIdx1Log := []int{1, 2, 3, 4, 6, 7, 8}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.isNeedCreateFile {
				os.Create(path)
			}
			defer os.Remove(path)

			got := strings.Split(fmt.Sprint(NewLogger(&tt.cfg)), " ")

			var expect []string
			if !tt.cfg.LogFileEnable && tt.cfg.LogStdoutEnable {
				expect = strings.Split(fmt.Sprint(li.new(li.newTee(tt.logStdout))), " ")
			} else if tt.cfg.LogFileEnable && !tt.cfg.LogStdoutEnable {
				expect = strings.Split(fmt.Sprint(li.new(li.newTee(tt.logFile))), " ")
			} else if !tt.cfg.LogFileEnable && !tt.cfg.LogStdoutEnable {
				expect = strings.Split(fmt.Sprint(nil), " ")
			} else {
				expect = strings.Split(fmt.Sprint(li.new(li.newTee(li.loggers(&tt.cfg, tt.logStdout,
					tt.logFile)...), li.zapOpts()...)), " ")
			}

			if len(expect) != len(got) {
				t.Errorf("expect %v\n\t\tgot    %v", expect, got)
			}

			if len(got) == 10 {
				for _, idx := range isEqIdx2Log {
					if got[idx] != expect[idx] {
						t.Errorf("expect %v\n\t\tgot    %v", expect[idx], got[idx])
					}
				}
			}
			if len(got) == 9 {
				for _, idx := range isEqIdx1Log {
					if got[idx] != expect[idx] {
						t.Errorf("expect %v\n\t\tgot    %v", expect[idx], got[idx])
					}
				}
			}
		})
	}
}
