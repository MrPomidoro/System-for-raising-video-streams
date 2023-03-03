package repository

import (
	"encoding/json"
	"fmt"
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

		stream := map[string]map[string]rtsp.SConf{
			"items": {
				"cam1": rtsp.SConf{Stream: "cam1", Id: 1, Conf: rtsp.Conf{
					SourceProtocol: "tcp",
					Source:         "source",
				}},
			}}

		json, _ := json.Marshal(stream)
		rw.Write(json)
	}))

	defer server.Close()

	client := server.Client()

	res, err := client.Get(server.URL + "/v1/paths/list")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var items map[string]map[string]rtsp.SConf
	if err := json.NewDecoder(res.Body).Decode(&items); err != nil {
		t.Error(err)
	}
	var streams map[string]rtsp.SConf
	for _, s := range items {
		streams = s
	}
	fmt.Println(streams)

	for stream, info := range streams {

		if stream != "cam1" && info.Stream != "cam1" {
			t.Errorf("unexpected stream: %s", stream)
		}

		if info.Id != 1 {
			t.Errorf("unexpected id: %d", info.Id)
		}

		if info.Conf.SourceProtocol != "tcp" {
			t.Errorf("unexpected protocol: %s", info.Conf.SourceProtocol)
		}

		if info.Conf.Source != "source" {
			t.Errorf("unexpected source: %s", info.Conf.Source)
		}
	}
}
