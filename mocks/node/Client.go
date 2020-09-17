// Code generated by mockery v2.2.1. DO NOT EDIT.

package mocks

import (
	hash "hash"

	node "github.com/NodeFactoryIo/vedran-daemon/internal/node"
	mock "github.com/stretchr/testify/mock"
)

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

// GetConfigHash provides a mock function with given fields:
func (_m *Client) GetConfigHash() (hash.Hash32, error) {
	ret := _m.Called()

	var r0 hash.Hash32
	if rf, ok := ret.Get(0).(func() hash.Hash32); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(hash.Hash32)
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

// GetMetrics provides a mock function with given fields:
func (_m *Client) GetMetrics() (*node.Metrics, error) {
	ret := _m.Called()

	var r0 *node.Metrics
	if rf, ok := ret.Get(0).(func() *node.Metrics); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*node.Metrics)
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

// GetMetricsURL provides a mock function with given fields:
func (_m *Client) GetMetricsURL() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetRPCURL provides a mock function with given fields:
func (_m *Client) GetRPCURL() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}