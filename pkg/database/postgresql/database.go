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
func NewDB(ctx context.Context, cfg config.Database, log *zap.Logger) (db *DB, err ce.IError) {

	err = ce.ErrorDatabase

	config := getConfig(cfg)

	pool, e := pgxpool.NewWithConfig(ctx, config)
	if e != nil {
		return nil, err.SetError(e)
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

func getConfig(cfg config.Database) *pgxpool.Config {
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
