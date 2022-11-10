package service

import (
	"context"
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

// Прототип приложения
type app struct {
	cfg                  *config.Config
	Log                  *logrus.Logger
	db                   *sql.DB
	SigChan              chan os.Signal
	refreshStreamUseCase refreshstream.RefreshStreamUseCase
}

// Функция, инициализирующая прототип приложения
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

// Алгоритм
func (a *app) Run() error {
	ctx := context.Background()
	logger.LogDebug(a.Log, "Context initializated")

	req, err := a.refreshStreamUseCase.Get(ctx)
	if err != nil {
		logger.LogError(a.Log, fmt.Sprintf("cannot get response from database: %v", err))
	} else {
		logger.LogDebug(a.Log, fmt.Sprintf("Response from database:\n%v", req))
	}

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
