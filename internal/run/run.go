package run

import (
	"log"
	"time"

	"github.com/NodeFactoryIo/vedran-daemon/internal/lb"
	"github.com/NodeFactoryIo/vedran-daemon/internal/telemetry"
	"github.com/go-co-op/gocron"
)

// Start registers to load balancer and starts sending telemetry
func Start(client *lb.Client, telemetry telemetry.TelemetryInterface, id string, nodeRPC string, nodeMetrics string, payoutAddress string) error {
	err := client.Register(id, nodeRPC, payoutAddress, "test-config-hash")
	if err != nil {
		return err
	}
	log.Printf("Registered to load balancer %s", client.BaseURL.String())

	scheduler := gocron.NewScheduler(time.UTC)
	telemetry.StartSendingTelemetry(scheduler, client, nodeMetrics)
	return nil
}
