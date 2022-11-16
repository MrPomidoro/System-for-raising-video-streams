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

func (s refreshStreamRepository) GetStatusTrue(ctx context.Context) ([]refreshstream.RefreshStream, error) {
	// Выполнение запроса
	rows, err := s.db.QueryContext(ctx, refreshstream.QUERY_STATUS_TRUE)
	if err != nil {
		logger.LogError(s.log, fmt.Sprintf("cannot complete Get request: %v", err))
		return nil, err
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
			logger.LogError(s.log, err)
			return nil, err
		}
		refreshStreamArr = append(refreshStreamArr, rs)
	}
	logger.LogInfo(s.log, "Received response from the database")
	return refreshStreamArr, nil
}

func (s refreshStreamRepository) GetStatusFalse(ctx context.Context) ([]refreshstream.RefreshStream, error) {
	// Выполнение запроса
	rows, err := s.db.QueryContext(ctx, refreshstream.QUERY_STATUS_FALSE)
	if err != nil {
		logger.LogError(s.log, fmt.Sprintf("cannot complete Get request: %v", err))
		return nil, err
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
			logger.LogError(s.log, err)
			return nil, err
		}
		refreshStreamArr = append(refreshStreamArr, rs)
	}
	logger.LogInfo(s.log, "Received response from the database")
	return refreshStreamArr, nil
}
