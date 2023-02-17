package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/logger"

	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
	_ "github.com/lib/pq"
)

// CreateDBConnection заполняет структуру данными из конфига и вызывает функцию connectToDB(),
// дающую подключение к базе данных
func CreateDBConnection(cfg *config.Config) (*DB, error) {
	var db DB
	db.err.SetLevel(ce.FatalLevel)
	db.err.SetCode("50.4.1")
	db.err.SetDesc("error at database operation level")

	db.port = cfg.Port
	db.host = cfg.Host
	db.dbName = cfg.DbName
	db.user = cfg.User
	db.password = cfg.Password

	db.driver = cfg.Driver
	db.dBConnectionTimeoutSecond = cfg.DbConnectionTimeoutSecond
	db.log = logger.NewLogger(cfg)

	var err error
	db.Db, err = db.connectToDB()

	return &db, db.err.SetError(err)
}

// connectToDB - функция, возвращающая открытое подключение к базе данных
func (db *DB) connectToDB() (*sql.DB, error) {
	var dbSQL *sql.DB

	sqlInfo := fmt.Sprintf(DBInfoConst,
		db.host, db.port, db.user, db.password,
		db.dbName)

	// Подключение
	dbSQL, err := sql.Open(db.driver, sqlInfo)
	if err != nil {
		return nil, db.err.SetError(err)
	}

	// Проверка подключения
	time.Sleep(time.Millisecond * 3)
	if err := dbSQL.Ping(); err == nil {
		db.log.Info(fmt.Sprintf("Success connect to database %s", db.dbName))
		return dbSQL, nil
	} else {
		db.log.Error(fmt.Sprintf("cannot connect to database: %v", err))
	}

	connLatency := time.Duration(10 * time.Millisecond)
	time.Sleep(connLatency)
	connTimeout := db.dBConnectionTimeoutSecond
	for t := connTimeout; t > 0; t-- {
		if dbSQL != nil {
			return dbSQL, nil
		}
		time.Sleep(time.Second * 3)
	}

	db.log.Warn(fmt.Sprintf("Time waiting of database connection exceeded limit: %v", connTimeout))
	return dbSQL, nil
}

// CloseDBConnection реализует отключение от базы данных
func (db *DB) CloseDBConnection(cfg *config.Config) error {

	if err := db.Db.Close(); err != nil {
		db.log.Error(fmt.Sprintf("cannot close database connection: %v", err))
		return db.err.SetError(err)
	}

	db.log.Debug("Established closing of connection to database")
	return nil
}

// DBPing реализует переподключение к базе данных при необходимости
// Происходит проверка контекста - если он закрыт, DBPing прекращаеи работу
func (db *DB) DBPing(ctx context.Context, cfg *config.Config) {

loop:
	for {
		errChan := make(chan error)
		defer close(errChan)
		db.ping(errChan)

		select {
		case <-ctx.Done():
			break loop
		case err := <-errChan:
			db.log.Debug(fmt.Sprintf("cannot connect to database %s", err))
			db.log.Debug("Try reconnect to database...")

			var db DB
			db.port = cfg.Port
			db.host = cfg.Host
			db.dbName = cfg.DbName
			db.user = cfg.User
			db.password = cfg.Password
			db.driver = cfg.Driver
			db.dBConnectionTimeoutSecond = cfg.DbConnectionTimeoutSecond

			db.connectToDB()
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (db *DB) ping(errChan chan error) {
	err := db.Db.Ping()
	if err != nil {
		errChan <- err
	}
}
