package run

import (
	"log"
	"time"

	"github.com/NodeFactoryIo/vedran-daemon/internal/lb"
	"github.com/NodeFactoryIo/vedran-daemon/internal/metrics"
	"github.com/go-co-op/gocron"
)

const (
	metricsSendInterval = 30
	pingSendInterval    = 5
)

// StartSendingTelemetry starts sending pings and metrics recurringly to load balancer
// and blocks current thread to prevent cmd exit
func StartSendingTelemetry(client *lb.Client, nodeMetrics string) error {
	scheduler := gocron.NewScheduler(time.UTC)

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
}
