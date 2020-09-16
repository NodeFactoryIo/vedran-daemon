package cmd

import (
	"fmt"
	"net/url"

	"github.com/NodeFactoryIo/vedran-daemon/internal/lb"
	"github.com/NodeFactoryIo/vedran-daemon/internal/run"
	"github.com/spf13/cobra"
)

var (
	nodeRPC       string
	nodeMetrics   string
	id            string
	lbBaseURL     string
	payoutAddress string
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register vedran-daemon with load balancer and start receiving requests",
	RunE:  register,
}

func init() {
	registerCmd.Flags().StringVar(&nodeRPC, "node-rpc", "", "Polkadot node rpc url (required)")
	registerCmd.Flags().StringVar(&nodeMetrics, "node-metrics", "", "Polkadot node metrics url (required)")
	registerCmd.Flags().StringVar(&id, "id", "", "Vedran-daemon id string (required)")
	registerCmd.Flags().StringVar(&lbBaseURL, "lb", "", "Target load balancer url (required)")
	registerCmd.Flags().StringVar(&payoutAddress, "payout-address", "", "Payout address to which reward tokens will be sent (required)")

	_ = registerCmd.MarkFlagRequired("node-rpc")
	_ = registerCmd.MarkFlagRequired("node-metrics")
	_ = registerCmd.MarkFlagRequired("id")
	_ = registerCmd.MarkFlagRequired("lb")
	_ = registerCmd.MarkFlagRequired("payout-address")

	RootCmd.AddCommand(registerCmd)
}

func register(_ *cobra.Command, _ []string) error {
	lbURL, err := url.Parse(lbBaseURL)
	if err != nil {
		return fmt.Errorf("Failed parsing load balancer url")
	}

	client := lb.NewClient(lbURL)
	err = run.Run(client, id, nodeRPC, nodeMetrics, payoutAddress)
	if err != nil {
		return fmt.Errorf("Failed registering to load balancer because of: %v", err)
	}

	return nil
}
