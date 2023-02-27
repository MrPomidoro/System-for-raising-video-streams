package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/logger"
	"go.uber.org/zap"

	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
	_ "github.com/lib/pq"
)

// CreateDBConnection заполняет структуру данными из конфига и вызывает функцию connectToDB(),
// дающую подключение к базе данных
func CreateDBConnection(ctx context.Context, cfg *config.Config) (*DB, ce.IError) {
	var db DB
	db.err = ce.ErrorDatabase

	db.port = cfg.Port
	db.host = cfg.Host
	db.dbName = cfg.DbName
	db.user = cfg.User
	db.password = cfg.Password

	db.driver = cfg.Driver
	db.dBConnectionTimeoutSecond = cfg.DbConnectionTimeoutSecond
	db.log = logger.NewLogger(cfg)

	var err error
	db.Db, err = db.connectToDB(*cfg)
	if err != nil {
		return nil, db.err.SetError(err)
	}

	return &db, nil
}

// connectToDB - функция, возвращающая открытое подключение к базе данных
func (db *DB) connectToDB(cfg config.Config) (*sql.DB, error) {
	var dbSQL *sql.DB

	sqlInfo := fmt.Sprintf(DBInfoConst,
		db.host, db.port, db.user, db.password,
		db.dbName)

	// Подключение
	dbSQL, err := sql.Open(db.driver, sqlInfo)
	if err != nil {
		return nil, err
	}

	// Проверка подключения
	time.Sleep(time.Millisecond * 3)
	if err := dbSQL.Ping(); err == nil {
		db.log.Info(fmt.Sprintf("Success connect to database %s", db.dbName))
		return dbSQL, nil
	} else {
		return nil, err
	}
}

// CloseDBConnection реализует отключение от базы данных
func (db *DB) CloseDBConnection(cfg *config.Config) *ce.Error {

	if err := db.Db.Close(); err != nil {
		return db.err.SetError(err)
	}

	db.log.Info("Established closing of connection to database")
	return nil
}

// DBPing реализует переподключение к базе данных при необходимости
// Происходит проверка контекста - если он закрыт, DBPing прекращаеи работу
func (db *DB) DBPing(ctx context.Context, cfg *config.Config, log *zap.Logger, errChan chan error) {

	defer close(errChan)

loop:
	for {
		if ctx.Err() != nil {
			break loop
		}
		go db.ping(ctx, errChan)

		time.Sleep(3 * time.Second)
		select {
		case <-ctx.Done():
			break loop
		case err := <-errChan:
			log.Debug(fmt.Sprintf("cannot connect to database: %s", err))
			log.Info("Try reconnect to database...")

			db.connectToDB(*cfg)
		default:
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (db *DB) ping(ctx context.Context, errChan chan error) {
	if ctx.Err() != nil {
		close(errChan)
		return
	}
	err := db.Db.Ping()
	if err != nil {
		errChan <- err
	}
}
