package controller

import (
	"database/sql"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/statusstream"
	"github.com/sirupsen/logrus"
)

type StatusStreamHandler struct {
	db      *sql.DB
	log     *logrus.Logger
	useCase statusstream.StatusStreamUseCase
}

func NewStatusStreamHandler(useCase statusstream.StatusStreamUseCase, db *sql.DB, log *logrus.Logger) *StatusStreamHandler {
	return &StatusStreamHandler{
		db:      db,
		log:     log,
		useCase: useCase,
	}
}
