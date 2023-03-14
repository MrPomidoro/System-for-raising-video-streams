package service

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"sync"
	"testing"

	appM "github.com/Kseniya-cha/System-for-raising-video-streams/cmd/mock"
	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	rtsp "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	"github.com/golang/mock/gomock"
)

func TestGetDBAndApi(t *testing.T) {

	var tests = []struct {
		name      string
		wantRTSP  map[string]rtsp.SConf
		wantDB    []refreshstream.Stream
		expectErr error
		isCanc    bool
	}{
		{
			name: "TestOK",
			wantRTSP: map[string]rtsp.SConf{
				"1": {Id: 1, Stream: "1", Conf: rtsp.Conf{
					Source: "rtsp://login:pass@1/1", SourceProtocol: "udp"}}},
			wantDB: []refreshstream.Stream{
				// {Id: 1, Stream: "1", Auth: sql.NullString{String: "login:pass", Valid: true},
				// 	Portsrv: "38652", Protocol: sql.NullString{String: "udp", Valid: true},
				// 	Ip: sql.NullString{String: "1", Valid: true}},
			},
			expectErr: nil,
			isCanc:    false,
		},
		{
			name:      "TestCtxCancel",
			wantRTSP:  nil,
			wantDB:    nil,
			expectErr: errors.New("context canceled"),
			isCanc:    true,
		},
	}

	mu := sync.Mutex{}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	a := appM.NewMockAppMock(ctrl)

	for _, tt := range tests {
		ctx, cancel := context.WithCancel(context.Background())
		t.Run(tt.name, func(t *testing.T) {

			if tt.isCanc {
				cancel()
			}
			a.EXPECT().GetDBAndApi(ctx, &mu)

			resDB, resRTSP, err := a.GetDBAndApi(ctx, &mu)

			if err != nil && tt.expectErr != nil {
				gotErrA := strings.Split(err.Error(), ": ")

				if tt.expectErr.Error() != gotErrA[len(gotErrA)-1] {
					t.Errorf("unexpected error %v", err)
				}
			} else if (err == nil && tt.expectErr != nil) || (err != nil && tt.expectErr == nil) {
				t.Errorf("unexpected error %v, expect %v", err, tt.expectErr)
			}

			if len(resDB) != len(tt.wantDB) || len(resRTSP) != len(tt.wantRTSP) {
				t.Errorf("expect %v, got %v", tt.wantDB, resDB)
			}
			for i := range resDB {
				if resDB[i] != tt.wantDB[i] {
					t.Errorf("expect %v, got %v", tt.wantDB, resDB)
				}
			}
			for k := range resRTSP {
				if resRTSP[k] != tt.wantRTSP[k] {
					t.Errorf("expect %v, got %v", tt.wantRTSP, resRTSP)
				}
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
			name:   "TestCamsAddOK",
			dataDB: []refreshstream.Stream{
				// {Id: 1, Stream: "1", Auth: sql.NullString{String: "login:pass", Valid: true},
				// 	Portsrv: "38652", Protocol: sql.NullString{String: "udp", Valid: true},
				// 	Ip: sql.NullString{String: "1", Valid: true}},
				// {Id: 2, Stream: "2", Auth: sql.NullString{String: "login:pass", Valid: true},
				// 	Portsrv: "38652", Protocol: sql.NullString{String: "udp", Valid: true},
				// 	Ip: sql.NullString{String: "1", Valid: true}},
			},
			dataRTSP: map[string]rtsp.SConf{
				// "1": {Id: 1, Stream: "1", Conf: rtsp.Conf{
				// 	Source: "rtsp://login:pass@1/1", SourceProtocol: "udp"}},
			},
			expect: map[string]rtsp.SConf{
				// "2": {Id: 2, Stream: "2", Conf: rtsp.Conf{
				// 	Source: "rtsp://login:pass@1/2", SourceProtocol: "udp"}},
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	a := appM.NewMockAppMock(ctrl)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a.EXPECT().GetCamsAdd(tt.dataDB, tt.dataRTSP)

			camsAdd := a.GetCamsAdd(tt.dataDB, tt.dataRTSP)
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
			name:   "TesOK",
			dataDB: []refreshstream.Stream{
				// {Id: 1, Stream: "1", Auth: sql.NullString{String: "login:pass", Valid: true},
				// 	Portsrv: "38652", Protocol: sql.NullString{String: "udp", Valid: true},
				// 	Ip: sql.NullString{String: "1", Valid: true}},
			},
			dataRTSP: map[string]rtsp.SConf{
				// "1": {Id: 1, Stream: "1", Conf: rtsp.Conf{
				// 	Source: "rtsp://login:pass@1/1", SourceProtocol: "udp"}},
				// "3": {Id: 3, Stream: "3", Conf: rtsp.Conf{
				// 	Source: "rtsp://login:pass@1/3", SourceProtocol: "udp"}},
			},
			expect: map[string]rtsp.SConf{
				// "3": {Id: 3, Stream: "3", Conf: rtsp.Conf{
				// 	Source: "rtsp://login:pass@1/3", SourceProtocol: "udp"}},
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	a := appM.NewMockAppMock(ctrl)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a.EXPECT().GetCamsRemove(tt.dataDB, tt.dataRTSP)

			camsRemove := tt.dataRTSP
			a.GetCamsRemove(tt.dataDB, camsRemove)
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
			name:   "TestOK",
			dataDB: []refreshstream.Stream{
				// {Id: 1, Stream: "1", Auth: sql.NullString{String: "login:pass2", Valid: true},
				// 	Portsrv: "38652", Protocol: sql.NullString{String: "udp", Valid: true},
				// 	Ip: sql.NullString{String: "1", Valid: true}},
			},
			dataRTSP: map[string]rtsp.SConf{
				// "1": {Id: 1, Stream: "1", Conf: rtsp.Conf{
				// 	Source: "rtsp://login:pass@1/1", SourceProtocol: "udp"}},
			},
			camsAdd:    map[string]rtsp.SConf{},
			camsRemove: map[string]rtsp.SConf{},

			expect: map[string]rtsp.SConf{
				// "1": {Id: 1, Stream: "1", Conf: rtsp.Conf{
				// 	Source: "rtsp://login:pass2@1/1", SourceProtocol: "udp"}},
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	a := appM.NewMockAppMock(ctrl)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a.EXPECT().GetCamsEdit(tt.dataDB, tt.dataRTSP, tt.camsAdd, tt.camsRemove)

			camsEdit := a.GetCamsEdit(tt.dataDB, tt.dataRTSP, tt.camsAdd, tt.camsRemove)
			if !reflect.DeepEqual(camsEdit, tt.expect) {
				t.Errorf("expect %v, got %v", tt.expect, camsEdit)
			}
		})
	}
}

func TestDbToCompare(t *testing.T) {

	var tests = []struct {
		name   string
		cfg    config.Config
		camDB  refreshstream.Stream
		expect rtsp.SConf
	}{
		{
			name:  "TestHaveRunField",
			cfg:   config.Config{Rtsp: config.Rtsp{Run: "usr/bin/av_reader-1.1.7/av_reader --config_file /etc/rss/rss-av_reader.yml --port %s --stream_path %s --camera_id %s"}},
			camDB: refreshstream.Stream{
				// Id: 1, Stream: "1", Auth: sql.NullString{String: "login:pass", Valid: true},
				// Portsrv: "1", Protocol: sql.NullString{String: "udp", Valid: true}, CamId: sql.NullString{String: "1", Valid: true},
				// Ip: sql.NullString{String: "1", Valid: true}, Sp: sql.NullString{String: "1", Valid: true}
			},
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
			name:  "TestHaveNotRunField",
			cfg:   config.Config{Rtsp: config.Rtsp{Run: ""}},
			camDB: refreshstream.Stream{
				// Id: 1, Stream: "1", Auth: sql.NullString{String: "login:pass", Valid: true},
				// Portsrv: "1", Protocol: sql.NullString{String: "udp", Valid: true}, CamId: sql.NullString{String: "1", Valid: true},
				// Ip: sql.NullString{String: "1", Valid: true}, Sp: sql.NullString{String: "1", Valid: true}
			},
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
			name:  "TestHaveNotProtocolField",
			cfg:   config.Config{Rtsp: config.Rtsp{Run: ""}},
			camDB: refreshstream.Stream{
				// Id: 1, Stream: "1", Auth: sql.NullString{String: "login:pass", Valid: true},
				// Portsrv: "1", Protocol: sql.NullString{String: "", Valid: false}, CamId: sql.NullString{String: "1", Valid: true},
				// Ip: sql.NullString{String: "1", Valid: true}, Sp: sql.NullString{String: "1", Valid: true}
			},
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

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func TestAddRemoveData(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	a := appM.NewMockAppMock(ctrl)

	var tests = []struct {
		name       string
		dataDB     []refreshstream.Stream
		dataRTSP   map[string]rtsp.SConf
		camsAdd    map[string]rtsp.SConf
		camsRemove map[string]rtsp.SConf
		isCanc     bool
		expectErr  error
	}{
		{
			name:   "TestOK",
			dataDB: []refreshstream.Stream{
				// {Id: 1, Stream: "1", Auth: sql.NullString{String: "login:pass", Valid: true},
				// 	Portsrv: "38652", Protocol: sql.NullString{String: "udp", Valid: true},
				// 	Ip: sql.NullString{String: "1", Valid: true}},
				// {Id: 2, Stream: "2", Auth: sql.NullString{String: "login:pass", Valid: true},
				// 	Portsrv: "38652", Protocol: sql.NullString{String: "udp", Valid: true},
				// 	Ip: sql.NullString{String: "1", Valid: true}},
			},
			dataRTSP: map[string]rtsp.SConf{
				// "1": {Id: 1, Stream: "1", Conf: rtsp.Conf{
				// 	Source: "rtsp://login:pass@1/1", SourceProtocol: "udp"}},
				// "3": {Id: 3, Stream: "3", Conf: rtsp.Conf{
				// 	Source: "rtsp://login:pass@1/3", SourceProtocol: "udp"}},
			},
			camsAdd: map[string]rtsp.SConf{
				// "2": {Id: 2, Stream: "2", Conf: rtsp.Conf{
				// 	Source: "rtsp://login:pass@1/2", SourceProtocol: "udp"}},
			},
			camsRemove: map[string]rtsp.SConf{
				// "3": {Id: 0, Stream: "3", Conf: rtsp.Conf{
				// 	Source: "rtsp://login:pass@1/3", SourceProtocol: "udp"}},
			},
			isCanc:    false,
			expectErr: nil,
		},
		{
			name:   "TestOKNothingChange",
			dataDB: []refreshstream.Stream{
				// {Id: 1, Stream: "1", Auth: sql.NullString{String: "login:pass", Valid: true},
				// 	Portsrv: "38652", Protocol: sql.NullString{String: "udp", Valid: true},
				// 	Ip: sql.NullString{String: "1", Valid: true}},
			},
			dataRTSP: map[string]rtsp.SConf{
				// "1": {Id: 1, Stream: "1", Conf: rtsp.Conf{
				// 	Source: "rtsp://login:pass@1/1", SourceProtocol: "udp"}},
			},
			isCanc:    false,
			expectErr: nil,
		},
		{
			name:      "TestCtxCancel",
			dataDB:    []refreshstream.Stream{},
			dataRTSP:  map[string]rtsp.SConf{},
			isCanc:    true,
			expectErr: errors.New("context canceled"),
		},
	}

	for _, tt := range tests {
		ctx, cancel := context.WithCancel(context.Background())

		t.Run(tt.name, func(t *testing.T) {

			if tt.isCanc {
				cancel()
			}
			a.EXPECT().AddRemoveData(ctx, tt.dataDB, tt.dataRTSP, tt.camsAdd, tt.camsRemove)

			err := a.AddRemoveData(ctx, tt.dataDB, tt.dataRTSP, tt.camsAdd, tt.camsRemove)

			if err != nil && tt.expectErr != nil {
				gotErrA := strings.Split(err.Error(), ": ")

				if tt.expectErr.Error() != gotErrA[len(gotErrA)-1] {
					t.Errorf("unexpected error %v", err)
				}
			} else if (err == nil && tt.expectErr != nil) || (err != nil && tt.expectErr == nil) {
				t.Errorf("unexpected error %v, expect %v", err, tt.expectErr)
			}
		})
	}
}

func TestAddData(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	a := appM.NewMockAppMock(ctrl)

	var tests = []struct {
		name      string
		camsAdd   map[string]rtsp.SConf
		isCanc    bool
		expectErr error
	}{
		{
			name:    "TestOK",
			camsAdd: map[string]rtsp.SConf{
				// "2": {Id: 0, Stream: "2", Conf: rtsp.Conf{
				// 	Source: "rtsp://login:pass@1/2", SourceProtocol: "udp"}},
			},
			isCanc:    false,
			expectErr: nil,
		},
		{
			name:      "TestCtxCancel",
			camsAdd:   map[string]rtsp.SConf{},
			isCanc:    true,
			expectErr: errors.New("context canceled"),
		},
	}

	for _, tt := range tests {
		ctx, cancel := context.WithCancel(context.Background())

		t.Run(tt.name, func(t *testing.T) {

			if tt.isCanc {
				cancel()
			}
			a.EXPECT().AddData(ctx, tt.camsAdd)

			err := a.AddData(ctx, tt.camsAdd)

			if err != nil && tt.expectErr != nil {
				gotErrA := strings.Split(err.Error(), ": ")

				if tt.expectErr.Error() != gotErrA[len(gotErrA)-1] {
					t.Errorf("unexpected error %v", err)
				}
			} else if (err == nil && tt.expectErr != nil) || (err != nil && tt.expectErr == nil) {
				t.Errorf("unexpected error %v, expect %v", err, tt.expectErr)
			}
		})
	}
}

func TestRemoveData(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	a := appM.NewMockAppMock(ctrl)

	var tests = []struct {
		name       string
		camsRemove map[string]rtsp.SConf
		isCanc     bool
		expectErr  error
	}{
		{
			name: "TestOK",
			camsRemove: map[string]rtsp.SConf{
				"2": {Id: 0, Stream: "2", Conf: rtsp.Conf{
					Source: "rtsp://login:pass@1/2", SourceProtocol: "udp"}},
			},
			isCanc:    false,
			expectErr: nil,
		},
		{
			name:       "TestCtxCancel",
			camsRemove: map[string]rtsp.SConf{},
			isCanc:     true,
			expectErr:  errors.New("context canceled"),
		},
	}

	for _, tt := range tests {
		ctx, cancel := context.WithCancel(context.Background())

		t.Run(tt.name, func(t *testing.T) {

			if tt.isCanc {
				cancel()
			}
			a.EXPECT().RemoveData(ctx, tt.camsRemove)

			err := a.RemoveData(ctx, tt.camsRemove)

			if err != nil && tt.expectErr != nil {
				gotErrA := strings.Split(err.Error(), ": ")

				if tt.expectErr.Error() != gotErrA[len(gotErrA)-1] {
					t.Errorf("unexpected error %v", err)
				}
			} else if (err == nil && tt.expectErr != nil) || (err != nil && tt.expectErr == nil) {
				t.Errorf("unexpected error %v, expect %v", err, tt.expectErr)
			}
		})
	}
}

func TestEditData(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	a := appM.NewMockAppMock(ctrl)

	var tests = []struct {
		name      string
		camsEdit  map[string]rtsp.SConf
		isCanc    bool
		expectErr error
	}{
		{
			name: "TestOK",
			camsEdit: map[string]rtsp.SConf{
				"2": {Id: 0, Stream: "2", Conf: rtsp.Conf{
					Source: "rtsp://login:pass@1/2", SourceProtocol: "udp"}},
			},
			isCanc:    false,
			expectErr: nil,
		},
		{
			name:      "TestCtxCancel",
			camsEdit:  map[string]rtsp.SConf{},
			isCanc:    true,
			expectErr: errors.New("context canceled"),
		},
	}

	for _, tt := range tests {
		ctx, cancel := context.WithCancel(context.Background())

		t.Run(tt.name, func(t *testing.T) {

			if tt.isCanc {
				cancel()
			}
			a.EXPECT().EditData(ctx, tt.camsEdit)

			err := a.EditData(ctx, tt.camsEdit)

			if err != nil && tt.expectErr != nil {
				gotErrA := strings.Split(err.Error(), ": ")

				if tt.expectErr.Error() != gotErrA[len(gotErrA)-1] {
					t.Errorf("unexpected error %v", err)
				}
			} else if (err == nil && tt.expectErr != nil) || (err != nil && tt.expectErr == nil) {
				t.Errorf("unexpected error %v, expect %v", err, tt.expectErr)
			}
		})
	}
}
