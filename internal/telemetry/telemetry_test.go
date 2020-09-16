package telemetry

import (
	"net/url"
	"testing"
	"time"

	"github.com/NodeFactoryIo/vedran-daemon/internal/lb"
	mocks "github.com/NodeFactoryIo/vedran-daemon/mocks/scheduler"
	"github.com/go-co-op/gocron"
	"github.com/stretchr/testify/mock"
)

func TestTelemetry_StartSendingTelemetry(t *testing.T) {
	lbURL, _ := url.Parse("localhost:4000")
	client := lb.NewClient(lbURL)

	type args struct {
		client      *lb.Client
		nodeMetrics string
	}

	tests := []struct {
		name               string
		t                  *Telemetry
		args               args
		expectedNumOfCalls int
	}{
		{
			name:               "Calls start blocking with ping and metrics jobs",
			args:               args{client, "localhost:9615"},
			expectedNumOfCalls: 1},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			mockScheduler := &mocks.Scheduler{}
			telemetry := &Telemetry{}
			mockScheduler.On("StartBlocking").Return()
			mockScheduler.On("Every", mock.Anything).Return(gocron.NewScheduler(time.UTC).Every(10))

			telemetry.StartSendingTelemetry(mockScheduler, tt.args.client, tt.args.nodeMetrics)

			mockScheduler.AssertNumberOfCalls(t, "StartBlocking", tt.expectedNumOfCalls)
		})
	}
}
