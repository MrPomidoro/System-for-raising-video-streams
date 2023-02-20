package refreshstream

import (
	"context"

	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
)

type RefreshStreamCommon interface {
	Get(ctx context.Context, status bool) ([]RefreshStream, ce.IError)
	Update(ctx context.Context, stream string) ce.IError
}

type RefreshStreamRepository interface {
	RefreshStreamCommon
}
