package cmd

import (
	"fmt"
	"net/url"
	"os"

	"github.com/NodeFactoryIo/vedran-daemon/internal/lb"
	"github.com/NodeFactoryIo/vedran-daemon/internal/node"
	"github.com/NodeFactoryIo/vedran-daemon/internal/run"
	"github.com/NodeFactoryIo/vedran-daemon/internal/telemetry"
	"github.com/NodeFactoryIo/vedran-daemon/pkg/logger"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	logLevel       string
	logFile        string
	nodeRPCURL     string
	nodeMetricsURL string
	id             string
	lbBaseURL      string
	payoutAddress  string
	lbURL          *url.URL
	metricsURL     *url.URL
	rpcURL         *url.URL
)

var startCmd = &cobra.Command{
	Use:   "vedran-daemon",
	Short: "Register vedran-daemon with load balancer and start sending telemetry",
	RunE:  start,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		level, err := log.ParseLevel(logLevel)
		if err != nil {
			log.Fatalf("Invalid log level %s", logLevel)
		}

		err = logger.SetupLogger(level, logFile)
		if err != nil {
			return err
		}

		return nil
	},
	Args: func(cmd *cobra.Command, args []string) error {
		var err error
		lbURL, err = url.Parse(lbBaseURL)
		if err != nil {
			return fmt.Errorf("Failed parsing load balancer url: %v", err)
		}

		metricsURL, err = url.Parse(nodeMetricsURL)
		if err != nil {
			return fmt.Errorf("Failed parsing metrics url: %v", err)
		}

		rpcURL, err = url.Parse(nodeRPCURL)
		if err != nil {
			return fmt.Errorf("Failed parsing rpc url: %v", err)
		}

		return nil
	},
}

func init() {
	startCmd.Flags().StringVar(&nodeRPCURL, "node-rpc", "http://localhost:9933", "Polkadot node rpc url")
	startCmd.Flags().StringVar(&nodeMetricsURL, "node-metrics", "http://localhost:9615", "Polkadot node metrics url")
	startCmd.Flags().StringVar(&id, "id", "", "Vedran-daemon id string (required)")
	startCmd.Flags().StringVar(&lbBaseURL, "lb", "", "Target load balancer url (required)")
	startCmd.Flags().StringVar(&payoutAddress, "payout-address", "", "Payout address to which reward tokens will be sent (required)")
	startCmd.Flags().StringVar(&logLevel, "log-level", "info", "Level of logging (eg. debug, info, warn, error)")
	startCmd.Flags().StringVar(&logFile, "log-file", "", "Path to logfile. If not set defaults to stdout")

	_ = startCmd.MarkFlagRequired("id")
	_ = startCmd.MarkFlagRequired("lb")
	_ = startCmd.MarkFlagRequired("payout-address")
}

func start(cmd *cobra.Command, _ []string) error {
	lbClient := lb.NewClient(lbURL)
	nodeClient := node.NewClient(rpcURL, metricsURL)
	telemetry := telemetry.NewTelemetry()

	err := run.Start(lbClient, nodeClient, telemetry, id, payoutAddress)
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
