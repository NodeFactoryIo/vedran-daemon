package telemetry

import (
	"log"

	"github.com/NodeFactoryIo/vedran-daemon/internal/lb"
	"github.com/NodeFactoryIo/vedran-daemon/internal/metrics"
	"github.com/NodeFactoryIo/vedran-daemon/internal/scheduler"
)

const (
	metricsSendInterval = 30
	pingSendInterval    = 5
)

type TelemetryInterface interface {
	StartSendingTelemetry(scheduler scheduler.Scheduler, client *lb.Client, nodeMetrics string)
}

type Telemetry struct{}

// StartSendingTelemetry start sending ping and metrics to load balancer and blocks current thread
func (t *Telemetry) StartSendingTelemetry(scheduler scheduler.Scheduler, client *lb.Client, nodeMetrics string) {
	fms := &metrics.FetchMetricsService{BaseURL: nodeMetrics}
	_, _ = scheduler.Every(metricsSendInterval).Seconds().Do(client.Metrics.Send, fms)
	log.Println("Started sending metrics to load balancer")

	_, _ = scheduler.Every(pingSendInterval).Seconds().Do(client.Ping.Send)
	log.Println("Started sending pings to load balancer")

	scheduler.StartBlocking()
}
