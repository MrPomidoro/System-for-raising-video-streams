package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB Эта структура будет хранить экземпляр подключения к базе данных.
type DB struct {
	Conn *pgxpool.Pool
}

type IDB interface {
	KeepAlive(ctx context.Context, errCh chan<- error)
	// DBPing(ctx context.Context, cfg *config.Config)
}
