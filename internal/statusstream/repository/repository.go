package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/statusstream"
)

type statusStreamRepository struct {
	db *sql.DB
}

func NewStatusStreamRepository(db *sql.DB) *statusStreamRepository {
	return &statusStreamRepository{
		db: db,
	}
}

func (s statusStreamRepository) Insert(ctx context.Context,
	ss *statusstream.StatusStream) error {

	query := fmt.Sprintf(statusstream.InsertToStatusStream, ss.StreamId, ss.StatusResponse)

	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("cannot insert: %v", err)
	}

	return nil
}
