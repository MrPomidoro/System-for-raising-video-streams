package statusstream

import (
	"context"

	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
)

type Common interface {
	Insert(ctx context.Context, rs *StatusStream) ce.IError
}

// Repository - интерфейс работы с запросами к базе данных
type Repository interface {
	Common
}
