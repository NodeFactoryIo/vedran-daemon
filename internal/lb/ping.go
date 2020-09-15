package lb

import (
	"log"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
)

// PingService is used for pinging load balancer to confirm daemon is alive
type PingService interface {
	Send() (*http.Response, error)
}

type pingService struct {
	client *Client
}

type pingRequest struct {
	timestamp int64
}

const (
	pingEndpoint = "/api/v1/nodes/pings"
)

func (ps *pingService) Send() (*http.Response, error) {
	body := &pingRequest{
		timestamp: time.Now().Unix(),
	}

	log.Println("Sending ping to load balancer")
	req, _ := ps.client.NewRequest(http.MethodPost, pingEndpoint, body)
	resp, err := ps.client.Do(req, nil)

	if err != nil {
		sentry.CaptureException(err)
		return nil, err
	}

	return resp, err
}
