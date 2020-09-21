package lb

import (
	"bytes"
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

func TestNewClient(t *testing.T) {
	expectedURL, _ := url.Parse("http://url.com")
	expectedClient := &Client{client: http.DefaultClient, BaseURL: expectedURL}
	expectedClient.Ping = &pingService{client: expectedClient}
	expectedClient.Metrics = &metricsService{client: expectedClient}

	type args struct {
		baseURL *url.URL
	}

	tests := []struct {
		name string
		args args
		want *Client
	}{
		{
			name: "Returns instance of client",
			args: args{expectedURL},
			want: expectedClient},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			got := NewClient(tt.args.baseURL)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}

		})

	}
}

func TestClient_Register(t *testing.T) {
	setup()
	defer teardown()

	type fields struct {
		client  *http.Client
		BaseURL string
		Token   string
	}
	type args struct {
		id            string
		nodeURL       string
		payoutAddress string
		configHash    string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		want       string
		handleFunc handleFnMock
	}{
		{
			name:    "Returns error if client does not return 200",
			args:    args{"test-id", "http://localhost:3000", "0xtestaddress", "config-hash"},
			fields:  fields{http.DefaultClient, "valid", ""},
			wantErr: true,
			want:    "",
			handleFunc: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Not Found", 404)
			}},
		{
			name:    "Sets token on client if request valid",
			args:    args{"test-id", "http://localhost:3000", "0xtestaddress", "config-hash"},
			fields:  fields{http.DefaultClient, "valid", ""},
			wantErr: false,
			want:    "test-token",
			handleFunc: func(w http.ResponseWriter, r *http.Request) {
				_, _ = io.WriteString(w, `{"token": "test-token"}`)
			}},
	}
	for _, tt := range tests {
		setup()

		t.Run(tt.name, func(t *testing.T) {
			var mockURL *url.URL
			if tt.fields.BaseURL == "valid" {
				mockURL, _ = url.Parse(server.URL + "/")
			} else {
				mockURL, _ = url.Parse("http://invalid:3000")
			}
			c := &Client{
				client:  tt.fields.client,
				BaseURL: mockURL,
				Token:   tt.fields.Token,
			}
			mux.HandleFunc("/api/v1/nodes", tt.handleFunc)

			err := c.Register(tt.args.id, tt.args.nodeURL, tt.args.payoutAddress, tt.args.configHash)

			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Register() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.want != c.Token {
				t.Errorf("Client.Register() token = %s, want %s", c.Token, tt.want)
			}

		})

		teardown()
	}
}

func TestClient_newRequest(t *testing.T) {
	baseURL, _ := url.Parse("http://url.com")
	expectedURL, _ := url.Parse("http://url.com" + "/url")
	expectedRequest, _ := http.NewRequest("GET", expectedURL.String(), new(bytes.Buffer))
	expectedRequest.Header.Add("Content-Type", "application/json")
	expectedRequest.Header.Add("Accept", "application/json")

	type fields struct {
		client  *http.Client
		BaseURL *url.URL
		Token   string
	}

	type args struct {
		method string
		urlStr string
		body   interface{}
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *http.Request
		wantErr bool
	}{
		{
			name:    "Returns error if url invalid",
			args:    args{http.MethodGet, "$%^&*(proper$#$%%^(password", nil},
			fields:  fields{http.DefaultClient, baseURL, ""},
			wantErr: true,
			want:    nil},
		{
			name:    "Returns error if body invalid",
			args:    args{http.MethodPost, "/url", make(chan int)},
			fields:  fields{http.DefaultClient, baseURL, ""},
			wantErr: true,
			want:    nil},
		{
			name:    "Returns error if method invalid",
			args:    args{"invalid method", "/url", nil},
			fields:  fields{http.DefaultClient, baseURL, ""},
			wantErr: true,
			want:    nil},
		{
			name:    "Creates request if data valid",
			args:    args{http.MethodGet, "/url", nil},
			fields:  fields{http.DefaultClient, baseURL, ""},
			wantErr: false,
			want:    expectedRequest},
		{
			name:    "Creates request with x-auth-header if client has token",
			args:    args{http.MethodGet, "/url", nil},
			fields:  fields{http.DefaultClient, baseURL, "token"},
			wantErr: false,
			want:    expectedRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			c := &Client{
				client:  tt.fields.client,
				BaseURL: tt.fields.BaseURL,
				Token:   tt.fields.Token,
			}

			got, err := c.newRequest(tt.args.method, tt.args.urlStr, tt.args.body)

			if (err != nil) != tt.wantErr {
				t.Errorf("Client.newRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != nil {
				if !reflect.DeepEqual(got.URL, tt.want.URL) && !reflect.DeepEqual(got.Method, tt.want.Method) {
					t.Errorf("Client.newRequest() = %v, want %v", got, tt.want)
				}
			}

			if c.Token != "" {
				if c.Token != got.Header.Get("X-Auth-Header") {
					t.Errorf("Client.newRequest() token = %s, want %s", got.Header.Get("X-Auth-Header"), c.Token)
					return
				}
			}

		})
	}
}

func TestClient_do(t *testing.T) {
	setup()
	defer teardown()

	type ExpectedResponse struct {
		Foo string `json:"foo"`
	}
	var testBody = new(ExpectedResponse)

	type fields struct {
		client  *http.Client
		BaseURL string
		Token   string
	}

	type args struct {
		url string
		v   interface{}
	}

	tests := []struct {
		name       string
		fields     fields
		args       args
		want       *ExpectedResponse
		wantErr    bool
		handleFunc handleFnMock
	}{
		{
			name:    "Returns error if status code not 200",
			args:    args{"invalid", nil},
			fields:  fields{http.DefaultClient, "valid", ""},
			wantErr: true,
			want:    nil,
			handleFunc: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Not Found", 404)
			}},
		{
			name:    "Returns error if server url invalid",
			args:    args{"invalid", nil},
			fields:  fields{http.DefaultClient, "invalid", ""},
			wantErr: true,
			want:    nil,
			handleFunc: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Not Found", 404)
			}},
		{
			name:    "Returns error if json unmarshal fails",
			args:    args{"/", testBody},
			fields:  fields{http.DefaultClient, "valid", ""},
			wantErr: true,
			want:    nil,
			handleFunc: func(w http.ResponseWriter, r *http.Request) {
				_, _ = io.WriteString(w, `invalid`)
			}},
		{
			name:    "Returns resp if request valid",
			args:    args{"/", testBody},
			fields:  fields{http.DefaultClient, "valid", ""},
			wantErr: false,
			want:    &ExpectedResponse{Foo: "bar"},
			handleFunc: func(w http.ResponseWriter, r *http.Request) {
				_, _ = io.WriteString(w, `{"foo": "bar"}`)
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
			c := &Client{
				client:  tt.fields.client,
				BaseURL: mockURL,
				Token:   tt.fields.Token,
			}
			req, _ := c.newRequest(http.MethodGet, tt.args.url, tt.args.v)
			mux.HandleFunc("/", tt.handleFunc)

			got, err := c.do(req, tt.args.v)

			if (err != nil) != tt.wantErr {
				t.Errorf("Client.do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != nil {
				if !reflect.DeepEqual(testBody, tt.want) {
					t.Errorf("Client.do() = %v, want %v", testBody, tt.want)
				}
			}
		})

		teardown()
	}
}
