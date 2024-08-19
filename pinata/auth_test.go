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
		require.Equal(t, "test_api_key", auth.APIKey)
		require.Equal(t, "test_api_secret", auth.APISecret)
		require.Equal(t, "test_jwt_token", auth.JWT)
	})

	t.Run("with only API key and secret", func(t *testing.T) {
		auth := NewAuth("test_api_key", "test_api_secret", "")
		require.NotNil(t, auth)
		require.Equal(t, "test_api_key", auth.APIKey)
		require.Equal(t, "test_api_secret", auth.APISecret)
		require.Empty(t, auth.JWT)
	})

	t.Run("with only JWT", func(t *testing.T) {
		auth := NewAuth("", "", "test_jwt_token")
		require.NotNil(t, auth)
		require.Empty(t, auth.APIKey)
		require.Empty(t, auth.APISecret)
		require.Equal(t, "test_jwt_token", auth.JWT)
	})

	t.Run("with empty fields", func(t *testing.T) {
		auth := NewAuth("", "", "")
		require.NotNil(t, auth)
		require.Empty(t, auth.APIKey)
		require.Empty(t, auth.APISecret)
		require.Empty(t, auth.JWT)
	})
}

func TestNewAuthWithJWT(t *testing.T) {
	t.Run("with valid JWT", func(t *testing.T) {
		jwt := "valid_jwt_token"
		auth := NewAuthWithJWT(jwt)
		require.NotNil(t, auth)
		require.Equal(t, jwt, auth.JWT)
		require.Empty(t, auth.APIKey)
		require.Empty(t, auth.APISecret)
	})

	t.Run("with empty JWT", func(t *testing.T) {
		auth := NewAuthWithJWT("")
		require.NotNil(t, auth)
		require.Empty(t, auth.JWT)
		require.Empty(t, auth.APIKey)
		require.Empty(t, auth.APISecret)
	})

	t.Run("with whitespace JWT", func(t *testing.T) {
		auth := NewAuthWithJWT("   ")
		require.NotNil(t, auth)
		require.Equal(t, "   ", auth.JWT)
		require.Empty(t, auth.APIKey)
		require.Empty(t, auth.APISecret)
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
		auth := &Auth{
			JWT: "test_jwt_token",
		}
		req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)

		auth.setAuthHeader(req)

		authHeader := req.Header.Get("Authorization")
		require.Equal(t, "Bearer test_jwt_token", authHeader)
		require.Empty(t, req.Header.Get("pinata_api_key"))
		require.Empty(t, req.Header.Get("pinata_secret_api_key"))
	})

	t.Run("with API key and secret", func(t *testing.T) {
		auth := &Auth{
			APIKey:    "test_api_key",
			APISecret: "test_api_secret",
		}
		req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)

		auth.setAuthHeader(req)

		require.Empty(t, req.Header.Get("Authorization"))
		require.Equal(t, "test_api_key", req.Header.Get("pinata_api_key"))
		require.Equal(t, "test_api_secret", req.Header.Get("pinata_secret_api_key"))
	})

	t.Run("with empty auth", func(t *testing.T) {
		auth := &Auth{}
		req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)

		auth.setAuthHeader(req)

		require.Empty(t, req.Header.Get("Authorization"))
		require.Empty(t, req.Header.Get("pinata_api_key"))
		require.Empty(t, req.Header.Get("pinata_secret_api_key"))
	})

	t.Run("JWT takes precedence over API key and secret", func(t *testing.T) {
		auth := &Auth{
			JWT:       "test_jwt_token",
			APIKey:    "test_api_key",
			APISecret: "test_api_secret",
		}
		req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)

		auth.setAuthHeader(req)

		authHeader := req.Header.Get("Authorization")
		require.Equal(t, "Bearer test_jwt_token", authHeader)
		require.Empty(t, req.Header.Get("pinata_api_key"))
		require.Empty(t, req.Header.Get("pinata_secret_api_key"))
	})
}
