package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/statusstream"
	mocks "github.com/Kseniya-cha/System-for-raising-video-streams/internal/statusstream/repository/mock"
	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/database/postgresql"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

func TestNewRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDB := postgresql.NewMockPgxIface(ctrl)
	mockLog := zap.NewNop()

	repo := NewRepository(mockDB, mockLog)
	repoS := strings.Split(fmt.Sprint(repo), " ")
	testRepoS := strings.Split(fmt.Sprint(&Repository{db: mockDB, log: mockLog, err: ce.ErrorStatusStream}), " ")
	for i := range repoS {
		if repoS[i] != testRepoS[i] {
			t.Errorf("Unexpected Repository struct: %v, expect %v", testRepoS, repoS)
		}
	}
}

func TestInsert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	mockDB := postgresql.NewMockPgxIface(ctrl) // сделать мок на интерфейс бд
	mockLog := zap.NewNop()
	repo := NewRepository(mockDB, mockLog)

	mockCommon := mocks.NewMockCommon(ctrl)
	repo.Common = mockCommon

	streamT := &statusstream.StatusStream{StreamId: 1, StatusResponse: true}
	streamF := &statusstream.StatusStream{StreamId: 1, StatusResponse: false}
	// Set up expectations
	// mockCommon.EXPECT().Insert(ctx, streamT)
	// mockCommon.EXPECT().Insert(ctx, streamF)

	t.Run("TestInsertTrue", func(t *testing.T) {
		err := repo.Insert(ctx, streamT)
		if err != ce.ErrorStatusStream.SetError(errors.New("ERROR: insert or update on table \"status_stream\" violates foreign key constraint \"stream_id\" (SQLSTATE 23503)")) {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("TestInsertFalse", func(t *testing.T) {
		cancel()
		err := repo.Insert(ctx, streamF)
		if err != ce.ErrorStatusStream.SetError(ctx.Err()) {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}
