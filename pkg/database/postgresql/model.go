package postgresql

import (
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

// DB Эта структура будет хранить экземпляр подключения к базе данных.
type DB struct {
	pgConfig *pgx.ConnConfig
	Conn     *pgx.Conn
	log      *zap.Logger
}
