package service

import (
	"database/sql"
	"os"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream/repository"
	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream/usecase"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/database"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/logger"
	"github.com/sirupsen/logrus"
)

type app struct {
	cfg                  *config.Config
	Log                  *logrus.Logger
	db                   *sql.DB
	SigChan              chan os.Signal
	refreshStreamUseCase refreshstream.RefreshStreamUseCase
}

// прототип приложения
func NewApp(cfg *config.Config) *app {
	log := logger.NewLog(cfg.LogLevel)
	db := database.CreateDBConnection(cfg)
	repo := repository.NewRefreshStreamRepository(db, log)
	sigChan := make(chan os.Signal, 1)

	return &app{
		cfg:                  cfg,
		db:                   db,
		Log:                  log,
		SigChan:              sigChan,
		refreshStreamUseCase: usecase.NewRefreshStreamUseCase(repo, db, log),
	}
}

func (a *app) Run() error {
	return nil
}
