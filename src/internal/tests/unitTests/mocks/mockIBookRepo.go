// Code generated by MockGen. DO NOT EDIT.
// Source: IBookRepo.go

// Package mock_repositories is a generated GoMock package.
package mock_repositories

import (
	models "BookSmart/internal/models"
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockIBookRepo is a mock of IBookRepo interfaces.
type MockIBookRepo struct {
	ctrl     *gomock.Controller
	recorder *MockIBookRepoMockRecorder
}

// MockIBookRepoMockRecorder is the mock recorder for MockIBookRepo.
type MockIBookRepoMockRecorder struct {
	mock *MockIBookRepo
}

// NewMockIBookRepo creates a new mock instance.
func NewMockIBookRepo(ctrl *gomock.Controller) *MockIBookRepo {
	mock := &MockIBookRepo{ctrl: ctrl}
	mock.recorder = &MockIBookRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIBookRepo) EXPECT() *MockIBookRepoMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockIBookRepo) Create(ctx context.Context, book *models.BookModel) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, book)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockIBookRepoMockRecorder) Create(ctx, book interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockIBookRepo)(nil).Create), ctx, book)
}

// DeleteByTitle mocks base method.
func (m *MockIBookRepo) DeleteByTitle(ctx context.Context, title string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByTitle", ctx, title)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByTitle indicates an expected call of DeleteByTitle.
func (mr *MockIBookRepoMockRecorder) DeleteByTitle(ctx, title interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByTitle", reflect.TypeOf((*MockIBookRepo)(nil).DeleteByTitle), ctx, title)
}

// GetByID mocks base method.
func (m *MockIBookRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.BookModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*models.BookModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockIBookRepoMockRecorder) GetByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockIBookRepo)(nil).GetByID), ctx, id)
}

// GetByTitle mocks base method.
func (m *MockIBookRepo) GetByTitle(ctx context.Context, title string) (*models.BookModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByTitle", ctx, title)
	ret0, _ := ret[0].(*models.BookModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByTitle indicates an expected call of GetByTitle.
func (mr *MockIBookRepoMockRecorder) GetByTitle(ctx, title interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByTitle", reflect.TypeOf((*MockIBookRepo)(nil).GetByTitle), ctx, title)
}

// Update mocks base method.
func (m *MockIBookRepo) Update(ctx context.Context, book *models.BookModel) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, book)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIBookRepoMockRecorder) Update(ctx, book interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIBookRepo)(nil).Update), ctx, book)
}