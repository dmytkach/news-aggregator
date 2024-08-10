// Code generated by MockGen. DO NOT EDIT.
// Source: feed.go

// Package mock_managers is a generated GoMock package.
package mock_managers

import (
	entity "news-aggregator/internal/entity"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockFeedManager is a mock of FeedManager interface.
type MockFeedManager struct {
	ctrl     *gomock.Controller
	recorder *MockFeedManagerMockRecorder
}

// MockFeedManagerMockRecorder is the mock recorder for MockFeedManager.
type MockFeedManagerMockRecorder struct {
	mock *MockFeedManager
}

// NewMockFeedManager creates a new mock instance.
func NewMockFeedManager(ctrl *gomock.Controller) *MockFeedManager {
	mock := &MockFeedManager{ctrl: ctrl}
	mock.recorder = &MockFeedManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFeedManager) EXPECT() *MockFeedManagerMockRecorder {
	return m.recorder
}

// FetchFeed mocks base method.
func (m *MockFeedManager) FetchFeed(path string) ([]entity.News, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchFeed", path)
	ret0, _ := ret[0].([]entity.News)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchFeed indicates an expected call of FetchFeed.
func (mr *MockFeedManagerMockRecorder) FetchFeed(path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchFeed", reflect.TypeOf((*MockFeedManager)(nil).FetchFeed), path)
}
