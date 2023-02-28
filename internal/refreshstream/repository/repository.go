package repository

import (
	"context"
	"database/sql"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
	"go.uber.org/zap"
)

type repository struct {
	db  *sql.DB
	log *zap.Logger
	err ce.IError
}

func NewRefreshStreamRepository(db *sql.DB, log *zap.Logger) *repository {
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
	s.log.Debug("Query to database:\n\t" + query)

	rows, err := s.db.QueryContext(ctx, query)
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

// Update отправляет запрос на изменение поля stream_status
// func (s repository) Update(ctx context.Context, stream string) ce.IError {

// 	query := fmt.Sprintf(refreshstream.QueryEditStatus, stream)
// 	s.log.Debug("Query to database:\n\t" + query)

// 	_, err := s.db.ExecContext(ctx, query)
// 	if err != nil {
// 		return s.err.SetError(err)
// 	}

// 	return nil
// }
