package lb

import (
	"log"
	"net/http"

	"github.com/NodeFactoryIo/vedran-daemon/internal/metrics"
	"github.com/getsentry/sentry-go"
)

// MetricsService is used for sending node metrics to load balancer
type MetricsService interface {
	Send(fm metrics.FetchMetrics) (*http.Response, error)
}

type metricsService struct {
	client *Client
}

const (
	metricsEndpoint = "/api/v1/nodes/metrics"
)

func (ms *metricsService) Send(fm metrics.FetchMetrics) (*http.Response, error) {
	metrics, err := fm.GetNodeMetrics()
	if err != nil {
		return nil, err
	}

	log.Println("Sending metrics to load balancer")
	req, _ := ms.client.newRequest(http.MethodPut, metricsEndpoint, metrics)
	resp, err := ms.client.do(req, nil)

	if err != nil {
		log.Printf("Falied sending metrics to load balancer because of: %v", err)
		sentry.CaptureException(err)
		return nil, err
	}

	return resp, err
}
