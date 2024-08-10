// Code generated by MockGen. DO NOT EDIT.
// Source: source.go

// Package mock_managers is a generated GoMock package.
package mock_managers

import (
	entity "news-aggregator/internal/entity"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockSourceManager is a mock of SourceManager interface.
type MockSourceManager struct {
	ctrl     *gomock.Controller
	recorder *MockSourceManagerMockRecorder
}

// MockSourceManagerMockRecorder is the mock recorder for MockSourceManager.
type MockSourceManagerMockRecorder struct {
	mock *MockSourceManager
}

// NewMockSourceManager creates a new mock instance.
func NewMockSourceManager(ctrl *gomock.Controller) *MockSourceManager {
	mock := &MockSourceManager{ctrl: ctrl}
	mock.recorder = &MockSourceManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSourceManager) EXPECT() *MockSourceManagerMockRecorder {
	return m.recorder
}

// CreateSource mocks base method.
func (m *MockSourceManager) CreateSource(name, url string) (entity.Source, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSource", name, url)
	ret0, _ := ret[0].(entity.Source)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSource indicates an expected call of CreateSource.
func (mr *MockSourceManagerMockRecorder) CreateSource(name, url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSource", reflect.TypeOf((*MockSourceManager)(nil).CreateSource), name, url)
}

// GetSource mocks base method.
func (m *MockSourceManager) GetSource(name string) (entity.Source, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSource", name)
	ret0, _ := ret[0].(entity.Source)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSource indicates an expected call of GetSource.
func (mr *MockSourceManagerMockRecorder) GetSource(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSource", reflect.TypeOf((*MockSourceManager)(nil).GetSource), name)
}

// GetSources mocks base method.
func (m *MockSourceManager) GetSources() ([]entity.Source, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSources")
	ret0, _ := ret[0].([]entity.Source)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSources indicates an expected call of GetSources.
func (mr *MockSourceManagerMockRecorder) GetSources() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSources", reflect.TypeOf((*MockSourceManager)(nil).GetSources))
}

// RemoveSourceByName mocks base method.
func (m *MockSourceManager) RemoveSourceByName(sourceName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveSourceByName", sourceName)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveSourceByName indicates an expected call of RemoveSourceByName.
func (mr *MockSourceManagerMockRecorder) RemoveSourceByName(sourceName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveSourceByName", reflect.TypeOf((*MockSourceManager)(nil).RemoveSourceByName), sourceName)
}

// UpdateSource mocks base method.
func (m *MockSourceManager) UpdateSource(name, newUrl string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSource", name, newUrl)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateSource indicates an expected call of UpdateSource.
func (mr *MockSourceManagerMockRecorder) UpdateSource(name, newUrl interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSource", reflect.TypeOf((*MockSourceManager)(nil).UpdateSource), name, newUrl)
}
