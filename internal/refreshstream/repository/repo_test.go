package repository

import (
	"context"
	"testing"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream/repository/mocks"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/database/postgresql"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"

	"github.com/pashagolub/pgxmock/v2"
)

func TestGet(t *testing.T) {

	// первая нерабочая хуйня
	//
	// mock, err := pgxmock.NewPool()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// defer mock.Close()
	// mock.ExpectBegin()
	// query := `SELECT *
	// FROM public."refresh_stream"
	// WHERE "stream" IS NOT null AND "stream_state" = true`
	// mock.ExpectQuery(query)
	// mock.ExpectQuery("select 1")

	// вторая нерабочая хуйня
	//
	// t.Parallel()
	// ctrl := gomock.NewController(t)
	// defer ctrl.Finish()
	// mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	// columns := []string{"id", "price"}
	// pgxRows := pgxpoolmock.NewRows(columns).AddRow(100, 100000.9).ToPgxRows()
	// mockPool.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxRows, nil)
	// repo := &repository{
	// }
	// NewRepository(mockPool, &zap.Logger{})

	// третья нерабочая хуйня
	//
	// mock, err := pgxmock.NewPool()
	// if err != nil {
	// 	t.Errorf("unexpected error: %v", err)
	// }
	// defer mock.Close()
	// mock.ExpectBegin()
	// mock.ExpectQuery(`SELECT *
	// 				FROM public."refresh_stream"
	// 				WHERE "stream" IS NOT null AND "stream_state" = true`).
	// 	WillReturnRows().
	// 	RowsWillBeClosed()
	// mock.ExpectExec(`UPDATE public."refresh_stream"
	// SET "stream_status"='true'
	// WHERE "stream"='1'`)
	// NewRepository(mock, &zap.Logger{})

	// четвёртая хуйня
	//
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Error("unexpected error:", err)
	}

	cols := []string{"id", "auth", "ip", "stream", "port_srv", "sp",
		"cam_id", "record_status", "stream_status", "record_state", "stream_state", "protocol"}

	// mock := Run(t)
	query := `SELECT *
		 		FROM public."refresh_stream"
		 		WHERE "stream" IS NOT null AND "stream_state" = true`

	// t.Run("а вдруг?..", func(t *testing.T) {
	mock.ExpectQuery(query).WillReturnRows(mock.NewRows(cols))
	// })

}

// func Run(t *testing.T) pgxmock.PgxConnIface {
// 	t.Helper()
// 	mock, err := pgxmock.NewConn()
// 	if err != nil {
// 		t.Error("unexpected error:", err)
// 	}
// 	defer mock.Close(context.Background())
// 	return mock
// }

func TestRepository_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.TODO()
	mockDB := postgresql.NewMockPgxIface(ctrl) // сделать мок на интерфейс бд
	mockLog := zap.NewNop()

	repo := NewRepository(postgresql.NewDB(ctx, mockDB, mockLog), mockLog)

	mockCommon := mocks.NewMockCommon(ctrl)
	repo.Common = mockCommon

	// Set up expectations
	mockCommon.EXPECT().Get(ctx, true).Return([]refreshstream.Stream{}, nil)

	// Call the method being tested
	_, err := repo.Get(ctx, true)

	// Check that the expectations were met
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}
