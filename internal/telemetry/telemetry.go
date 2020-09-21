package telemetry

import (
	"github.com/NodeFactoryIo/vedran-daemon/internal/lb"
	"github.com/NodeFactoryIo/vedran-daemon/internal/node"
	"github.com/NodeFactoryIo/vedran-daemon/internal/scheduler"
	log "github.com/sirupsen/logrus"
)

const (
	metricsSendInterval = 30
	pingSendInterval    = 5
)

// Telemetry is used for scheduling telemetry reports
type Telemetry interface {
	StartSendingTelemetry(scheduler scheduler.Scheduler, lbClient *lb.Client, nodeClient node.Client)
}

// NewTelemetry returns instance of Telemetry
func NewTelemetry() Telemetry {
	return &telemetry{}
}

type telemetry struct{}

// StartSendingTelemetry start sending ping and metrics to load balancer and blocks current thread
func (t *telemetry) StartSendingTelemetry(scheduler scheduler.Scheduler, lbClient *lb.Client, nodeClient node.Client) {
	_, _ = scheduler.Every(metricsSendInterval).Seconds().Do(lbClient.Metrics.Send, nodeClient)
	log.Info("Started sending metrics to load balancer")

	_, _ = scheduler.Every(pingSendInterval).Seconds().Do(lbClient.Ping.Send)
	log.Info("Started sending pings to load balancer")

	scheduler.StartBlocking()
}
