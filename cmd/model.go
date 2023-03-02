package service

//go:generate mockgen -destination=../mocks/mock_app.go -package=mocks github.com/Kseniya-cha/System-for-raising-video-streams/cmd App

import (
	"context"
	"os"
	"sync"

	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/database/postgresql"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	rsrepo "github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream/repository"
	rtsp "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server"
	rtsprepo "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server/repository"
	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/statusstream"
	ssrepo "github.com/Kseniya-cha/System-for-raising-video-streams/internal/statusstream/repository"
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
	db  postgresql.IDB

	sigChan  chan os.Signal
	doneChan chan struct{}

	refreshStreamRepo refreshstream.Repository
	statusStreamRepo  statusstream.Repository
	rtspRepo          rtsp.Repository

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
	idb := postgresql.IDB(db)

	sigChan := make(chan os.Signal, 1)
	doneChan := make(chan struct{})

	return &app{
		cfg: cfg,
		db:  idb,
		log: log,

		sigChan:  sigChan,
		doneChan: doneChan,

		refreshStreamRepo: rsrepo.NewRepository(db, log),
		statusStreamRepo:  ssrepo.NewRepository(db, log),
		rtspRepo:          rtsprepo.NewRepository(cfg, log),

		err: err,
	}, nil
}

type appIn interface {
	// addData(ctx context.Context, camsAdd map[string]rtsp.SConf) ce.IError
	// removeData(ctx context.Context, dataRTSP map[string]rtsp.SConf) ce.IError
	// editData(ctx context.Context, camsEdit map[string]rtsp.SConf) ce.IError

	//
	getDB(ctx context.Context, mu *sync.Mutex) ([]refreshstream.Stream, ce.IError)
	getRTSP(ctx context.Context) (map[string]rtsp.SConf, ce.IError)
}
