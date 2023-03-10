package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
			jsonStruct := map[string]map[string]map[string]map[string]string{
				"items": {
					"cam1": {
						"conf": {
							"source":         "rtsp://login:pass@1:123/cam1",
							"sourceProtocol": "udp",
							"runOnReady":     "",
						},
					},
					"cam2": {
						"conf": {
							"source":         "rtsp://login:pass@1:321/cam2",
							"sourceProtocol": "tcp",
							"runOnReady":     "",
						},
					},
				},
			}

			jsonString, err := json.Marshal(jsonStruct)
			if err != nil {
				rw.WriteHeader(500)
			}
			rw.Write(jsonString)

		case "/v1/config/paths/add/1":
		case "/v1/config/paths/remove/1":
		case "/v1/config/paths/edit/1":
		default:
			rw.WriteHeader(404)
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

	ctx, cancel := context.WithCancel(context.Background())
	t.Run("TestGetOK", func(t *testing.T) {

		res, err := repo.GetRtsp(ctx)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(res) == 0 {
			t.Error("have not response")
		}
	})
	t.Run("TestGetInvalidURL", func(t *testing.T) {
		cfg.Rtsp.Api.UrlGet = "/v/lists"
		repo = NewRepository(cfg, log)
		_, err := repo.GetRtsp(context.Background())
		if err != customError.ErrorRTSP.SetError(errors.New("unexpected status code: 404")) {
			t.Errorf("unexpected error: %v", err)
		}
	})
	cancel()
	t.Run("TestGetCtxCancel", func(t *testing.T) {
		_, err := repo.GetRtsp(ctx)
		if err != customError.ErrorRTSP.SetError(ctx.Err()) {
			t.Errorf("unexpected error: %v", err)
		}
	})

	ctx, cancel = context.WithCancel(context.Background())
	t.Run("TestAddOK", func(t *testing.T) {
		err := repo.PostAddRTSP(ctx, sconf)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	t.Run("TestAddInvalidURL", func(t *testing.T) {
		cfg.Rtsp.Api.UrlAdd = "/v1/config/adddd/"
		repo = NewRepository(cfg, log)
		err := repo.PostAddRTSP(context.Background(), sconf)
		if err != customError.ErrorRTSP.SetError(fmt.Errorf("unexpected status code: 404")) {
			t.Errorf("unexpected error: %v", err)
		}
	})
	cancel()
	t.Run("TestAddCtxCancel", func(t *testing.T) {
		err := repo.PostAddRTSP(ctx, sconf)
		if err != customError.ErrorRTSP.SetError(ctx.Err()) {
			t.Errorf("unexpected error: %v", err)
		}
	})

	ctx, cancel = context.WithCancel(context.Background())
	t.Run("TestEditOK", func(t *testing.T) {
		err := repo.PostEditRTSP(ctx, sconf)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	t.Run("TestEditInvalidURL", func(t *testing.T) {
		cfg.Rtsp.Api.UrlEdit = "/v1/config/ediiiit/"
		repo = NewRepository(cfg, log)
		err := repo.PostEditRTSP(context.Background(), sconf)
		if err != customError.ErrorRTSP.SetError(fmt.Errorf("unexpected status code: 404")) {
			t.Errorf("unexpected error: %v", err)
		}
	})
	cancel()
	t.Run("TestEditCtxCancel", func(t *testing.T) {
		err := repo.PostEditRTSP(ctx, sconf)
		if err != customError.ErrorRTSP.SetError(ctx.Err()) {
			t.Errorf("unexpected error: %v", err)
		}
	})

	ctx, cancel = context.WithCancel(context.Background())
	t.Run("TestRemoveOK", func(t *testing.T) {
		err := repo.PostRemoveRTSP(ctx, sconf)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	t.Run("TestRemoveInvalidURL", func(t *testing.T) {
		cfg.Rtsp.Api.UrlRemove = "/v1/config/rem/"
		repo = NewRepository(cfg, log)
		err := repo.PostRemoveRTSP(context.Background(), sconf)
		if err != customError.ErrorRTSP.SetError(fmt.Errorf("unexpected status code: 404")) {
			t.Errorf("unexpected error: %v", err)
		}
	})
	cancel()
	t.Run("TestRemoveCtxCancel", func(t *testing.T) {
		err := repo.PostRemoveRTSP(ctx, sconf)
		if err != customError.ErrorRTSP.SetError(ctx.Err()) {
			t.Errorf("unexpected error: %v", err)
		}
	})

}
