package node

import (
	"hash"
	"hash/fnv"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfigHash(t *testing.T) {

	type args struct {
		baseURL string
	}

	type Test struct {
		name       string
		args       args
		want       hash.Hash32
		wantErr    bool
		handleFunc handleFnMock
	}

	hash := fnv.New32()
	_, _ = hash.Write([]byte("author_rotateKeyschain_getBlock"))
	tests := []Test{
		{
			name:    "Returns error if node not started",
			args:    args{"invalid"},
			want:    nil,
			wantErr: true,
			handleFunc: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Not Found", 404)
			}},
		{
			name:    "Returns error if rpc method does not exist",
			args:    args{"valid"},
			want:    nil,
			wantErr: true,
			handleFunc: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				http.Error(w, "Not Found", 404)
			}},
		{
			name:    "Returns error if parsing rpc methods fails",
			args:    args{"valid"},
			want:    nil,
			wantErr: true,
			handleFunc: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				_, _ = io.WriteString(w, `invalid`)
			}},
		{
			name:    "Returns error if no rpc methods found",
			args:    args{"valid"},
			want:    nil,
			wantErr: true,
			handleFunc: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				w.Header().Set("Content-Type", "application/json")
				_, _ = io.WriteString(
					w,
					`{
						"jsonrpc": 2.0,
						"result": {
							"methods": [
							],
							"version": 1
						},
						"id": 1
					}`)
			}},
		{
			name:    "Returns hash if rpc response valid",
			args:    args{"valid"},
			want:    hash,
			wantErr: false,
			handleFunc: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				w.Header().Set("Content-Type", "application/json")
				_, _ = io.WriteString(
					w,
					`{
						"jsonrpc": 2.0,
						"result": {
							"methods": [
								"chain_getBlock",
								"author_rotateKeys"
							],
							"version": 1
						},
						"id": 1
					}`)
			}},
	}

	for _, tt := range tests {
		setup()
		defer teardown()

		t.Run(tt.name, func(t *testing.T) {
			var baseURL *url.URL
			if tt.args.baseURL == "valid" {
				baseURL, _ = url.Parse(server.URL)
			} else {
				baseURL, _ = url.Parse("http://invalid:3003")
			}
			mux.HandleFunc("/", tt.handleFunc)

			client := NewClient(baseURL, baseURL)
			got, err := client.GetConfigHash()

			if (err != nil) != tt.wantErr {
				t.Errorf("GetConfigHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConfigHash() = %v, want %v", got, tt.want)
			}

		})
	}
}
