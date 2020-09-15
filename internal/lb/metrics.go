package lb

import (
	"log"
	"net/http"
	"time"

	"github.com/NodeFactoryIo/vedran-daemon/internal/metrics"
	"github.com/getsentry/sentry-go"
)

// MetricsService is used for sending node metrics to load balancer
type MetricsService interface {
	Send(metricsBaseURL string) (*http.Response, error)
}

type metricsService struct {
	client *Client
}

type metricsRequest struct {
	timestamp time.Time
}

const (
	metricsEndpoint = "/api/v1/nodes/metrics"
)

func (ms *metricsService) Send(metricsBaseURL string) (*http.Response, error) {
	metrics, err := metrics.GetNodeMetrics(metricsBaseURL)
	if err != nil {
		return nil, err
	}

	req, err := ms.client.NewRequest(http.MethodPut, metricsEndpoint, metrics)
	if err != nil {
		return nil, err
	}

	log.Println("Sending metrics to load balancer")
	resp, err := ms.client.Do(req, new(interface{}))

	if err != nil {
		sentry.CaptureException(err)
		return nil, err
	}

	return resp, err
}
