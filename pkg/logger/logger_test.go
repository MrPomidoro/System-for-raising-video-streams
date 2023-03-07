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
	cfg := config.Config{Logger: config.Logger{LogLevel: "INFO"}}
	log := logger{}

	if log.logLevel(&cfg) != zapcore.InfoLevel {
		t.Errorf("expect %v, got %v", log.logLevel(&cfg), zapcore.InfoLevel)
	}
}

func TestNewLogger(t *testing.T) {
	l := &logger{}
	li := Logger(l)

	conf := li.newProductionEncoderConfig()
	conf.EncodeTime = zapcore.ISO8601TimeEncoder
	conf.EncodeLevel = li.customEncoderLevel
	conf.MessageKey = "message"
	conf.CallerKey = "caller"
	conf.TimeKey = "time"
	jsonEncoder := li.newJSONEncode(conf)
	textEncoder := li.newConsoleEncoder(conf)

	tests := []struct {
		name      string
		cfg       config.Config
		logFile   zapcore.Core
		logStdout zapcore.Core
	}{
		{
			name: "TestFileStdout",
			cfg: config.Config{Logger: config.Logger{
				LogLevel: "INFO", RewriteLog: false, LogFile: "./pkg/logger/out.log",
				LogFileEnable: true, LogStdoutEnable: true, MaxSize: 100,
				MaxAge: 28, MaxBackups: 7,
			}},
			logFile: li.newCore(
				jsonEncoder,
				li.addSync(&lumberjack.Logger{
					Filename:   "./pkg/logger/out.log",
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
		},
		{
			name: "TestFile",
			cfg: config.Config{Logger: config.Logger{
				LogLevel: "INFO", RewriteLog: false, LogFile: "./pkg/logger/out.log",
				LogFileEnable: true, LogStdoutEnable: false, MaxSize: 100,
				MaxAge: 28, MaxBackups: 7,
			}},
			logFile: li.newCore(
				jsonEncoder,
				li.addSync(&lumberjack.Logger{
					Filename:   "./pkg/logger/out.log",
					MaxSize:    100,
					MaxAge:     28,
					MaxBackups: 7,
				}),
				zapcore.InfoLevel,
			),
		},
		{
			name: "TestStdout",
			cfg: config.Config{Logger: config.Logger{
				LogLevel: "INFO", RewriteLog: false, LogFile: "./pkg/logger/out.log",
				LogFileEnable: false, LogStdoutEnable: true, MaxSize: 100,
				MaxAge: 28, MaxBackups: 7,
			}},
			logStdout: li.newCore(
				textEncoder,
				zapcore.AddSync(os.Stdout),
				zapcore.InfoLevel,
			),
		},
		{
			name: "TestBothFalse",
			cfg: config.Config{Logger: config.Logger{
				LogLevel: "INFO", RewriteLog: false, LogFile: "./pkg/logger/out.log",
				LogFileEnable: false, LogStdoutEnable: false, MaxSize: 100,
				MaxAge: 28, MaxBackups: 7,
			}},
		},
	}

	isEqIdx2Log := []int{2, 3, 4, 5, 7, 8, 9}
	isEqIdx1Log := []int{1, 2, 3, 4, 6, 7, 8}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

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
