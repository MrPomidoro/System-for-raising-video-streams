package refreshstream

import (
	"context"
)

type RefreshStreamCommon interface {
	Get(ctx context.Context, status bool) ([]RefreshStream, error)
	Update(ctx context.Context, stream string) error
}

type RefreshStreamRepository interface {
	RefreshStreamCommon
}
