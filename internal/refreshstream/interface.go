package refreshstream

import (
	"context"
)

type RefreshStreamCommon interface {
	GetStatusTrue(ctx context.Context) ([]RefreshStream, error)
	GetStatusFalse(ctx context.Context) ([]RefreshStream, error)
	// GetId(ctx context.Context, id interface{}) (RefreshStream, error)
	// Insert(ctx context.Context, rs *RefreshStream) error
	// Update(ctx context.Context, rs *RefreshStream) error
	// Delete(ctx context.Context, id interface{}) error
}

type RefreshStreamUseCase interface {
	RefreshStreamCommon
}

type RefreshStreamRepository interface {
	RefreshStreamCommon
}
