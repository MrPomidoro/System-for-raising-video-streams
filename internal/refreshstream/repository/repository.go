package repository

import (
	"context"
	"fmt"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/database/postgresql"
	"go.uber.org/zap"
)

type Repository struct {
	Common refreshstream.Common
	db     postgresql.IDB
	log    *zap.Logger
	cfg    *config.Database
	err    ce.IError
}

func NewRepository(db postgresql.IDB, cfg *config.Database, log *zap.Logger) *Repository {
	return &Repository{
		db:  db,
		log: log,
		cfg: cfg,
		err: ce.ErrorRefreshStream,
	}
}

// Get отправляет запрос на получение данных из таблицы
func (s Repository) Get(ctx context.Context, status bool) ([]refreshstream.Stream, ce.IError) {
	var query string
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	switch status {
	case true:
		query = fmt.Sprintf(refreshstream.QueryStateTrue, s.cfg.TableName)
	case false:
		query = fmt.Sprintf(refreshstream.QueryStateFalse, s.cfg.TableName)
	}

	if ctx.Err() != nil {
		return nil, s.err.SetError(ctx.Err())
	}

	s.log.Debug("Query to database:\n\t" + query)

	rows, err := s.db.GetConn().Query(ctx, query)
	if err != nil {
		return nil, s.err.SetError(err)
	}
	defer rows.Close()

	// Слайс копий структур
	var res []refreshstream.Stream

	// pgxscan.Select(ctx, s.db, &res, query)
	for rows.Next() {
		rs := refreshstream.Stream{}
		err = rows.Scan(&rs.Id, &rs.Login, &rs.Pass, &rs.Ip,
			&rs.CamPath, &rs.CodeMp, &rs.StatePublic, &rs.StatusPublic)
		if err != nil {
			return nil, s.err.SetError(err)
		}
		rs.Port = "554"
		res = append(res, rs)
	}

	return res, nil
}
