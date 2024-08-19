package pinata

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddCidSignature(t *testing.T) {
	t.Run("successful signature addition", func(t *testing.T) {
		auth := &auth{jwt: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/v3/ipfs/signature/test_cid", r.URL.Path)
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "Bearer valid_jwt_token", r.Header.Get("Authorization"))
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))

			var payload map[string]string
			err := json.NewDecoder(r.Body).Decode(&payload)
			require.NoError(t, err)
			require.Equal(t, "test_signature", payload["signature"])

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"data":{"cid":"test_cid","signature":"test_signature"}}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		cidSignature, err := client.AddCidSignature("test_cid", "test_signature")

		require.NoError(t, err)
		require.NotNil(t, cidSignature)
		require.Equal(t, "test_cid", cidSignature.Data.Cid)
		require.Equal(t, "test_signature", cidSignature.Data.Signature)
	})

	t.Run("empty cid", func(t *testing.T) {
		auth := &auth{jwt: "valid_jwt_token"}
		client := New(auth)

		cidSignature, err := client.AddCidSignature("", "test_signature")

		require.Error(t, err)
		require.Nil(t, cidSignature)
		require.Contains(t, err.Error(), "cid and signature is required")
	})

	t.Run("empty signature", func(t *testing.T) {
		auth := &auth{jwt: "valid_jwt_token"}
		client := New(auth)

		cidSignature, err := client.AddCidSignature("test_cid", "")

		require.Error(t, err)
		require.Nil(t, cidSignature)
		require.Contains(t, err.Error(), "cid and signature is required")
	})

	t.Run("server error", func(t *testing.T) {
		auth := &auth{jwt: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"Internal server error"}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		cidSignature, err := client.AddCidSignature("test_cid", "test_signature")

		require.Error(t, err)
		require.Nil(t, cidSignature)
		require.Contains(t, err.Error(), "Internal server error")
	})

	t.Run("invalid JSON response", func(t *testing.T) {
		auth := &auth{jwt: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"cid":"test_cid","signature":}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		cidSignature, err := client.AddCidSignature("test_cid", "test_signature")

		require.Error(t, err)
		require.Nil(t, cidSignature)
		require.Contains(t, err.Error(), "invalid character")
	})
}

func TestGetCidSignature(t *testing.T) {
	t.Run("successful signature retrieval", func(t *testing.T) {
		auth := &auth{jwt: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/v3/ipfs/signature/test_cid", r.URL.Path)
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "Bearer valid_jwt_token", r.Header.Get("Authorization"))

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"data":{"cid":"test_cid","signature":"test_signature"}}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		cidSignature, err := client.GetCidSignature("test_cid")

		require.NoError(t, err)
		require.NotNil(t, cidSignature)
		require.Equal(t, "test_cid", cidSignature.Data.Cid)
		require.Equal(t, "test_signature", cidSignature.Data.Signature)
	})

	t.Run("empty cid", func(t *testing.T) {
		auth := &auth{jwt: "valid_jwt_token"}
		client := New(auth)

		cidSignature, err := client.GetCidSignature("")

		require.Error(t, err)
		require.Nil(t, cidSignature)
		require.Contains(t, err.Error(), "cid is required")
	})

	t.Run("server error", func(t *testing.T) {
		auth := &auth{jwt: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"Internal server error"}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		cidSignature, err := client.GetCidSignature("test_cid")

		require.Error(t, err)
		require.Nil(t, cidSignature)
		require.Contains(t, err.Error(), "Internal server error")
	})

	t.Run("invalid JSON response", func(t *testing.T) {
		auth := &auth{jwt: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"data":{"cid":"test_cid","signature":}}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		cidSignature, err := client.GetCidSignature("test_cid")

		require.Error(t, err)
		require.Nil(t, cidSignature)
		require.Contains(t, err.Error(), "invalid character")
	})

	t.Run("not found error", func(t *testing.T) {
		auth := &auth{jwt: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":"Signature not found"}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		cidSignature, err := client.GetCidSignature("non_existent_cid")

		require.Error(t, err)
		require.Nil(t, cidSignature)
		require.Contains(t, err.Error(), "Signature not found")
	})
}

func TestRemoveCidSignature(t *testing.T) {
	t.Run("successful signature removal", func(t *testing.T) {
		auth := &auth{jwt: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/v3/ipfs/signature/test_cid", r.URL.Path)
			require.Equal(t, http.MethodDelete, r.Method)
			require.Equal(t, "Bearer valid_jwt_token", r.Header.Get("Authorization"))
			w.WriteHeader(http.StatusOK)
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		err := client.RemoveCidSignature("test_cid")

		require.NoError(t, err)
	})

	t.Run("empty cid", func(t *testing.T) {
		auth := &auth{jwt: "valid_jwt_token"}
		client := New(auth)

		err := client.RemoveCidSignature("")

		require.Error(t, err)
		require.Contains(t, err.Error(), "cid is required")
	})

	t.Run("server error", func(t *testing.T) {
		auth := &auth{jwt: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"Internal server error"}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		err := client.RemoveCidSignature("test_cid")

		require.Error(t, err)
		require.Contains(t, err.Error(), "Internal server error")
	})

	t.Run("not found error", func(t *testing.T) {
		auth := &auth{jwt: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":"Signature not found"}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		err := client.RemoveCidSignature("non_existent_cid")

		require.Error(t, err)
		require.Contains(t, err.Error(), "Signature not found")
	})

	t.Run("unauthorized error", func(t *testing.T) {
		auth := &auth{jwt: "invalid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"Unauthorized"}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		err := client.RemoveCidSignature("test_cid")

		require.Error(t, err)
		require.Contains(t, err.Error(), "Unauthorized")
	})
}
