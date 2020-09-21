package node

import (
	"encoding/json"
	"hash"
	"hash/fnv"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfigHash(t *testing.T) {

	type Test struct {
		name       string
		want       hash.Hash32
		wantErr    bool
		handleFunc handleFnMock
	}

	methodsResponse := `{
		"jsonrpc": "2.0",
			"result": {
			"methods": [
				"system_chain",
				"a_chain"
			]
		},
		"id": 1
	}`
	nodeRolesResponse := `{
		"jsonrpc": "2.0",
		"result": [
			"Full",
			"Archive"
		],
		"id": 1
	}`
	chainResponse := `{
		"jsonrpc": "2.0",
		"result": "Polkadot",
		"id": 1
	}`
	chainTypeResponse := `{
		"jsonrpc": "2.0",
		"result": "Live",
		"id": 1
	}`
	propertiesResponse := `{
		"jsonrpc": "2.0",
		"result": {
			"ss58format": 5,
			"tokenDecimals": 0,
			"tokenSymbol": "Dot"
		},
		"id": 1
	}`

	expectedHash := fnv.New32()
	_, _ = expectedHash.Write([]byte("Polkadot"))
	_, _ = expectedHash.Write([]byte("Live"))
	_, _ = expectedHash.Write([]byte("a_chainsystem_chain"))
	_, _ = expectedHash.Write([]byte("ArchiveFull"))
	_, _ = expectedHash.Write([]byte("5"))
	_, _ = expectedHash.Write([]byte("0"))
	_, _ = expectedHash.Write([]byte("Dot"))

	tests := []Test{
		{
			name:    "Returns error if getNodeRPCMethods fails",
			want:    nil,
			wantErr: true,
			handleFunc: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				var rpcRequest RPCRequest
				defer r.Body.Close()
				body, _ := ioutil.ReadAll((r.Body))
				_ = json.Unmarshal(body, &rpcRequest)

				if rpcRequest.Method == "rpc_methods" {
					http.Error(w, "error", 404)
				}
			}},
		{
			name:    "Returns error if getNodeRoles fails",
			want:    nil,
			wantErr: true,
			handleFunc: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				var rpcRequest RPCRequest
				defer r.Body.Close()
				body, _ := ioutil.ReadAll((r.Body))
				_ = json.Unmarshal(body, &rpcRequest)

				if rpcRequest.Method == "rpc_methods" {
					_, _ = io.WriteString(w, methodsResponse)
				} else if rpcRequest.Method == "system_nodeRoles" {
					http.Error(w, "error", 404)
				}
			}},
		{
			name:    "Returns error if getChain fails",
			want:    nil,
			wantErr: true,
			handleFunc: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				var rpcRequest RPCRequest
				defer r.Body.Close()
				body, _ := ioutil.ReadAll((r.Body))
				_ = json.Unmarshal(body, &rpcRequest)

				if rpcRequest.Method == "rpc_methods" {
					_, _ = io.WriteString(w, methodsResponse)
				} else if rpcRequest.Method == "system_nodeRoles" {
					_, _ = io.WriteString(w, nodeRolesResponse)
				} else if rpcRequest.Method == "system_chain" {
					http.Error(w, "error", 404)
				}
			}},
		{
			name:    "Returns error if getChainType fails",
			want:    nil,
			wantErr: true,
			handleFunc: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				var rpcRequest RPCRequest
				defer r.Body.Close()
				body, _ := ioutil.ReadAll((r.Body))
				_ = json.Unmarshal(body, &rpcRequest)

				if rpcRequest.Method == "rpc_methods" {
					_, _ = io.WriteString(w, methodsResponse)
				} else if rpcRequest.Method == "system_nodeRoles" {
					_, _ = io.WriteString(w, nodeRolesResponse)
				} else if rpcRequest.Method == "system_chain" {
					_, _ = io.WriteString(w, chainResponse)
				} else if rpcRequest.Method == "system_chainType" {
					http.Error(w, "error", 404)
				}
			}},
		{
			name:    "Returns error if getNodeProperties fails",
			want:    nil,
			wantErr: true,
			handleFunc: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				var rpcRequest RPCRequest
				defer r.Body.Close()
				body, _ := ioutil.ReadAll((r.Body))
				_ = json.Unmarshal(body, &rpcRequest)

				if rpcRequest.Method == "rpc_methods" {
					_, _ = io.WriteString(w, methodsResponse)
				} else if rpcRequest.Method == "system_nodeRoles" {
					_, _ = io.WriteString(w, nodeRolesResponse)
				} else if rpcRequest.Method == "system_chain" {
					_, _ = io.WriteString(w, chainResponse)
				} else if rpcRequest.Method == "system_chainType" {
					_, _ = io.WriteString(w, chainTypeResponse)
				} else if rpcRequest.Method == "system_properties" {
					http.Error(w, "error", 404)
				}
			}},
		{
			name:    "Returns valid hash if rpc calls succeed",
			want:    expectedHash,
			wantErr: false,
			handleFunc: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				var rpcRequest RPCRequest
				defer r.Body.Close()
				body, _ := ioutil.ReadAll((r.Body))
				_ = json.Unmarshal(body, &rpcRequest)

				if rpcRequest.Method == "rpc_methods" {
					_, _ = io.WriteString(w, methodsResponse)
				} else if rpcRequest.Method == "system_nodeRoles" {
					_, _ = io.WriteString(w, nodeRolesResponse)
				} else if rpcRequest.Method == "system_chain" {
					_, _ = io.WriteString(w, chainResponse)
				} else if rpcRequest.Method == "system_chainType" {
					_, _ = io.WriteString(w, chainTypeResponse)
				} else if rpcRequest.Method == "system_properties" {
					_, _ = io.WriteString(w, propertiesResponse)
				}
			}},
	}

	for _, tt := range tests {
		setup()
		defer teardown()

		t.Run(tt.name, func(t *testing.T) {
			baseURL, _ := url.Parse(server.URL)
			mux.HandleFunc("/", tt.handleFunc)
			client := NewClient(baseURL, baseURL)

			got, err := client.GetConfigHash()

			if (err != nil) != tt.wantErr {
				t.Errorf("GetConfigHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConfigHash() = %v, want %v", got, tt.want)
			}

		})
	}
}
