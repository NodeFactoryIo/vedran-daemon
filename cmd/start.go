package cmd

import (
	"fmt"
	"net/url"
	"os"

	"github.com/NodeFactoryIo/vedran-daemon/internal/lb"
	"github.com/NodeFactoryIo/vedran-daemon/internal/metrics"
	"github.com/NodeFactoryIo/vedran-daemon/internal/run"
	"github.com/NodeFactoryIo/vedran-daemon/internal/telemetry"
	"github.com/spf13/cobra"
)

var (
	nodeRPCURL     string
	nodeMetricsURL string
	id             string
	lbBaseURL      string
	payoutAddress  string
	lbURL          *url.URL
	metricsURL     *url.URL
)

var startCmd = &cobra.Command{
	Use:   "vedran-daemon",
	Short: "Register vedran-daemon with load balancer and start sending telemetry",
	RunE:  start,
	Args: func(cms *cobra.Command, args []string) error {
		var err error
		lbURL, err = url.Parse(lbBaseURL)
		if err != nil {
			return fmt.Errorf("Failed parsing load balancer url")
		}

		metricsURL, err = url.Parse(nodeMetricsURL)
		if err != nil {
			return fmt.Errorf("Failed parsing load balancer url")
		}

		return nil
	},
}

func init() {
	startCmd.Flags().StringVar(&nodeRPCURL, "node-rpc", "localhost:9933", "Polkadot node rpc url")
	startCmd.Flags().StringVar(&nodeMetricsURL, "node-metrics", "localhost:9615", "Polkadot node metrics url")
	startCmd.Flags().StringVar(&id, "id", "", "Vedran-daemon id string (required)")
	startCmd.Flags().StringVar(&lbBaseURL, "lb", "", "Target load balancer url (required)")
	startCmd.Flags().StringVar(&payoutAddress, "payout-address", "", "Payout address to which reward tokens will be sent (required)")

	_ = startCmd.MarkFlagRequired("id")
	_ = startCmd.MarkFlagRequired("lb")
	_ = startCmd.MarkFlagRequired("payout-address")
}

func start(cmd *cobra.Command, _ []string) error {
	lbClient := lb.NewClient(lbURL)
	metricsClient := metrics.NewClient(metricsURL)
	telemetry := &telemetry.Telemetry{}

	err := run.Start(lbClient, metricsClient, telemetry, id, nodeRPCURL, payoutAddress)
	if err != nil {
		return fmt.Errorf("Failed starting vedran daemon because: %v", err)
	}

	return nil
}

func Execute() {
	if err := startCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
