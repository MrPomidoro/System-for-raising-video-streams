package repository

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/statusstream"
	mocks "github.com/Kseniya-cha/System-for-raising-video-streams/internal/statusstream/repository/mock"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
	sqlMock "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/database/postgresql/mocks"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

func TestNewRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDB := sqlMock.NewMockIDB(ctrl)
	mockLog := zap.NewNop()

	repo := NewRepository(mockDB, &config.Database{}, mockLog)
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
	mockDB := sqlMock.NewMockIDB(ctrl)
	mockLog := zap.NewNop()
	repo := NewRepository(mockDB, &config.Database{}, mockLog)

	mockCommon := mocks.NewMockCommon(ctrl)
	repo.Common = mockCommon

	streamT := &statusstream.StatusStream{StreamId: 1, StatusResponse: true}
	streamF := &statusstream.StatusStream{StreamId: 0, StatusResponse: false}
	// Set up expectations
	mockCommon.EXPECT().Insert(ctx, streamT)
	mockCommon.EXPECT().Insert(ctx, streamF).Times(2)

	t.Run("TestInsertTrue", func(t *testing.T) {
		err := repo.Common.Insert(ctx, streamT)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("TestInsertFalse", func(t *testing.T) {
		err := repo.Common.Insert(ctx, streamF)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("TestInsertCtxCancel", func(t *testing.T) {
		cancel()
		err := repo.Common.Insert(ctx, streamF)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}
