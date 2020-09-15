package lb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	mediaType        = "application/json"
	registerEndpoint = "/api/v1/nodes"
)

// RegisterRequest contains data needed to connect daemon with lb
type RegisterRequest struct {
	ID            string `json:"id"`
	ConfigHash    string `json:"config_hash"`
	NodeURL       string `json:"node_url"`
	PayoutAddress string `json:"payout_address"`
}

// TokenResponse from lb register endpoint
type TokenResponse struct {
	Token string `json:"token"`
}

// Client used to communicate with vedran load balancer
type Client struct {
	client  *http.Client
	BaseURL *url.URL
	Token   string
}

// NewClient creates vedran load balancer client instance
func NewClient(baseURL *url.URL) *Client {
	httpClient := http.DefaultClient
	c := &Client{client: httpClient, BaseURL: baseURL}

	return c
}

// Register daemon with load balancer and store token in client
func (c *Client) Register(id string, nodeURL string, payoutAddress string, configHash string) error {
	body := &RegisterRequest{
		ID:            id,
		NodeURL:       nodeURL,
		PayoutAddress: payoutAddress,
		ConfigHash:    configHash,
	}
	req, _ := c.NewRequest(http.MethodPost, registerEndpoint, body)
	tokenResponse := new(TokenResponse)
	_, err := c.Do(req, tokenResponse)

	if tokenResponse.Token != "" {
		c.Token = tokenResponse.Token
	}

	return err
}

// NewRequest creates an API request. A relative URL can be provided in urlStr, which will be resolved to the
// BaseURL of the Client. If Client contains token X-Auth-Header will be added to request.
func (c *Client) NewRequest(method string, urlStr string, body interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(c.BaseURL.Path + urlStr)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if body != nil {
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)
	if c.Token != "" {
		req.Header.Add("X-Auth-Header", c.Token)
	}
	return req, nil
}

// Do sends an API request and returns the API response. The API response is JSON decoded and stored in the value
// pointed to by v, or returned as an error if an API error has occurred
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Request %v returned invalid status code %d", req, resp.StatusCode)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll((resp.Body))
	err = json.Unmarshal(body, &v)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
