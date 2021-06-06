// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/code_layer_generator (interfaces: CodeLayerGenerator)

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	code_layer_generator "gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/code_layer_generator"
	model_frame_path "gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/model_frame_path"
	reflect "reflect"
)

// MockCodeLayerGenerator is a mock of CodeLayerGenerator interface
type MockCodeLayerGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockCodeLayerGeneratorMockRecorder
}

// MockCodeLayerGeneratorMockRecorder is the mock recorder for MockCodeLayerGenerator
type MockCodeLayerGeneratorMockRecorder struct {
	mock *MockCodeLayerGenerator
}

// NewMockCodeLayerGenerator creates a new mock instance
func NewMockCodeLayerGenerator(ctrl *gomock.Controller) *MockCodeLayerGenerator {
	mock := &MockCodeLayerGenerator{ctrl: ctrl}
	mock.recorder = &MockCodeLayerGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCodeLayerGenerator) EXPECT() *MockCodeLayerGeneratorMockRecorder {
	return m.recorder
}

// GenerateCodeLayersForFramePath mocks base method
func (m *MockCodeLayerGenerator) GenerateCodeLayersForFramePath(arg0 model_frame_path.ModelFramePath) (code_layer_generator.ModuleCodeLayers, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateCodeLayersForFramePath", arg0)
	ret0, _ := ret[0].(code_layer_generator.ModuleCodeLayers)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateCodeLayersForFramePath indicates an expected call of GenerateCodeLayersForFramePath
func (mr *MockCodeLayerGeneratorMockRecorder) GenerateCodeLayersForFramePath(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateCodeLayersForFramePath", reflect.TypeOf((*MockCodeLayerGenerator)(nil).GenerateCodeLayersForFramePath), arg0)
}
