package pinata

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateApiKey(t *testing.T) {
	t.Run("successful API key generation", func(t *testing.T) {
		auth := &Auth{JWT: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/users/generateApiKey", r.URL.Path)
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "Bearer valid_jwt_token", r.Header.Get("Authorization"))
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))

			var payload GenerateApiKeyOptions
			err := json.NewDecoder(r.Body).Decode(&payload)
			require.NoError(t, err)
			require.Equal(t, "test_key", payload.KeyName)

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"pinata_api_key":"generated_api_key","pinata_api_secret":"generated_api_secret"}`))
		}))
		defer mockServer.Close()
		client.BaseURL = mockServer.URL

		options := &GenerateApiKeyOptions{KeyName: "test_key"}
		secret, err := client.GenerateApiKey(options)

		require.NoError(t, err)
		require.NotNil(t, secret)
		require.Equal(t, "generated_api_key", secret.PinataApiKey)
		require.Equal(t, "generated_api_secret", secret.PinataApiSecret)
	})

	t.Run("nil options", func(t *testing.T) {
		auth := &Auth{JWT: "valid_jwt_token"}
		client := New(auth)

		secret, err := client.GenerateApiKey(nil)

		require.Error(t, err)
		require.Nil(t, secret)
		require.Contains(t, err.Error(), "options cannot be nil")
	})

	t.Run("server error", func(t *testing.T) {
		auth := &Auth{JWT: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"Internal server error"}`))
		}))
		defer mockServer.Close()
		client.BaseURL = mockServer.URL

		options := &GenerateApiKeyOptions{KeyName: "test_key"}
		secret, err := client.GenerateApiKey(options)

		require.Error(t, err)
		require.Nil(t, secret)
		require.Contains(t, err.Error(), "Internal server error")
	})

	t.Run("invalid JSON response", func(t *testing.T) {
		auth := &Auth{JWT: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"API_KEY":"generated_api_key","API_SECRET":}`))
		}))
		defer mockServer.Close()
		client.BaseURL = mockServer.URL

		options := &GenerateApiKeyOptions{KeyName: "test_key"}
		secret, err := client.GenerateApiKey(options)

		require.Error(t, err)
		require.Nil(t, secret)
		require.Contains(t, err.Error(), "invalid character")
	})
}

func TestGenerateApiKeyV3(t *testing.T) {
	t.Run("successful API key generation", func(t *testing.T) {
		auth := &Auth{JWT: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/v3/pinata/keys", r.URL.Path)
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "Bearer valid_jwt_token", r.Header.Get("Authorization"))
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))

			var payload GenerateApiKeyOptions
			err := json.NewDecoder(r.Body).Decode(&payload)
			require.NoError(t, err)
			require.Equal(t, "test_key_v3", payload.KeyName)

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"pinata_api_key":"generated_api_key_v3","pinata_api_secret":"generated_api_secret_v3"}`))
		}))
		defer mockServer.Close()
		client.BaseURL = mockServer.URL

		options := &GenerateApiKeyOptions{KeyName: "test_key_v3"}
		secret, err := client.GenerateApiKeyV3(options)

		require.NoError(t, err)
		require.NotNil(t, secret)
		require.Equal(t, "generated_api_key_v3", secret.PinataApiKey)
		require.Equal(t, "generated_api_secret_v3", secret.PinataApiSecret)
	})

	t.Run("nil options", func(t *testing.T) {
		auth := &Auth{JWT: "valid_jwt_token"}
		client := New(auth)

		secret, err := client.GenerateApiKeyV3(nil)

		require.Error(t, err)
		require.Nil(t, secret)
		require.Contains(t, err.Error(), "options cannot be nil")
	})

	t.Run("server error response", func(t *testing.T) {
		auth := &Auth{JWT: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"Internal server error"}`))
		}))
		defer mockServer.Close()
		client.BaseURL = mockServer.URL

		options := &GenerateApiKeyOptions{KeyName: "test_key_v3"}
		secret, err := client.GenerateApiKeyV3(options)

		require.Error(t, err)
		require.Nil(t, secret)
		require.Contains(t, err.Error(), "Internal server error")
	})

	t.Run("invalid JSON response", func(t *testing.T) {
		auth := &Auth{JWT: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"API_KEY":"generated_api_key_v3","API_SECRET":}`))
		}))
		defer mockServer.Close()
		client.BaseURL = mockServer.URL

		options := &GenerateApiKeyOptions{KeyName: "test_key_v3"}
		secret, err := client.GenerateApiKeyV3(options)

		require.Error(t, err)
		require.Nil(t, secret)
		require.Contains(t, err.Error(), "invalid character")
	})
}

func TestListApiKeys(t *testing.T) {
	t.Run("successful API key listing", func(t *testing.T) {
		auth := &Auth{JWT: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/users/apiKeys", r.URL.Path)
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "Bearer valid_jwt_token", r.Header.Get("Authorization"))

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"keys": [{"key": "api_key_1"}, {"key": "api_key_2"}]}`))
		}))
		defer mockServer.Close()
		client.BaseURL = mockServer.URL

		response, err := client.ListApiKeys()

		require.NoError(t, err)
		require.NotNil(t, response)
		require.Len(t, response.Keys, 2)
		require.Equal(t, "api_key_1", response.Keys[0].Key)
		require.Equal(t, "api_key_2", response.Keys[1].Key)
	})

	t.Run("empty API key list", func(t *testing.T) {
		auth := &Auth{JWT: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"keys": []}`))
		}))
		defer mockServer.Close()
		client.BaseURL = mockServer.URL

		response, err := client.ListApiKeys()

		require.NoError(t, err)
		require.NotNil(t, response)
		require.Empty(t, response.Keys)
	})

	t.Run("server error", func(t *testing.T) {
		auth := &Auth{JWT: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"Internal server error"}`))
		}))
		defer mockServer.Close()
		client.BaseURL = mockServer.URL

		response, err := client.ListApiKeys()

		require.Error(t, err)
		require.Nil(t, response)
		require.Contains(t, err.Error(), "Internal server error")
	})

	t.Run("invalid JSON response", func(t *testing.T) {
		auth := &Auth{JWT: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"keys": [{"key": "api_key_1"}, {]}`))
		}))
		defer mockServer.Close()
		client.BaseURL = mockServer.URL

		response, err := client.ListApiKeys()

		require.Error(t, err)
		require.Nil(t, response)
		require.Contains(t, err.Error(), "invalid character")
	})

	t.Run("unauthorized request", func(t *testing.T) {
		auth := &Auth{JWT: "invalid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"Unauthorized"}`))
		}))
		defer mockServer.Close()
		client.BaseURL = mockServer.URL

		response, err := client.ListApiKeys()

		require.Error(t, err)
		require.Nil(t, response)
		require.Contains(t, err.Error(), "Unauthorized")
	})
}

func TestListApiKeyV3(t *testing.T) {
	t.Run("successful API key listing with options", func(t *testing.T) {
		auth := &Auth{JWT: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/v3/pinata/keys", r.URL.Path)
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "Bearer valid_jwt_token", r.Header.Get("Authorization"))
			require.Equal(t, "20", r.URL.Query().Get("offset"))

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"keys": [{"key": "api_key_1"}, {"key": "api_key_2"}], "count": 2}`))
		}))
		defer mockServer.Close()
		client.BaseURL = mockServer.URL

		options := &ListApiKeysOptions{
			Offset: 20,
		}
		response, err := client.ListApiKeyV3(options)

		require.NoError(t, err)
		require.NotNil(t, response)
		require.Len(t, response.Keys, 2)
		require.Equal(t, "api_key_1", response.Keys[0].Key)
		require.Equal(t, "api_key_2", response.Keys[1].Key)
		require.Equal(t, 2, response.Count)
	})

	t.Run("successful API key listing without options", func(t *testing.T) {
		auth := &Auth{JWT: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/v3/pinata/keys", r.URL.Path)
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "Bearer valid_jwt_token", r.Header.Get("Authorization"))
			require.Empty(t, r.URL.Query().Get("offset"))

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"keys": [{"key": "api_key_1"}], "count": 1}`))
		}))
		defer mockServer.Close()
		client.BaseURL = mockServer.URL

		response, err := client.ListApiKeyV3(nil)

		require.NoError(t, err)
		require.NotNil(t, response)
		require.Len(t, response.Keys, 1)
		require.Equal(t, "api_key_1", response.Keys[0].Key)
		require.Equal(t, 1, response.Count)
	})

	t.Run("server error response", func(t *testing.T) {
		auth := &Auth{JWT: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"Internal server error"}`))
		}))
		defer mockServer.Close()
		client.BaseURL = mockServer.URL

		response, err := client.ListApiKeyV3(nil)

		require.Error(t, err)
		require.Nil(t, response)
		require.Contains(t, err.Error(), "Internal server error")
	})

	t.Run("invalid JSON response", func(t *testing.T) {
		auth := &Auth{JWT: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"keys": [{"key": "api_key_1"}, {]}`))
		}))
		defer mockServer.Close()
		client.BaseURL = mockServer.URL

		response, err := client.ListApiKeyV3(nil)

		require.Error(t, err)
		require.Nil(t, response)
		require.Contains(t, err.Error(), "invalid character")
	})

	t.Run("empty API key list", func(t *testing.T) {
		auth := &Auth{JWT: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"keys": [], "count": 0}`))
		}))
		defer mockServer.Close()
		client.BaseURL = mockServer.URL

		response, err := client.ListApiKeyV3(nil)

		require.NoError(t, err)
		require.NotNil(t, response)
		require.Empty(t, response.Keys)
		require.Equal(t, 0, response.Count)
	})
}

func TestRevokeApiKey(t *testing.T) {
	t.Run("successful API key revocation", func(t *testing.T) {
		auth := &Auth{JWT: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/users/revokeApiKey", r.URL.Path)
			require.Equal(t, http.MethodPut, r.Method)
			require.Equal(t, "Bearer valid_jwt_token", r.Header.Get("Authorization"))
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))

			var payload map[string]string
			err := json.NewDecoder(r.Body).Decode(&payload)
			require.NoError(t, err)
			require.Equal(t, "test_api_key", payload["apiKey"])

			w.WriteHeader(http.StatusOK)
		}))
		defer mockServer.Close()
		client.BaseURL = mockServer.URL

		err := client.RevokeApiKey("test_api_key")

		require.NoError(t, err)
	})

	t.Run("empty API key", func(t *testing.T) {
		auth := &Auth{JWT: "valid_jwt_token"}
		client := New(auth)

		err := client.RevokeApiKey("")

		require.Error(t, err)
		require.Contains(t, err.Error(), "api key is required")
	})

	t.Run("server error", func(t *testing.T) {
		auth := &Auth{JWT: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"Internal server error"}`))
		}))
		defer mockServer.Close()
		client.BaseURL = mockServer.URL

		err := client.RevokeApiKey("test_api_key")

		require.Error(t, err)
		require.Contains(t, err.Error(), "Internal server error")
	})

	t.Run("network error", func(t *testing.T) {
		auth := &Auth{JWT: "valid_jwt_token"}
		client := New(auth)
		client.BaseURL = "http://invalid-url"

		err := client.RevokeApiKey("test_api_key")

		require.Error(t, err) 
	})

	t.Run("unauthorized request", func(t *testing.T) {
		auth := &Auth{JWT: "invalid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"Unauthorized"}`))
		}))
		defer mockServer.Close()
		client.BaseURL = mockServer.URL

		err := client.RevokeApiKey("test_api_key")

		require.Error(t, err)
		require.Contains(t, err.Error(), "Unauthorized")
	})
}

func TestRevokeApiKeyV3(t *testing.T) {
	t.Run("successful API key revocation", func(t *testing.T) {
		auth := &Auth{JWT: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/v3/pinata/keys/test_api_key", r.URL.Path)
			require.Equal(t, http.MethodPut, r.Method)
			require.Equal(t, "Bearer valid_jwt_token", r.Header.Get("Authorization"))
			w.WriteHeader(http.StatusOK)
		}))
		defer mockServer.Close()
		client.BaseURL = mockServer.URL

		err := client.RevokeApiKeyV3("test_api_key")

		require.NoError(t, err)
	})

	t.Run("empty API key", func(t *testing.T) {
		auth := &Auth{JWT: "valid_jwt_token"}
		client := New(auth)

		err := client.RevokeApiKeyV3("")

		require.Error(t, err)
		require.Contains(t, err.Error(), "key is required")
	})

	t.Run("server error", func(t *testing.T) {
		auth := &Auth{JWT: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"Internal server error"}`))
		}))
		defer mockServer.Close()
		client.BaseURL = mockServer.URL

		err := client.RevokeApiKeyV3("test_api_key")

		require.Error(t, err)
		require.Contains(t, err.Error(), "Internal server error")
	})

	t.Run("network error", func(t *testing.T) {
		auth := &Auth{JWT: "valid_jwt_token"}
		client := New(auth)
		client.BaseURL = "http://invalid-url"

		err := client.RevokeApiKeyV3("test_api_key")

		require.Error(t, err) 
	})

	t.Run("unauthorized request", func(t *testing.T) {
		auth := &Auth{JWT: "invalid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"Unauthorized"}`))
		}))
		defer mockServer.Close()
		client.BaseURL = mockServer.URL

		err := client.RevokeApiKeyV3("test_api_key")

		require.Error(t, err)
		require.Contains(t, err.Error(), "Unauthorized")
	})
}
