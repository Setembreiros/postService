// Code generated by MockGen. DO NOT EDIT.
// Source: database.go

// Package mock_database is a generated GoMock package.
package mock_database

import (
	context "context"
	database "postservice/internal/db"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockDatabaseClient is a mock of DatabaseClient interface.
type MockDatabaseClient struct {
	ctrl     *gomock.Controller
	recorder *MockDatabaseClientMockRecorder
}

// MockDatabaseClientMockRecorder is the mock recorder for MockDatabaseClient.
type MockDatabaseClientMockRecorder struct {
	mock *MockDatabaseClient
}

// NewMockDatabaseClient creates a new mock instance.
func NewMockDatabaseClient(ctrl *gomock.Controller) *MockDatabaseClient {
	mock := &MockDatabaseClient{ctrl: ctrl}
	mock.recorder = &MockDatabaseClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDatabaseClient) EXPECT() *MockDatabaseClientMockRecorder {
	return m.recorder
}

// CreateIndexesOnTable mocks base method.
func (m *MockDatabaseClient) CreateIndexesOnTable(tableName, indexName string, inndexes *[]database.TableAttributes, ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateIndexesOnTable", tableName, indexName, inndexes, ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateIndexesOnTable indicates an expected call of CreateIndexesOnTable.
func (mr *MockDatabaseClientMockRecorder) CreateIndexesOnTable(tableName, indexName, inndexes, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateIndexesOnTable", reflect.TypeOf((*MockDatabaseClient)(nil).CreateIndexesOnTable), tableName, indexName, inndexes, ctx)
}

// CreateTable mocks base method.
func (m *MockDatabaseClient) CreateTable(tableName string, keys *[]database.TableAttributes, ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTable", tableName, keys, ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTable indicates an expected call of CreateTable.
func (mr *MockDatabaseClientMockRecorder) CreateTable(tableName, keys, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTable", reflect.TypeOf((*MockDatabaseClient)(nil).CreateTable), tableName, keys, ctx)
}

// GetData mocks base method.
func (m *MockDatabaseClient) GetData(tableName string, key, result any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetData", tableName, key, result)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetData indicates an expected call of GetData.
func (mr *MockDatabaseClientMockRecorder) GetData(tableName, key, result interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetData", reflect.TypeOf((*MockDatabaseClient)(nil).GetData), tableName, key, result)
}

// GetPostsByIds mocks base method.
func (m *MockDatabaseClient) GetPostsByIds(postIds []string) ([]*database.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPostsByIds", postIds)
	ret0, _ := ret[0].([]*database.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPostsByIds indicates an expected call of GetPostsByIds.
func (mr *MockDatabaseClientMockRecorder) GetPostsByIds(postIds interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPostsByIds", reflect.TypeOf((*MockDatabaseClient)(nil).GetPostsByIds), postIds)
}

// GetPostsByIndexUser mocks base method.
func (m *MockDatabaseClient) GetPostsByIndexUser(username, lastPostId, lastPostCreatedAt string, limit int) ([]*database.Post, string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPostsByIndexUser", username, lastPostId, lastPostCreatedAt, limit)
	ret0, _ := ret[0].([]*database.Post)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(string)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// GetPostsByIndexUser indicates an expected call of GetPostsByIndexUser.
func (mr *MockDatabaseClientMockRecorder) GetPostsByIndexUser(username, lastPostId, lastPostCreatedAt, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPostsByIndexUser", reflect.TypeOf((*MockDatabaseClient)(nil).GetPostsByIndexUser), username, lastPostId, lastPostCreatedAt, limit)
}

// IndexExists mocks base method.
func (m *MockDatabaseClient) IndexExists(tableName, indexName string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IndexExists", tableName, indexName)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IndexExists indicates an expected call of IndexExists.
func (mr *MockDatabaseClientMockRecorder) IndexExists(tableName, indexName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IndexExists", reflect.TypeOf((*MockDatabaseClient)(nil).IndexExists), tableName, indexName)
}

// InsertData mocks base method.
func (m *MockDatabaseClient) InsertData(tableName string, attributes any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertData", tableName, attributes)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertData indicates an expected call of InsertData.
func (mr *MockDatabaseClientMockRecorder) InsertData(tableName, attributes interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertData", reflect.TypeOf((*MockDatabaseClient)(nil).InsertData), tableName, attributes)
}

// RemoveData mocks base method.
func (m *MockDatabaseClient) RemoveData(tableName string, key any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveData", tableName, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveData indicates an expected call of RemoveData.
func (mr *MockDatabaseClientMockRecorder) RemoveData(tableName, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveData", reflect.TypeOf((*MockDatabaseClient)(nil).RemoveData), tableName, key)
}

// RemoveMultipleData mocks base method.
func (m *MockDatabaseClient) RemoveMultipleData(tableName string, keys []any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveMultipleData", tableName, keys)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveMultipleData indicates an expected call of RemoveMultipleData.
func (mr *MockDatabaseClientMockRecorder) RemoveMultipleData(tableName, keys interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveMultipleData", reflect.TypeOf((*MockDatabaseClient)(nil).RemoveMultipleData), tableName, keys)
}

// TableExists mocks base method.
func (m *MockDatabaseClient) TableExists(tableName string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TableExists", tableName)
	ret0, _ := ret[0].(bool)
	return ret0
}

// TableExists indicates an expected call of TableExists.
func (mr *MockDatabaseClientMockRecorder) TableExists(tableName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TableExists", reflect.TypeOf((*MockDatabaseClient)(nil).TableExists), tableName)
}
