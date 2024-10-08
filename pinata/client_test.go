package pinata

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("with default settings", func(t *testing.T) {
		auth := &Auth{
			jwt: "test_jwt_token",
		}
		client := New(auth)

		require.Equal(t, BaseURL, client.baseURL)
		require.Equal(t, auth, client.auth)
		require.NotNil(t, client.httpClient)
		require.NotNil(t, client.transport)

		require.Equal(t, 90*time.Second, client.httpClient.Timeout)

		transport, ok := client.httpClient.Transport.(*http.Transport)
		require.True(t, ok)
		require.Equal(t, 100, transport.MaxIdleConns)
		require.Equal(t, 100, transport.MaxIdleConnsPerHost)
		require.Equal(t, 90*time.Second, transport.IdleConnTimeout)
	})

	t.Run("with custom base URL", func(t *testing.T) {
		auth := &Auth{
			apiKey:    "test_api_key",
			apiSecret: "test_api_secret",
		}
		client := New(auth)
		client.baseURL = "https://custom.pinata.cloud"

		require.Equal(t, "https://custom.pinata.cloud", client.baseURL)
	})

	t.Run("with nil auth", func(t *testing.T) {
		client := New(nil)

		require.NotNil(t, client)
		require.Nil(t, client.auth)
	})

	t.Run("transport equality", func(t *testing.T) {
		client := New(&Auth{})

		require.Equal(t, client.transport, client.httpClient.Transport)
	})
}

func TestNewRequest(t *testing.T) {
	t.Run("basic request creation", func(t *testing.T) {
		client := New(&Auth{jwt: "test_jwt"})
		rb := client.NewRequest(http.MethodGet, "/test/path")

		require.NotNil(t, rb)
		require.Equal(t, client, rb.client)
		require.Equal(t, http.MethodGet, rb.method)
		require.Equal(t, "/test/path", rb.path)
		require.Empty(t, rb.pathParams)
		require.Empty(t, rb.queryParams)
		require.Empty(t, rb.headers)
	})

	t.Run("different HTTP methods", func(t *testing.T) {
		client := New(&Auth{jwt: "test_jwt"})
		methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}

		for _, method := range methods {
			rb := client.NewRequest(method, "/test/path")
			require.Equal(t, method, rb.method)
		}
	})

	t.Run("path with special characters", func(t *testing.T) {
		client := New(&Auth{jwt: "test_jwt"})
		rb := client.NewRequest(http.MethodGet, "/test/path with spaces/and/special-chars!@#$%^&*()")

		require.Equal(t, "/test/path with spaces/and/special-chars!@#$%^&*()", rb.path)
	})

	t.Run("empty path", func(t *testing.T) {
		client := New(&Auth{jwt: "test_jwt"})
		rb := client.NewRequest(http.MethodGet, "")

		require.Equal(t, "", rb.path)
	})
}

func TestTestAuthentication(t *testing.T) {
	t.Run("successful authentication", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/data/testAuthentication", r.URL.Path)
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "Bearer valid_jwt_token", r.Header.Get("Authorization"))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message":"Congratulations! You are authenticated"}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		response, err := client.TestAuthentication()

		require.NoError(t, err)
		require.NotNil(t, response)
		require.Equal(t, "Congratulations! You are authenticated", response.Message)
	})

	t.Run("authentication failure", func(t *testing.T) {
		auth := &Auth{jwt: "invalid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"Invalid authentication credentials"}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		response, err := client.TestAuthentication()

		require.Error(t, err)
		require.Nil(t, response)
		require.Contains(t, err.Error(), "Invalid authentication credentials")
	})

	t.Run("network error", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)
		client.baseURL = "http://non-existent-url.com"

		response, err := client.TestAuthentication()

		require.Error(t, err)
		require.Nil(t, response)
	})
}
