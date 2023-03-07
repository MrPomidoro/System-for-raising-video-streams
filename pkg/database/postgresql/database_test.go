package postgresql

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// func TestNewDB(t *testing.T) {
// 	cfg := config.Config{
// 		Database: config.Database{
// 			User:     "sysadmin",
// 			Password: "w3X{77PpCR",
// 			Host:     "192.168.0.32",
// 			DbName:   "www",
// 			Driver:   "postgres",
// 			Port:     5432,
// 		},
// 	}
// 	log := logger.NewLogger(&cfg)
// 	ctx := context.Background()
// 	pool, err := pgxpool.NewWithConfig(ctx, getConfig(cfg.Database))
// 	if err != nil {
// 		t.Error("unexpected error:", err)
// 	}
// 	defer pool.Close()

// 	poolNew, err := NewDB(ctx, cfg.Database, log)
// 	if err != nil {
// 		t.Error("unexpected error:", err)
// 	}

// 	assert.NotNil(t, poolNew)
// 	assert.Equal(t, poolNew.Conn, pool)
// }

func TestNewDB(t *testing.T) {
	cfg := config.Config{
		Database: config.Database{User: "sysadmin",
			Password: "w3X{77PpCR",
			Host:     "192.168.0.32",
			DbName:   "www",
			Driver:   "postgres",
			Port:     5432,
		},
	} // задаем конфигурацию для базы данных
	log := logger.NewLogger(&cfg)
	// создаем mock объект для pgxpool.Pool
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// pool, err := pgxpool.NewWithConfig(context.Background(), getConfig(cfg.Database))
	// pool, err := pgxpool.Connect(context.Background(), "postgres://user:password@localhost/dbname")
	// require.NoError(t, err)

	// создаем объект для тестирования
	// testDB := &DB{pool}

	// имитируем вызов pgxpool.NewWithConfig()
	mock.ExpectQuery("SELECT 1").WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow("1"))
	// newPool, err := testDB.NewPool(ctx, cfg)
	newPool, err := NewDB(context.Background(), cfg.Database, log)
	assert.NoError(t, err)
	assert.NotNil(t, newPool)

	// проверяем, что все вызовы mock объекта были выполнены
	assert.NoError(t, mock.ExpectationsWereMet())
}
