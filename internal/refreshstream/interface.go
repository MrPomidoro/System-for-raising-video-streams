package refreshstream

import (
	"context"

	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
)

// Common ...
type Common interface {
	Get(ctx context.Context, status bool) ([]Stream, ce.IError)
}

// Repository - интерфейс работы с запросами к базе данных
type Repository interface {
	Common
}
