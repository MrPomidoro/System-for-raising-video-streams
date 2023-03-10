package mocks

import (
	"context"
	"database/sql"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
	"github.com/golang/mock/gomock"
)

// MockCommon is a mock implementation of Common interface
type MockCommon struct {
	ctrl     *gomock.Controller
	recorder *MockCommonMockRecorder
}

// MockCommonMockRecorder is the mock recorder for MockCommon
type MockCommonMockRecorder struct {
	mock *MockCommon
}

// NewMockCommon creates a new mock instance
func NewMockCommon(ctrl *gomock.Controller) *MockCommon {
	mock := &MockCommon{ctrl: ctrl}
	mock.recorder = &MockCommonMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCommon) EXPECT() *MockCommonMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockCommon) Get(ctx context.Context, status bool) ([]refreshstream.Stream, ce.IError) {
	ret := m.ctrl.Call(m, "Get", ctx, status)

	if ctx.Err() != nil {
		return nil, ce.ErrorRefreshStream.SetError(ctx.Err())
	}

	var ret0 []refreshstream.Stream
	if status {
		ret0 = []refreshstream.Stream{{
			Id:           1,
			Auth:         sql.NullString{String: "login:pass", Valid: true},
			Ip:           sql.NullString{String: "ip", Valid: true},
			Stream:       "1",
			Portsrv:      "123",
			Sp:           sql.NullString{String: "sp", Valid: true},
			CamId:        sql.NullString{String: "cam1", Valid: true},
			Protocol:     sql.NullString{String: "tcp", Valid: true},
			RecordStatus: sql.NullBool{Bool: true, Valid: true},
			StreamStatus: sql.NullBool{Bool: true, Valid: true},
			RecordState:  sql.NullBool{Bool: true, Valid: true},
			StreamState:  sql.NullBool{Bool: true, Valid: true},
		}}
	} else {
		ret0 = []refreshstream.Stream{{
			Id:           1,
			Auth:         sql.NullString{String: "login:pass", Valid: true},
			Ip:           sql.NullString{String: "ip", Valid: true},
			Stream:       "1",
			Portsrv:      "123",
			Sp:           sql.NullString{String: "sp", Valid: true},
			CamId:        sql.NullString{String: "cam1", Valid: true},
			Protocol:     sql.NullString{String: "tcp", Valid: true},
			RecordStatus: sql.NullBool{Bool: true, Valid: true},
			StreamStatus: sql.NullBool{Bool: true, Valid: true},
			RecordState:  sql.NullBool{Bool: false, Valid: true},
			StreamState:  sql.NullBool{Bool: true, Valid: true},
		}}
	}
	ret1, _ := ret[1].(ce.IError)

	return ret0, ret1
}

// Close mocks base method
func (m *MockCommon) Close() {
	m.ctrl.Call(m, "Close")
}

func (mr *MockCommonMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCall(mr.mock, "Get", arg0, arg1)
}
