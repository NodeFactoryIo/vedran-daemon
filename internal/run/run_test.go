package run

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/NodeFactoryIo/vedran-daemon/internal/lb"
	mocks "github.com/NodeFactoryIo/vedran-daemon/mocks/telemetry"
	"github.com/stretchr/testify/mock"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
)

type handleFnMock func(http.ResponseWriter, *http.Request)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
}

func teardown() {
	server.Close()
}

func TestStart(t *testing.T) {
	setup()
	defer teardown()

	lbURL, _ := url.Parse(server.URL)
	lbClient := lb.NewClient(lbURL)
	type args struct {
		client        *lb.Client
		id            string
		nodeRPC       string
		nodeMetrics   string
		payoutAddress string
	}

	tests := []struct {
		name                        string
		args                        args
		wantErr                     bool
		handleFunc                  handleFnMock
		startSendingTelemetryResult error
	}{
		{
			name:    "Returns error if lb register fails",
			args:    args{lbClient, "test-id", "localhost:9933", "localhost:9615", "0xtestpayoutaddress"},
			wantErr: true,
			handleFunc: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Not Found", 404)
			},
			startSendingTelemetryResult: nil},
		{
			name:    "Returns error if startSendingTelemetry fails",
			args:    args{lbClient, "test-id", "localhost:9933", "localhost:9615", "0xtestpayoutaddress"},
			wantErr: true,
			handleFunc: func(w http.ResponseWriter, r *http.Request) {
				_, _ = io.WriteString(w, `{"token": "test-token"}`)
			},
			startSendingTelemetryResult: fmt.Errorf("Errpr")},
		{
			name:    "Returns nil if startSendingTelemetry succeeds",
			args:    args{lbClient, "test-id", "localhost:9933", "localhost:9615", "0xtestpayoutaddress"},
			wantErr: false,
			handleFunc: func(w http.ResponseWriter, r *http.Request) {
				_, _ = io.WriteString(w, `{"token": "test-token"}`)
			},
			startSendingTelemetryResult: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup()
			telemetryMock := &mocks.TelemetryInterface{}
			telemetryMock.On("StartSendingTelemetry", mock.Anything, mock.Anything, mock.Anything).Return(tt.startSendingTelemetryResult)
			url, _ := url.Parse(server.URL)
			lbClient.BaseURL = url
			mux.HandleFunc("/api/v1/nodes", tt.handleFunc)

			err := Start(tt.args.client, telemetryMock, tt.args.id, tt.args.nodeRPC, tt.args.nodeMetrics, tt.args.payoutAddress)

			if (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}

			teardown()
		})
	}
}
