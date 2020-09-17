// Code generated by mockery v2.2.1. DO NOT EDIT.

package mocks

import (
	lb "github.com/NodeFactoryIo/vedran-daemon/internal/lb"
	metrics "github.com/NodeFactoryIo/vedran-daemon/internal/metrics"

	mock "github.com/stretchr/testify/mock"

	scheduler "github.com/NodeFactoryIo/vedran-daemon/internal/scheduler"
)

// TelemetryInterface is an autogenerated mock type for the TelemetryInterface type
type TelemetryInterface struct {
	mock.Mock
}

// StartSendingTelemetry provides a mock function with given fields: _a0, lbClient, metricsClient
func (_m *TelemetryInterface) StartSendingTelemetry(_a0 scheduler.Scheduler, lbClient *lb.Client, metricsClient metrics.Client) {
	_m.Called(_a0, lbClient, metricsClient)
}
