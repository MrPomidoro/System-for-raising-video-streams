package service

import (
	"context"
	"os"

	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/database/postgresql"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	rsrepository "github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream/repository"
	rtspsimpleserver "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server"
	rtsprepository "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server/repository"
	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/statusstream"
	ssrepository "github.com/Kseniya-cha/System-for-raising-video-streams/internal/statusstream/repository"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/logger"
	"go.uber.org/zap"
)

type App interface {
	Run(context.Context)
	GracefulShutdown(cancel context.CancelFunc)
}

// app - прототип приложения
type app struct {
	cfg *config.Config
	log *zap.Logger
	db  *postgresql.DB

	sigChan  chan os.Signal
	doneChan chan struct{}

	refreshStreamRepo refreshstream.RefreshStreamRepository
	statusStreamRepo  statusstream.Repository
	rtspRepo          rtspsimpleserver.RTSPRepository

	err ce.IError
}

// NewApp инициализирует прототип приложения
func NewApp(ctx context.Context, cfg *config.Config) (*app, ce.IError) {
	err := ce.ErrorApp
	log := logger.NewLogger(cfg)

	db, e := postgresql.NewDB(ctx, &cfg.Database, log)
	if e != nil {
		return nil, err.SetError(e)
	}

	sigChan := make(chan os.Signal, 1)
	doneChan := make(chan struct{})

	repoRS := rsrepository.NewRefreshStreamRepository(db, log)
	repoSS := ssrepository.NewStatusStreamRepository(db, log)
	repoRTSP := rtsprepository.NewRTSPRepository(cfg, log)

	return &app{
		cfg: cfg,
		db:  db,
		log: log,

		sigChan:  sigChan,
		doneChan: doneChan,

		refreshStreamRepo: repoRS,
		statusStreamRepo:  repoSS,
		rtspRepo:          repoRTSP,

		err: err,
	}, nil
}
