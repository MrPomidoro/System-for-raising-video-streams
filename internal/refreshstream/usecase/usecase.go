package usecase

import (
	"context"
	"database/sql"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	"github.com/sirupsen/logrus"
)

type refreshStreamUseCase struct {
	db   *sql.DB
	log  *logrus.Logger
	repo refreshstream.RefreshStreamRepository // интерфейс

}

func NewRefreshStreamUseCase(repo refreshstream.RefreshStreamRepository,
	db *sql.DB, log *logrus.Logger) *refreshStreamUseCase {
	return &refreshStreamUseCase{
		db:   db,
		log:  log,
		repo: repo,
	}
}

func (s *refreshStreamUseCase) Get(ctx context.Context) ([]refreshstream.RefreshStream, error) {
	return s.repo.Get(ctx)
}
