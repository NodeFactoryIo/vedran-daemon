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

// Start registers to load balancer and starts sending telemetry
func Start(client *lb.Client, id string, nodeRPC string, nodeMetrics string, payoutAddress string) error {
	err := client.Register(id, nodeRPC, payoutAddress, "test-config-hash")
	if err != nil {
		return err
	}
	log.Printf("Registered to load balancer %s", client.BaseURL.String())

	scheduler := gocron.NewScheduler(time.UTC)
	err = startSendingTelemetry(scheduler, client, nodeMetrics)
	return err
}

var startSendingTelemetry = func(scheduler Scheduler, client *lb.Client, nodeMetrics string) error {
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
