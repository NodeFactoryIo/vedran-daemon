package run

import (
	"log"
	"time"

	"github.com/NodeFactoryIo/vedran-daemon/internal/lb"
	"github.com/NodeFactoryIo/vedran-daemon/internal/node"
	"github.com/NodeFactoryIo/vedran-daemon/internal/telemetry"
	"github.com/go-co-op/gocron"
)

// Start registers to load balancer and starts sending telemetry
func Start(lbClient *lb.Client, nodeClient node.Client, telemetry telemetry.TelemetryInterface, id string, payoutAddress string) error {
	err := lbClient.Register(id, nodeClient.GetRPCURL(), payoutAddress, "test-config-hash")
	if err != nil {
		return err
	}
	log.Printf("Registered to load balancer %s", lbClient.BaseURL.String())

	scheduler := gocron.NewScheduler(time.UTC)
	telemetry.StartSendingTelemetry(scheduler, lbClient, nodeClient)
	return nil
}
