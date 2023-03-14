package repository

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	mocks "github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream/repository/mock"
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
	testRepoS := strings.Split(fmt.Sprint(&Repository{db: mockDB, log: mockLog, err: ce.ErrorRefreshStream}), " ")
	for i := range repoS {
		if repoS[i] != testRepoS[i] {
			t.Errorf("Unexpected Repository struct: %v, expect %v", testRepoS, repoS)
		}
	}
}

func TestRepository_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	mockDB := sqlMock.NewMockIDB(ctrl)
	mockLog := zap.NewNop()
	repo := NewRepository(mockDB, &config.Database{}, mockLog)

	mockCommon := mocks.NewMockCommon(ctrl)
	repo.Common = mockCommon

	// Set up expectations
	mockCommon.EXPECT().Get(ctx, true).Return([]refreshstream.Stream{}, nil)
	mockCommon.EXPECT().Get(ctx, false).Return([]refreshstream.Stream{}, nil)

	expectT := []refreshstream.Stream{{}}

	t.Run("TestGetTrue", func(t *testing.T) {
		// Call the method being tested
		res, err := repo.Common.Get(ctx, true)
		// Check that the expectations were met
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !reflect.DeepEqual(res, expectT) {
			t.Errorf("unexpected result from query Get: %v", res)
		}
	})

	t.Run("TestGetCtxCancel", func(t *testing.T) {
		cancel()
		_, err := repo.Common.Get(ctx, false)
		// Check that the expectations were met
		if err != ce.ErrorRefreshStream.SetError(ctx.Err()) {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}
