package tunnel

import (
	"net/url"
	"time"

	"github.com/NodeFactoryIo/vedran/pkg/http-tunnel/client"
	log "github.com/sirupsen/logrus"
)

const (
	DefaultBackoffInterval    = 500 * time.Millisecond
	DefaultBackoffMultiplier  = 1.5
	DefaultBackoffMaxInterval = 60 * time.Second
	DefaultBackoffMaxTime     = 15 * time.Minute
	Protocol                  = "tcp"
)

type Tunnel struct {
	NodeRPCURL *url.URL
}

type Tunneler interface {
	// StartTunnel connects to load balancer tunnel port and creates connection
	StartTunnel(nodeID string, tunnelServerAddress string, token string)
}

func (t *Tunnel) StartTunnel(nodeID string, tunnelServerAddress string, token string) {
	c, err := client.NewClient(&client.ClientConfig{
		ServerAddress: tunnelServerAddress,
		Tunnels: map[string]*client.Tunnel{
			"default": {
				Protocol:   Protocol,
				Addr:       t.NodeRPCURL.Host,
				RemoteAddr: "0.0.0.0:AUTO",
			},
		},
		Logger:    log.NewEntry(log.New()),
		AuthToken: token,
		IdName:    nodeID,
	})
	if err != nil {
		log.Fatal("Failed to connect to tunnel: ", err)
	}

	err = c.Start()
	if err != nil {
		log.Fatal("Failed to start tunnels: ", err)
	}
}
