package lb

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// PingService is used for pinging load balancer to confirm daemon is alive
type PingService interface {
	Send() (*http.Response, error)
}

type pingService struct {
	client *Client
}

type pingRequest struct {
	Timestamp int64 `json:"timestamp"`
}

const (
	pingEndpoint = "/api/v1/nodes/pings"
)

func (ps *pingService) Send() (*http.Response, error) {
	body := &pingRequest{
		Timestamp: time.Now().Unix(),
	}

	log.Debug("Sending ping to load balancer")
	req, _ := ps.client.newRequest(http.MethodPost, pingEndpoint, body)
	resp, err := ps.client.do(req, nil)

	if err != nil {
		log.Errorf("Failed sending ping to load balancer because of: %v", err)
		return nil, err
	}

	return resp, err
}
