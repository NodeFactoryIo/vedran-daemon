package lb

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/NodeFactoryIo/vedran-daemon/internal/metrics"
	"github.com/stretchr/testify/assert"
)

func Test_metricsService_Send(t *testing.T) {
	setup()
	defer teardown()

	peerCount := float64(19)
	bestBlockHeight := float64(432933)
	finalizedBlockHeight := float64(432640)
	readyTransactionCount := float64(0)

	type args struct {
		metricsBaseURL string
	}
	tests := []struct {
		name         string
		args         args
		wantErr      bool
		lbHandleFunc handleFnMock
		want         int
	}{
		{
			name:    "Returns error if sending metrics fails",
			args:    args{""},
			wantErr: true,
			want:    0,
			lbHandleFunc: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Not Found", 404)
			}},
		{
			name:    "Returns resp if sending metrics succedes",
			args:    args{""},
			wantErr: false,
			want:    200,
			lbHandleFunc: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPut, r.Method)

				var expectedMetrics metrics.Metrics
				defer r.Body.Close()
				body, _ := ioutil.ReadAll((r.Body))
				_ = json.Unmarshal(body, &expectedMetrics)

				assert.Equal(t, expectedMetrics, metrics.Metrics{
					PeerCount:             &peerCount,
					BestBlockHeight:       &bestBlockHeight,
					FinalizedBlockHeight:  &finalizedBlockHeight,
					ReadyTransactionCount: &readyTransactionCount})

				_, _ = io.WriteString(w, `{"status": "ok"}`)
			}},
	}
	for _, tt := range tests {
		setup()

		t.Run(tt.name, func(t *testing.T) {
			mux.HandleFunc("/api/v1/nodes/metrics", tt.lbHandleFunc)
			mockURL, _ := url.Parse(server.URL)
			client := NewClient(mockURL)
			ms := &metricsService{client: client}

			got, err := ms.Send(
				&metrics.Metrics{
					PeerCount:             &peerCount,
					BestBlockHeight:       &bestBlockHeight,
					FinalizedBlockHeight:  &finalizedBlockHeight,
					ReadyTransactionCount: &readyTransactionCount})

			if (err != nil) != tt.wantErr {
				t.Errorf("metricsService.Send() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != nil && got.StatusCode != 200 {
				t.Errorf("metricsService.Send() statusCode = %d, want %d", got.StatusCode, tt.want)
				return
			}
		})

		teardown()
	}
}
