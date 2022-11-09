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
	Repo refreshstream.RefreshStreamRepository // интерфейс

}

func NewRefreshStreamUseCase(Repo refreshstream.RefreshStreamRepository,
	db *sql.DB, log *logrus.Logger) *refreshStreamUseCase {
	return &refreshStreamUseCase{
		db:   db,
		log:  log,
		Repo: Repo,
	}
}

func (s *refreshStreamUseCase) Get(ctx context.Context) ([]refreshstream.RefreshStream, error) {
	return s.Repo.Get(ctx)
}
