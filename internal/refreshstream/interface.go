package refreshstream

import (
	"context"
)

type RefreshStreamCommon interface {
	Get(ctx context.Context, status bool) ([]RefreshStream, error)
	Update(ctx context.Context, stream string) error
}

type RefreshStreamUseCase interface {
	RefreshStreamCommon
}

type RefreshStreamRepository interface {
	RefreshStreamCommon
}
