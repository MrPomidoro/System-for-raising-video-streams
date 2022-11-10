package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/statusstream"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/logger"
	"github.com/sirupsen/logrus"
)

type statusStreamRepository struct {
	db  *sql.DB
	log *logrus.Logger
}

func NewStatusStreamRepository(db *sql.DB, log *logrus.Logger) *statusStreamRepository {
	return &statusStreamRepository{
		db:  db,
		log: log,
	}
}

func (s statusStreamRepository) Insert(ctx context.Context,
	ss *statusstream.StatusStream) error {

	query := fmt.Sprintf(statusstream.InsertToStatusStream, ss.StreamId, ss.StatusResponse)

	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("cannot insert: %v", err)
	}

	logger.LogDebug(s.log, "Success insert")
	return nil
}
