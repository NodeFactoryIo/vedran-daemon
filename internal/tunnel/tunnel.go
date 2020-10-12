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

type Tunnel struct {
	TunnelURL  *url.URL
	NodeRPCURL *url.URL
}

type Tunneler interface {
	// StartTunnel connects to load balancer tunnel port and creates connection
	StartTunnel(nodeID string, token string)
}

func (t *Tunnel) StartTunnel(nodeID string, token string) {
	c, err := client.NewClient(&client.ClientConfig{
		ServerAddress: t.TunnelURL.Host,
		Tunnels: map[string]*client.Tunnel{
			"default": {
				Protocol:   Protocol,
				Addr:       t.NodeRPCURL.Host,
				RemoteAddr: RemoteAddr,
			},
		},
		Logger:    log.NewEntry(log.New()),
		AuthToken: token,
	})
	if err != nil {
		log.Fatalf("Failed to connect to tunnel: ", err)
	}

	err = c.Start()
	if err != nil {
		log.Fatal("Failed to start tunnels: ", err)
	}
}
