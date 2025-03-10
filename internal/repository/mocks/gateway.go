// Code generated by MockGen. DO NOT EDIT.
// Source: gateway.go

// Package mocks is a generated GoMock package.
package mocks

import (
	models "payment-gateway/internal/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockGatewayRepository is a mock of GatewayRepository interface.
type MockGatewayRepository struct {
	ctrl     *gomock.Controller
	recorder *MockGatewayRepositoryMockRecorder
}

// MockGatewayRepositoryMockRecorder is the mock recorder for MockGatewayRepository.
type MockGatewayRepositoryMockRecorder struct {
	mock *MockGatewayRepository
}

// NewMockGatewayRepository creates a new mock instance.
func NewMockGatewayRepository(ctrl *gomock.Controller) *MockGatewayRepository {
	mock := &MockGatewayRepository{ctrl: ctrl}
	mock.recorder = &MockGatewayRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGatewayRepository) EXPECT() *MockGatewayRepositoryMockRecorder {
	return m.recorder
}

// CreateGateway mocks base method.
func (m *MockGatewayRepository) CreateGateway(gateway models.Gateway) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateGateway", gateway)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateGateway indicates an expected call of CreateGateway.
func (mr *MockGatewayRepositoryMockRecorder) CreateGateway(gateway interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGateway", reflect.TypeOf((*MockGatewayRepository)(nil).CreateGateway), gateway)
}

// GetAvailableGateways mocks base method.
func (m *MockGatewayRepository) GetAvailableGateways(countryID int) ([]models.Gateway, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAvailableGateways", countryID)
	ret0, _ := ret[0].([]models.Gateway)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAvailableGateways indicates an expected call of GetAvailableGateways.
func (mr *MockGatewayRepositoryMockRecorder) GetAvailableGateways(countryID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAvailableGateways", reflect.TypeOf((*MockGatewayRepository)(nil).GetAvailableGateways), countryID)
}

// GetGateways mocks base method.
func (m *MockGatewayRepository) GetGateways() ([]models.Gateway, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGateways")
	ret0, _ := ret[0].([]models.Gateway)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGateways indicates an expected call of GetGateways.
func (mr *MockGatewayRepositoryMockRecorder) GetGateways() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGateways", reflect.TypeOf((*MockGatewayRepository)(nil).GetGateways))
}
