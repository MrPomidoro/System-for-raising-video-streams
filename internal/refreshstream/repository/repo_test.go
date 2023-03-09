package repository

import (
	"context"
	"testing"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream/repository/mocks"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/database/postgresql"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

func TestRepository_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	mockDB := postgresql.NewMockPgxIface(ctrl) // сделать мок на интерфейс бд
	mockLog := zap.NewNop()
	// testdb, err := postgresql.NewDB(ctx, config.Database{}, mockLog)
	// if err != nil {
	// 	t.Fatalf("Unexpected error: %v", err)
	// }
	repo := NewRepository(mockDB, mockLog)

	mockCommon := mocks.NewMockCommon(ctrl)
	repo.Common = mockCommon

	// Set up expectations
	mockCommon.EXPECT().Get(ctx, true).Return([]refreshstream.Stream{}, nil)

	// Call the method being tested
	_, err := repo.Common.Get(ctx, true)

	// Check that the expectations were met
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
