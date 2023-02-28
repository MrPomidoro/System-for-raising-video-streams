package service

import (
	"context"
	"fmt"
	"os"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	rsrepository "github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream/repository"
	rtspsimpleserver "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server"
	rtsprepository "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server/repository"
	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/statusstream"
	ssrepository "github.com/Kseniya-cha/System-for-raising-video-streams/internal/statusstream/repository"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/database"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/logger"
	"go.uber.org/zap"
)

// app - прототип приложения
type app struct {
	cfg *config.Config
	log *zap.Logger
	db  *database.DB

	sigChan  chan os.Signal
	doneChan chan struct{}

	refreshStreamRepo refreshstream.RefreshStreamRepository
	statusStreamRepo  statusstream.Repository
	rtspRepo          rtspsimpleserver.RTSPRepository

	err ce.IError
}

// NewApp инициализирует прототип приложения
func NewApp(ctx context.Context, cfg *config.Config) (*app, ce.IError) {
	err := ce.ErrorStorage
	log := logger.NewLogger(cfg)

	if !cfg.DatabaseConnect {
		return nil, err.SetError(fmt.Errorf("no permission to connect to database"))
	}

	db, e := database.CreateDBConnection(ctx, cfg)
	if e != nil {
		return nil, err
	}

	sigChan := make(chan os.Signal, 1)
	doneChan := make(chan struct{})

	repoRS := rsrepository.NewRefreshStreamRepository(db.Db, log)
	repoSS := ssrepository.NewStatusStreamRepository(db.Db, log)
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
