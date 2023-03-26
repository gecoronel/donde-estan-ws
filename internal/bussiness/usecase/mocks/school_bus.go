// Code generated by MockGen. DO NOT EDIT.
// Source: school_bus.go

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	reflect "reflect"

	gateway "github.com/gecoronel/donde-estan-ws/internal/bussiness/gateway"
	model "github.com/gecoronel/donde-estan-ws/internal/bussiness/model"
	gomock "github.com/golang/mock/gomock"
)

// MockSchoolBusUseCase is a mock of SchoolBusUseCase interface.
type MockSchoolBusUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockSchoolBusUseCaseMockRecorder
}

// MockSchoolBusUseCaseMockRecorder is the mock recorder for MockSchoolBusUseCase.
type MockSchoolBusUseCaseMockRecorder struct {
	mock *MockSchoolBusUseCase
}

// NewMockSchoolBusUseCase creates a new mock instance.
func NewMockSchoolBusUseCase(ctrl *gomock.Controller) *MockSchoolBusUseCase {
	mock := &MockSchoolBusUseCase{ctrl: ctrl}
	mock.recorder = &MockSchoolBusUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSchoolBusUseCase) EXPECT() *MockSchoolBusUseCaseMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockSchoolBusUseCase) Delete(arg0 uint64, arg1 gateway.ServiceLocator) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockSchoolBusUseCaseMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockSchoolBusUseCase)(nil).Delete), arg0, arg1)
}

// Get mocks base method.
func (m *MockSchoolBusUseCase) Get(arg0 uint64, arg1 gateway.ServiceLocator) (*model.SchoolBus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*model.SchoolBus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockSchoolBusUseCaseMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockSchoolBusUseCase)(nil).Get), arg0, arg1)
}

// Save mocks base method.
func (m *MockSchoolBusUseCase) Save(arg0 model.SchoolBus, arg1 gateway.ServiceLocator) (*model.SchoolBus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0, arg1)
	ret0, _ := ret[0].(*model.SchoolBus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Save indicates an expected call of Save.
func (mr *MockSchoolBusUseCaseMockRecorder) Save(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockSchoolBusUseCase)(nil).Save), arg0, arg1)
}

// Update mocks base method.
func (m *MockSchoolBusUseCase) Update(arg0 model.SchoolBus, arg1 gateway.ServiceLocator) (*model.SchoolBus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(*model.SchoolBus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockSchoolBusUseCaseMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockSchoolBusUseCase)(nil).Update), arg0, arg1)
}
