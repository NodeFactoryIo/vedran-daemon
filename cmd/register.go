package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var (
	node          string
	id            string
	lb            string
	payoutAddress string
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register vedran-daemon with load balancer and start receiving requests",
	RunE:  register,
}

func init() {
	registerCmd.Flags().StringVar(&node, "node", "", "Polkadot node url (required)")
	registerCmd.Flags().StringVar(&id, "id", "", "Vedran-daemon id string (required)")
	registerCmd.Flags().StringVar(&lb, "lb", "", "Target load balancer url (required)")
	registerCmd.Flags().StringVar(&payoutAddress, "payout-address", "", "Payout address to which reward tokens will be sent (required)")

	_ = registerCmd.MarkFlagRequired("node")
	_ = registerCmd.MarkFlagRequired("id")
	_ = registerCmd.MarkFlagRequired("lb")
	_ = registerCmd.MarkFlagRequired("payout-address")

	RootCmd.AddCommand(registerCmd)
}

func register(_ *cobra.Command, _ []string) error {
	resp, err := http.Get(node + "/metrics")
	if err != nil && resp.StatusCode != 200 {
		return fmt.Errorf("Error polling polkadot node: %v", err)
	}

	values := map[string]string{
		"id":             id,
		"config_hash":    "config_hash",
		"node_url":       node,
		"payout_address": payoutAddress}
	jsonValue, _ := json.Marshal(values)
	resp, err = http.Post(lb+"/api/v1/nodes", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return fmt.Errorf("Error retrieving load balancer token: %v", err)
	}

	fmt.Printf("Load balancer response: %v", resp)
	return nil
}
