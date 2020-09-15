// Code generated by mockery v2.2.1. DO NOT EDIT.

package mocks

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// PingService is an autogenerated mock type for the PingService type
type PingService struct {
	mock.Mock
}

// Send provides a mock function with given fields:
func (_m *PingService) Send() (*http.Response, error) {
	ret := _m.Called()

	var r0 *http.Response
	if rf, ok := ret.Get(0).(func() *http.Response); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Response)
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
