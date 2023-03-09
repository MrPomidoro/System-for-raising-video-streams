package mocks

import (
	"context"

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
	ret0, _ := ret[0].([]refreshstream.Stream)
	ret1, _ := ret[1].(ce.IError)
	return ret0, ret1
}

// Close mocks base method
func (m *MockCommon) Close() {
	m.ctrl.Call(m, "Close")
}

func (mr MockCommonMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCall(mr.mock, "Get", arg0, arg1)
}
