package repository

import (
	"context"
	"database/sql"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
	"go.uber.org/zap"
)

type refreshStreamRepository struct {
	db  *sql.DB
	log *zap.Logger
	err ce.IError
}

func NewRefreshStreamRepository(db *sql.DB, log *zap.Logger) *refreshStreamRepository {
	return &refreshStreamRepository{
		db:  db,
		log: log,
		err: ce.ErrorRefreshStream,
	}
}

// Get отправляет запрос на получение данных из таблицы
func (s refreshStreamRepository) Get(ctx context.Context, status bool) ([]refreshstream.RefreshStream, ce.IError) {
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
	res := []refreshstream.RefreshStream{}
	for rows.Next() {
		rs := refreshstream.RefreshStream{}
		err := rows.Scan(&rs.Id, &rs.Auth, &rs.Ip, &rs.Stream,
			&rs.Portsrv, &rs.Sp, &rs.CamId, &rs.Record_status,
			&rs.Stream_status, &rs.Record_state, &rs.Stream_state, &rs.Protocol)
		if err != nil {
			return nil, s.err.SetError(err)
		}
		res = append(res, rs)
	}
	return res, nil
}

// Update отправляет запрос на изменение поля stream_status
// func (s refreshStreamRepository) Update(ctx context.Context, stream string) ce.IError {

// 	query := fmt.Sprintf(refreshstream.QueryEditStatus, stream)
// 	s.log.Debug("Query to database:\n\t" + query)

// 	_, err := s.db.ExecContext(ctx, query)
// 	if err != nil {
// 		return s.err.SetError(err)
// 	}

// 	return nil
// }
