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

// type MockPgxPoolIface interface {
// 	// using pgxpool interface
// 	Begin(context.Context) (pgx.Tx, error)
// 	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
// 	QueryRow(context.Context, string, ...interface{}) pgx.Row
// 	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
// 	Ping(context.Context) error
// 	Close()
// }

// type MockPgxIface struct {
// 	Conn MockPgxPoolIface
// }

// func NewMockPgxIface(*gomock.Controller) *MockPgxIface {

// 	config := getConfig(config.Database{
// 		Host:     "192.168.0.32",
// 		Port:     5432,
// 		DbName:   "www",
// 		User:     "sysadmin",
// 		Password: "w3X{77PpCR",
// 	})

// 	pool, err := pgxpool.NewWithConfig(context.Background(), config)
// 	if err != nil {
// 		return nil
// 	}

// 	return &MockPgxIface{pool}
// }
