package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
)

type refreshStreamRepository struct {
	db *sql.DB
}

func NewRefreshStreamRepository(db *sql.DB) *refreshStreamRepository {
	return &refreshStreamRepository{
		db: db,
	}
}

func (s refreshStreamRepository) Get(ctx context.Context, status bool) ([]refreshstream.RefreshStream, error) {
	var query string

	switch status {
	case true:
		query = refreshstream.QueryStatusTrue
	case false:
		query = refreshstream.QueryStatusFalse
	}

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("cannot complete Get request: %v", err)
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
			return nil, err
		}
		refreshStreamArr = append(refreshStreamArr, rs)
	}
	return refreshStreamArr, nil
}
