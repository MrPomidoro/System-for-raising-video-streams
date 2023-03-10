package postgresql

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// В этом тесте можно обойтись и без моков, так как мы проверяем фактическое подключение к бд
func TestDatabaseConnection(t *testing.T) {
	// задаем параметры подключения к базе данных
	conf, _ := pgxpool.ParseConfig("")
	conf.ConnConfig.User = "sysadmin"
	conf.ConnConfig.Password = "w3X{77PpCR"
	conf.ConnConfig.Host = "192.168.0.32"
	conf.ConnConfig.Port = 5432
	conf.ConnConfig.Database = "www"
	conf.MaxConns = 2

	// создаем пул соединений к базе данных
	pool, err := pgxpool.NewWithConfig(context.Background(), conf)
	require.NoError(t, err)
	defer pool.Close()

	// делаем запрос к базе данных
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		t.Errorf("error acquiring connection from pool: %s", err)
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), "SELECT 1")
	if err != nil {
		t.Errorf("error executing query: %s", err)
	}
	cfg := config.Config{Database: config.Database{
		DbName:   "www",
		User:     "sysadmin",
		Port:     5432,
		Host:     "192.168.0.32",
		Password: "w3X{77PpCR",
	}}

	newdb, err := NewDB(context.Background(), cfg.Database, zap.NewNop())
	if err != nil {
		t.Errorf("error executing query: %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	t.Run("KeepAliveConnOk", func(t *testing.T) {
		errCh := make(chan error)
		defer close(errCh)
		go newdb.ping(ctx, errCh)
		time.Sleep(500 * time.Millisecond)

		select {
		case <-ctx.Done():
			close(errCh)
			t.Error("context is closed")
		case err := <-errCh:
			t.Errorf("got value from errChan: %v", err)
		default:
		}
	})

	t.Run("TestNewDBandGetConn", func(t *testing.T) {
		newPool := newdb.GetConn()
		newdbS := strings.Split(fmt.Sprint(newPool), " ")
		dbS := strings.Split(fmt.Sprint(pool), " ")
		indexes := []int{0, 1, 2, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 16, 17}

		for _, idx := range indexes {
			if newdbS[idx] != dbS[idx] {
				t.Errorf("expect: %v, got: %v", dbS[idx], newdbS[idx])
			}
		}
	})

	t.Run("TestIsConn", func(t *testing.T) {
		isConn := newdb.IsConn(context.Background())
		if !isConn {
			t.Errorf("expect connection to database")
		}
	})

	cancel()

	t.Run("TestIsConnCtxCancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		isConn := newdb.IsConn(ctx)
		if isConn {
			t.Errorf("expect no connection to database")
		}
	})

	t.Run("KeepAliveCtxCancel", func(t *testing.T) {
		errCh := make(chan error)
		defer close(errCh)
		go newdb.ping(ctx, errCh)
		time.Sleep(500 * time.Millisecond)

		after := time.After(500 * time.Millisecond)
		select {
		case <-ctx.Done():
		case err := <-errCh:
			t.Errorf("got value from errChan: %v", err)
		case <-after:
			t.Errorf("expexct error: closed pool; got nil")
		}
	})

	t.Run("TestCloseConnection", func(t *testing.T) {
		newdb.Close()
		_, err = newdb.Conn.Exec(context.Background(), "SELECT 1")
		if err == nil {
			t.Errorf("error closing connection: %s", err)
		}
	})

	t.Run("KeepAliveConnClosed", func(t *testing.T) {
		errCh := make(chan error)
		ctx = context.Background()
		go newdb.ping(ctx, errCh)
		time.Sleep(500 * time.Millisecond)

		after := time.After(500 * time.Millisecond)
		select {
		case <-ctx.Done():
			t.Error("context is closed")
		case <-errCh:
		case <-after:
			t.Errorf("expexct error: closed pool; got nil")
		}
		close(errCh)
	})

	t.Run("TestIsConnClosed", func(t *testing.T) {
		isConn := newdb.IsConn(context.Background())
		if isConn {
			t.Errorf("expect no connection to database")
		}
	})
}
