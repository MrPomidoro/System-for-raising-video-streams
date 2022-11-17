package usecase

import (
	"context"
	"database/sql"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/statusstream"
)

type statusStreamUseCase struct {
	db   *sql.DB
	repo statusstream.StatusStreamRepository // интерфейс

}

func NewStatusStreamUseCase(repo statusstream.StatusStreamRepository,
	db *sql.DB) *statusStreamUseCase {
	return &statusStreamUseCase{
		db:   db,
		repo: repo,
	}
}

func (s *statusStreamUseCase) Insert(ctx context.Context,
	ss *statusstream.StatusStream) error {

	return s.repo.Insert(ctx, ss)
}
