package node

import (
	"hash"
	"net/url"
)

// Client is used to interact with polkadot node
type Client interface {
	GetMetrics() (*Metrics, error)
	GetConfigHash() (hash.Hash32, error)
	GetRPCURL() string
	GetMetricsURL() string
}

// NewClient creates node client instance
func NewClient(rpcURL *url.URL, metricsURL *url.URL) Client {
	return &client{MetricsBaseURL: metricsURL, RPCBaseURL: rpcURL}
}

type client struct {
	MetricsBaseURL *url.URL
	RPCBaseURL     *url.URL
}

// GetRPCURL returns rpc base url as string
func (client *client) GetRPCURL() string {
	return client.RPCBaseURL.String()
}

// GetMetricsURL returns metrics base url as string
func (client *client) GetMetricsURL() string {
	return client.MetricsBaseURL.String()
}
