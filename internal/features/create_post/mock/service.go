// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_create_post is a generated GoMock package.
package mock_create_post

import (
	create_post "postservice/internal/features/create_post"
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

// AddNewPostMetaData mocks base method.
func (m *MockRepository) AddNewPostMetaData(data *create_post.Post) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddNewPostMetaData", data)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddNewPostMetaData indicates an expected call of AddNewPostMetaData.
func (mr *MockRepositoryMockRecorder) AddNewPostMetaData(data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNewPostMetaData", reflect.TypeOf((*MockRepository)(nil).AddNewPostMetaData), data)
}

// GetPostMetadata mocks base method.
func (m *MockRepository) GetPostMetadata(postId string) (*create_post.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPostMetadata", postId)
	ret0, _ := ret[0].(*create_post.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPostMetadata indicates an expected call of GetPostMetadata.
func (mr *MockRepositoryMockRecorder) GetPostMetadata(postId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPostMetadata", reflect.TypeOf((*MockRepository)(nil).GetPostMetadata), postId)
}

// GetPresignedUrlForUploadingText mocks base method.
func (m *MockRepository) GetPresignedUrlForUploadingText(data *create_post.Post) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPresignedUrlForUploadingText", data)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPresignedUrlForUploadingText indicates an expected call of GetPresignedUrlForUploadingText.
func (mr *MockRepositoryMockRecorder) GetPresignedUrlForUploadingText(data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPresignedUrlForUploadingText", reflect.TypeOf((*MockRepository)(nil).GetPresignedUrlForUploadingText), data)
}

// RemoveUnconfirmedPost mocks base method.
func (m *MockRepository) RemoveUnconfirmedPost(postId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveUnconfirmedPost", postId)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveUnconfirmedPost indicates an expected call of RemoveUnconfirmedPost.
func (mr *MockRepositoryMockRecorder) RemoveUnconfirmedPost(postId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveUnconfirmedPost", reflect.TypeOf((*MockRepository)(nil).RemoveUnconfirmedPost), postId)
}
