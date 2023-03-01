package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// NewDB Эта функция создает новый экземпляр DB.
func NewDB(ctx context.Context, cfg *config.Database, log *zap.Logger) (db *DB, err error) {

	e := ce.ErrorDatabase

	config := getConfig(cfg, log)

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, e.SetError(err)
	}

	return &DB{pool}, nil
}

func (db *DB) KeepAlive(ctx context.Context, log *zap.Logger, errCh chan error) {

	for {
		if ctx.Err() != nil {
			close(errCh)
			return
		}

		go db.ping(ctx, errCh)

		time.Sleep(3 * time.Second)
		select {
		case <-ctx.Done():
			close(errCh)
			return
		case err := <-errCh:
			log.Debug(fmt.Sprintf("cannot connect to database: %s", err))
			log.Info("Try reconnect to database...")

		default:
		}
	}
}

func (db *DB) ping(ctx context.Context, errCh chan error) {
	if ctx.Err() != nil {
		return
	}

	conn, err := db.Conn.Acquire(context.Background())
	if err != nil {
		select {
		case <-ctx.Done():
			return
		case errCh <- fmt.Errorf("failed to acquire connection: %w", err):
		}
		return
	}

	tx, _ := conn.Begin(ctx)
	defer conn.Release()
	// defer tx.Rollback(ctx)

	if _, err = tx.Exec(context.Background(), "SELECT 1"); err != nil {
		select {
		case <-ctx.Done():
			return
		case errCh <- fmt.Errorf("failed to execute test query: %w", err):
		}
		return
	}
}

func (db *DB) IsConn(ctx context.Context) bool {

	if ctx.Err() != nil {
		return false
	}

	conn, err := db.Conn.Acquire(context.Background())
	if err != nil {
		return false
	}

	tx, _ := conn.Begin(ctx)
	defer conn.Release()
	// defer tx.Rollback(ctx)

	if _, err = tx.Exec(context.Background(), "SELECT 1"); err != nil {
		return false
	}

	return true
}

func getConfig(cfg *config.Database, log *zap.Logger) *pgxpool.Config {
	// Настраиваем конфигурацию пула подключений к базе данных
	config, _ := pgxpool.ParseConfig("")
	config.ConnConfig.User = cfg.User
	config.ConnConfig.Password = cfg.Password
	config.ConnConfig.Host = cfg.Host
	config.ConnConfig.Port = uint16(cfg.Port)
	config.ConnConfig.Database = cfg.DbName

	// Устанавливаем максимальное количество соединений в пуле
	config.MaxConns = 2

	return config
}

func (db *DB) Close() {
	db.Conn.Close()
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
