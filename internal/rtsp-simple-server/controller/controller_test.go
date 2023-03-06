package controller

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	rtsp "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/logger"
)

func TestRTSP(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/v1/paths/list":
			json := []byte(`
			{
				"items": {
					"cam1": {
						"confName": "cam1",
						"conf": {
							"source": "rtsp://login:pass@1:123/cam1",
							"sourceProtocol": "automatic",
							"sourceAnyPortEnable": false,
							"sourceOnDemand": false,
							"sourceOnDemandStartTimeout": "10s",
							"sourceOnDemandCloseAfter": "10s",
							"runOnDemandStartTimeout": "10s",
							"runOnDemandCloseAfter": "10s",
							"runOnReady": "",
							"runOnReadyRestart": true,
							"runOnRead": "",
							"runOnReadRestart": false
						}
					}
				}
			}`)

			rw.Write(json)

		case "/v1/config/paths/add/1":
		case "/v1/config/paths/remove/1":
		case "/v1/config/paths/edit/1":
		default:
			t.Errorf("unexpected request URL: %s", req.URL.Path)
		}
	}))
	defer server.Close()

	cfg := &config.Config{
		Rtsp: config.Rtsp{
			Url: server.URL,
		},
		Logger: config.Logger{
			LogLevel:        "INFO",
			LogFileEnable:   true,
			LogStdoutEnable: true,
			RewriteLog:      true,
		},
	}
	cfg.Rtsp.Api.UrlGet = "/v1/paths/list"
	cfg.Rtsp.Api.UrlAdd = "/v1/config/paths/add/"
	cfg.Rtsp.Api.UrlRemove = "/v1/config/paths/remove/"
	cfg.Rtsp.Api.UrlEdit = "/v1/config/paths/edit/"

	log := logger.NewLogger(cfg)
	repo := NewRepository(cfg, log)
	sconf := rtsp.SConf{
		Stream: "1",
		Id:     1,
		Conf: rtsp.Conf{
			Source:         "source",
			SourceProtocol: "tcp",
		},
	}

	t.Run("TestGet", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		res, err := repo.GetRtsp(ctx)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(res) == 0 {
			t.Error("have not response")
		}

		cancel()
		_, err = repo.GetRtsp(ctx)
		if err != customError.ErrorRTSP.SetError(ctx.Err()) {
			t.Errorf("unexpected error: %v", err)
		}

		cfg.Rtsp.Api.UrlGet = "/v1/paths/listsss"
		repo = NewRepository(cfg, log)
		_, err = repo.GetRtsp(context.Background())
		if err != customError.ErrorRTSP.SetError(errors.New("unexpected end of JSON input")) {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("TestAdd", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		err := repo.PostAddRTSP(ctx, sconf)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		cancel()
		err = repo.PostAddRTSP(ctx, sconf)
		if err != customError.ErrorRTSP.SetError(ctx.Err()) {
			t.Errorf("unexpected error: %v", err)
		}

		// cfg.Rtsp.Api.UrlAdd = "/v1/config/adddd/"
		// repo = NewRepository(cfg, log)
		// err = repo.PostAddRTSP(context.Background(), sconf)
		// fmt.Println("err add", err)
		// if err != customError.ErrorRTSP.SetError(fmt.Errorf("unexpected request URL: %s", cfg.Rtsp.Api.UrlAdd)) {
		// 	t.Errorf("unexpected error: %v", err)
		// }
	})

	t.Run("TestEdit", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		err := repo.PostEditRTSP(ctx, sconf)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		cancel()
		err = repo.PostEditRTSP(ctx, sconf)
		if err != customError.ErrorRTSP.SetError(ctx.Err()) {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("TestRemove", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		err := repo.PostRemoveRTSP(ctx, sconf)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		cancel()
		err = repo.PostRemoveRTSP(ctx, sconf)
		if err != customError.ErrorRTSP.SetError(ctx.Err()) {
			t.Errorf("unexpected error: %v", err)
		}
	})

}
