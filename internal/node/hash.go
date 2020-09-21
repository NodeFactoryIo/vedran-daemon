package node

import (
	"fmt"
	"hash"
	"hash/fnv"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
)

// MethodsResult contains all available rpc methods of node
type MethodsResult struct {
	Version int      `json:"version"`
	Methods []string `json:"methods"`
}

// Properties are node properties
type Properties struct {
	SS58Format    int    `json:"ss58Format"`
	TokenDecimals int    `json:"tokenDecimals"`
	TokenSymbol   string `json:"tokenSymbol"`
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

	sort.Strings(methods)
	sort.Strings(nodeRoles)
	log.Infof(`
		Created config hash with:
			Chain: %s
			Chain Type: %s
			Methods: %v
			Node Roles: %v
			SS58Format: %d
			Token Decimals: %d
			Token Symbol: %s
	`,
		chain, chainType, methods, nodeRoles, properties.SS58Format,
		properties.TokenDecimals, properties.TokenSymbol)
	hash := fnv.New32()
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
	var chainType string
	_, err := client.sendRPCRequest("system_chainType", []string{}, &chainType)

	if err != nil {
		return "", err
	}

	return chainType, nil
}

func (client *client) getChain() (string, error) {
	var chain string
	_, err := client.sendRPCRequest("system_chain", []string{}, &chain)

	if err != nil {
		return "", err
	}

	return chain, nil
}

func (client *client) getNodeRoles() ([]string, error) {
	var nodeRoles []string
	_, err := client.sendRPCRequest("system_nodeRoles", []string{}, &nodeRoles)

	if err != nil {
		return nil, err
	}

	return nodeRoles, nil
}

func (client *client) getNodeProperties() (Properties, error) {
	var properties Properties
	_, err := client.sendRPCRequest("system_properties", []string{}, &properties)

	if err != nil {
		return Properties{}, err
	}

	return properties, nil
}

func (client *client) getNodeRPCMethods() ([]string, error) {
	var methods MethodsResult
	_, err := client.sendRPCRequest("rpc_methods", []string{}, &methods)

	if err != nil {
		return nil, err
	}

	return methods.Methods, nil
}
