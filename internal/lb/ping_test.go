package lb

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_pingService_Send(t *testing.T) {
	setup()
	defer teardown()

	tests := []struct {
		name         string
		wantErr      bool
		lbHandleFunc handleFnMock
		want         int
	}{
		{
			name:    "Returns error if sending ping fails",
			wantErr: true,
			want:    0,
			lbHandleFunc: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Not Found", 404)
			}},
		{
			name:    "Returns resp if sending ping succedes",
			wantErr: false,
			want:    200,
			lbHandleFunc: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)

				var pingRequest pingRequest
				defer r.Body.Close()
				body, _ := ioutil.ReadAll((r.Body))
				_ = json.Unmarshal(body, &pingRequest)

				assert.NotNil(t, pingRequest.Timestamp)
				assert.Greater(t, pingRequest.Timestamp, int64(1000))

				_, _ = io.WriteString(w, `{"status": "ok"}`)
			}},
	}
	for _, tt := range tests {
		setup()

		t.Run(tt.name, func(t *testing.T) {
			mux.HandleFunc("/api/v1/nodes/pings", tt.lbHandleFunc)
			mockURL, _ := url.Parse(server.URL)
			client := NewClient(mockURL)
			ps := &pingService{client: client}

			got, err := ps.Send()

			if (err != nil) != tt.wantErr {
				t.Errorf("pingService.Send() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != nil && got.StatusCode != 200 {
				t.Errorf("pingService.Send() statusCode = %d, want %d", got.StatusCode, tt.want)
				return
			}
		})

		teardown()
	}
}
