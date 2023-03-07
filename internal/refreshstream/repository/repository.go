package repository

import (
	"context"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/database/postgresql"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx"
	"go.uber.org/zap"
)

type PgxIface interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	Ping(context.Context) error
	Close()
}

// Database is wrapper for PgxIface
type Database struct {
	DB PgxIface
}

// NewSelector is an initializer for Selector
func NewDatabase(ds PgxIface) Database {
	return Database{DB: ds}
}

//
//
//

type repository struct {
	db  *postgresql.DB
	log *zap.Logger
	err ce.IError
}

func NewRepository(db *postgresql.DB, log *zap.Logger) *repository {
	return &repository{
		db:  db,
		log: log,
		err: ce.ErrorRefreshStream,
	}
}

// Get отправляет запрос на получение данных из таблицы
func (s repository) Get(ctx context.Context, status bool) ([]refreshstream.Stream, ce.IError) {
	var query string
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	switch status {
	case true:
		query = refreshstream.QueryStateTrue
	case false:
		query = refreshstream.QueryStateFalse
	}

	if ctx.Err() != nil {
		return nil, s.err.SetError(ctx.Err())
	}

	s.log.Debug("Query to database:\n\t" + query)

	rows, err := s.db.Conn.Query(ctx, query)
	if err != nil {
		return nil, s.err.SetError(err)
	}
	defer rows.Close()

	// Слайс копий структур
	res := []refreshstream.Stream{}
	for rows.Next() {
		rs := refreshstream.Stream{}
		err = rows.Scan(&rs.Id, &rs.Auth, &rs.Ip, &rs.Stream,
			&rs.Portsrv, &rs.Sp, &rs.CamId, &rs.RecordStatus,
			&rs.StreamStatus, &rs.RecordState, &rs.StreamState, &rs.Protocol)
		if err != nil {
			return nil, s.err.SetError(err)
		}
		res = append(res, rs)
	}

	return res, nil
}
