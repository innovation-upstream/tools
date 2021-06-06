// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/module/registry (interfaces: ModuleHeader)

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockModuleHeader is a mock of ModuleHeader interface
type MockModuleHeader struct {
	ctrl     *gomock.Controller
	recorder *MockModuleHeaderMockRecorder
}

// MockModuleHeaderMockRecorder is the mock recorder for MockModuleHeader
type MockModuleHeaderMockRecorder struct {
	mock *MockModuleHeader
}

// NewMockModuleHeader creates a new mock instance
func NewMockModuleHeader(ctrl *gomock.Controller) *MockModuleHeader {
	mock := &MockModuleHeader{ctrl: ctrl}
	mock.recorder = &MockModuleHeaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockModuleHeader) EXPECT() *MockModuleHeaderMockRecorder {
	return m.recorder
}

// GetJSON mocks base method
func (m *MockModuleHeader) GetJSON() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetJSON")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetJSON indicates an expected call of GetJSON
func (mr *MockModuleHeaderMockRecorder) GetJSON() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetJSON", reflect.TypeOf((*MockModuleHeader)(nil).GetJSON))
}

// GetLocation mocks base method
func (m *MockModuleHeader) GetLocation() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLocation")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetLocation indicates an expected call of GetLocation
func (mr *MockModuleHeaderMockRecorder) GetLocation() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLocation", reflect.TypeOf((*MockModuleHeader)(nil).GetLocation))
}

// GetName mocks base method
func (m *MockModuleHeader) GetName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetName indicates an expected call of GetName
func (mr *MockModuleHeaderMockRecorder) GetName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetName", reflect.TypeOf((*MockModuleHeader)(nil).GetName))
}