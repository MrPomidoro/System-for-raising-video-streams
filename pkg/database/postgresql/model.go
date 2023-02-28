package postgresql

import (
	"github.com/jackc/pgx/v4"
)

// DB Эта структура будет хранить экземпляр подключения к базе данных.
type DB struct {
	pgConfig *pgx.ConnConfig
	Conn     *pgx.Conn
}
