package rtspsimpleserver

import (
	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
)

type RTSPCommon interface {
	GetRtsp() (map[string]interface{}, *ce.Error)
	PostAddRTSP(camDB refreshstream.RefreshStream) *ce.Error
	PostRemoveRTSP(camRTSP string) *ce.Error
	PostEditRTSP(camDB refreshstream.RefreshStream, conf Conf) *ce.Error
}

type RTSPUseCase interface {
	RTSPCommon
}

type RTSPRepository interface {
	RTSPCommon
}
