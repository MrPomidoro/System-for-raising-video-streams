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
	err *ce.Error
}

func NewStatusStreamRepository(db *sql.DB, log *zap.Logger) *statusStreamRepository {
	return &statusStreamRepository{
		db:  db,
		log: log,
		err: ce.NewError(ce.ErrorLevel, "50.4.4", "status stream entity error at database operation level"),
	}
}

// Insert отправляет запрос на добавление лога
func (s statusStreamRepository) Insert(ctx context.Context,
	ss *statusstream.StatusStream) *ce.Error {

	query := fmt.Sprintf(statusstream.InsertToStatusStream, ss.StreamId, ss.StatusResponse)
	s.log.Debug(query)

	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return s.err.SetError(err)
	}

	return nil
}
