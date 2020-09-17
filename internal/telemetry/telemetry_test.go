package telemetry

import (
	"net/url"
	"testing"
	"time"

	"github.com/NodeFactoryIo/vedran-daemon/internal/lb"
	nodeMocks "github.com/NodeFactoryIo/vedran-daemon/mocks/node"
	schedulerMocks "github.com/NodeFactoryIo/vedran-daemon/mocks/scheduler"
	"github.com/go-co-op/gocron"
	"github.com/stretchr/testify/mock"
)

func TestTelemetry_StartSendingTelemetry(t *testing.T) {
	lbURL, _ := url.Parse("localhost:4000")
	lbClient := lb.NewClient(lbURL)
	nodeClient := &nodeMocks.Client{}

	tests := []struct {
		name               string
		t                  *Telemetry
		expectedNumOfCalls int
	}{
		{
			name:               "Calls start blocking with ping and metrics jobs",
			expectedNumOfCalls: 1},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			mockScheduler := &schedulerMocks.Scheduler{}
			telemetry := &Telemetry{}
			mockScheduler.On("StartBlocking").Return()
			mockScheduler.On("Every", mock.Anything).Return(gocron.NewScheduler(time.UTC).Every(10))

			telemetry.StartSendingTelemetry(mockScheduler, lbClient, nodeClient)

			mockScheduler.AssertNumberOfCalls(t, "StartBlocking", tt.expectedNumOfCalls)
		})
	}
}
