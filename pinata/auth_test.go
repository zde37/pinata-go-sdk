package pinata

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewAuth(t *testing.T) {
	t.Run("with all fields provided", func(t *testing.T) {
		auth := NewAuth("test_api_key", "test_api_secret", "test_jwt_token")
		require.NotNil(t, auth)
		require.Equal(t, "test_api_key", auth.apiKey)
		require.Equal(t, "test_api_secret", auth.apiSecret)
		require.Equal(t, "test_jwt_token", auth.jwt)
	})

	t.Run("with only API key and secret", func(t *testing.T) {
		auth := NewAuth("test_api_key", "test_api_secret", "")
		require.NotNil(t, auth)
		require.Equal(t, "test_api_key", auth.apiKey)
		require.Equal(t, "test_api_secret", auth.apiSecret)
		require.Empty(t, auth.jwt)
	})

	t.Run("with only JWT", func(t *testing.T) {
		auth := NewAuth("", "", "test_jwt_token")
		require.NotNil(t, auth)
		require.Empty(t, auth.apiKey)
		require.Empty(t, auth.apiSecret)
		require.Equal(t, "test_jwt_token", auth.jwt)
	})

	t.Run("with empty fields", func(t *testing.T) {
		auth := NewAuth("", "", "")
		require.NotNil(t, auth)
		require.Empty(t, auth.apiKey)
		require.Empty(t, auth.apiSecret)
		require.Empty(t, auth.jwt)
	})
}

func TestNewAuthWithJWT(t *testing.T) {
	t.Run("with valid JWT", func(t *testing.T) {
		jwt := "valid_jwt_token"
		auth := NewAuthWithJWT(jwt)
		require.NotNil(t, auth)
		require.Equal(t, jwt, auth.jwt)
		require.Empty(t, auth.apiKey)
		require.Empty(t, auth.apiSecret)
	})

	t.Run("with empty JWT", func(t *testing.T) {
		auth := NewAuthWithJWT("")
		require.NotNil(t, auth)
		require.Empty(t, auth.jwt)
		require.Empty(t, auth.apiKey)
		require.Empty(t, auth.apiSecret)
	})

	t.Run("with whitespace JWT", func(t *testing.T) {
		auth := NewAuthWithJWT("   ")
		require.NotNil(t, auth)
		require.Equal(t, "   ", auth.jwt)
		require.Empty(t, auth.apiKey)
		require.Empty(t, auth.apiSecret)
	})

	t.Run("setAuthHeader with JWT from NewAuthWithJWT", func(t *testing.T) {
		jwt := "test_jwt_from_new_auth"
		auth := NewAuthWithJWT(jwt)
		req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)

		auth.setAuthHeader(req)

		authHeader := req.Header.Get("Authorization")
		require.Equal(t, "Bearer "+jwt, authHeader)
		require.Empty(t, req.Header.Get("pinata_api_key"))
		require.Empty(t, req.Header.Get("pinata_secret_api_key"))
	})
}

func TestSetAuthHeader(t *testing.T) {
	t.Run("with JWT", func(t *testing.T) {
		auth := &auth{
			jwt: "test_jwt_token",
		}
		req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)

		auth.setAuthHeader(req)

		authHeader := req.Header.Get("Authorization")
		require.Equal(t, "Bearer test_jwt_token", authHeader)
		require.Empty(t, req.Header.Get("pinata_api_key"))
		require.Empty(t, req.Header.Get("pinata_secret_api_key"))
	})

	t.Run("with API key and secret", func(t *testing.T) {
		auth := &auth{
			apiKey:    "test_api_key",
			apiSecret: "test_api_secret",
		}
		req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)

		auth.setAuthHeader(req)

		require.Empty(t, req.Header.Get("Authorization"))
		require.Equal(t, "test_api_key", req.Header.Get("pinata_api_key"))
		require.Equal(t, "test_api_secret", req.Header.Get("pinata_secret_api_key"))
	})

	t.Run("with empty auth", func(t *testing.T) {
		auth := &auth{}
		req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)

		auth.setAuthHeader(req)

		require.Empty(t, req.Header.Get("Authorization"))
		require.Empty(t, req.Header.Get("pinata_api_key"))
		require.Empty(t, req.Header.Get("pinata_secret_api_key"))
	})

	t.Run("JWT takes precedence over API key and secret", func(t *testing.T) {
		auth := &auth{
			jwt:       "test_jwt_token",
			apiKey:    "test_api_key",
			apiSecret: "test_api_secret",
		}
		req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)

		auth.setAuthHeader(req)

		authHeader := req.Header.Get("Authorization")
		require.Equal(t, "Bearer test_jwt_token", authHeader)
		require.Empty(t, req.Header.Get("pinata_api_key"))
		require.Empty(t, req.Header.Get("pinata_secret_api_key"))
	})
}
