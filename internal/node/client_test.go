package node

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
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

func Test_client_GetRPCURL(t *testing.T) {
	rpcURL, _ := url.Parse("http://localhost:9933")
	metricsURL, _ := url.Parse("http://localhost:9615")

	type fields struct {
		MetricsBaseURL *url.URL
		RPCBaseURL     *url.URL
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Returns string representation of rpc url",
			fields: fields{metricsURL, rpcURL},
			want:   "http://localhost:9933"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &client{
				MetricsBaseURL: tt.fields.MetricsBaseURL,
				RPCBaseURL:     tt.fields.RPCBaseURL,
			}

			if got := client.GetRPCURL(); got != tt.want {
				t.Errorf("client.GetRPCURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_client_GetMetricsURL(t *testing.T) {
	rpcURL, _ := url.Parse("http://localhost:9933")
	metricsURL, _ := url.Parse("http://localhost:9615")

	type fields struct {
		MetricsBaseURL *url.URL
		RPCBaseURL     *url.URL
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Returns string representation of rpc url",
			fields: fields{metricsURL, rpcURL},
			want:   "http://localhost:9615"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &client{
				MetricsBaseURL: tt.fields.MetricsBaseURL,
				RPCBaseURL:     tt.fields.RPCBaseURL,
			}

			if got := client.GetMetricsURL(); got != tt.want {
				t.Errorf("client.GetMetricsURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
