package usecase

import (
	"context"
	"database/sql"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
)

type refreshStreamUseCase struct {
	db   *sql.DB
	repo refreshstream.RefreshStreamRepository
}

func NewRefreshStreamUseCase(repo refreshstream.RefreshStreamRepository,
	db *sql.DB) *refreshStreamUseCase {
	return &refreshStreamUseCase{
		db:   db,
		repo: repo,
	}
}

func (s *refreshStreamUseCase) Get(ctx context.Context, status bool) ([]refreshstream.RefreshStream, error) {
	return s.repo.Get(ctx, status)
}
