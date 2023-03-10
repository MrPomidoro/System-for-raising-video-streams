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
	rtspcontr "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server/controller"
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

	getDB(ctx context.Context, mu *sync.Mutex) ([]refreshstream.Stream, ce.IError)
	getRTSP(ctx context.Context) (map[string]rtsp.SConf, ce.IError)
	addData(ctx context.Context, camsAdd map[string]rtsp.SConf) ce.IError
	getCamsEdit(dataDB []refreshstream.Stream, dataRTSP map[string]rtsp.SConf,
		camsAdd map[string]rtsp.SConf, camsRemove map[string]rtsp.SConf) map[string]rtsp.SConf
	addRemoveData(ctx context.Context, dataDB []refreshstream.Stream,
		dataRTSP map[string]rtsp.SConf, camsAdd map[string]rtsp.SConf,
		camsRemove map[string]rtsp.SConf) ce.IError
	getCamsAdd(dataDB []refreshstream.Stream,
		dataRTSP map[string]rtsp.SConf) map[string]rtsp.SConf
	getCamsRemove(dataDB []refreshstream.Stream,
		dataRTSP map[string]rtsp.SConf)
}

// app - прототип приложения
type app struct {
	cfg *config.Config
	log *zap.Logger
	db  postgresql.IDB

	sigChan chan os.Signal

	refreshStreamRepo refreshstream.Repository
	statusStreamRepo  statusstream.Repository
	rtspRepo          rtsp.Repository

	err ce.IError
}

// NewApp инициализирует прототип приложения
func NewApp(ctx context.Context, cfg *config.Config) (*app, ce.IError) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	log := logger.NewLogger(cfg)

	db, err := postgresql.NewDB(ctx, cfg.Database, log)
	if err != nil {
		return nil, err
	}

	sigChan := make(chan os.Signal, 1)

	return &app{
		cfg: cfg,
		db:  db,
		log: log,

		sigChan: sigChan,

		refreshStreamRepo: rsrepo.NewRepository(db, log),
		statusStreamRepo:  ssrepo.NewRepository(db, log),
		rtspRepo:          rtspcontr.NewRepository(cfg, log),

		err: ce.ErrorApp,
	}, nil
}

type appIn interface {
	getDB(ctx context.Context, mu *sync.Mutex) ([]refreshstream.Stream, ce.IError)
	getRTSP(ctx context.Context) (map[string]rtsp.SConf, ce.IError)
}
