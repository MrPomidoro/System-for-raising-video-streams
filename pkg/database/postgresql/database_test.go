package postgresql

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// В этом тесте можно обойтись и без моков, так как мы проверяем фактическое подключение к бд
func TestDatabaseConnection(t *testing.T) {
	// задаем параметры подключения к базе данных
	config, _ := pgxpool.ParseConfig("")
	config.ConnConfig.User = "sysadmin"
	config.ConnConfig.Password = "w3X{77PpCR"
	config.ConnConfig.Host = "192.168.0.32"
	config.ConnConfig.Port = 5432
	config.ConnConfig.Database = "www"

	// создаем пул соединений к базе данных
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	require.NoError(t, err)
	defer pool.Close()

	db := DB{Conn: pool}

	// делаем запрос к базе данных
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		t.Fatalf("error acquiring connection from pool: %s", err)
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), "SELECT 1")
	if err != nil {
		t.Fatalf("error executing query: %s", err)
	}

	fmt.Println("1")
	if db.GetConn() != pool {
		t.Fatalf("expect: %v, got: %v", pool, db.GetConn())
	}
	fmt.Println("2")

	// Test close connection

	db.Close()
	fmt.Println("3")
	_, err = db.Conn.Exec(context.Background(), "SELECT 1")
	if err == nil {
		t.Fatalf("error executing query: %s", err)
	}
}
