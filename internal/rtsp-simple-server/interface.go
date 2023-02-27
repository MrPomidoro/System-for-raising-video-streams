package rtspsimpleserver

import (
	"context"

	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
)

type RTSPCommon interface {
	// GetRtsp отправляет GET запрос на получение данных
	GetRtsp(ctx context.Context, dataRTSPchan chan SConf) ([]SConf, ce.IError)
	// PostAddRTSP отправляет POST запрос на добавление потока
	PostAddRTSP(camDB SConf) ce.IError
	// PostRemoveRTSP отправляет POST запрос на удаление потока
	PostRemoveRTSP(camRTSP SConf) ce.IError
	// PostEditRTSP отправляет POST запрос на изменение потока
	PostEditRTSP(cam SConf) ce.IError
}

type RTSPUseCase interface {
	RTSPCommon
}

type RTSPRepository interface {
	RTSPCommon
}
