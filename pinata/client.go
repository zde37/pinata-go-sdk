package pinata

import (
	"net/http"
	"time"
)

const BaseURL = "https://api.pinata.cloud"

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Auth       *Auth
	Transport  *http.Transport
}

type AuthTestResponse struct {
	Message string `json:"message"`
}

func NewClient(auth *Auth) *Client {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	}

	return &Client{
		BaseURL: BaseURL,
		HTTPClient: &http.Client{
			Timeout:   time.Second * 30,
			Transport: transport,
		},
		Auth:      auth,
		Transport: transport,
	}
}

func (c *Client) NewRequest(method, path string) *requestBuilder {
	return &requestBuilder{
		client:      c,
		method:      method,
		path:        path,
		pathParams:  make(map[string]string),
		queryParams: make(map[string]string),
		headers:     make(map[string]string),
	}
}

func (c *Client) TestAuthentication() (*AuthTestResponse, error) {
    var response AuthTestResponse
    err := c.NewRequest("GET", "/data/testAuthentication").
        Send(&response)

    if err != nil {
        return nil, err
    }

    return &response, nil
}