package service

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

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

// Метод структуры app, закрывающий канал
func (a *app) StopChan(sig chan<- os.Signal) {
	close(sig)
}

// Метод для корректного завершения работы программы
// при получении прерывающего сигнала
func (a *app) GracefulShutdown(sig chan os.Signal) {
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sign := <-sig:
		logger.LogWarn(a.Log, fmt.Sprintf("Got signal: %v, exiting", sign))
		database.CloseDBConnection(a.cfg, a.db)
		time.Sleep(time.Second * 2)
		a.StopChan(sig)
	}
}
