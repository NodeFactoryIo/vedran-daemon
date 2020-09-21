package run

import (
	"encoding/base64"
	"time"

	"github.com/NodeFactoryIo/vedran-daemon/internal/lb"
	"github.com/NodeFactoryIo/vedran-daemon/internal/node"
	"github.com/NodeFactoryIo/vedran-daemon/internal/telemetry"
	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
)

// Start registers to load balancer and starts sending telemetry
func Start(lbClient *lb.Client, nodeClient node.Client, telemetry telemetry.Telemetry, id string, payoutAddress string) error {
	configHash, err := nodeClient.GetConfigHash()
	if err != nil {
		return err
	}

	err = lbClient.Register(id, nodeClient.GetRPCURL(), payoutAddress, base64.StdEncoding.EncodeToString(configHash.Sum(nil)[:]))
	if err != nil {
		return err
	}
	log.Infof("Registered to load balancer %s", lbClient.BaseURL.String())

	scheduler := gocron.NewScheduler(time.UTC)
	telemetry.StartSendingTelemetry(scheduler, lbClient, nodeClient)
	return nil
}
