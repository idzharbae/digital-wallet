// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/idzharbae/digital-wallet/src/internal/repository (interfaces: UserBalanceRepository)
//
// Generated by this command:
//
//	mockgen -destination=repomock/userbalancerepo_mock.go -package=repomock github.com/idzharbae/digital-wallet/src/internal/repository UserBalanceRepository
//

// Package repomock is a generated GoMock package.
package repomock

import (
	context "context"
	reflect "reflect"

	repository "github.com/idzharbae/digital-wallet/src/internal/repository"
	pgx "github.com/jackc/pgx/v5"
	gomock "go.uber.org/mock/gomock"
)

// MockUserBalanceRepository is a mock of UserBalanceRepository interface.
type MockUserBalanceRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserBalanceRepositoryMockRecorder
}

// MockUserBalanceRepositoryMockRecorder is the mock recorder for MockUserBalanceRepository.
type MockUserBalanceRepositoryMockRecorder struct {
	mock *MockUserBalanceRepository
}

// NewMockUserBalanceRepository creates a new mock instance.
func NewMockUserBalanceRepository(ctrl *gomock.Controller) *MockUserBalanceRepository {
	mock := &MockUserBalanceRepository{ctrl: ctrl}
	mock.recorder = &MockUserBalanceRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserBalanceRepository) EXPECT() *MockUserBalanceRepositoryMockRecorder {
	return m.recorder
}

// CreateUserBalance mocks base method.
func (m *MockUserBalanceRepository) CreateUserBalance(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUserBalance", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUserBalance indicates an expected call of CreateUserBalance.
func (mr *MockUserBalanceRepositoryMockRecorder) CreateUserBalance(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUserBalance", reflect.TypeOf((*MockUserBalanceRepository)(nil).CreateUserBalance), arg0, arg1)
}

// GetUserBalance mocks base method.
func (m *MockUserBalanceRepository) GetUserBalance(arg0 context.Context, arg1 string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserBalance", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserBalance indicates an expected call of GetUserBalance.
func (mr *MockUserBalanceRepositoryMockRecorder) GetUserBalance(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserBalance", reflect.TypeOf((*MockUserBalanceRepository)(nil).GetUserBalance), arg0, arg1)
}

// GetUserBalanceForUpdate mocks base method.
func (m *MockUserBalanceRepository) GetUserBalanceForUpdate(arg0 context.Context, arg1 string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserBalanceForUpdate", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserBalanceForUpdate indicates an expected call of GetUserBalanceForUpdate.
func (mr *MockUserBalanceRepositoryMockRecorder) GetUserBalanceForUpdate(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserBalanceForUpdate", reflect.TypeOf((*MockUserBalanceRepository)(nil).GetUserBalanceForUpdate), arg0, arg1)
}

// UpdateBalance mocks base method.
func (m *MockUserBalanceRepository) UpdateBalance(arg0 context.Context, arg1 string, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBalance", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateBalance indicates an expected call of UpdateBalance.
func (mr *MockUserBalanceRepositoryMockRecorder) UpdateBalance(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBalance", reflect.TypeOf((*MockUserBalanceRepository)(nil).UpdateBalance), arg0, arg1, arg2)
}

// WithTransaction mocks base method.
func (m *MockUserBalanceRepository) WithTransaction(arg0 pgx.Tx) repository.UserBalanceRepository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithTransaction", arg0)
	ret0, _ := ret[0].(repository.UserBalanceRepository)
	return ret0
}

// WithTransaction indicates an expected call of WithTransaction.
func (mr *MockUserBalanceRepositoryMockRecorder) WithTransaction(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithTransaction", reflect.TypeOf((*MockUserBalanceRepository)(nil).WithTransaction), arg0)
}
