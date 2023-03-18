// Code generated by MockGen. DO NOT EDIT.
// Source: user.go

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	reflect "reflect"

	gateway "github.com/gecoronel/donde-estan-ws/internal/bussiness/gateway"
	model "github.com/gecoronel/donde-estan-ws/internal/bussiness/model"
	gomock "github.com/golang/mock/gomock"
)

// MockUserUseCase is a mock of UserUseCase interface.
type MockUserUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockUserUseCaseMockRecorder
}

// MockUserUseCaseMockRecorder is the mock recorder for MockUserUseCase.
type MockUserUseCaseMockRecorder struct {
	mock *MockUserUseCase
}

// NewMockUserUseCase creates a new mock instance.
func NewMockUserUseCase(ctrl *gomock.Controller) *MockUserUseCase {
	mock := &MockUserUseCase{ctrl: ctrl}
	mock.recorder = &MockUserUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserUseCase) EXPECT() *MockUserUseCaseMockRecorder {
	return m.recorder
}

// AddObservedUserInObserverUser mocks base method.
func (m *MockUserUseCase) AddObservedUserInObserverUser(arg0 string, arg1 uint64, arg2 gateway.ServiceLocator) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddObservedUserInObserverUser", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddObservedUserInObserverUser indicates an expected call of AddObservedUserInObserverUser.
func (mr *MockUserUseCaseMockRecorder) AddObservedUserInObserverUser(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddObservedUserInObserverUser", reflect.TypeOf((*MockUserUseCase)(nil).AddObservedUserInObserverUser), arg0, arg1, arg2)
}

// CreateObservedUser mocks base method.
func (m *MockUserUseCase) CreateObservedUser(arg0 model.ObservedUser, arg1 gateway.ServiceLocator) (*model.ObservedUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateObservedUser", arg0, arg1)
	ret0, _ := ret[0].(*model.ObservedUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateObservedUser indicates an expected call of CreateObservedUser.
func (mr *MockUserUseCaseMockRecorder) CreateObservedUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateObservedUser", reflect.TypeOf((*MockUserUseCase)(nil).CreateObservedUser), arg0, arg1)
}

// CreateObserverUser mocks base method.
func (m *MockUserUseCase) CreateObserverUser(arg0 model.ObserverUser, arg1 gateway.ServiceLocator) (*model.ObserverUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateObserverUser", arg0, arg1)
	ret0, _ := ret[0].(*model.ObserverUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateObserverUser indicates an expected call of CreateObserverUser.
func (mr *MockUserUseCaseMockRecorder) CreateObserverUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateObserverUser", reflect.TypeOf((*MockUserUseCase)(nil).CreateObserverUser), arg0, arg1)
}

// DeleteObservedUser mocks base method.
func (m *MockUserUseCase) DeleteObservedUser(arg0 uint64, arg1 gateway.ServiceLocator) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteObservedUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteObservedUser indicates an expected call of DeleteObservedUser.
func (mr *MockUserUseCaseMockRecorder) DeleteObservedUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteObservedUser", reflect.TypeOf((*MockUserUseCase)(nil).DeleteObservedUser), arg0, arg1)
}

// DeleteObservedUserInObserverUser mocks base method.
func (m *MockUserUseCase) DeleteObservedUserInObserverUser(arg0, arg1 uint64, arg2 gateway.ServiceLocator) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteObservedUserInObserverUser", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteObservedUserInObserverUser indicates an expected call of DeleteObservedUserInObserverUser.
func (mr *MockUserUseCaseMockRecorder) DeleteObservedUserInObserverUser(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteObservedUserInObserverUser", reflect.TypeOf((*MockUserUseCase)(nil).DeleteObservedUserInObserverUser), arg0, arg1, arg2)
}

// DeleteObserverUser mocks base method.
func (m *MockUserUseCase) DeleteObserverUser(arg0 uint64, arg1 gateway.ServiceLocator) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteObserverUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteObserverUser indicates an expected call of DeleteObserverUser.
func (mr *MockUserUseCaseMockRecorder) DeleteObserverUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteObserverUser", reflect.TypeOf((*MockUserUseCase)(nil).DeleteObserverUser), arg0, arg1)
}

// FindByEmail mocks base method.
func (m *MockUserUseCase) FindByEmail(arg0 string, arg1 gateway.ServiceLocator) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByEmail", arg0, arg1)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByEmail indicates an expected call of FindByEmail.
func (mr *MockUserUseCaseMockRecorder) FindByEmail(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByEmail", reflect.TypeOf((*MockUserUseCase)(nil).FindByEmail), arg0, arg1)
}

// FindByUsername mocks base method.
func (m *MockUserUseCase) FindByUsername(arg0 string, arg1 gateway.ServiceLocator) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByUsername", arg0, arg1)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByUsername indicates an expected call of FindByUsername.
func (mr *MockUserUseCaseMockRecorder) FindByUsername(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByUsername", reflect.TypeOf((*MockUserUseCase)(nil).FindByUsername), arg0, arg1)
}

// Get mocks base method.
func (m *MockUserUseCase) Get(arg0 uint64, arg1 gateway.ServiceLocator) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockUserUseCaseMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockUserUseCase)(nil).Get), arg0, arg1)
}

// Login mocks base method.
func (m *MockUserUseCase) Login(arg0 model.Login, arg1 gateway.ServiceLocator) (model.IUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", arg0, arg1)
	ret0, _ := ret[0].(model.IUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockUserUseCaseMockRecorder) Login(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockUserUseCase)(nil).Login), arg0, arg1)
}

// UpdateObservedUser mocks base method.
func (m *MockUserUseCase) UpdateObservedUser(arg0 model.ObservedUser, arg1 gateway.ServiceLocator) (*model.ObservedUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateObservedUser", arg0, arg1)
	ret0, _ := ret[0].(*model.ObservedUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateObservedUser indicates an expected call of UpdateObservedUser.
func (mr *MockUserUseCaseMockRecorder) UpdateObservedUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateObservedUser", reflect.TypeOf((*MockUserUseCase)(nil).UpdateObservedUser), arg0, arg1)
}

// UpdateObserverUser mocks base method.
func (m *MockUserUseCase) UpdateObserverUser(arg0 model.ObserverUser, arg1 gateway.ServiceLocator) (*model.ObserverUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateObserverUser", arg0, arg1)
	ret0, _ := ret[0].(*model.ObserverUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateObserverUser indicates an expected call of UpdateObserverUser.
func (mr *MockUserUseCaseMockRecorder) UpdateObserverUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateObserverUser", reflect.TypeOf((*MockUserUseCase)(nil).UpdateObserverUser), arg0, arg1)
}
