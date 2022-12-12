package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/logger"
	_ "github.com/lib/pq"
)

// CreateDBConnection заполняет структуру данными из конфига и вызывает функцию connectToDB(),
// дающую подключение к базе данных
func CreateDBConnection(cfg *config.Config) *sql.DB {
	var dbcfg Database

	dbcfg.Port = cfg.Port
	dbcfg.Host = cfg.Host
	dbcfg.Db_name = cfg.Db_Name
	dbcfg.User = cfg.User
	dbcfg.Password = cfg.Password

	dbcfg.Driver = cfg.Driver
	dbcfg.DBConnectionTimeoutSecond = cfg.Db_Connection_Timeout_Second
	dbcfg.Log = logger.NewLog(cfg.LogLevel)

	return connectToDB(&dbcfg)
}

// connectToDB - функция, возвращающая открытое подключение к базе данных
func connectToDB(dbcfg *Database) *sql.DB {
	var dbSQL *sql.DB

	sqlInfo := fmt.Sprintf(DBInfoConst,
		dbcfg.Host, dbcfg.Port, dbcfg.User, dbcfg.Password,
		dbcfg.Db_name)

	// Подключение
	dbSQL, err := sql.Open(dbcfg.Driver, sqlInfo)
	if err != nil {
		logger.LogError(dbcfg.Log, fmt.Sprintf("cannot get connect to database: %v", err))
	}

	// Проверка подключения
	time.Sleep(time.Millisecond * 3)
	if err := dbSQL.Ping(); err == nil {
		logger.LogInfo(dbcfg.Log, fmt.Sprintf("Success connect to database %s", dbcfg.Db_name))
		return dbSQL
	} else {
		logger.LogError(dbcfg.Log, fmt.Sprintf("cannot connect to database: %s", err))
	}

	connLatency := time.Duration(10 * time.Millisecond)
	time.Sleep(connLatency * time.Millisecond)
	connTimeout := dbcfg.DBConnectionTimeoutSecond
	for t := connTimeout; t > 0; t-- {
		if dbSQL != nil {
			return dbSQL
		}
		time.Sleep(time.Second * 3)
	}

	logger.LogError(dbcfg.Log, fmt.Sprintf("Time waiting of database connection exceeded limit: %v", connTimeout))
	return dbSQL
}

// CloseDBConnection реализует отключение от базы данных
func CloseDBConnection(cfg *config.Config, dbSQL *sql.DB) {
	log := logger.NewLog(cfg.LogLevel)
	if err := dbSQL.Close(); err != nil {
		logger.LogError(log, fmt.Sprintf("cannot close database connection: %v", err))
		return
	}
	logger.LogDebug(log, "Established closing of connection to database")
}

// DBPing реализует переподключение к базе данных при необходимости
// Происходит проверка контекста - если он закрыт, DBPing прекращаеи работу
func DBPing(ctx context.Context, cfg *config.Config, db *sql.DB) {
	log := logger.NewLog(cfg.LogLevel)

loop:
	for {
		errChan := make(chan error)
		defer close(errChan)
		ping(db, errChan)

		select {
		case <-ctx.Done():
			break loop
		case err := <-errChan:
			logger.LogDebug(log, fmt.Sprintf("cannot connect to database %s", err))
			logger.LogDebug(log, "try connect to database...")

			var dbcfg Database
			dbcfg.Port = cfg.Port
			dbcfg.Host = cfg.Host
			dbcfg.Db_name = cfg.Db_Name
			dbcfg.User = cfg.User
			dbcfg.Password = cfg.Password
			dbcfg.Driver = cfg.Driver
			dbcfg.DBConnectionTimeoutSecond = cfg.Db_Connection_Timeout_Second
			dbcfg.Log = log

			connectToDB(&dbcfg)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func ping(db *sql.DB, errChan chan error) {
	err := db.Ping()
	if err != nil {
		errChan <- err
	}
}
