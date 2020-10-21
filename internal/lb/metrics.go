package lb

import (
	"net/http"

	"github.com/NodeFactoryIo/vedran-daemon/internal/node"
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
		log.Errorf("Failed sending metrics to load balancer because of: %v", err)
		return nil, err
	}

	log.Debugf("Sending metrics to load balancer: %+v", metrics)
	req, _ := ms.client.newRequest(http.MethodPut, metricsEndpoint, metrics)
	resp, err := ms.client.do(req, nil)

	if err != nil {
		log.Errorf("Failed sending metrics to load balancer because of: %v", err)
		return nil, err
	}

	return resp, err
}
