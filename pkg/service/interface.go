package service

import (
	"context"
)

type App interface {
	Run(context.Context)
	GracefulShutdown(ctx context.Context, cancel context.CancelFunc)
}
