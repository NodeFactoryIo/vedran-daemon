package node

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash"
	"hash/fnv"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

// MethodsResult contains all available rpc methods of node
type MethodsResult struct {
	Version int      `json:"version"`
	Methods []string `json:"methods"`
}

// RPCMethodsResponse is available rpc methods request response
type RPCMethodsResponse struct {
	JSONRPC float32       `json:"jsonrpc"`
	Result  MethodsResult `json:"result"`
	ID      int           `json:"id"`
}

// GetConfigHash hashes sorted available rpc methods from node
func (client *client) GetConfigHash() (hash.Hash32, error) {
	methods, err := client.getNodeRPCMethods()
	if err != nil {
		return nil, err
	}

	hash := fnv.New32()
	sort.Strings(methods)
	_, _ = hash.Write([]byte(strings.Join(methods, "")))

	return hash, nil
}

func (client *client) getNodeRPCMethods() ([]string, error) {
	var jsonStr = []byte(`{
		"id": 1,
		"jsonrpc": "2.0",
		"method": "rpc_methods"
	}`)
	resp, err := http.Post(client.GetRPCURL(), "application/json", bytes.NewBuffer(jsonStr))

	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Node rpc request returned invalid status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll((resp.Body))
	var rpcMethods RPCMethodsResponse
	err = json.Unmarshal(body, &rpcMethods)

	if err != nil {
		return nil, err
	} else if len(rpcMethods.Result.Methods) == 0 {
		return nil, fmt.Errorf("No registered rpc methods: %v", err)
	}

	return rpcMethods.Result.Methods, nil
}
