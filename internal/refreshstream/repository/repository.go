package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
)

type refreshStreamRepository struct {
	db  *sql.DB
	err *ce.Error
}

func NewRefreshStreamRepository(db *sql.DB) *refreshStreamRepository {
	return &refreshStreamRepository{
		db:  db,
		err: ce.NewError(ce.ErrorLevel, "50.4.2", "refresh stream entity error at database operation level"),
	}
}

// Get отправляет запрос на получение данных из таблицы
func (s refreshStreamRepository) Get(ctx context.Context, status bool) ([]refreshstream.RefreshStream, *ce.Error) {
	var query string

	switch status {
	case true:
		query = refreshstream.QueryStateTrue
	case false:
		query = refreshstream.QueryStateFalse
	}

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, s.err.SetError(err)
	}
	defer rows.Close()

	// Слайс копий структур
	refreshStreamArr := []refreshstream.RefreshStream{}
	for rows.Next() {
		rs := refreshstream.RefreshStream{}
		err := rows.Scan(&rs.Id, &rs.Auth, &rs.Ip, &rs.Stream,
			&rs.Portsrv, &rs.Sp, &rs.CamId, &rs.Record_status,
			&rs.Stream_status, &rs.Record_state, &rs.Stream_state, &rs.Protocol)
		if err != nil {
			return nil, s.err.SetError(err)
		}
		refreshStreamArr = append(refreshStreamArr, rs)
	}
	return refreshStreamArr, nil
}

// Update отправляет запрос на изменение поля stream_status
func (s refreshStreamRepository) Update(ctx context.Context, stream string) *ce.Error {

	query := fmt.Sprintf(refreshstream.QueryEditStatus, stream)

	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return s.err.SetError(err)
	}

	return nil
}
