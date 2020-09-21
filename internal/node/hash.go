package node

import (
	"fmt"
	"hash"
	"hash/fnv"
	"sort"
	"strings"
)

type JSONRPCResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
}

// MethodsResult contains all available rpc methods of node
type MethodsResult struct {
	Version int      `json:"version"`
	Methods []string `json:"methods"`
}

// RPCMethodsResponse is rpc response od rpc_methods call
type RPCMethodsResponse struct {
	JSONRPCResponse
	Result MethodsResult `json:"result"`
}

// RPCChainTypeResponse is rpc response of system_chainType call
type RPCChainTypeResponse struct {
	JSONRPCResponse
	Result string `json:"result"`
}

// RPCChainResponse is rpc response of system_chain call
type RPCChainResponse struct {
	JSONRPCResponse
	Result string `json:"result"`
}

// RPCNodeRolesResponse is rpc response of system_nodeRoles call
type RPCNodeRolesResponse struct {
	JSONRPCResponse
	Result []string `json:"result"`
}

// Properties are node properties
type Properties struct {
	SS58Format    int    `json:"ss58Format"`
	TokenDecimals int    `json:"tokenDecimals"`
	TokenSymbol   string `json:"tokenSymbol"`
}

// RPCPropertiesResponse is rpc response of system_properties call
type RPCPropertiesResponse struct {
	JSONRPCResponse
	Result Properties `json:"result"`
}

// GetConfigHash hashes sorted node configuration
func (client *client) GetConfigHash() (hash.Hash32, error) {
	methods, err := client.getNodeRPCMethods()
	if err != nil {
		return nil, err
	}

	nodeRoles, err := client.getNodeRoles()
	if err != nil {
		return nil, err
	}

	chain, err := client.getChain()
	if err != nil {
		return nil, err
	}

	chainType, err := client.getChainType()
	if err != nil {
		return nil, err
	}

	properties, err := client.getNodeProperties()
	if err != nil {
		return nil, err
	}

	hash := fnv.New32()
	sort.Strings(methods)
	sort.Strings(nodeRoles)
	_, _ = hash.Write([]byte(chain))
	_, _ = hash.Write([]byte(chainType))
	_, _ = hash.Write([]byte(strings.Join(methods, "")))
	_, _ = hash.Write([]byte(strings.Join(nodeRoles, "")))
	_, _ = hash.Write([]byte(fmt.Sprint(properties.SS58Format)))
	_, _ = hash.Write([]byte(fmt.Sprint(properties.TokenDecimals)))
	_, _ = hash.Write([]byte(string(properties.TokenSymbol)))

	return hash, nil
}

func (client *client) getChainType() (string, error) {
	var chainType RPCChainTypeResponse
	_, err := client.sendRPCRequest("system_chainType", []string{}, &chainType)

	if err != nil {
		return "", err
	}

	return chainType.Result, nil
}

func (client *client) getChain() (string, error) {
	var chain RPCChainResponse
	_, err := client.sendRPCRequest("system_chain", []string{}, &chain)

	if err != nil {
		return "", err
	}

	return chain.Result, nil
}

func (client *client) getNodeRoles() ([]string, error) {
	var nodeRoles RPCNodeRolesResponse
	_, err := client.sendRPCRequest("system_nodeRoles", []string{}, &nodeRoles)

	if err != nil {
		return nil, err
	}

	return nodeRoles.Result, nil
}

func (client *client) getNodeProperties() (Properties, error) {
	var properties RPCPropertiesResponse
	_, err := client.sendRPCRequest("system_properties", []string{}, &properties)

	if err != nil {
		return Properties{}, err
	}

	return properties.Result, nil
}

func (client *client) getNodeRPCMethods() ([]string, error) {
	var rpcMethods RPCMethodsResponse
	_, err := client.sendRPCRequest("rpc_methods", []string{}, &rpcMethods)

	if err != nil {
		return nil, err
	}

	return rpcMethods.Result.Methods, nil
}
