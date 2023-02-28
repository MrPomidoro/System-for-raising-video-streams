package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/statusstream"
	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
	"go.uber.org/zap"
)

type statusStreamRepository struct {
	db  *sql.DB
	log *zap.Logger
	err ce.IError
}

func NewStatusStreamRepository(db *sql.DB, log *zap.Logger) *statusStreamRepository {
	return &statusStreamRepository{
		db:  db,
		log: log,
		err: ce.ErrorStatusStream,
	}
}

// Insert отправляет запрос на добавление лога
func (s statusStreamRepository) Insert(ctx context.Context,
	ss *statusstream.StatusStream) ce.IError {

	query := fmt.Sprintf(statusstream.InsertToStatusStream, ss.StreamId, ss.StatusResponse)
	s.log.Debug("Query to database:\n\t" + query)

	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return s.err.SetError(err)
	}

	return nil
}
