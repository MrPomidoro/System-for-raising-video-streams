package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	rtsp "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server"
)

func TestGetRtsp(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/v1/paths/list" {
			t.Errorf("unexpected request URL: %s", req.URL.Path)
		}

		json := []byte(`"items": {
			"cam1": {
				"confName": "cam1",
				"conf": {
					"source": "rtsp://login:pass@1:123/cam1",
					"sourceProtocol": "automatic",
					"sourceAnyPortEnable": false,
					"sourceOnDemand": false,
					"sourceOnDemandStartTimeout": "10s",
					"sourceOnDemandCloseAfter": "10s",
					"runOnDemand": "",
					"runOnDemandRestart": false,
					"runOnDemandStartTimeout": "10s",
					"runOnDemandCloseAfter": "10s",
					"runOnReady": "",
					"runOnReadyRestart": true,
					"runOnRead": "",
					"runOnReadRestart": false
				},
			},"
		}`)

		rw.Write(json)
	}))

	defer server.Close()

	tests := []struct {
		name string
		path string
	}{
		{
			name: "Test for correct",
			path: "/v1/paths/list",
		},
		{
			name: "Test got wrong request",
			path: "/v1/path/list",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			client := server.Client()

			res, err := client.Get(server.URL + tt.path)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				t.Errorf("unexpected status code: %d", res.StatusCode)
			}

			body, _ := io.ReadAll(res.Body)

			var item map[string]map[string]rtsp.SConf
			json.Unmarshal(body, &item)

			streams := make(map[string]rtsp.SConf)
			for _, ress := range item {
				for stream, i := range ress {
					cam := rtsp.SConf{}
					cam.Stream = stream
					cam.Conf = i.Conf

					streams[stream] = cam
				}
			}

			for stream, info := range streams {
				fmt.Println(stream, info)

				if stream != "cam1" && info.Stream != "cam1" {
					t.Errorf("unexpected stream: %s", stream)
				}

				if info.Conf.SourceProtocol != "tcp" {
					t.Errorf("unexpected protocol: %s", info.Conf.SourceProtocol)
				}

				if info.Conf.Source != "source" {
					t.Errorf("unexpected source: %s", info.Conf.Source)
				}
			}
		})
	}
}
