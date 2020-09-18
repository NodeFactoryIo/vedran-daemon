package node

import (
	"fmt"
	"hash"
	"hash/fnv"
	"net/http"
	"net/url"

	"github.com/prometheus/common/expfmt"
)

// Client is used to interact with polkadot node
type Client interface {
	GetMetrics() (*Metrics, error)
	GetConfigHash() (hash.Hash32, error)
	GetRPCURL() string
	GetMetricsURL() string
}

// Metrics required to be sent to load balancer
type Metrics struct {
	PeerCount             *float64 `json:"peer_count"`
	BestBlockHeight       *float64 `json:"best_block_height"`
	FinalizedBlockHeight  *float64 `json:"finalized_block_height"`
	ReadyTransactionCount *float64 `json:"read_transaction_count"`
}

// NewClient creates node client instance
func NewClient(rpcURL *url.URL, metricsURL *url.URL) Client {
	return &client{MetricsBaseURL: metricsURL, RPCBaseURL: rpcURL}
}

type client struct {
	MetricsBaseURL *url.URL
	RPCBaseURL     *url.URL
}

const (
	metricsEndpoint = "/metrics"
)

// GetRPCURL returns rpc base url as string
func (client *client) GetRPCURL() string {
	return client.RPCBaseURL.String()
}

// GetMetricsURL returns metrics base url as string
func (client *client) GetMetricsURL() string {
	return client.MetricsBaseURL.String()
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

	metrics := &Metrics{
		metricFamilies["polkadot_sync_peers"].GetMetric()[0].Gauge.Value,
		metricFamilies["polkadot_block_height"].GetMetric()[0].Gauge.Value,
		metricFamilies["polkadot_block_height"].GetMetric()[1].Gauge.Value,
		metricFamilies["polkadot_ready_transactions_number"].GetMetric()[0].Gauge.Value,
	}
	return metrics, nil
}

// GetConfigHash returns sorted hash of supported rpc methods3
func (client *client) GetConfigHash() (hash.Hash32, error) {
	hash := fnv.New32()

	_, err := hash.Write([]byte("TODO:hash"))
	if err != nil {
		return nil, err
	}

	return hash, nil
}
