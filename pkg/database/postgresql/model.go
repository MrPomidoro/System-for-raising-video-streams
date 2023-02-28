package postgresql

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// DB Эта структура будет хранить экземпляр подключения к базе данных.
type DB struct {
	Conn *pgxpool.Pool
}
