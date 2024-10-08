package pinata

import (
	"net/http"
	"time"
)

const BaseURL = "https://api.pinata.cloud"

// Client is the main struct for interacting with the Pinata API. It contains the necessary
// configuration and authentication details to make requests to the API.
type Client struct {
	baseURL    string
	httpClient *http.Client
	auth       *Auth
	transport  *http.Transport
}

// authTestResponse represents the response from the Pinata API's test authentication endpoint.
// It contains a message field with the result of the authentication test.
type authTestResponse struct {
	Message string `json:"message"`
}

// New creates a new Pinata API client with the provided authentication credentials.
// It configures the HTTP client with a transport that has a maximum of 100 idle connections,
// a maximum of 100 idle connections per host, and an idle connection timeout of 90 seconds.
// The HTTP client also has a timeout of 30 seconds.
func New(auth *Auth) *Client {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	}

	return &Client{
		baseURL: BaseURL,
		httpClient: &http.Client{
			Timeout:   time.Second * 90,
			Transport: transport,
		},
		auth:      auth,
		transport: transport,
	}
}

// NewRequest creates a new request builder for the Pinata API. The request builder
// allows for configuring the HTTP method, path, path parameters, query parameters,
// and headers before sending the request.
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

// TestAuthentication tests the authentication credentials configured in the Pinata API client.
// It sends a GET request to the "/data/testAuthentication" endpoint and returns the response
// message indicating whether the authentication was successful or not.
func (c *Client) TestAuthentication() (*authTestResponse, error) {
	var response authTestResponse
	err := c.NewRequest(http.MethodGet, "/data/testAuthentication").
		Send(&response)

	if err != nil {
		return nil, err
	}

	return &response, nil
}
