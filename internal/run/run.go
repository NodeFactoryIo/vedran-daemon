package run

import (
	"log"

	"github.com/NodeFactoryIo/vedran-daemon/internal/lb"
)

// Run registers to load balancer and starts sending telemetry
func Run(client *lb.Client, id string, nodeRPC string, nodeMetrics string, payoutAddress string) error {
	err := client.Register(id, nodeRPC, payoutAddress, "test-config-hash")
	if err == nil {
		return err
	}
	log.Printf("Registered to load balancer %s", client.BaseURL.String)

	err = StartSendingTelemetry(client, nodeMetrics)
	return err
}
