package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/logger"
	"github.com/sirupsen/logrus"
)

type refreshStreamRepository struct {
	db  *sql.DB
	log *logrus.Logger
}

func NewRefreshStreamRepository(db *sql.DB, log *logrus.Logger) *refreshStreamRepository {
	return &refreshStreamRepository{
		db:  db,
		log: log,
	}
}

func (s refreshStreamRepository) Get(ctx context.Context) ([]refreshstream.RefreshStream, error) {

	template := refreshstream.SELECT_COL_FROM_TBL
	chose := "*"
	tbl := `public."refresh_stream"`
	query := fmt.Sprintf(template, chose, tbl)

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		logger.LogError(s.log, fmt.Sprintf("cannot get: %v", err))
		return nil, err
	}
	defer rows.Close()

	// Слайс копий структур
	refreshStreamArr := []refreshstream.RefreshStream{}
	for rows.Next() {
		rs := refreshstream.RefreshStream{}
		err := rows.Scan(&rs.Id, &rs.Auth, &rs.Ip, &rs.Stream,
			&rs.Portsrv, &rs.Sp, &rs.Camid, &rs.Record_status,
			&rs.Stream_status, &rs.Record_state, &rs.Stream_state)
		if err != nil {
			logger.LogError(s.log, err)
			return nil, err
		}
		refreshStreamArr = append(refreshStreamArr, rs)
	}
	logger.LogDebug(s.log, "Received response from the database")
	return refreshStreamArr, nil
}
