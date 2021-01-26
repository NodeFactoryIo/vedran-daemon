package run

import (
	"encoding/base64"
	"hash"
	"time"

	"github.com/NodeFactoryIo/vedran-daemon/internal/lb"
	"github.com/NodeFactoryIo/vedran-daemon/internal/node"
	"github.com/NodeFactoryIo/vedran-daemon/internal/telemetry"
	"github.com/NodeFactoryIo/vedran-daemon/internal/tunnel"
	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
)

var sleep = time.Sleep

// Start registers to load balancer and starts sending telemetry
func Start(
	tunnel tunnel.Tunneler,
	lbClient *lb.Client,
	nodeClient node.Client,
	telemetry telemetry.Telemetry,
	id string,
	payoutAddress string,
) error {
	var configHash hash.Hash32
	for {
		var err error
		configHash, err = nodeClient.GetConfigHash()
		if err == nil {
			break
		}

		log.Errorf("Failed retrieving node metrics because of %v. Retrying in 5 seconds...", err)
		sleep(time.Second * 5)
	}

	registerResponse, err := lbClient.Register(
		id, payoutAddress, base64.StdEncoding.EncodeToString(configHash.Sum(nil)[:]),
	)
	if err != nil {
		return err
	}
	log.Infof("Registered to load balancer %s", lbClient.BaseURL.String())

	go tunnel.StartTunnel(id, registerResponse.TunnelServerAddress, registerResponse.Token)

	scheduler := gocron.NewScheduler(time.UTC)
	telemetry.StartSendingTelemetry(scheduler, lbClient, nodeClient)
	return nil
}
