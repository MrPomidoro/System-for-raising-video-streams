package rtspsimpleserver

import (
	"context"

	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
)

type Common interface {
	GetRtsp(ctx context.Context) (map[string]SConf, ce.IError)
	PostAddRTSP(ctx context.Context, camDB SConf) ce.IError
	PostRemoveRTSP(ctx context.Context, camRTSP SConf) ce.IError
	PostEditRTSP(ctx context.Context, cam SConf) ce.IError
}

// Repository - интерфейс работы с запросами к базе данных
type Repository interface {
	Common
}
