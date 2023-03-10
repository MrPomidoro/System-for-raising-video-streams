package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// DB Эта структура будет хранить экземпляр подключения к базе данных.
type DB struct {
	Conn *pgxpool.Pool
}

type IDB interface {
	KeepAlive(ctx context.Context, log *zap.Logger, errCh chan error)
	IsConn(ctx context.Context) bool
	Close()
	GetConn() *pgxpool.Pool
}
