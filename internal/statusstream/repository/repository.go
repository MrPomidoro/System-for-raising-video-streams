package repository

import (
	"context"
	"fmt"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/statusstream"
	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/database/postgresql"
	"go.uber.org/zap"
)

type Repository struct {
	Common statusstream.Common
	db     postgresql.IDB
	log    *zap.Logger
	err    ce.IError
}

func NewRepository(db postgresql.IDB, log *zap.Logger) *Repository {
	return &Repository{
		db:  db,
		log: log,
		err: ce.ErrorStatusStream,
	}
}

// Insert отправляет запрос на добавление лога
func (s Repository) Insert(ctx context.Context,
	ss *statusstream.StatusStream) ce.IError {

	if ss.StreamId == 0 {
		ss.StreamId = 1
	}

	query := fmt.Sprintf(statusstream.InsertToStatusStream, ss.StreamId, ss.StatusResponse)
	s.log.Debug("Query to database:\n\t" + query)

	_, err := s.db.GetConn().Exec(ctx, query)
	if err != nil {
		return s.err.SetError(err)
	}

	return nil
}
