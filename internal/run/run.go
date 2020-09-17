package run

import (
	"log"
	"time"

	"github.com/NodeFactoryIo/vedran-daemon/internal/lb"
	"github.com/NodeFactoryIo/vedran-daemon/internal/metrics"
	"github.com/NodeFactoryIo/vedran-daemon/internal/telemetry"
	"github.com/go-co-op/gocron"
)

// Start registers to load balancer and starts sending telemetry
func Start(lbClient *lb.Client, metricsClient metrics.Client, telemetry telemetry.TelemetryInterface, id string, nodeRPC string, payoutAddress string) error {
	err := lbClient.Register(id, nodeRPC, payoutAddress, "test-config-hash")
	if err != nil {
		return err
	}
	log.Printf("Registered to load balancer %s", lbClient.BaseURL.String())

	scheduler := gocron.NewScheduler(time.UTC)
	telemetry.StartSendingTelemetry(scheduler, lbClient, metricsClient)
	return nil
}
