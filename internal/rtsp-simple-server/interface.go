package rtspsimpleserver

import (
	"context"

	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
)

type Common interface {
	// GetRtsp отправляет GET запрос на получение данных
	GetRtsp(ctx context.Context) (map[string]SConf, ce.IError)
	// PostAddRTSP отправляет POST запрос на добавление потока
	PostAddRTSP(ctx context.Context, camDB SConf) ce.IError
	// PostRemoveRTSP отправляет POST запрос на удаление потока
	PostRemoveRTSP(ctx context.Context, camRTSP SConf) ce.IError
	// PostEditRTSP отправляет POST запрос на изменение потока
	PostEditRTSP(ctx context.Context, cam SConf) ce.IError
}

type RTSPUseCase interface {
	Common
}

type RTSPRepository interface {
	Common
}
