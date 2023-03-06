package postgresql

import (
	"context"
	"testing"

	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	cfg := config.Config{
		Database: config.Database{
			User:     "sysadmin",
			Password: "w3X{77PpCR",
			Host:     "192.168.0.32",
			DbName:   "www",
			Driver:   "postgres",
			Port:     5432,
		},
	}
	log := logger.NewLogger(&cfg)
	ctx := context.Background()
	pool, err := pgxpool.NewWithConfig(ctx, getConfig(cfg.Database))
	if err != nil {
		t.Error("unexpected error:", err)
	}
	defer pool.Close()

	poolNew, err := NewDB(ctx, cfg.Database, log)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	assert.NotNil(t, poolNew)
	assert.Equal(t, poolNew.Conn, pool)
}
