package refreshstream

import (
	"context"
)

type Common interface {
	Get(ctx context.Context, status bool) ([]RefreshStream, error)
}

type UseCase interface {
	Common
}

type Repository interface {
	Common
}
