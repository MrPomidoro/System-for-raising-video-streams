package statusstream

import (
	"context"
)

type StatusStreamCommon interface {
	Get(ctx context.Context) ([]StatusStream, error)
	// GetId(ctx context.Context, id interface{}) (StatusStream, error)
	// Insert(ctx context.Context, rs *StatusStream) error
	// Update(ctx context.Context, rs *StatusStream) error
	// Delete(ctx context.Context, id interface{}) error
}

type StatusStreamUseCase interface {
	StatusStreamCommon
}

type StatusStreamRepository interface {
	StatusStreamCommon
}
