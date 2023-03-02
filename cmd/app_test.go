package service

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"sync"
	"testing"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	rtsp "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
	"github.com/stretchr/testify/mock"
)

type mockClient struct{ mock.Mock }

func newMockClient() *mockClient { return &mockClient{} }

func (c *mockClient) getDB(ctx context.Context, mu *sync.Mutex) ([]refreshstream.Stream, ce.IError) {
	args := c.Called(ctx, mu)
	if args.Get(0) == nil {
		e, ok := args.Error(1).(ce.Error)
		if !ok {
			fmt.Println("blyat")
		} else {
			return nil, &e
		}
	}
	return args.Get(0).([]refreshstream.Stream), nil
}

func (c *mockClient) getRTSP(ctx context.Context) (map[string]rtsp.SConf, ce.IError) {
	args := c.Called(ctx)
	if args.Get(0) == nil {
		e, ok := args.Error(1).(ce.Error)
		if !ok {
			fmt.Println("blyat")
		} else {
			return nil, &e
		}
	}
	return args.Get(0).(map[string]rtsp.SConf), nil
}

func TestGetDBAndApi(t *testing.T) {

	var tests = []struct {
		name     string
		wantRTSP map[string]rtsp.SConf
		wantDB   []refreshstream.Stream
	}{
		{
			name: "Test for case correct work",
			wantRTSP: map[string]rtsp.SConf{
				"1": {Id: 1, Stream: "1", Conf: rtsp.Conf{
					Source: "rtsp://login:pass@1/1", SourceProtocol: "udp"}}},
			wantDB: []refreshstream.Stream{
				{Id: 1, Stream: "1", Auth: sql.NullString{"login:pass", true},
					Portsrv: "38652", Protocol: sql.NullString{"udp", true},
					Ip: sql.NullString{"1", true}}},
		},
	}

	ctx := context.Background()
	mu := sync.Mutex{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			c := newMockClient()

			c.On("getDB", ctx, &mu).Return(tt.wantDB, nil)
			c.On("getRTSP", ctx).Return(tt.wantRTSP, nil)

			resDB, resRTSP, _ := getDBAndApi(ctx, c, &mu)

			if !reflect.DeepEqual(resRTSP, tt.wantRTSP) {
				t.Errorf("expect %v, got %v", tt.wantRTSP, resRTSP)
			}
			if !reflect.DeepEqual(resDB, tt.wantDB) {
				t.Errorf("expect %v, got %v", tt.wantDB, resDB)
			}
		})
	}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func TestGetCamsAdd(t *testing.T) {
	var tests = []struct {
		name     string
		dataDB   []refreshstream.Stream
		dataRTSP map[string]rtsp.SConf
		expect   map[string]rtsp.SConf
	}{
		{
			name: "Test for get streams to add",
			dataDB: []refreshstream.Stream{
				{Id: 1, Stream: "1", Auth: sql.NullString{"login:pass", true},
					Portsrv: "38652", Protocol: sql.NullString{"udp", true},
					Ip: sql.NullString{"1", true}},
				{Id: 2, Stream: "2", Auth: sql.NullString{"login:pass", true},
					Portsrv: "38652", Protocol: sql.NullString{"udp", true},
					Ip: sql.NullString{"1", true}},
			},
			dataRTSP: map[string]rtsp.SConf{
				"1": {Id: 1, Stream: "1", Conf: rtsp.Conf{
					Source: "rtsp://login:pass@1/1", SourceProtocol: "udp"}},
			},
			expect: map[string]rtsp.SConf{
				"2": {Id: 2, Stream: "2", Conf: rtsp.Conf{
					Source: "rtsp://login:pass@1/2", SourceProtocol: "udp"}},
			},
		},
	}

	a, err := NewApp(context.Background(), &config.Config{})
	if err != nil {
		fmt.Println("cannot create app", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			camsAdd := a.getCamsAdd(tt.dataDB, tt.dataRTSP)
			if !reflect.DeepEqual(camsAdd, tt.expect) {
				t.Errorf("expect %v, got %v", tt.expect, camsAdd)
			}
		})
	}
}

func TestGetCamsRemove(t *testing.T) {
	var tests = []struct {
		name     string
		dataDB   []refreshstream.Stream
		dataRTSP map[string]rtsp.SConf
		expect   map[string]rtsp.SConf
	}{
		{
			name: "Test for get streams to remove",
			dataDB: []refreshstream.Stream{
				{Id: 1, Stream: "1", Auth: sql.NullString{"login:pass", true},
					Portsrv: "38652", Protocol: sql.NullString{"udp", true},
					Ip: sql.NullString{"1", true}},
			},
			dataRTSP: map[string]rtsp.SConf{
				"1": {Id: 1, Stream: "1", Conf: rtsp.Conf{
					Source: "rtsp://login:pass@1/1", SourceProtocol: "udp"}},
				"3": {Id: 3, Stream: "3", Conf: rtsp.Conf{
					Source: "rtsp://login:pass@1/3", SourceProtocol: "udp"}},
			},
			expect: map[string]rtsp.SConf{
				"3": {Id: 3, Stream: "3", Conf: rtsp.Conf{
					Source: "rtsp://login:pass@1/3", SourceProtocol: "udp"}},
			},
		},
	}

	a, err := NewApp(context.Background(), &config.Config{})
	if err != nil {
		fmt.Println("cannot create app", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			camsRemove := tt.dataRTSP
			a.getCamsRemove(tt.dataDB, camsRemove)
			if !reflect.DeepEqual(camsRemove, tt.expect) {
				t.Errorf("expect %v, got %v", tt.expect, camsRemove)
			}
		})
	}
}

func TestGetCamsEdit(t *testing.T) {
	var tests = []struct {
		name       string
		dataDB     []refreshstream.Stream
		dataRTSP   map[string]rtsp.SConf
		camsAdd    map[string]rtsp.SConf
		camsRemove map[string]rtsp.SConf
		expect     map[string]rtsp.SConf
	}{
		{
			name: "Test for get streams to edit",
			dataDB: []refreshstream.Stream{
				{Id: 1, Stream: "1", Auth: sql.NullString{"login:pass2", true},
					Portsrv: "38652", Protocol: sql.NullString{"udp", true},
					Ip: sql.NullString{"1", true}},
			},
			dataRTSP: map[string]rtsp.SConf{
				"1": {Id: 1, Stream: "1", Conf: rtsp.Conf{
					Source: "rtsp://login:pass@1/1", SourceProtocol: "udp"}},
			},
			camsAdd:    map[string]rtsp.SConf{},
			camsRemove: map[string]rtsp.SConf{},

			expect: map[string]rtsp.SConf{
				"1": {Id: 1, Stream: "1", Conf: rtsp.Conf{
					Source: "rtsp://login:pass2@1/1", SourceProtocol: "udp"}},
			},
		},
	}

	a, err := NewApp(context.Background(), &config.Config{})
	if err != nil {
		fmt.Println("cannot create app", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			camsEdit := a.getCamsEdit(tt.dataDB, tt.dataRTSP, tt.camsAdd, tt.camsRemove)
			if !reflect.DeepEqual(camsEdit, tt.expect) {
				t.Errorf("expect %v, got %v", tt.expect, camsEdit)
			}
		})
	}
}

// dbToCompare(cfg *config.Config, camDB refreshstream.Stream) rtsp.SConf
func TestDbToCompare(t *testing.T) {

	var tests = []struct {
		name   string
		cfg    config.Config
		camDB  refreshstream.Stream
		expect rtsp.SConf
	}{
		{
			name: "Test have Run",
			cfg:  config.Config{Rtsp: config.Rtsp{Run: "usr/bin/av_reader-1.1.7/av_reader --config_file /etc/rss/rss-av_reader.yml --port %s --stream_path %s --camera_id %s"}},
			camDB: refreshstream.Stream{Id: 1, Stream: "1", Auth: sql.NullString{"login:pass", true},
				Portsrv: "1", Protocol: sql.NullString{"udp", true}, CamId: sql.NullString{"1", true},
				Ip: sql.NullString{"1", true}, Sp: sql.NullString{"1", true}},
			expect: rtsp.SConf{
				Stream: "1",
				Id:     1,
				Conf: rtsp.Conf{
					SourceProtocol: "udp",
					RunOnReady:     "usr/bin/av_reader-1.1.7/av_reader --config_file /etc/rss/rss-av_reader.yml --port 1 --stream_path 1 --camera_id 1",
					Source:         "rtsp://login:pass@1/1",
				},
			},
		},
		{
			name: "Test have not Run",
			cfg:  config.Config{Rtsp: config.Rtsp{Run: ""}},
			camDB: refreshstream.Stream{Id: 1, Stream: "1", Auth: sql.NullString{"login:pass", true},
				Portsrv: "1", Protocol: sql.NullString{"udp", true}, CamId: sql.NullString{"1", true},
				Ip: sql.NullString{"1", true}, Sp: sql.NullString{"1", true}},
			expect: rtsp.SConf{
				Stream: "1",
				Id:     1,
				Conf: rtsp.Conf{
					SourceProtocol: "udp",
					RunOnReady:     "",
					Source:         "rtsp://login:pass@1/1",
				},
			},
		},
		{
			name: "Test have not protocol",
			cfg:  config.Config{Rtsp: config.Rtsp{Run: ""}},
			camDB: refreshstream.Stream{Id: 1, Stream: "1", Auth: sql.NullString{"login:pass", true},
				Portsrv: "1", Protocol: sql.NullString{"", false}, CamId: sql.NullString{"1", true},
				Ip: sql.NullString{"1", true}, Sp: sql.NullString{"1", true}},
			expect: rtsp.SConf{
				Stream: "1",
				Id:     1,
				Conf: rtsp.Conf{
					SourceProtocol: "tcp",
					RunOnReady:     "",
					Source:         "rtsp://login:pass@1/1",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newDB := dbToCompare(&tt.cfg, tt.camDB)
			if !reflect.DeepEqual(newDB, tt.expect) {
				t.Errorf("expect %+v\ngot %+v", tt.expect, newDB)
			}
		})
	}
}
