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
func (datab *DB) CreateDBConnection(cfg *config.Config) *DB {
	// var dbcfg DB

	datab.Port = cfg.Port
	datab.Host = cfg.Host
	datab.Db_name = cfg.Db_Name
	datab.User = cfg.User
	datab.Password = cfg.Password

	datab.Driver = cfg.Driver
	datab.DBConnectionTimeoutSecond = cfg.Db_Connection_Timeout_Second
	datab.Log = logger.NewLogger(cfg)

	datab.Db = datab.connectToDB()

	return datab
}

// connectToDB - функция, возвращающая открытое подключение к базе данных
func (datab *DB) connectToDB() *sql.DB {
	var dbSQL *sql.DB

	sqlInfo := fmt.Sprintf(DBInfoConst,
		datab.Host, datab.Port, datab.User, datab.Password,
		datab.Db_name)

	// Подключение
	dbSQL, err := sql.Open(datab.Driver, sqlInfo)
	if err != nil {
		datab.Log.Error(fmt.Sprintf("cannot get connect to database: %v", err))
	}

	// Проверка подключения
	time.Sleep(time.Millisecond * 3)
	if err := dbSQL.Ping(); err == nil {
		datab.Log.Info(fmt.Sprintf("Success connect to database %s", datab.Db_name))
		return dbSQL
	} else {
		datab.Log.Error(fmt.Sprintf("cannot connect to database: %s", err))
	}

	connLatency := time.Duration(10 * time.Millisecond)
	time.Sleep(connLatency * time.Millisecond)
	connTimeout := datab.DBConnectionTimeoutSecond
	for t := connTimeout; t > 0; t-- {
		if dbSQL != nil {
			return dbSQL
		}
		time.Sleep(time.Second * 3)
	}

	datab.Log.Warn(fmt.Sprintf("Time waiting of database connection exceeded limit: %v", connTimeout))
	return dbSQL
}

// CloseDBConnection реализует отключение от базы данных
func (datab *DB) CloseDBConnection(cfg *config.Config) {
	log := logger.NewLogger(cfg)
	if err := datab.Db.Close(); err != nil {
		log.Error(fmt.Sprintf("cannot close database connection: %v", err))
		return
	}
	log.Debug("Established closing of connection to database")
}

// DBPing реализует переподключение к базе данных при необходимости
// Происходит проверка контекста - если он закрыт, DBPing прекращаеи работу
func (datab *DB) DBPing(ctx context.Context, cfg *config.Config) {
	log := logger.NewLogger(cfg)

loop:
	for {
		errChan := make(chan error)
		defer close(errChan)
		datab.ping(errChan)

		select {
		case <-ctx.Done():
			break loop
		case err := <-errChan:
			log.Debug(fmt.Sprintf("cannot connect to database %s", err))
			log.Debug("try connect to database...")

			var datab DB
			datab.Port = cfg.Port
			datab.Host = cfg.Host
			datab.Db_name = cfg.Db_Name
			datab.User = cfg.User
			datab.Password = cfg.Password
			datab.Driver = cfg.Driver
			datab.DBConnectionTimeoutSecond = cfg.Db_Connection_Timeout_Second
			datab.Log = log

			datab.connectToDB()
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (datab *DB) ping(errChan chan error) {
	err := datab.Db.Ping()
	if err != nil {
		errChan <- err
	}
}
