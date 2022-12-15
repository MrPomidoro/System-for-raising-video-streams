package database

import (
	"context"

	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
)

type DBI interface {
	CreateDBConnection(cfg *config.Config) *DB
	CloseDBConnection(cfg *config.Config)
	DBPing(ctx context.Context, cfg *config.Config)
}
