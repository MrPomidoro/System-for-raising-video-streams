package postgresql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

// NewDB Эта функция создает новый экземпляр DB.
func NewDB(ctx context.Context, cfg *config.Database, log *zap.Logger) (db *DB, err error) {

	if cfg.Driver != "postgres" {
		return nil, errors.New("this driver not has use in connection postgresql database")
	}

	c := GetConfig(cfg, log)
	conn, err := pgx.ConnectConfig(ctx, c)
	if err != nil {
		return nil, err
	}

	db = &DB{c, conn}

	go db.keepAlive()

	return db, nil
}

func (db *DB) keepAlive() {
	for {
		time.Sleep(5 * time.Second)

		if err := db.Conn.Ping(context.Background()); err != nil {
			fmt.Printf("lost database connection: %v\n", err)

			if err = db.reconnect(); err != nil {
				fmt.Printf("failed to reconnect to database: %v\n", err)
			}
		}
	}
}

func (db *DB) reconnect() error {
	var conn *pgx.Conn
	var err error

	for {
		conn, err = pgx.ConnectConfig(context.Background(), db.Conn.Config())
		if err != nil {
			fmt.Printf("failed to reconnect to database: %v\n", err)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	// if err = db.Conn.Close(context.Background()); err != nil {
	// 	return err
	// }
	db.Conn = conn

	fmt.Println("successfully reconnected to database")

	return nil
}

func GetConfig(cfg *config.Database, log *zap.Logger) *pgx.ConnConfig {
	// c := &pgx.ConnConfig{}
	c, _ := pgx.ParseConfig("")
	c.Host = cfg.Host
	c.Port = uint16(cfg.Port)
	c.Database = cfg.DbName
	c.User = cfg.User
	c.Password = cfg.Password
	c.ConnectTimeout = cfg.ConnectionTimeout

	return c
}

//// Connection заполняет структуру данными из конфига и вызывает функцию db(),
//// дающую подключение к базе данных
//func Connection(cfg *config.Database, log *zap.Logger) (*DB, ce.IError) {
//	var db DB
//	db.err = ce.ErrorDatabase
//	db.driver = cfg.Driver
//	db.dBConnectionTimeoutSecond = cfg.DbConnectionTimeoutSecond
//	db.log = log
//
//	sqlInfo := fmt.Sprintf(DBInfoConst, cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DbName)
//
//	var err error
//	db.Db, err = db.db(sqlInfo)
//	if err != nil {
//		return nil, db.err.SetError(err)
//	}
//
//	return &db, nil
//}
//
//// db - функция, возвращающая открытое подключение к базе данных
//func (db *DB) db(sqlInfo string) (dbSQL *sql.DB, err error) {
//
//	// Подключение
//	dbSQL, err = sql.Open(db.driver, sqlInfo)
//	if err != nil {
//		return nil, err
//	}
//
//	// Проверка подключения
//	time.Sleep(time.Millisecond * 3)
//	if err = dbSQL.Ping(); err == nil {
//		db.log.Info(fmt.Sprintf("Success connect to database"))
//		return dbSQL, nil
//	} else {
//		return nil, err
//	}
//}
//
//// Close реализует отключение от базы данных
//func (db *DB) Close() *ce.Error {
//
//	if err := db.Db.Close(); err != nil {
//		return db.err.SetError(err)
//	}
//
//	db.log.Info("Established closing of connection to database")
//	return nil
//}
//
//// Ping реализует переподключение к базе данных при необходимости
//// Происходит проверка контекста - если он закрыт, Ping прекращаеи работу
//func (db *DB) Ping(ctx context.Context, log *zap.Logger, errChan chan error) {
//
//	defer close(errChan)
//
//loop:
//	for {
//		if ctx.Err() != nil {
//			break loop
//		}
//		go db.ping(ctx, errChan)
//
//		time.Sleep(3 * time.Second)
//		select {
//		case <-ctx.Done():
//			break loop
//		case err := <-errChan:
//			log.Debug(fmt.Sprintf("cannot connect to database: %s", err))
//			log.Info("Try reconnect to database...")
//
//			_, err = db.db()
//			if err != nil {
//				log.Error(db.err.SetError(err).Error())
//				continue
//			}
//		default:
//		}
//	}
//}
//
//func (db *DB) ping(ctx context.Context, errChan chan error) {
//	if ctx.Err() != nil {
//		return
//	}
//	err := db.Db.Ping()
//	if err != nil {
//		errChan <- err
//	}
//}
