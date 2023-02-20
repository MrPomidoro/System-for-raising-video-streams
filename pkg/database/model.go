package database

import (
	"database/sql"
	"time"

	"go.uber.org/zap"

	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
)

// Database - структура с параметрами для базы данных
type database struct {
	port     string
	host     string
	dbName   string
	user     string
	password string

	driver                    string
	dBConnectionTimeoutSecond time.Duration
	log                       *zap.Logger

	err ce.IError
}

type DB struct {
	database
	Db *sql.DB
}
