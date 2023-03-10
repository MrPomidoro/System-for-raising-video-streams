// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Kseniya-cha/System-for-raising-video-streams/cmd (interfaces: AppMock)

// Package mock is a generated GoMock package.
package mock

import (
	"context"
	"reflect"
	"sync"
	// "fmt"
	"database/sql"

	refreshstream "github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	rtspsimpleserver "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server"
	customError "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
	gomock "github.com/golang/mock/gomock"
	rtsp "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server"
)

// MockAppMock is a mock of AppMock interface.
type MockAppMock struct {
	ctrl     *gomock.Controller
	recorder *MockAppMockMockRecorder
}

// MockAppMockMockRecorder is the mock recorder for MockAppMock.
type MockAppMockMockRecorder struct {
	mock *MockAppMock
}

// NewMockAppMock creates a new mock instance.
func NewMockAppMock(ctrl *gomock.Controller) *MockAppMock {
	mock := &MockAppMock{ctrl: ctrl}
	mock.recorder = &MockAppMockMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAppMock) EXPECT() *MockAppMockMockRecorder {
	return m.recorder
}

// AddData mocks base method.
func (m *MockAppMock) AddData(arg0 context.Context, arg1 map[string]rtspsimpleserver.SConf) customError.IError {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddData", arg0, arg1)

	err:= customError.ErrorRefreshStream.SetError(arg0.Err())
	if arg0.Err()!=nil{
		return err
	}

	ret0, _ := ret[0].(customError.IError)
	return ret0
}

// AddData indicates an expected call of AddData.
func (mr *MockAppMockMockRecorder) AddData(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddData", reflect.TypeOf((*MockAppMock)(nil).AddData), arg0, arg1)
}

// AddRemoveData mocks base method.
func (m *MockAppMock) AddRemoveData(arg0 context.Context, arg1 []refreshstream.Stream, arg2, arg3, arg4 map[string]rtspsimpleserver.SConf) customError.IError {

	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddRemoveData", arg0, arg1, arg2, arg3, arg4)

	if arg0.Err()!=nil{
		return customError.ErrorApp.SetError(arg0.Err())
	}

	ret0, _ := ret[0].(customError.IError)
	return ret0
}

// AddRemoveData indicates an expected call of AddRemoveData.
func (mr *MockAppMockMockRecorder) AddRemoveData(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddRemoveData", reflect.TypeOf((*MockAppMock)(nil).AddRemoveData), arg0, arg1, arg2, arg3, arg4)
}

// EditData mocks base method.
func (m *MockAppMock) EditData(arg0 context.Context, arg1 map[string]rtspsimpleserver.SConf) customError.IError {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditData", arg0, arg1)

	err:= customError.ErrorApp.SetError(arg0.Err())
	if arg0.Err()!=nil{
		return err
	}

	ret0, _ := ret[0].(customError.IError)
	return ret0
}

// EditData indicates an expected call of EditData.
func (mr *MockAppMockMockRecorder) EditData(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditData", reflect.TypeOf((*MockAppMock)(nil).EditData), arg0, arg1)
}

// GetCamsAdd mocks base method.
func (m *MockAppMock) GetCamsAdd(arg0 []refreshstream.Stream, arg1 map[string]rtspsimpleserver.SConf) map[string]rtspsimpleserver.SConf {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetCamsAdd", arg0, arg1)
	// ret0, _ := ret[0].(map[string]rtspsimpleserver.SConf)
	return map[string]rtsp.SConf{
		"2": {Id: 2, Stream: "2", Conf: rtsp.Conf{
			Source: "rtsp://login:pass@1/2", SourceProtocol: "udp"}}}
	// return ret0
}

// GetCamsAdd indicates an expected call of GetCamsAdd.
func (mr *MockAppMockMockRecorder) GetCamsAdd(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCamsAdd", reflect.TypeOf((*MockAppMock)(nil).GetCamsAdd), arg0, arg1)
}

// GetCamsEdit mocks base method.
func (m *MockAppMock) GetCamsEdit(arg0 []refreshstream.Stream, arg1, arg2, arg3 map[string]rtspsimpleserver.SConf) map[string]rtspsimpleserver.SConf {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetCamsEdit", arg0, arg1, arg2, arg3)
	// ret0, _ := ret[0].(map[string]rtspsimpleserver.SConf)
	return map[string]rtsp.SConf{
		"1": {Id: 1, Stream: "1", Conf: rtsp.Conf{
			Source: "rtsp://login:pass2@1/1", SourceProtocol: "udp"}}}
	// return ret0
}

// GetCamsEdit indicates an expected call of GetCamsEdit.
func (mr *MockAppMockMockRecorder) GetCamsEdit(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCamsEdit", reflect.TypeOf((*MockAppMock)(nil).GetCamsEdit), arg0, arg1, arg2, arg3)
}

// GetCamsRemove mocks base method.
func (m *MockAppMock) GetCamsRemove(arg0 []refreshstream.Stream, arg1 map[string]rtspsimpleserver.SConf) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetCamsRemove", arg0, arg1)

	for _, camDB := range arg0 {
		delete(arg1, camDB.Stream)
	}
}

// GetCamsRemove indicates an expected call of GetCamsRemove.
func (mr *MockAppMockMockRecorder) GetCamsRemove(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCamsRemove", reflect.TypeOf((*MockAppMock)(nil).GetCamsRemove), arg0, arg1)
}

// GetDBAndApi mocks base method.
func (m *MockAppMock) GetDBAndApi(arg0 context.Context, arg1 *sync.Mutex) ([]refreshstream.Stream, map[string]rtspsimpleserver.SConf, customError.IError) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetDBAndApi", arg0, arg1)

	err := customError.ErrorRefreshStream.SetError(arg0.Err())
	if arg0.Err() != nil{
		return []refreshstream.Stream{}, make(map[string]rtspsimpleserver.SConf), err
	}

	// ret0, _ := ret[0].([]refreshstream.Stream)
	// ret1, _ := ret[1].(map[string]rtspsimpleserver.SConf)
	// ret2, _ := ret[2].(customError.IError)

	return []refreshstream.Stream{
		{Id: 1, Stream: "1", Auth: sql.NullString{String: "login:pass", Valid: true},
			Portsrv: "38652", Protocol: sql.NullString{String: "udp", Valid: true},
			Ip: sql.NullString{String: "1", Valid: true}}},
		map[string]rtsp.SConf{
			"1": {Id: 1, Stream: "1", Conf: rtsp.Conf{
				Source: "rtsp://login:pass@1/1", SourceProtocol: "udp"}}},
		nil
}

// GetDBAndApi indicates an expected call of GetDBAndApi.
func (mr *MockAppMockMockRecorder) GetDBAndApi(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDBAndApi", reflect.TypeOf((*MockAppMock)(nil).GetDBAndApi), arg0, arg1)
}

// GracefulShutdown mocks base method.
func (m *MockAppMock) GracefulShutdown(arg0 context.CancelFunc) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GracefulShutdown", arg0)
}

// GracefulShutdown indicates an expected call of GracefulShutdown.
func (mr *MockAppMockMockRecorder) GracefulShutdown(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GracefulShutdown", reflect.TypeOf((*MockAppMock)(nil).GracefulShutdown), arg0)
}

// RemoveData mocks base method.
func (m *MockAppMock) RemoveData(arg0 context.Context, arg1 map[string]rtspsimpleserver.SConf) customError.IError {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveData", arg0, arg1)

	err:= customError.ErrorApp.SetError(arg0.Err())
	if arg0.Err()!=nil{
		return err
	}

	ret0, _ := ret[0].(customError.IError)
	return ret0
}

// RemoveData indicates an expected call of RemoveData.
func (mr *MockAppMockMockRecorder) RemoveData(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveData", reflect.TypeOf((*MockAppMock)(nil).RemoveData), arg0, arg1)
}

// Run mocks base method.
func (m *MockAppMock) Run(arg0 context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Run", arg0)
}

// Run indicates an expected call of Run.
func (mr *MockAppMockMockRecorder) Run(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockAppMock)(nil).Run), arg0)
}
