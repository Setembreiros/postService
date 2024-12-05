// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_get_post is a generated GoMock package.
package mock_get_post

import (
	get_post "postservice/internal/features/get_post"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// GetPresignedUrlsForDownloading mocks base method.
func (m *MockRepository) GetPresignedUrlsForDownloading(username, lastCreatedAt string, limit int) ([]get_post.PostUrl, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPresignedUrlsForDownloading", username, lastCreatedAt, limit)
	ret0, _ := ret[0].([]get_post.PostUrl)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetPresignedUrlsForDownloading indicates an expected call of GetPresignedUrlsForDownloading.
func (mr *MockRepositoryMockRecorder) GetPresignedUrlsForDownloading(username, lastCreatedAt, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPresignedUrlsForDownloading", reflect.TypeOf((*MockRepository)(nil).GetPresignedUrlsForDownloading), username, lastCreatedAt, limit)
}
