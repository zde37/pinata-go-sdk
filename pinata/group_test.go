package pinata

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateGroup(t *testing.T) {
	t.Run("successful group creation", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/groups", r.URL.Path)
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "Bearer valid_jwt_token", r.Header.Get("Authorization"))
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))

			var payload map[string]string
			err := json.NewDecoder(r.Body).Decode(&payload)
			require.NoError(t, err)
			require.Equal(t, "test_group", payload["name"])

			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"id":"group123","name":"test_group"}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		group, err := client.CreateGroup("test_group")

		require.NoError(t, err)
		require.NotNil(t, group)
		require.Equal(t, "group123", group.ID)
		require.Equal(t, "test_group", group.GroupName)
	})

	t.Run("empty group name", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)

		group, err := client.CreateGroup("")

		require.Error(t, err)
		require.Nil(t, group)
		require.Contains(t, err.Error(), "group name is required")
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

		group, err := client.CreateGroup("test_group")

		require.Error(t, err)
		require.Nil(t, group)
		require.Contains(t, err.Error(), "Internal server error")
	})

	t.Run("invalid JSON response", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"id":"group123","name":}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		group, err := client.CreateGroup("test_group")

		require.Error(t, err)
		require.Nil(t, group)
		require.Contains(t, err.Error(), "invalid character")
	})
}

func TestGetGroup(t *testing.T) {
	t.Run("successful group retrieval", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/groups/group123", r.URL.Path)
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "Bearer valid_jwt_token", r.Header.Get("Authorization"))

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id":"group123","name":"test_group"}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		group, err := client.GetGroup("group123")

		require.NoError(t, err)
		require.NotNil(t, group)
		require.Equal(t, "group123", group.ID)
		require.Equal(t, "test_group", group.GroupName)
	})

	t.Run("empty group ID", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)

		group, err := client.GetGroup("")

		require.Error(t, err)
		require.Nil(t, group)
		require.Contains(t, err.Error(), "group id is required")
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

		group, err := client.GetGroup("group123")

		require.Error(t, err)
		require.Nil(t, group)
		require.Contains(t, err.Error(), "Internal server error")
	})

	t.Run("not found error", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":"Group not found"}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		group, err := client.GetGroup("nonexistent_group")

		require.Error(t, err)
		require.Nil(t, group)
		require.Contains(t, err.Error(), "Group not found")
	})

	t.Run("invalid JSON response", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id":"group123","name":}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		group, err := client.GetGroup("group123")

		require.Error(t, err)
		require.Nil(t, group)
		require.Contains(t, err.Error(), "invalid character")
	})
}

func TestListGroups(t *testing.T) {
	t.Run("successful groups listing", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/groups", r.URL.Path)
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "Bearer valid_jwt_token", r.Header.Get("Authorization"))

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"id":"group1","name":"test_group1"},{"id":"group2","name":"test_group2"}]`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		groups, err := client.ListGroups(nil)

		require.NoError(t, err)
		require.NotNil(t, groups)
		require.Len(t, groups, 2)
		require.Equal(t, "group1", groups[0].ID)
		require.Equal(t, "test_group1", groups[0].GroupName)
		require.Equal(t, "group2", groups[1].ID)
		require.Equal(t, "test_group2", groups[1].GroupName)
	})

	t.Run("with query parameters", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/groups", r.URL.Path)
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "Bearer valid_jwt_token", r.Header.Get("Authorization"))
			require.Equal(t, "10", r.URL.Query().Get("limit"))
			require.Equal(t, "5", r.URL.Query().Get("offset"))

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"id":"group3","name":"test_group3"}]`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		options := &ListGroupsOptions{
			Limit:  10,
			Offset: 5,
		}
		groups, err := client.ListGroups(options)

		require.NoError(t, err)
		require.NotNil(t, groups)
		require.Len(t, groups, 1)
		require.Equal(t, "group3", groups[0].ID)
		require.Equal(t, "test_group3", groups[0].GroupName)
	})

	t.Run("empty response", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[]`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		groups, err := client.ListGroups(nil)

		require.NoError(t, err)
		require.NotNil(t, groups)
		require.Len(t, groups, 0)
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

		groups, err := client.ListGroups(nil)

		require.Error(t, err)
		require.Nil(t, groups)
		require.Contains(t, err.Error(), "Internal server error")
	})

	t.Run("invalid JSON response", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"id":"group1","name":}]`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		groups, err := client.ListGroups(nil)

		require.Error(t, err)
		require.Nil(t, groups)
		require.Contains(t, err.Error(), "invalid character")
	})
}

func TestUpdateGroup(t *testing.T) {
	t.Run("successful group update", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/groups/group123", r.URL.Path)
			require.Equal(t, http.MethodPut, r.Method)
			require.Equal(t, "Bearer valid_jwt_token", r.Header.Get("Authorization"))
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))

			var payload map[string]string
			err := json.NewDecoder(r.Body).Decode(&payload)
			require.NoError(t, err)
			require.Equal(t, "new_group_name", payload["name"])

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id":"group123","name":"new_group_name"}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		group, err := client.UpdateGroup("group123", "new_group_name")

		require.NoError(t, err)
		require.NotNil(t, group)
		require.Equal(t, "group123", group.ID)
		require.Equal(t, "new_group_name", group.GroupName)
	})

	t.Run("empty group ID", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)

		group, err := client.UpdateGroup("", "new_group_name")

		require.Error(t, err)
		require.Nil(t, group)
		require.Contains(t, err.Error(), "group id and new group name are required")
	})

	t.Run("empty new group name", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)

		group, err := client.UpdateGroup("group123", "")

		require.Error(t, err)
		require.Nil(t, group)
		require.Contains(t, err.Error(), "group id and new group name are required")
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

		group, err := client.UpdateGroup("group123", "new_group_name")

		require.Error(t, err)
		require.Nil(t, group)
		require.Contains(t, err.Error(), "Internal server error")
	})

	t.Run("not found error", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":"Group not found"}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		group, err := client.UpdateGroup("nonexistent_group", "new_group_name")

		require.Error(t, err)
		require.Nil(t, group)
		require.Contains(t, err.Error(), "Group not found")
	})

	t.Run("invalid JSON response", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id":"group123","name":}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		group, err := client.UpdateGroup("group123", "new_group_name")

		require.Error(t, err)
		require.Nil(t, group)
		require.Contains(t, err.Error(), "invalid character")
	})
}

func TestAddCidToGroup(t *testing.T) {
	t.Run("successful add CID to group", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/groups/group123/cids", r.URL.Path)
			require.Equal(t, http.MethodPut, r.Method)
			require.Equal(t, "Bearer valid_jwt_token", r.Header.Get("Authorization"))
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))

			var payload map[string][]string
			err := json.NewDecoder(r.Body).Decode(&payload)
			require.NoError(t, err)
			require.Equal(t, []string{"cid1", "cid2"}, payload["cids"])

			w.WriteHeader(http.StatusOK)
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		err := client.AddCidToGroup("group123", []string{"cid1", "cid2"})

		require.NoError(t, err)
	})

	t.Run("empty group ID", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)

		err := client.AddCidToGroup("", []string{"cid1", "cid2"})

		require.Error(t, err)
		require.Contains(t, err.Error(), "group id and at least one cid is required")
	})

	t.Run("empty CIDs list", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)

		err := client.AddCidToGroup("group123", []string{})

		require.Error(t, err)
		require.Contains(t, err.Error(), "group id and at least one cid is required")
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

		err := client.AddCidToGroup("group123", []string{"cid1"})

		require.Error(t, err)
		require.Contains(t, err.Error(), "Internal server error")
	})
}

func TestRemoveCidFromGroup(t *testing.T) {
	t.Run("successful remove CID from group", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/groups/group123/cids", r.URL.Path)
			require.Equal(t, http.MethodDelete, r.Method)
			require.Equal(t, "Bearer valid_jwt_token", r.Header.Get("Authorization"))
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))

			var payload map[string][]string
			err := json.NewDecoder(r.Body).Decode(&payload)
			require.NoError(t, err)
			require.Equal(t, []string{"cid1", "cid2"}, payload["cids"])

			w.WriteHeader(http.StatusOK)
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		err := client.RemoveCidFromGroup("group123", []string{"cid1", "cid2"})

		require.NoError(t, err)
	})

	t.Run("empty group ID", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)

		err := client.RemoveCidFromGroup("", []string{"cid1", "cid2"})

		require.Error(t, err)
		require.Contains(t, err.Error(), "group id and at least one cid is required")
	})

	t.Run("empty CIDs list", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)

		err := client.RemoveCidFromGroup("group123", []string{})

		require.Error(t, err)
		require.Contains(t, err.Error(), "group id and at least one cid is required")
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

		err := client.RemoveCidFromGroup("group123", []string{"cid1"})

		require.Error(t, err)
		require.Contains(t, err.Error(), "Internal server error")
	})

	t.Run("multiple CIDs removal", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/groups/group123/cids", r.URL.Path)
			require.Equal(t, http.MethodDelete, r.Method)

			var payload map[string][]string
			err := json.NewDecoder(r.Body).Decode(&payload)
			require.NoError(t, err)
			require.Equal(t, []string{"cid1", "cid2", "cid3"}, payload["cids"])

			w.WriteHeader(http.StatusOK)
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		err := client.RemoveCidFromGroup("group123", []string{"cid1", "cid2", "cid3"})

		require.NoError(t, err)
	})
}

func TestRemoveGroup(t *testing.T) {
	t.Run("successful group removal", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/groups/group123", r.URL.Path)
			require.Equal(t, http.MethodDelete, r.Method)
			require.Equal(t, "Bearer valid_jwt_token", r.Header.Get("Authorization"))
			w.WriteHeader(http.StatusOK)
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		err := client.RemoveGroup("group123")

		require.NoError(t, err)
	})

	t.Run("empty group ID", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)

		err := client.RemoveGroup("")

		require.Error(t, err)
		require.Contains(t, err.Error(), "group id is required")
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

		err := client.RemoveGroup("group123")

		require.Error(t, err)
		require.Contains(t, err.Error(), "Internal server error")
	})

	t.Run("not found error", func(t *testing.T) {
		auth := &Auth{jwt: "valid_jwt_token"}
		client := New(auth)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":"Group not found"}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		err := client.RemoveGroup("nonexistent_group")

		require.Error(t, err)
		require.Contains(t, err.Error(), "Group not found")
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

		err := client.RemoveGroup("group123")

		require.Error(t, err)
		require.Contains(t, err.Error(), "Unauthorized")
	})
}
