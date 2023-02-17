package service

import (
	"context"
	"os"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	rsrepository "github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream/repository"
	rtspsimpleserver "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server"
	rtsprepository "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server/repository"
	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/statusstream"
	ssrepository "github.com/Kseniya-cha/System-for-raising-video-streams/internal/statusstream/repository"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/database"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/logger"
	"go.uber.org/zap"
)

// app - прототип приложения
type app struct {
	cfg               *config.Config
	log               *zap.Logger
	db                *database.DB
	sigChan           chan os.Signal
	refreshStreamRepo refreshstream.RefreshStreamRepository
	statusStreamRepo  statusstream.StatusStreamRepository
	rtspRepo          rtspsimpleserver.RTSPRepository
}

// NewApp инициализирует прототип приложения
func NewApp(ctx context.Context, cfg *config.Config) *app {
	log := logger.NewLogger(cfg)

	if !cfg.DatabaseConnect {
		log.Error("no permission to connect to database")
		return &app{}
	}

	db := database.CreateDBConnection(cfg)
	sigChan := make(chan os.Signal, 1)
	repoRS := rsrepository.NewRefreshStreamRepository(db.Db)
	repoSS := ssrepository.NewStatusStreamRepository(db.Db)
	repoRTSP := rtsprepository.NewRTSPRepository(cfg, log)

	return &app{
		cfg:               cfg,
		db:                db,
		log:               log,
		sigChan:           sigChan,
		refreshStreamRepo: repoRS,
		statusStreamRepo:  repoSS,
		rtspRepo:          repoRTSP,
	}
}
