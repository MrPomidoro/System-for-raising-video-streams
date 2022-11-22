package usecase

import (
	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	rtspsimpleserver "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server"
)

type rtspUseCase struct {
	repo rtspsimpleserver.RTSPRepository
}

func NewRTSPUseCase(repo rtspsimpleserver.RTSPRepository) *rtspUseCase {
	return &rtspUseCase{
		repo: repo,
	}
}

func (rtsp *rtspUseCase) GetRtsp() (map[string]interface{}, error) {
	return rtsp.repo.GetRtsp()
}

func (rtsp *rtspUseCase) PostAddRTSP(camDB refreshstream.RefreshStream) error {
	return rtsp.repo.PostAddRTSP(camDB)
}

func (rtsp *rtspUseCase) PostRemoveRTSP(camRTSP string) error {
	return rtsp.repo.PostRemoveRTSP(camRTSP)
}

func (rtsp *rtspUseCase) PostEditRTSP(camDB refreshstream.RefreshStream, conf rtspsimpleserver.Conf) error {
	return rtsp.repo.PostEditRTSP(camDB, conf)
}
