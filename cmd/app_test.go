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
		return nil, ce.NewError(0, "", "").SetError(args.Error(1))
	}
	return args.Get(0).([]refreshstream.Stream), nil
}
func (c *mockClient) getRTSP(ctx context.Context) (map[string]rtsp.SConf, ce.IError) {
	args := c.Called(ctx)
	if args.Get(0) == nil {
		return nil, ce.NewError(0, "", "").SetError(args.Error(1))
	}
	return args.Get(0).(map[string]rtsp.SConf), nil
}

func TestGetDBAndApi(t *testing.T) {

	mapRTSP := map[string]rtsp.SConf{
		"1": {Id: 1, Stream: "1", Conf: rtsp.Conf{
			Source: "rtsp://login:pass@1/1", SourceProtocol: "udp"}}}
	sliceDB := []refreshstream.Stream{
		{Id: 1, Stream: "1", Auth: sql.NullString{"login:pass", true},
			Portsrv: "38652", Protocol: sql.NullString{"udp", true},
			Ip: sql.NullString{"1", true}},
	}
	ctx := context.Background()
	mu := sync.Mutex{}

	c := newMockClient()
	c.On("getDB", ctx, &mu).Return(sliceDB, nil)
	c.On("getRTSP", ctx).Return(mapRTSP, nil)

	resDB, resRTSP, err := getDBAndApi(ctx, c, &mu)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(resRTSP, mapRTSP) {
		t.Errorf("expect %v, got %v", mapRTSP, resRTSP)
	}
	if !reflect.DeepEqual(resDB, sliceDB) {
		t.Errorf("expect %v, got %v", sliceDB, resDB)
	}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func TestGetCamsAdd(t *testing.T) {
	var test = struct {
		dataDB   []refreshstream.Stream
		dataRTSP map[string]rtsp.SConf
		expect   map[string]rtsp.SConf
	}{
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
	}

	a, err := NewApp(context.Background(), &config.Config{})
	if err != nil {
		fmt.Println("cannot create app", err)
	}

	camsAdd := a.getCamsAdd(test.dataDB, test.dataRTSP)
	if !reflect.DeepEqual(camsAdd, test.expect) {
		t.Errorf("expect %v, got %v", test.expect, camsAdd)
	}
}

func TestGetCamsRemove(t *testing.T) {
	var test = struct {
		dataDB   []refreshstream.Stream
		dataRTSP map[string]rtsp.SConf
		expect   map[string]rtsp.SConf
	}{
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
	}

	a, err := NewApp(context.Background(), &config.Config{})
	if err != nil {
		fmt.Println("cannot create app", err)
	}

	camsRemove := test.dataRTSP
	a.getCamsRemove(test.dataDB, camsRemove)
	if !reflect.DeepEqual(camsRemove, test.expect) {
		t.Errorf("expect %v, got %v", test.expect, camsRemove)
	}
}

func TestGetCamsEdit(t *testing.T) {
	var test = struct {
		dataDB     []refreshstream.Stream
		dataRTSP   map[string]rtsp.SConf
		camsAdd    map[string]rtsp.SConf
		camsRemove map[string]rtsp.SConf
		expect     map[string]rtsp.SConf
	}{
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
	}

	a, err := NewApp(context.Background(), &config.Config{})
	if err != nil {
		fmt.Println("cannot create app", err)
	}

	camsEdit := a.getCamsEdit(test.dataDB, test.dataRTSP, test.camsAdd, test.camsRemove)
	if !reflect.DeepEqual(camsEdit, test.expect) {
		t.Errorf("expect %v, got %v", test.expect, camsEdit)
	}
}

// func (c *mockClient) ListGetCamsAdd(q AnimalFactsQuery) (*AnimalFactsResponse, error) {
// 	args := c.Called(q)
// 	if args.Get(0) == nil {
// 		return nil, args.Error(1)
// 	}

// 	return args.Get(0).(*AnimalFactsResponse), args.Error(1)
// }
