package repository

import (
	"context"
	"database/sql"

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

func (s statusStreamRepository) Insert(ctx context.Context) error {

	return nil
}
