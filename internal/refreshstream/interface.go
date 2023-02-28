package refreshstream

import (
	"context"

	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
)

type RefreshStreamCommon interface {
	// Get отправляет запрос на получение данных из таблицы
	Get(ctx context.Context, status bool) ([]RefreshStream, ce.IError)
	// Update отправляет запрос на изменение поля stream_status
	// Update(ctx context.Context, stream string) ce.IError
}

type RefreshStreamRepository interface {
	RefreshStreamCommon
}
