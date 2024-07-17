// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/idzharbae/digital-wallet/src/internal/usecase (interfaces: UserUC)
//
// Generated by this command:
//
//	mockgen -destination=ucmock/useruc_mock.go -package=ucmock github.com/idzharbae/digital-wallet/src/internal/usecase UserUC
//

// Package ucmock is a generated GoMock package.
package ucmock

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockUserUC is a mock of UserUC interface.
type MockUserUC struct {
	ctrl     *gomock.Controller
	recorder *MockUserUCMockRecorder
}

// MockUserUCMockRecorder is the mock recorder for MockUserUC.
type MockUserUCMockRecorder struct {
	mock *MockUserUC
}

// NewMockUserUC creates a new mock instance.
func NewMockUserUC(ctrl *gomock.Controller) *MockUserUC {
	mock := &MockUserUC{ctrl: ctrl}
	mock.recorder = &MockUserUCMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserUC) EXPECT() *MockUserUCMockRecorder {
	return m.recorder
}

// RegisterUser mocks base method.
func (m *MockUserUC) RegisterUser(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterUser", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterUser indicates an expected call of RegisterUser.
func (mr *MockUserUCMockRecorder) RegisterUser(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterUser", reflect.TypeOf((*MockUserUC)(nil).RegisterUser), arg0)
}