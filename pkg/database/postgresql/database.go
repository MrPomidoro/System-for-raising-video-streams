package postgresql

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// NewDB Эта функция создает новый экземпляр DB.
func NewDB(ctx context.Context, cfg *config.Database, log *zap.Logger) (db *DB, err error) {

	config := getConfig(cfg, log)

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	return &DB{pool}, nil
}

func (db *DB) KeepAlive(ctx context.Context, log *zap.Logger, errCh chan error) {
	defer close(errCh)

loop:
	for {
		if ctx.Err() != nil {
			break loop
		}
		go db.ping(ctx, errCh)

		time.Sleep(3 * time.Second)
		select {
		case <-ctx.Done():
			break loop
		case err := <-errCh:
			log.Debug(fmt.Sprintf("cannot connect to database: %s", err))
			log.Info("Try reconnect to database...")

		default:
		}
		// for {
		// 	if ctx.Err() != nil {
		// 		return
		// 	}
		// 	go db.ping(ctx, errCh)
		// 	time.Sleep(3 * time.Second)
		// 	fmt.Println("Time after 3 second")
		// 	// Выполняем тестовый запрос, чтобы убедиться, что соединение работает
		// 	fmt.Println("1")
		// 	err := db.Conn.Ping(ctx)
		// 	fmt.Println("2")
		// 	if err != nil {
		// 		errCh <- fmt.Errorf("failed to acquire connection: %w", err)
		// 		continue
		// 	}
		//
		//
		//
		// conn, err := db.Conn.Acquire(context.Background())
		// if err != nil {
		// errCh <- fmt.Errorf("failed to acquire connection: %w", err)
		// continue
		// }
		// fmt.Println("2")
		// defer conn.Release()
		// if _, err = conn.Exec(context.Background(), "SELECT 1"); err != nil {
		// 	fmt.Println("3")
		// 	errCh <- fmt.Errorf("failed to execute test query: %w", err)
		// 	continue
		// }
		// fmt.Println("4")
	}
}

func (db *DB) ping(ctx context.Context, errChan chan error) {
	if ctx.Err() != nil {
		return
	}
	fmt.Println("try ping")
	err := db.Conn.Ping(ctx)
	fmt.Println("err of ping:", err)
	if err != nil {
		errChan <- err
	}
}

// func (db *DB) reconnect() error {
// 	var conn *pgx.Conn
// 	var err error
// 	for {
// 		conn, err = pgx.ConnectConfig(context.Background(), db.Conn.Config())
// 		if err != nil {
// 			fmt.Printf("failed to reconnect to database: %v\n", err)
// 			time.Sleep(5 * time.Second)
// 			continue
// 		}
// 		break
// 	}
// 	// if err = db.Conn.Close(context.Background()); err != nil {
// 	// 	return err
// 	// }
// 	db.Conn = conn
// 	fmt.Println("successfully reconnected to database")
// 	return nil
// }

func getConfig(cfg *config.Database, log *zap.Logger) *pgxpool.Config {
	// Настраиваем конфигурацию пула подключений к базе данных
	config, _ := pgxpool.ParseConfig("")
	config.ConnConfig.User = cfg.User
	config.ConnConfig.Password = cfg.Password
	config.ConnConfig.Host = cfg.Host
	config.ConnConfig.Port = uint16(cfg.Port)
	config.ConnConfig.Database = cfg.DbName

	// Устанавливаем максимальное количество соединений в пуле
	config.MaxConns = 1

	return config
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
