package usecase

import (
	"context"
	"database/sql"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/statusstream"
	"github.com/sirupsen/logrus"
)

type statusStreamUseCase struct {
	db   *sql.DB
	log  *logrus.Logger
	repo statusstream.StatusStreamRepository // интерфейс

}

func NewStatusStreamUseCase(repo statusstream.StatusStreamRepository,
	db *sql.DB, log *logrus.Logger) *statusStreamUseCase {
	return &statusStreamUseCase{
		db:   db,
		log:  log,
		repo: repo,
	}
}

func (s *statusStreamUseCase) Get(ctx context.Context) ([]statusstream.StatusStream, error) {
	return s.repo.Get(ctx)
}
