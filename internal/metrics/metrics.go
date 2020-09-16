package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/common/expfmt"
)

// Metrics required to be sent to load balancer
type Metrics struct {
	PeerCount             *float64
	BestBlockHeight       *float64
	FinalizedBlockHeight  *float64
	ReadyTransactionCount *float64
}

const (
	metricsEndpoint = "/metrics"
)

// GetNodeMetrics retrieves polkadot metrics from prometheus server
func GetNodeMetrics(baseURL string) (*Metrics, error) {
	resp, err := http.Get(baseURL + metricsEndpoint)
	if err != nil {
		return nil, fmt.Errorf("Metrics endpoint returned error: %v", err)
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Metrics endpoint returned invalid status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	var parser expfmt.TextParser
	metricFamilies, err := parser.TextToMetricFamilies(resp.Body)
	if err != nil {
		return nil, err
	}

	metrics := &Metrics{
		metricFamilies["polkadot_sync_peers"].GetMetric()[0].Gauge.Value,
		metricFamilies["polkadot_block_height"].GetMetric()[0].Gauge.Value,
		metricFamilies["polkadot_block_height"].GetMetric()[1].Gauge.Value,
		metricFamilies["polkadot_ready_transactions_number"].GetMetric()[0].Gauge.Value,
	}
	return metrics, nil
}
