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
	RemoteAddr                = "0.0.0.0:AUTO"
)

// Tunnel is tunnel connection with load balancer
type Tunnel struct {
	NodeRPCURL *url.URL
}

// Tunneler defines methods for connecting to load balancer tunnel
type Tunneler interface {
	// StartTunnel connects to load balancer tunnel port and creates connection.
	// nodeID is id that is passed in daemon,
	// tunnelServerAddress is public address of load balancer tunnel server and
	// token is jwt token given when registering with load balancer.
	StartTunnel(nodeID string, tunnelServerAddress string, token string)
}

func (t *Tunnel) StartTunnel(nodeID string, tunnelServerAddress string, token string) {
	c, err := client.NewClient(&client.ClientConfig{
		ServerAddress: tunnelServerAddress,
		Tunnels: map[string]*client.Tunnel{
			"default": {
				Protocol:   Protocol,
				Addr:       t.NodeRPCURL.Host,
				RemoteAddr: RemoteAddr,
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
