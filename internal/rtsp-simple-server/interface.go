package rtspsimpleserver

import (
	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
)

type RTSPCommon interface {
	GetRtsp() ([]SConf, ce.IError)
	PostAddRTSP(camDB refreshstream.RefreshStream) ce.IError
	PostRemoveRTSP(camRTSP string) ce.IError
	PostEditRTSP(camDB refreshstream.RefreshStream, conf Conf) ce.IError
}

type RTSPUseCase interface {
	RTSPCommon
}

type RTSPRepository interface {
	RTSPCommon
}
