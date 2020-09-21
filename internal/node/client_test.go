package node

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
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

func TestClient_SendRPCRequest(t *testing.T) {
	setup()
	defer teardown()

	type InvalidResult []string
	type ValidResult string

	invalidResult := new(InvalidResult)
	validResult := new(ValidResult)

	type fields struct {
		client  *http.Client
		BaseURL string
	}

	type args struct {
		method string
		params []string
		v      interface{}
	}

	tests := []struct {
		name       string
		fields     fields
		args       args
		want       ValidResult
		wantErr    bool
		handleFunc handleFnMock
	}{
		{
			name:    "Returns error if server url invalid",
			args:    args{"system_chain", []string{}, nil},
			fields:  fields{http.DefaultClient, "invalid"},
			wantErr: true,
			want:    "",
			handleFunc: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Not Found", 404)
			}},
		{
			name:    "Returns error if response not 200",
			args:    args{"system_chain", []string{}, nil},
			fields:  fields{http.DefaultClient, "valid"},
			wantErr: true,
			want:    "",
			handleFunc: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Not Found", 404)
			}},
		{
			name:    "Returns error if rpc code not 200",
			args:    args{"system_chain", []string{}, validResult},
			fields:  fields{http.DefaultClient, "valid"},
			wantErr: true,
			want:    "",
			handleFunc: func(w http.ResponseWriter, r *http.Request) {
				_, _ = io.WriteString(w, `{
					"jsonrpc": "2.0",
					"error": {
						"code": -32600,
						"message": "Error"
					},
					"id": 1
				}`)
			}},
		{
			name:    "Returns error if json unmarshal fails",
			args:    args{"system_chain", []string{}, validResult},
			fields:  fields{http.DefaultClient, "valid"},
			wantErr: true,
			want:    "",
			handleFunc: func(w http.ResponseWriter, r *http.Request) {
				_, _ = io.WriteString(w, `invalid`)
			}},
		{
			name:    "Returns error if map structure decode fails",
			args:    args{"system_chain", []string{}, invalidResult},
			fields:  fields{http.DefaultClient, "valid"},
			wantErr: true,
			want:    "",
			handleFunc: func(w http.ResponseWriter, r *http.Request) {
				_, _ = io.WriteString(w, `{
					"jsonrpc": "2.0",
					"result": "Live",
					"id": 1
				}`)
			}},
		{
			name:    "Returns resp if request valid",
			args:    args{"system_chain", []string{}, validResult},
			fields:  fields{http.DefaultClient, "valid"},
			wantErr: false,
			want:    "Live",
			handleFunc: func(w http.ResponseWriter, r *http.Request) {
				_, _ = io.WriteString(w, `{
					"jsonrpc": "2.0",
					"result": "Live",
					"id": 1
				}`)
			}},
	}

	for _, tt := range tests {
		setup()

		t.Run(tt.name, func(t *testing.T) {
			var mockURL *url.URL
			if tt.fields.BaseURL == "valid" {
				mockURL, _ = url.Parse(server.URL)
			} else {
				mockURL, _ = url.Parse("http://invalid:3000")
			}
			c := &client{mockURL, mockURL}
			mux.HandleFunc("/", tt.handleFunc)

			got, err := c.sendRPCRequest(tt.args.method, tt.args.params, tt.args.v)

			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SendRPCRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != nil {
				if !reflect.DeepEqual(validResult, &tt.want) {
					t.Errorf("Client.SendRPCRequest() = %v, want %v", validResult, tt.want)
				}
			}
		})

		teardown()
	}
}
