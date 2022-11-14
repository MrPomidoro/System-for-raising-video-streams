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
	db     *sql.DB
	logStC *logrus.Logger
}

func NewRefreshStreamRepository(db *sql.DB, logStC *logrus.Logger) *refreshStreamRepository {
	return &refreshStreamRepository{
		db:     db,
		logStC: logStC,
	}
}

func (s refreshStreamRepository) Get(ctx context.Context) ([]refreshstream.RefreshStream, error) {

	/*
		template := refreshstream.SELECT_COL_FROM_TBL
		chose := "*"
		tbl := `public."refresh_stream"`
		query := fmt.Sprintf(template, chose, tbl)
	*/

	// Выполнение запроса
	rows, err := s.db.QueryContext(ctx, refreshstream.QUERY)
	if err != nil {
		logger.LogErrorStatusCode(s.logStC, fmt.Sprintf("cannot get: %v", err), "Get", "400")
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
			logger.LogErrorStatusCode(s.logStC, err, "Get", "400")
			return nil, err
		}
		refreshStreamArr = append(refreshStreamArr, rs)
	}
	logger.LogInfoStatusCode(s.logStC, "Received response from the database", "Get", "200")
	return refreshStreamArr, nil
}
