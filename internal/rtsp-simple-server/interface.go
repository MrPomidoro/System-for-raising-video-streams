package rtspsimpleserver

import "github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"

type RTSPCommon interface {
	GetRtsp() (map[string]interface{}, error)
	PostAddRTSP(camDB refreshstream.RefreshStream) error
	PostRemoveRTSP(camRTSP string) error
	PostEditRTSP(camDB refreshstream.RefreshStream, conf Conf) error
}

type RTSPUseCase interface {
	RTSPCommon
}

type RTSPRepository interface {
	RTSPCommon
}
