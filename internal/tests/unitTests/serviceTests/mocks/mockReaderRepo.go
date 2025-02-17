// Code generated by MockGen. DO NOT EDIT.
// Source: ./components/component-services/intfRepo/IReaderRepo.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	models "github.com/nikitalystsev/BookSmart-services/core/models"
)

// MockIReaderRepo is a mock of IReaderRepo interface.
type MockIReaderRepo struct {
	ctrl     *gomock.Controller
	recorder *MockIReaderRepoMockRecorder
}

// MockIReaderRepoMockRecorder is the mock recorder for MockIReaderRepo.
type MockIReaderRepoMockRecorder struct {
	mock *MockIReaderRepo
}

// NewMockIReaderRepo creates a new mock instance.
func NewMockIReaderRepo(ctrl *gomock.Controller) *MockIReaderRepo {
	mock := &MockIReaderRepo{ctrl: ctrl}
	mock.recorder = &MockIReaderRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIReaderRepo) EXPECT() *MockIReaderRepoMockRecorder {
	return m.recorder
}

// AddToFavorites mocks base method.
func (m *MockIReaderRepo) AddToFavorites(ctx context.Context, readerID, bookID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddToFavorites", ctx, readerID, bookID)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddToFavorites indicates an expected call of AddToFavorites.
func (mr *MockIReaderRepoMockRecorder) AddToFavorites(ctx, readerID, bookID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddToFavorites", reflect.TypeOf((*MockIReaderRepo)(nil).AddToFavorites), ctx, readerID, bookID)
}

// Create mocks base method.
func (m *MockIReaderRepo) Create(ctx context.Context, reader *models.ReaderModel) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, reader)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockIReaderRepoMockRecorder) Create(ctx, reader interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockIReaderRepo)(nil).Create), ctx, reader)
}

// GetByID mocks base method.
func (m *MockIReaderRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.ReaderModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*models.ReaderModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockIReaderRepoMockRecorder) GetByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockIReaderRepo)(nil).GetByID), ctx, id)
}

// GetByPhoneNumber mocks base method.
func (m *MockIReaderRepo) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*models.ReaderModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByPhoneNumber", ctx, phoneNumber)
	ret0, _ := ret[0].(*models.ReaderModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByPhoneNumber indicates an expected call of GetByPhoneNumber.
func (mr *MockIReaderRepoMockRecorder) GetByPhoneNumber(ctx, phoneNumber interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByPhoneNumber", reflect.TypeOf((*MockIReaderRepo)(nil).GetByPhoneNumber), ctx, phoneNumber)
}

// GetByRefreshToken mocks base method.
func (m *MockIReaderRepo) GetByRefreshToken(ctx context.Context, token string) (*models.ReaderModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByRefreshToken", ctx, token)
	ret0, _ := ret[0].(*models.ReaderModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByRefreshToken indicates an expected call of GetByRefreshToken.
func (mr *MockIReaderRepoMockRecorder) GetByRefreshToken(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByRefreshToken", reflect.TypeOf((*MockIReaderRepo)(nil).GetByRefreshToken), ctx, token)
}

// IsFavorite mocks base method.
func (m *MockIReaderRepo) IsFavorite(ctx context.Context, readerID, bookID uuid.UUID) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsFavorite", ctx, readerID, bookID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsFavorite indicates an expected call of IsFavorite.
func (mr *MockIReaderRepoMockRecorder) IsFavorite(ctx, readerID, bookID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsFavorite", reflect.TypeOf((*MockIReaderRepo)(nil).IsFavorite), ctx, readerID, bookID)
}

// SaveRefreshToken mocks base method.
func (m *MockIReaderRepo) SaveRefreshToken(ctx context.Context, id uuid.UUID, token string, ttl time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveRefreshToken", ctx, id, token, ttl)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveRefreshToken indicates an expected call of SaveRefreshToken.
func (mr *MockIReaderRepoMockRecorder) SaveRefreshToken(ctx, id, token, ttl interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveRefreshToken", reflect.TypeOf((*MockIReaderRepo)(nil).SaveRefreshToken), ctx, id, token, ttl)
}
