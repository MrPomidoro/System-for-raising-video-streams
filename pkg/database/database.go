package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/logger"
	_ "github.com/lib/pq"
)

// заполняется структура из конфига, вызывается функция connectToDB()
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

// подключение к базе данных
func connectToDB(dbcfg *Database) *sql.DB {
	var dbSQL *sql.DB

	sqlInfo := fmt.Sprintf(DBInfoConst,
		dbcfg.Host, dbcfg.Port, dbcfg.User, dbcfg.Password,
		dbcfg.Db_name)

	// подключение
	dbSQL, err := sql.Open(dbcfg.Driver, sqlInfo)
	if err != nil {
		logger.LogError(dbcfg.Log, fmt.Sprintf("cannot get connect to database: %v", err))
	}

	// проверка подключения
	time.Sleep(time.Millisecond * 3)
	if err := dbSQL.Ping(); err == nil {
		logger.LogDebug(dbcfg.Log, fmt.Sprintf("Success connect to database %s", dbcfg.Db_name))
		return dbSQL
	} else {
		logger.LogError(dbcfg.Log, "cannot connect to database %s")
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

	logger.LogError(dbcfg.Log, fmt.Sprintf("Time waiting of DB connection exceeded limit: %v", connTimeout))
	return dbSQL
}

// отключение от базы данных
func CloseDBConnection(cfg *config.Config, dbSQL *sql.DB) {
	log := logger.NewLog(cfg.LogLevel)
	if err := dbSQL.Close(); err != nil {
		logger.LogError(log, fmt.Sprintf("cannot close DB connection. Error: %v", err))
		return
	}
	logger.LogDebug(log, "Established closing of connection to DB")
}
