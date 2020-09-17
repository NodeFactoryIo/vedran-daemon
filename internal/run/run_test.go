package run

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/NodeFactoryIo/vedran-daemon/internal/lb"
	nodeMocks "github.com/NodeFactoryIo/vedran-daemon/mocks/node"
	telemetryMocks "github.com/NodeFactoryIo/vedran-daemon/mocks/telemetry"
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
	nodeClient := &nodeMocks.Client{}

	type args struct {
		client        *lb.Client
		id            string
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
			args:    args{lbClient, "test-id", "0xtestpayoutaddress"},
			wantErr: true,
			handleFunc: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Not Found", 404)
			}},
		{
			name:    "Returns nil if startSendingTelemetry succeeds",
			args:    args{lbClient, "test-id", "0xtestpayoutaddress"},
			wantErr: false,
			handleFunc: func(w http.ResponseWriter, r *http.Request) {
				_, _ = io.WriteString(w, `{"token": "test-token"}`)
			}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup()

			telemetryMock := &telemetryMocks.TelemetryInterface{}
			telemetryMock.On("StartSendingTelemetry", mock.Anything, mock.Anything, mock.Anything).Return()
			nodeClient.On("GetRPCURL").Return("http://localhost:9933")
			url, _ := url.Parse(server.URL)
			lbClient.BaseURL = url
			mux.HandleFunc("/api/v1/nodes", tt.handleFunc)

			err := Start(tt.args.client, nodeClient, telemetryMock, tt.args.id, tt.args.payoutAddress)

			if (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}

			teardown()
		})
	}
}