package tunnel

import (
	"net/url"
	"strconv"
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
	StartTunnel(nodeID string, tunnelURL string, token string, port int)
}

func (t *Tunnel) StartTunnel(nodeID string, tunnelURL string, token string, port int) {
	c, err := client.NewClient(&client.ClientConfig{
		ServerAddress: tunnelURL,
		Tunnels: map[string]*client.Tunnel{
			"default": {
				Protocol:   Protocol,
				Addr:       t.NodeRPCURL.Host,
				RemoteAddr: "0.0.0.0:" + strconv.Itoa(port),
			},
		},
		Logger:    log.NewEntry(log.New()),
		AuthToken: token,
	})
	if err != nil {
		log.Fatal("Failed to connect to tunnel: ", err)
	}

	err = c.Start()
	if err != nil {
		log.Fatal("Failed to start tunnels: ", err)
	}
}
