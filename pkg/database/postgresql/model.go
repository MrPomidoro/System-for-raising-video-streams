package postgresql

import (
	"context"

	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	"github.com/golang/mock/gomock"
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

type MockPgxIface struct {
	Conn *pgxpool.Pool
}

func NewMockPgxIface(*gomock.Controller) *DB {

	config := getConfig(config.Database{
		Host:     "192.168.0.32",
		Port:     5432,
		DbName:   "www",
		User:     "sysadmin",
		Password: "w3X{77PpCR",
	})

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil
	}

	return &DB{pool}
}
