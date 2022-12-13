package database

import (
	"time"

	"go.uber.org/zap"
)

// Database - структура с параметрами для базы данных
type Database struct {
	Port     string
	Host     string
	Db_name  string
	User     string
	Password string

	Driver                    string
	DBConnectionTimeoutSecond time.Duration
	Log                       *zap.Logger
}
