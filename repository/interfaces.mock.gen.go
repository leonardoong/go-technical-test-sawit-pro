// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package repository is a generated GoMock package.
package repository

import (
	context "context"
	reflect "reflect"

	model "github.com/SawitProRecruitment/UserService/model"
	gomock "github.com/golang/mock/gomock"
)

// MockRepositoryInterface is a mock of RepositoryInterface interface.
type MockRepositoryInterface struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryInterfaceMockRecorder
}

// MockRepositoryInterfaceMockRecorder is the mock recorder for MockRepositoryInterface.
type MockRepositoryInterfaceMockRecorder struct {
	mock *MockRepositoryInterface
}

// NewMockRepositoryInterface creates a new mock instance.
func NewMockRepositoryInterface(ctrl *gomock.Controller) *MockRepositoryInterface {
	mock := &MockRepositoryInterface{ctrl: ctrl}
	mock.recorder = &MockRepositoryInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepositoryInterface) EXPECT() *MockRepositoryInterfaceMockRecorder {
	return m.recorder
}

// GetLoginData mocks base method.
func (m *MockRepositoryInterface) GetLoginData(ctx context.Context, input GetLoginDataInput) (GetLoginDataOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLoginData", ctx, input)
	ret0, _ := ret[0].(GetLoginDataOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLoginData indicates an expected call of GetLoginData.
func (mr *MockRepositoryInterfaceMockRecorder) GetLoginData(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLoginData", reflect.TypeOf((*MockRepositoryInterface)(nil).GetLoginData), ctx, input)
}

// GetUserDataByUserID mocks base method.
func (m *MockRepositoryInterface) GetUserDataByUserID(ctx context.Context, input GetUserDataByUserIDInput) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserDataByUserID", ctx, input)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserDataByUserID indicates an expected call of GetUserDataByUserID.
func (mr *MockRepositoryInterfaceMockRecorder) GetUserDataByUserID(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserDataByUserID", reflect.TypeOf((*MockRepositoryInterface)(nil).GetUserDataByUserID), ctx, input)
}

// InsertUser mocks base method.
func (m *MockRepositoryInterface) InsertUser(ctx context.Context, in InsertUserInput) (InsertUserOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertUser", ctx, in)
	ret0, _ := ret[0].(InsertUserOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertUser indicates an expected call of InsertUser.
func (mr *MockRepositoryInterfaceMockRecorder) InsertUser(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertUser", reflect.TypeOf((*MockRepositoryInterface)(nil).InsertUser), ctx, in)
}

// UpdateSuccessfulLogin mocks base method.
func (m *MockRepositoryInterface) UpdateSuccessfulLogin(ctx context.Context, in UpdateSuccessfulLoginInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSuccessfulLogin", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateSuccessfulLogin indicates an expected call of UpdateSuccessfulLogin.
func (mr *MockRepositoryInterfaceMockRecorder) UpdateSuccessfulLogin(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSuccessfulLogin", reflect.TypeOf((*MockRepositoryInterface)(nil).UpdateSuccessfulLogin), ctx, in)
}

// UpdateUserData mocks base method.
func (m *MockRepositoryInterface) UpdateUserData(ctx context.Context, in UpdateUserDataInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserData", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserData indicates an expected call of UpdateUserData.
func (mr *MockRepositoryInterfaceMockRecorder) UpdateUserData(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserData", reflect.TypeOf((*MockRepositoryInterface)(nil).UpdateUserData), ctx, in)
}