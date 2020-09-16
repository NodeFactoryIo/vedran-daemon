package telemetry

import (
	"log"

	"github.com/NodeFactoryIo/vedran-daemon/internal/lb"
	"github.com/NodeFactoryIo/vedran-daemon/internal/metrics"
	"github.com/go-co-op/gocron"
)

const (
	metricsSendInterval = 30
	pingSendInterval    = 5
)

type Scheduler interface {
	Every(interval uint64) *gocron.Scheduler
	StartBlocking()
}

type TelemetryInterface interface {
	StartSendingTelemetry(scheduler Scheduler, client *lb.Client, nodeMetrics string) error
}

type Telemetry struct{}

// StartSendingTelemetry start sending ping and metrics to load balancer and blocks current thread
func (t *Telemetry) StartSendingTelemetry(scheduler Scheduler, client *lb.Client, nodeMetrics string) error {
	fms := &metrics.FetchMetricsService{BaseURL: nodeMetrics}
	_, err := scheduler.Every(metricsSendInterval).Seconds().Do(client.Metrics.Send, fms)
	if err != nil {
		return err
	}
	log.Println("Started sending metrics to load balancer")

	_, err = scheduler.Every(pingSendInterval).Seconds().Do(client.Ping.Send)
	if err != nil {
		return err
	}
	log.Println("Started sending pings to load balancer")

	scheduler.StartBlocking()
	return nil

}
