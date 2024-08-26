package pinata

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddSwap(t *testing.T) {
	t.Run("successful swap addition", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/v3/ipfs/swap/test_cid", r.URL.Path)
			require.Equal(t, http.MethodPut, r.Method)
			require.Equal(t, "Bearer valid_jwt_token", r.Header.Get("Authorization"))

			var payload map[string]string
			err := json.NewDecoder(r.Body).Decode(&payload)
			require.NoError(t, err)
			require.Equal(t, "test_swap_cid", payload["swapCid"])

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"data" : {"mappedCid": "test_swap_cid"}}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		response, err := client.AddSwap("test_cid", "test_swap_cid")

		require.NoError(t, err)
		require.NotNil(t, response)
		require.Equal(t, "test_swap_cid", response.Data.MappedCid) 
	})

	t.Run("empty cid", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)

		response, err := client.AddSwap("", "test_swap_cid")

		require.Error(t, err)
		require.Nil(t, response)
		require.Contains(t, err.Error(), "cid and swapcid are required")
	})

	t.Run("empty swap cid", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)

		response, err := client.AddSwap("test_cid", "")

		require.Error(t, err)
		require.Nil(t, response)
		require.Contains(t, err.Error(), "cid and swapcid are required")
	})

	t.Run("server error", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"Internal server error"}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		response, err := client.AddSwap("test_cid", "test_swap_cid")

		require.Error(t, err)
		require.Nil(t, response)
		require.Contains(t, err.Error(), "Internal server error")
	})

	t.Run("unauthorized error", func(t *testing.T) {
		auth := &Auth{jwt: "invalid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"Unauthorized"}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		response, err := client.AddSwap("test_cid", "test_swap_cid")

		require.Error(t, err)
		require.Nil(t, response)
		require.Contains(t, err.Error(), "Unauthorized")
	})
}

func TestGetSwapHistory(t *testing.T) {
	t.Run("successful swap history retrieval", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/v3/ipfs/swap/test_cid", r.URL.Path)
			require.Equal(t, http.MethodDelete, r.Method)
			require.Equal(t, "Bearer valid_jwt_token", r.Header.Get("Authorization"))
			require.Equal(t, "test_domain", r.URL.Query().Get("domain"))

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"data": [{"createdAt": "2023-05-01T12:00:00Z", "mappedCid": "swap_cid_1"}]}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		response, err := client.GetSwapHistory("test_cid", "test_domain")

		require.NoError(t, err)
		require.NotNil(t, response)
		require.Len(t, response.Data, 1)
		require.Equal(t, "swap_cid_1", response.Data[0].MappedCid)
		require.Equal(t, "2023-05-01 12:00:00 +0000 UTC", response.Data[0].CreatedAt.String())
	})

	t.Run("empty cid", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)

		response, err := client.GetSwapHistory("", "test_domain")

		require.Error(t, err)
		require.Nil(t, response)
		require.Contains(t, err.Error(), "cid and domain are required")
	})

	t.Run("empty domain", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)

		response, err := client.GetSwapHistory("test_cid", "")

		require.Error(t, err)
		require.Nil(t, response)
		require.Contains(t, err.Error(), "cid and domain are required")
	})

	t.Run("server error", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"Internal server error"}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		response, err := client.GetSwapHistory("test_cid", "test_domain")

		require.Error(t, err)
		require.Nil(t, response)
		require.Contains(t, err.Error(), "Internal server error")
	})

	t.Run("not found error", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":"Swap history not found"}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		response, err := client.GetSwapHistory("non_existent_cid", "test_domain")

		require.Error(t, err)
		require.Nil(t, response)
		require.Contains(t, err.Error(), "Swap history not found")
	})
}

func TestRemoveSwap(t *testing.T) {
	t.Run("successful swap removal", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/v3/ipfs/swap/test_cid", r.URL.Path)
			require.Equal(t, http.MethodDelete, r.Method)
			require.Equal(t, "Bearer valid_jwt_token", r.Header.Get("Authorization"))

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"data" : {"message": "Swap removed successfully"}}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		response, err := client.RemoveSwap("test_cid")

		require.NoError(t, err)
		require.NotNil(t, response) 
	})

	t.Run("empty cid", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)

		response, err := client.RemoveSwap("")

		require.Error(t, err)
		require.Nil(t, response)
		require.Contains(t, err.Error(), "cid is required")
	})

	t.Run("server error", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"Internal server error"}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		response, err := client.RemoveSwap("test_cid")

		require.Error(t, err)
		require.Nil(t, response)
		require.Contains(t, err.Error(), "Internal server error")
	})

	t.Run("not found error", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":"Swap not found"}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		response, err := client.RemoveSwap("non_existent_cid")

		require.Error(t, err)
		require.Nil(t, response)
		require.Contains(t, err.Error(), "Swap not found")
	})
}
