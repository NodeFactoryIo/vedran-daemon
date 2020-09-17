package telemetry

import (
	"log"

	"github.com/NodeFactoryIo/vedran-daemon/internal/lb"
	"github.com/NodeFactoryIo/vedran-daemon/internal/node"
	"github.com/NodeFactoryIo/vedran-daemon/internal/scheduler"
)

const (
	metricsSendInterval = 30
	pingSendInterval    = 5
)

type TelemetryInterface interface {
	StartSendingTelemetry(scheduler scheduler.Scheduler, lbClient *lb.Client, nodeClient node.Client)
}

type Telemetry struct{}

// StartSendingTelemetry start sending ping and metrics to load balancer and blocks current thread
func (t *Telemetry) StartSendingTelemetry(scheduler scheduler.Scheduler, lbClient *lb.Client, nodeClient node.Client) {
	_, _ = scheduler.Every(metricsSendInterval).Seconds().Do(lbClient.Metrics.Send, nodeClient)
	log.Println("Started sending metrics to load balancer")

	_, _ = scheduler.Every(pingSendInterval).Seconds().Do(lbClient.Ping.Send)
	log.Println("Started sending pings to load balancer")

	scheduler.StartBlocking()
}
