package postgresql

import (
	"context"

	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
)

type DBI interface {
	DBPing(ctx context.Context, cfg *config.Config)
}
