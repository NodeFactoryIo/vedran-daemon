// Code generated by mockery v2.2.1. DO NOT EDIT.

package mocks

import (
	metrics "github.com/NodeFactoryIo/vedran-daemon/internal/metrics"
	mock "github.com/stretchr/testify/mock"
)

// FetchMetrics is an autogenerated mock type for the FetchMetrics type
type FetchMetrics struct {
	mock.Mock
}

// GetNodeMetrics provides a mock function with given fields:
func (_m *FetchMetrics) GetNodeMetrics() (*metrics.Metrics, error) {
	ret := _m.Called()

	var r0 *metrics.Metrics
	if rf, ok := ret.Get(0).(func() *metrics.Metrics); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*metrics.Metrics)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
