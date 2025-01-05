// Code generated by MockGen. DO NOT EDIT.
// Source: object_storage.go

// Package mock_objectstorage is a generated GoMock package.
package mock_objectstorage

import (
	objectstorage "postservice/internal/objectStorage"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockObjectStorageClient is a mock of ObjectStorageClient interface.
type MockObjectStorageClient struct {
	ctrl     *gomock.Controller
	recorder *MockObjectStorageClientMockRecorder
}

// MockObjectStorageClientMockRecorder is the mock recorder for MockObjectStorageClient.
type MockObjectStorageClientMockRecorder struct {
	mock *MockObjectStorageClient
}

// NewMockObjectStorageClient creates a new mock instance.
func NewMockObjectStorageClient(ctrl *gomock.Controller) *MockObjectStorageClient {
	mock := &MockObjectStorageClient{ctrl: ctrl}
	mock.recorder = &MockObjectStorageClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockObjectStorageClient) EXPECT() *MockObjectStorageClientMockRecorder {
	return m.recorder
}

// CompleteMultipartUpload mocks base method.
func (m *MockObjectStorageClient) CompleteMultipartUpload(multipartobject objectstorage.MultipartObject) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CompleteMultipartUpload", multipartobject)
	ret0, _ := ret[0].(error)
	return ret0
}

// CompleteMultipartUpload indicates an expected call of CompleteMultipartUpload.
func (mr *MockObjectStorageClientMockRecorder) CompleteMultipartUpload(multipartobject interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CompleteMultipartUpload", reflect.TypeOf((*MockObjectStorageClient)(nil).CompleteMultipartUpload), multipartobject)
}

// DeleteObjects mocks base method.
func (m *MockObjectStorageClient) DeleteObjects(objectKeys []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteObjects", objectKeys)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteObjects indicates an expected call of DeleteObjects.
func (mr *MockObjectStorageClientMockRecorder) DeleteObjects(objectKeys interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteObjects", reflect.TypeOf((*MockObjectStorageClient)(nil).DeleteObjects), objectKeys)
}

// GetPreSignedUrlForGettingObject mocks base method.
func (m *MockObjectStorageClient) GetPreSignedUrlForGettingObject(objectKey string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPreSignedUrlForGettingObject", objectKey)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPreSignedUrlForGettingObject indicates an expected call of GetPreSignedUrlForGettingObject.
func (mr *MockObjectStorageClientMockRecorder) GetPreSignedUrlForGettingObject(objectKey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPreSignedUrlForGettingObject", reflect.TypeOf((*MockObjectStorageClient)(nil).GetPreSignedUrlForGettingObject), objectKey)
}

// GetPreSignedUrlsForPuttingObject mocks base method.
func (m *MockObjectStorageClient) GetPreSignedUrlsForPuttingObject(objectKey string, size int) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPreSignedUrlsForPuttingObject", objectKey, size)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPreSignedUrlsForPuttingObject indicates an expected call of GetPreSignedUrlsForPuttingObject.
func (mr *MockObjectStorageClientMockRecorder) GetPreSignedUrlsForPuttingObject(objectKey, size interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPreSignedUrlsForPuttingObject", reflect.TypeOf((*MockObjectStorageClient)(nil).GetPreSignedUrlsForPuttingObject), objectKey, size)
}
