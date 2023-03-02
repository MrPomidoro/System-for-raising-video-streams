package repository

import (
	"context"
	"fmt"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/statusstream"
	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/database/postgresql"
	"go.uber.org/zap"
)

type repository struct {
	db  *postgresql.DB
	log *zap.Logger
	err ce.IError
}

func NewRepository(db *postgresql.DB, log *zap.Logger) *repository {
	return &repository{
		db:  db,
		log: log,
		err: ce.ErrorStatusStream,
	}
}

// Insert отправляет запрос на добавление лога
func (s repository) Insert(ctx context.Context,
	ss *statusstream.StatusStream) ce.IError {

	query := fmt.Sprintf(statusstream.InsertToStatusStream, ss.StreamId, ss.StatusResponse)
	s.log.Debug("Query to database:\n\t" + query)

	_, err := s.db.Conn.Exec(ctx, query)
	if err != nil {
		return s.err.SetError(err)
	}

	return nil
}
