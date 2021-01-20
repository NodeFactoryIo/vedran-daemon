package node

import (
	"fmt"
	"net/http"

	"github.com/prometheus/common/expfmt"
)

const (
	metricsEndpoint = "/metrics"
)

// Metrics required to be sent to load balancer
type Metrics struct {
	PeerCount             *float64 `json:"peer_count"`
	BestBlockHeight       *float64 `json:"best_block_height"`
	FinalizedBlockHeight  *float64 `json:"finalized_block_height"`
	TargetBlockHeight     *float64 `json:"target_block_height"`
	ReadyTransactionCount *float64 `json:"ready_transaction_count"`
}

// GetNodeMetrics retrieves polkadot metrics from prometheus server
func (client *client) GetMetrics() (*Metrics, error) {
	metricsURL, _ := client.MetricsBaseURL.Parse(metricsEndpoint)
	resp, err := http.Get(metricsURL.String())
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

	blockMetrics := metricFamilies["polkadot_block_height"].GetMetric()
	bestBlockHeight := blockMetrics[0].Gauge.Value
	finalizedBlockHeight := blockMetrics[1].Gauge.Value
	targetBlockHeight := bestBlockHeight
	// if node not synced there is additional metric - target block
	if len(blockMetrics) == 3 {
		targetBlockHeight = blockMetrics[2].Gauge.Value
	}

	metrics := &Metrics{
		metricFamilies["polkadot_sync_peers"].GetMetric()[0].Gauge.Value,
		bestBlockHeight,
		finalizedBlockHeight,
		targetBlockHeight,
		metricFamilies["polkadot_ready_transactions_number"].GetMetric()[0].Gauge.Value,
	}
	return metrics, nil
}
