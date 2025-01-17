// Code generated by MockGen. DO NOT EDIT.
// Source: flux-meter.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	selectorv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/selector/v1"
	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	iface "github.com/fluxninja/aperture/pkg/policies/dataplane/iface"
	gomock "github.com/golang/mock/gomock"
	prometheus "github.com/prometheus/client_golang/prometheus"
)

// MockFluxMeter is a mock of FluxMeter interface.
type MockFluxMeter struct {
	ctrl     *gomock.Controller
	recorder *MockFluxMeterMockRecorder
}

// MockFluxMeterMockRecorder is the mock recorder for MockFluxMeter.
type MockFluxMeterMockRecorder struct {
	mock *MockFluxMeter
}

// NewMockFluxMeter creates a new mock instance.
func NewMockFluxMeter(ctrl *gomock.Controller) *MockFluxMeter {
	mock := &MockFluxMeter{ctrl: ctrl}
	mock.recorder = &MockFluxMeterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFluxMeter) EXPECT() *MockFluxMeterMockRecorder {
	return m.recorder
}

// GetAttributeKey mocks base method.
func (m *MockFluxMeter) GetAttributeKey() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAttributeKey")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetAttributeKey indicates an expected call of GetAttributeKey.
func (mr *MockFluxMeterMockRecorder) GetAttributeKey() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAttributeKey", reflect.TypeOf((*MockFluxMeter)(nil).GetAttributeKey))
}

// GetFluxMeterID mocks base method.
func (m *MockFluxMeter) GetFluxMeterID() iface.FluxMeterID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFluxMeterID")
	ret0, _ := ret[0].(iface.FluxMeterID)
	return ret0
}

// GetFluxMeterID indicates an expected call of GetFluxMeterID.
func (mr *MockFluxMeterMockRecorder) GetFluxMeterID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFluxMeterID", reflect.TypeOf((*MockFluxMeter)(nil).GetFluxMeterID))
}

// GetFluxMeterName mocks base method.
func (m *MockFluxMeter) GetFluxMeterName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFluxMeterName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetFluxMeterName indicates an expected call of GetFluxMeterName.
func (mr *MockFluxMeterMockRecorder) GetFluxMeterName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFluxMeterName", reflect.TypeOf((*MockFluxMeter)(nil).GetFluxMeterName))
}

// GetHistogram mocks base method.
func (m *MockFluxMeter) GetHistogram(decisionType flowcontrolv1.CheckResponse_DecisionType, statusCode, featureStatus string) prometheus.Observer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHistogram", decisionType, statusCode, featureStatus)
	ret0, _ := ret[0].(prometheus.Observer)
	return ret0
}

// GetHistogram indicates an expected call of GetHistogram.
func (mr *MockFluxMeterMockRecorder) GetHistogram(decisionType, statusCode, featureStatus interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHistogram", reflect.TypeOf((*MockFluxMeter)(nil).GetHistogram), decisionType, statusCode, featureStatus)
}

// GetSelector mocks base method.
func (m *MockFluxMeter) GetSelector() *selectorv1.Selector {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSelector")
	ret0, _ := ret[0].(*selectorv1.Selector)
	return ret0
}

// GetSelector indicates an expected call of GetSelector.
func (mr *MockFluxMeterMockRecorder) GetSelector() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSelector", reflect.TypeOf((*MockFluxMeter)(nil).GetSelector))
}
