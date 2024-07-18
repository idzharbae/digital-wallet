// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/idzharbae/digital-wallet/src/internal/repository (interfaces: UserTransactionRepository)
//
// Generated by this command:
//
//	mockgen -destination=repomock/usertransaction_mock.go -package=repomock github.com/idzharbae/digital-wallet/src/internal/repository UserTransactionRepository
//

// Package repomock is a generated GoMock package.
package repomock

import (
	context "context"
	reflect "reflect"

	entity "github.com/idzharbae/digital-wallet/src/internal/entity"
	repository "github.com/idzharbae/digital-wallet/src/internal/repository"
	pgx "github.com/jackc/pgx/v5"
	gomock "go.uber.org/mock/gomock"
)

// MockUserTransactionRepository is a mock of UserTransactionRepository interface.
type MockUserTransactionRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserTransactionRepositoryMockRecorder
}

// MockUserTransactionRepositoryMockRecorder is the mock recorder for MockUserTransactionRepository.
type MockUserTransactionRepositoryMockRecorder struct {
	mock *MockUserTransactionRepository
}

// NewMockUserTransactionRepository creates a new mock instance.
func NewMockUserTransactionRepository(ctrl *gomock.Controller) *MockUserTransactionRepository {
	mock := &MockUserTransactionRepository{ctrl: ctrl}
	mock.recorder = &MockUserTransactionRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserTransactionRepository) EXPECT() *MockUserTransactionRepositoryMockRecorder {
	return m.recorder
}

// GetTopTransactingUsers mocks base method.
func (m *MockUserTransactionRepository) GetTopTransactingUsers(arg0 context.Context) ([]entity.TotalDebit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTopTransactingUsers", arg0)
	ret0, _ := ret[0].([]entity.TotalDebit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTopTransactingUsers indicates an expected call of GetTopTransactingUsers.
func (mr *MockUserTransactionRepositoryMockRecorder) GetTopTransactingUsers(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTopTransactingUsers", reflect.TypeOf((*MockUserTransactionRepository)(nil).GetTopTransactingUsers), arg0)
}

// GetUserTopTransactions mocks base method.
func (m *MockUserTransactionRepository) GetUserTopTransactions(arg0 context.Context, arg1 string) ([]entity.UserTransaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserTopTransactions", arg0, arg1)
	ret0, _ := ret[0].([]entity.UserTransaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserTopTransactions indicates an expected call of GetUserTopTransactions.
func (mr *MockUserTransactionRepositoryMockRecorder) GetUserTopTransactions(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserTopTransactions", reflect.TypeOf((*MockUserTransactionRepository)(nil).GetUserTopTransactions), arg0, arg1)
}

// InsertTransaction mocks base method.
func (m *MockUserTransactionRepository) InsertTransaction(arg0 context.Context, arg1, arg2 string, arg3 entity.TransactionType, arg4 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertTransaction", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertTransaction indicates an expected call of InsertTransaction.
func (mr *MockUserTransactionRepositoryMockRecorder) InsertTransaction(arg0, arg1, arg2, arg3, arg4 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertTransaction", reflect.TypeOf((*MockUserTransactionRepository)(nil).InsertTransaction), arg0, arg1, arg2, arg3, arg4)
}

// RefreshTopTransactingUsers mocks base method.
func (m *MockUserTransactionRepository) RefreshTopTransactingUsers(arg0 context.Context) ([]entity.TotalDebit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshTopTransactingUsers", arg0)
	ret0, _ := ret[0].([]entity.TotalDebit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefreshTopTransactingUsers indicates an expected call of RefreshTopTransactingUsers.
func (mr *MockUserTransactionRepositoryMockRecorder) RefreshTopTransactingUsers(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshTopTransactingUsers", reflect.TypeOf((*MockUserTransactionRepository)(nil).RefreshTopTransactingUsers), arg0)
}

// UpsertTotalDebit mocks base method.
func (m *MockUserTransactionRepository) UpsertTotalDebit(arg0 context.Context, arg1 string, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertTotalDebit", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertTotalDebit indicates an expected call of UpsertTotalDebit.
func (mr *MockUserTransactionRepositoryMockRecorder) UpsertTotalDebit(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertTotalDebit", reflect.TypeOf((*MockUserTransactionRepository)(nil).UpsertTotalDebit), arg0, arg1, arg2)
}

// WithTransaction mocks base method.
func (m *MockUserTransactionRepository) WithTransaction(arg0 pgx.Tx) repository.UserTransactionRepository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithTransaction", arg0)
	ret0, _ := ret[0].(repository.UserTransactionRepository)
	return ret0
}

// WithTransaction indicates an expected call of WithTransaction.
func (mr *MockUserTransactionRepositoryMockRecorder) WithTransaction(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithTransaction", reflect.TypeOf((*MockUserTransactionRepository)(nil).WithTransaction), arg0)
}