package lb

import (
	"fmt"
	"net/http"

	"github.com/NodeFactoryIo/vedran-daemon/internal/node"
	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
)

// MetricsService is used for sending node metrics to load balancer
type MetricsService interface {
	Send(client node.Client) (*http.Response, error)
}

type metricsService struct {
	client *Client
}

const (
	metricsEndpoint = "/api/v1/nodes/metrics"
)

func (ms *metricsService) Send(client node.Client) (*http.Response, error) {
	metrics, err := client.GetMetrics()
	if err != nil {
		return nil, err
	}

	log.Info("Sending metrics to load balancer")
	req, _ := ms.client.newRequest(http.MethodPut, metricsEndpoint, metrics)
	resp, err := ms.client.do(req, nil)

	if err != nil {
		log.Error(fmt.Sprintf("Failed sending metrics to load balancer because of: %v", err))
		sentry.CaptureException(err)
		return nil, err
	}

	return resp, err
}
