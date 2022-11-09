package controller

import (
	"database/sql"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	"github.com/sirupsen/logrus"
)

type refreshStreamHandler struct {
	db      *sql.DB
	log     *logrus.Logger
	useCase refreshstream.RefreshStreamUseCase
}

func NewRefreshStreamHandler(useCase refreshstream.RefreshStreamUseCase, db *sql.DB, log *logrus.Logger) *refreshStreamHandler {
	return &refreshStreamHandler{
		db:      db,
		log:     log,
		useCase: useCase,
	}
}
