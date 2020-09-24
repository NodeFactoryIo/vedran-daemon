package run

import (
	"fmt"
	"hash"
	"hash/fnv"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

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
	sleep = func(d time.Duration) {}

	setup()
	defer teardown()

	lbURL, _ := url.Parse(server.URL)
	lbClient := lb.NewClient(lbURL)
	nodeClient := &nodeMocks.Client{}
	testHash := fnv.New32()

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
		firstGetConfigHashResult    hash.Hash32
		firstGetConfigHashError     error
		secondGetConfigHashResult   hash.Hash32
		secondGetConfigHashError    error
	}{
		{
			name:    "Retries get config hash if get config hash fails and returns error if register fails",
			args:    args{lbClient, "test-id", "0xtestpayoutaddress"},
			wantErr: true,
			handleFunc: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Not Found", 404)
			},
			firstGetConfigHashError:   fmt.Errorf("Error"),
			firstGetConfigHashResult:  nil,
			secondGetConfigHashError:  nil,
			secondGetConfigHashResult: testHash},
		{
			name:    "Returns nil if startSendingTelemetry succeeds",
			args:    args{lbClient, "test-id", "0xtestpayoutaddress"},
			wantErr: false,
			handleFunc: func(w http.ResponseWriter, r *http.Request) {
				_, _ = io.WriteString(w, `{"token": "test-token"}`)
			},
			firstGetConfigHashError:  nil,
			firstGetConfigHashResult: testHash},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup()

			telemetryMock := &telemetryMocks.Telemetry{}
			telemetryMock.On("StartSendingTelemetry", mock.Anything, mock.Anything, mock.Anything).Return()
			nodeClient.On("GetRPCURL").Return("http://localhost:9933")
			nodeClient.On("GetConfigHash").Once().Return(tt.firstGetConfigHashResult, tt.firstGetConfigHashError)
			nodeClient.On("GetConfigHash").Once().Return(tt.secondGetConfigHashResult, tt.secondGetConfigHashError)

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
