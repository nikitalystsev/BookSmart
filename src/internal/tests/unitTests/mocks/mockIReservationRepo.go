// Code generated by MockGen. DO NOT EDIT.
// Source: IReservationRepo.go

// Package mock_repositories is a generated GoMock package.
package mock_repositories

import (
	models "BookSmart/internal/models"
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockIReservationRepo is a mock of IReservationRepo interfaces.
type MockIReservationRepo struct {
	ctrl     *gomock.Controller
	recorder *MockIReservationRepoMockRecorder
}

// MockIReservationRepoMockRecorder is the mock recorder for MockIReservationRepo.
type MockIReservationRepoMockRecorder struct {
	mock *MockIReservationRepo
}

// NewMockIReservationRepo creates a new mock instance.
func NewMockIReservationRepo(ctrl *gomock.Controller) *MockIReservationRepo {
	mock := &MockIReservationRepo{ctrl: ctrl}
	mock.recorder = &MockIReservationRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIReservationRepo) EXPECT() *MockIReservationRepoMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockIReservationRepo) Create(ctx context.Context, reservation *models.ReservationModel) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, reservation)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockIReservationRepoMockRecorder) Create(ctx, reservation interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockIReservationRepo)(nil).Create), ctx, reservation)
}

// GetActiveByReaderID mocks base method.
func (m *MockIReservationRepo) GetActiveByReaderID(ctx context.Context, readerID uuid.UUID) ([]*models.ReservationModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActiveByReaderID", ctx, readerID)
	ret0, _ := ret[0].([]*models.ReservationModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetActiveByReaderID indicates an expected call of GetActiveByReaderID.
func (mr *MockIReservationRepoMockRecorder) GetActiveByReaderID(ctx, readerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActiveByReaderID", reflect.TypeOf((*MockIReservationRepo)(nil).GetActiveByReaderID), ctx, readerID)
}

// GetByID mocks base method.
func (m *MockIReservationRepo) GetByID(ctx context.Context, reservationID uuid.UUID) (*models.ReservationModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, reservationID)
	ret0, _ := ret[0].(*models.ReservationModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockIReservationRepoMockRecorder) GetByID(ctx, reservationID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockIReservationRepo)(nil).GetByID), ctx, reservationID)
}

// GetByReaderAndBook mocks base method.
func (m *MockIReservationRepo) GetByReaderAndBook(ctx context.Context, readerID, bookID uuid.UUID) (*models.ReservationModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByReaderAndBook", ctx, readerID, bookID)
	ret0, _ := ret[0].(*models.ReservationModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByReaderAndBook indicates an expected call of GetByReaderAndBook.
func (mr *MockIReservationRepoMockRecorder) GetByReaderAndBook(ctx, readerID, bookID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByReaderAndBook", reflect.TypeOf((*MockIReservationRepo)(nil).GetByReaderAndBook), ctx, readerID, bookID)
}

// GetOverdueByReaderID mocks base method.
func (m *MockIReservationRepo) GetOverdueByReaderID(ctx context.Context, readerID uuid.UUID) ([]*models.ReservationModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOverdueByReaderID", ctx, readerID)
	ret0, _ := ret[0].([]*models.ReservationModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOverdueByReaderID indicates an expected call of GetOverdueByReaderID.
func (mr *MockIReservationRepoMockRecorder) GetOverdueByReaderID(ctx, readerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOverdueByReaderID", reflect.TypeOf((*MockIReservationRepo)(nil).GetOverdueByReaderID), ctx, readerID)
}

// Update mocks base method.
func (m *MockIReservationRepo) Update(ctx context.Context, reservation *models.ReservationModel) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, reservation)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIReservationRepoMockRecorder) Update(ctx, reservation interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIReservationRepo)(nil).Update), ctx, reservation)
}