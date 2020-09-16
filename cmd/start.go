package cmd

import (
	"fmt"
	"net/url"
	"os"

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

var startCmd = &cobra.Command{
	Use:   "vedran-daemon",
	Short: "Register vedran-daemon with load balancer and start sending telemetry",
	RunE:  start,
}

func init() {
	startCmd.Flags().StringVar(&nodeRPC, "node-rpc", "localhost:9933", "Polkadot node rpc url")
	startCmd.Flags().StringVar(&nodeMetrics, "node-metrics", "localhost:9615", "Polkadot node metrics url")
	startCmd.Flags().StringVar(&id, "id", "", "Vedran-daemon id string (required)")
	startCmd.Flags().StringVar(&lbBaseURL, "lb", "", "Target load balancer url (required)")
	startCmd.Flags().StringVar(&payoutAddress, "payout-address", "", "Payout address to which reward tokens will be sent (required)")

	_ = startCmd.MarkFlagRequired("id")
	_ = startCmd.MarkFlagRequired("lb")
	_ = startCmd.MarkFlagRequired("payout-address")
}

func start(cmd *cobra.Command, _ []string) error {
	lbURL, err := url.Parse(lbBaseURL)
	if err != nil {
		return fmt.Errorf("Failed parsing load balancer url")
	}

	client := lb.NewClient(lbURL)
	err = run.Start(client, id, nodeRPC, nodeMetrics, payoutAddress)

	if err != nil {
		return fmt.Errorf("Failed registering to load balancer because of: %v", err)
	}

	return nil
}

func Execute() {
	if err := startCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
