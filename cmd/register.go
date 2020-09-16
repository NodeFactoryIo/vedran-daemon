package cmd

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/NodeFactoryIo/vedran-daemon/internal/lb"
	"github.com/NodeFactoryIo/vedran-daemon/internal/metrics"
	"github.com/go-co-op/gocron"
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
	lbBaseURL, err := url.Parse(lbBaseURL)
	if err != nil {
		return fmt.Errorf("Failed parsing load balancer url")
	}

	client := lb.NewClient(lbBaseURL)
	err = client.Register(id, nodeRPC, payoutAddress, "test-config-hash")
	if err == nil {
		return fmt.Errorf("Failed registering to load balancer on url %s because of: %v", lbBaseURL.String(), err)
	}
	log.Printf("Registered to load balancer %s", lbBaseURL)

	scheduler := gocron.NewScheduler(time.UTC)

	fms := &metrics.FetchMetricsService{BaseURL: nodeMetrics}
	_, err = scheduler.Every(30).Seconds().Do(client.Metrics.Send, fms)
	if err != nil {
		return fmt.Errorf("Failed starting sending metrics becuase of: %v", err)
	}
	log.Println("Started sending metrics to load balancer")

	_, err = scheduler.Every(5).Seconds().Do(client.Ping.Send)
	if err != nil {
		return fmt.Errorf("Failed starting sending pings becuase of: %v", err)
	}
	log.Println("Started sending pings to load balancer")

	scheduler.StartBlocking()
	return nil
}
