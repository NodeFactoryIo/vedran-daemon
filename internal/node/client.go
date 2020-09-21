package node

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/mitchellh/mapstructure"
)

// Client is used to interact with polkadot node
type Client interface {
	GetMetrics() (*Metrics, error)
	GetConfigHash() (hash.Hash32, error)
	GetRPCURL() string
	GetMetricsURL() string
}

// RPCError is error returned if rpc request fails
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// RPCRequest is used for retrieving data from node
type RPCRequest struct {
	ID      int      `json:"id"`
	JSONRPC string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
}

// RPCResponse is response from rpc request
type RPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Error   RPCError    `json:"error"`
	ID      int         `json:"id"`
	Result  interface{} `json:"result"`
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

// SendsRPCRequest sends rpc request to node rpc url and decodes result to v
func (client *client) sendRPCRequest(method string, params []string, v interface{}) (*http.Response, error) {
	rpcReq := &RPCRequest{
		ID:      1,
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}
	buf := new(bytes.Buffer)
	_ = json.NewEncoder(buf).Encode(rpcReq)

	resp, err := http.Post(client.GetRPCURL(), "application/json", buf)
	if err != nil {
		return nil, err
	}

	var rpcResponse RPCResponse
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll((resp.Body))
	err = json.Unmarshal(body, &rpcResponse)

	if err != nil {
		return nil, err
	} else if rpcResponse.Error != (RPCError{}) {
		return nil, fmt.Errorf("Node rpc request returned invalid status code: %v", rpcResponse.Error)
	}

	err = mapstructure.Decode(rpcResponse.Result, &v)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
