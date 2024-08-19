package pinata

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPinFileToIPFS(t *testing.T) {
	t.Run("successful file pinning", func(t *testing.T) {
		auth := &auth{jwt: "valid_jwt_token"}
		client := New(auth)
		tempFile, err := os.CreateTemp("", "test_file_*.txt")
		require.NoError(t, err)
		defer os.Remove(tempFile.Name())

		_, err = tempFile.WriteString("Test content")
		require.NoError(t, err)
		tempFile.Close()

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/pinning/pinFileToIPFS", r.URL.Path)
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "Bearer valid_jwt_token", r.Header.Get("Authorization"))
			require.Contains(t, r.Header.Get("Content-Type"), "multipart/form-data")

			err := r.ParseMultipartForm(10 << 20)
			require.NoError(t, err)

			file, _, err := r.FormFile("file")
			require.NoError(t, err)
			defer file.Close()

			content, err := io.ReadAll(file)
			require.NoError(t, err)
			require.Equal(t, "Test content", string(content))

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"IpfsHash":"Qm123456","PinSize":123,"Timestamp":"2023-05-01T12:00:00Z"}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		response, err := client.PinFileToIPFS(tempFile.Name(), nil)

		require.NoError(t, err)
		require.NotNil(t, response)
		require.Equal(t, "Qm123456", response.IpfsHash)
		require.Equal(t, 123, response.PinSize)
		require.Equal(t, "2023-05-01T12:00:00Z", response.Timestamp)
	})

	t.Run("empty file path", func(t *testing.T) {
		auth := &auth{jwt: "valid_jwt_token"}
		client := New(auth)

		response, err := client.PinFileToIPFS("", nil)

		require.Error(t, err)
		require.Nil(t, response)
		require.Contains(t, err.Error(), "filepath is required")
	})

	t.Run("non-existent file", func(t *testing.T) {
		auth := &auth{jwt: "valid_jwt_token"}
		client := New(auth)

		response, err := client.PinFileToIPFS("/path/to/non/existent/file.txt", nil)

		require.Error(t, err)
		require.Nil(t, response)
		require.Contains(t, err.Error(), "failed to open file")
	})

	t.Run("with pin options", func(t *testing.T) {
		auth := &auth{jwt: "valid_jwt_token"}
		client := New(auth)
		tempFile, err := os.CreateTemp("", "test_file_*.txt")
		require.NoError(t, err)
		defer os.Remove(tempFile.Name())

		_, err = tempFile.WriteString("Test content")
		require.NoError(t, err)
		tempFile.Close()

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := r.ParseMultipartForm(10 << 20)
			require.NoError(t, err)

			pinataOptions := r.FormValue("pinataOptions")
			require.NotEmpty(t, pinataOptions)

			var options PinOptions
			err = json.Unmarshal([]byte(pinataOptions), &options)
			require.NoError(t, err)
			require.Equal(t, "test_name", options.PinataMetadata.Name)

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"IpfsHash":"Qm789012","PinSize":456,"Timestamp":"2023-05-02T12:00:00Z"}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		options := &PinOptions{
			PinataMetadata: PinataMetadata{
				Name: "test_name",
			},
		}
		response, err := client.PinFileToIPFS(tempFile.Name(), options)

		require.NoError(t, err)
		require.NotNil(t, response)
		require.Equal(t, "Qm789012", response.IpfsHash)
		require.Equal(t, 456, response.PinSize)
		require.Equal(t, "2023-05-02T12:00:00Z", response.Timestamp)
	})

	t.Run("server error", func(t *testing.T) {
		auth := &auth{jwt: "valid_jwt_token"}
		client := New(auth)
		tempFile, err := os.CreateTemp("", "test_file_*.txt")
		require.NoError(t, err)
		defer os.Remove(tempFile.Name())

		_, err = tempFile.WriteString("Test content")
		require.NoError(t, err)
		tempFile.Close()

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"Internal server error"}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		response, err := client.PinFileToIPFS(tempFile.Name(), nil)

		require.Error(t, err)
		require.Nil(t, response)
		require.Contains(t, err.Error(), "Internal server error")
	})
}

func TestPinJSONToIPFS(t *testing.T) {
	t.Run("successful JSON pinning", func(t *testing.T) {
		auth := &auth{jwt: "valid_jwt_token"}
		client := New(auth)

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/pinning/pinJSONToIPFS", r.URL.Path)
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "Bearer valid_jwt_token", r.Header.Get("Authorization"))
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))

			var payload map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&payload)
			require.NoError(t, err)
			require.Contains(t, payload, "pinataContent")

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"IpfsHash":"Qm987654","PinSize":789,"Timestamp":"2023-05-03T12:00:00Z"}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		data := map[string]string{"key": "value"}
		response, err := client.PinJSONToIPFS(data, nil)

		require.NoError(t, err)
		require.NotNil(t, response)
		require.Equal(t, "Qm987654", response.IpfsHash)
		require.Equal(t, 789, response.PinSize)
		require.Equal(t, "2023-05-03T12:00:00Z", response.Timestamp)
	})

	t.Run("nil data", func(t *testing.T) {
		auth := &auth{jwt: "valid_jwt_token"}
		client := New(auth)

		response, err := client.PinJSONToIPFS(nil, nil)

		require.Error(t, err)
		require.Nil(t, response)
		require.Contains(t, err.Error(), "jsonData is required")
	})

	t.Run("with pin options", func(t *testing.T) {
		auth := &auth{jwt: "valid_jwt_token"}
		client := New(auth)

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var payload map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&payload)
			require.NoError(t, err)
			require.Contains(t, payload, "pinataContent")
			require.Contains(t, payload, "pinataOptions")
			require.Contains(t, payload, "pinataMetadata")

			options, ok := payload["pinataOptions"].(map[string]interface{})
			require.True(t, ok)
			require.Equal(t, 5, int(options["cidVersion"].(float64)))

			metadata, ok := payload["pinataMetadata"].(map[string]interface{})
			require.True(t, ok)
			require.Equal(t, "test_json", metadata["name"])

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"IpfsHash":"Qm135790","PinSize":246,"Timestamp":"2023-05-04T12:00:00Z"}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		data := map[string]int{"number": 42}
		options := &PinOptions{
			PinataOptions: struct {
				CidVersion int "json:\"cidVersion,omitempty\""
			}{
				CidVersion: 5,
			},
			PinataMetadata: PinataMetadata{
				Name: "test_json",
			},
		}
		response, err := client.PinJSONToIPFS(data, options)

		require.NoError(t, err)
		require.NotNil(t, response)
		require.Equal(t, "Qm135790", response.IpfsHash)
		require.Equal(t, 246, response.PinSize)
		require.Equal(t, "2023-05-04T12:00:00Z", response.Timestamp)
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

		data := map[string]bool{"flag": true}
		response, err := client.PinJSONToIPFS(data, nil)

		require.Error(t, err)
		require.Nil(t, response)
		require.Contains(t, err.Error(), "Internal server error")
	})
}

func TestPinByCid(t *testing.T) {
	t.Run("successful pin by CID", func(t *testing.T) {
		auth := &auth{jwt: "valid_jwt_token"}
		client := New(auth)

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/pinning/pinByHash", r.URL.Path)
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "Bearer valid_jwt_token", r.Header.Get("Authorization"))
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))

			var payload map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&payload)
			require.NoError(t, err)
			require.Equal(t, "QmTestHash123", payload["hashToPin"])

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id":"test_id","ipfsHash":"QmTestHash123","status":"pinned"}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		response, err := client.PinByCid("QmTestHash123", nil)

		require.NoError(t, err)
		require.NotNil(t, response)
		require.Equal(t, "test_id", response.ID)
		require.Equal(t, "QmTestHash123", response.IpfsHash)
		require.Equal(t, "pinned", response.Status)
	})

	t.Run("empty hash to pin", func(t *testing.T) {
		auth := &auth{jwt: "valid_jwt_token"}
		client := New(auth)

		response, err := client.PinByCid("", nil)

		require.Error(t, err)
		require.Nil(t, response)
		require.Contains(t, err.Error(), "hashToPin is required")
	})

	t.Run("with pin options", func(t *testing.T) {
		auth := &auth{jwt: "valid_jwt_token"}
		client := New(auth)

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var payload map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&payload)
			require.NoError(t, err)
			require.Equal(t, "QmTestHash456", payload["hashToPin"])
			require.Contains(t, payload, "pinataOptions")
			require.Contains(t, payload, "pinataMetadata")

			options := payload["pinataOptions"].(map[string]interface{})
			require.Equal(t, "test_group", options["groupId"])
			require.Len(t, options["hostNodes"], 2)

			metadata := payload["pinataMetadata"].(map[string]interface{})
			require.Equal(t, "test_pin", metadata["name"])

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id":"test_id_2","ipfsHash":"QmTestHash456","status":"pinned","created":"2023-05-06T12:00:00Z"}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		options := &PinByCidOptions{
			PinataOptions: struct {
				GroupId   string   "json:\"groupId,omitempty\""
				HostNodes []string "json:\"hostNodes,omitempty\""
			}{
				GroupId:   "test_group",
				HostNodes: []string{"node1", "node2"},
			},
			PinataMetadata: PinataMetadata{
				Name: "test_pin",
			},
		}
		response, err := client.PinByCid("QmTestHash456", options)

		require.NoError(t, err)
		require.NotNil(t, response)
		require.Equal(t, "test_id_2", response.ID)
		require.Equal(t, "QmTestHash456", response.IpfsHash)
		require.Equal(t, "pinned", response.Status)
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

		response, err := client.PinByCid("QmTestHash789", nil)

		require.Error(t, err)
		require.Nil(t, response)
		require.Contains(t, err.Error(), "Internal server error")
	})
}

func TestListFiles(t *testing.T) {
	t.Run("successful list files without options", func(t *testing.T) {
		auth := &auth{jwt: "valid_jwt_token"}
		client := New(auth)

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/data/pinList", r.URL.Path)
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "Bearer valid_jwt_token", r.Header.Get("Authorization"))

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"count":2,"rows":[{"id":"file1","ipfs_pin_hash":"Qm123","size":100,"user_id":"user1","date_pinned":"2023-05-07T12:00:00Z"},{"id":"file2","ipfs_pin_hash":"Qm456","size":200,"user_id":"user1","date_pinned":"2023-05-08T12:00:00Z"}]}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		response, err := client.ListFiles(nil)

		require.NoError(t, err)
		require.NotNil(t, response)
		require.Equal(t, 2, response.Count)
		require.Len(t, response.Rows, 2)
		require.Equal(t, "file1", response.Rows[0].ID)
		require.Equal(t, "Qm123", response.Rows[0].IPFSPinHash)
		require.Equal(t, 100, response.Rows[0].Size)
		require.Equal(t, "user1", response.Rows[0].UserID)
		require.Equal(t, "2023-05-07T12:00:00Z", response.Rows[0].DatePinned)
	})

	t.Run("list files with options", func(t *testing.T) {
		auth := &auth{jwt: "valid_jwt_token"}
		client := New(auth)

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/data/pinList", r.URL.Path)
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "Bearer valid_jwt_token", r.Header.Get("Authorization"))

			query := r.URL.Query()
			require.Equal(t, "10", query.Get("pageLimit"))
			require.Equal(t, "20", query.Get("pageOffset"))
			// require.Equal(t, "test", query.Get("metadata[name]"))

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"count":1,"rows":[{"id":"file3","ipfs_pin_hash":"Qm789","size":300,"user_id":"user1","date_pinned":"2023-05-09T12:00:00Z"}]}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		options := &ListFilesOptions{
			PageLimit:  10,
			PageOffset: 20,
			Metadata: map[string]interface{}{
				"name": "test",
			},
		}
		response, err := client.ListFiles(options)

		require.NoError(t, err)
		require.NotNil(t, response)
		require.Equal(t, 1, response.Count)
		require.Len(t, response.Rows, 1)
		require.Equal(t, "file3", response.Rows[0].ID)
		require.Equal(t, "Qm789", response.Rows[0].IPFSPinHash)
		require.Equal(t, 300, response.Rows[0].Size)
		require.Equal(t, "user1", response.Rows[0].UserID)
		require.Equal(t, "2023-05-09T12:00:00Z", response.Rows[0].DatePinned)
	})

	t.Run("empty response", func(t *testing.T) {
		auth := &auth{jwt: "valid_jwt_token"}
		client := New(auth)

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"count":0,"rows":[]}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		response, err := client.ListFiles(nil)

		require.NoError(t, err)
		require.NotNil(t, response)
		require.Equal(t, 0, response.Count)
		require.Empty(t, response.Rows)
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

		response, err := client.ListFiles(nil)

		require.Error(t, err)
		require.Nil(t, response)
		require.Contains(t, err.Error(), "Internal server error")
	})
}

func TestListPinByCidJobs(t *testing.T) {
	t.Run("successful list pin by CID jobs without options", func(t *testing.T) {
		auth := &auth{jwt: "valid_jwt_token"}
		client := New(auth)

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/pinning/pinJobs", r.URL.Path)
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "Bearer valid_jwt_token", r.Header.Get("Authorization"))

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"count":2,"rows":[{"id":"job1","ipfs_pin_hash":"Qm123","status":"retrieving","date_queued":"2023-05-10T12:00:00Z"},{"id":"job2","ipfs_pin_hash":"Qm456","status":"retrieving","date_queued":"2023-05-11T12:00:00Z"}]}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		response, err := client.ListPinByCidJobs(nil)

		require.NoError(t, err)
		require.NotNil(t, response)
		require.Equal(t, 2, response.Count)
		require.Len(t, response.Rows, 2)
		require.Equal(t, "job1", response.Rows[0].ID)
		require.Equal(t, "Qm123", response.Rows[0].IPFSPinHash)
		require.Equal(t, "retrieving", response.Rows[0].Status)
		require.Equal(t, "2023-05-10T12:00:00Z", response.Rows[0].DateQueued)
	})

	t.Run("list pin by CID jobs with options", func(t *testing.T) {
		auth := &auth{jwt: "valid_jwt_token"}
		client := New(auth)

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/pinning/pinJobs", r.URL.Path)
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "Bearer valid_jwt_token", r.Header.Get("Authorization"))

			query := r.URL.Query()
			require.Equal(t, "5", query.Get("limit"))
			require.Equal(t, "10", query.Get("offset"))
			require.Equal(t, "retrieving", query.Get("status"))

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"count":1,"rows":[{"id":"job3","ipfs_pin_hash":"Qm789","status":"retrieving","date_queued":"2023-05-12T12:00:00Z"}]}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		options := &ListPinByCidOptions{
			Limit:  5,
			Offset: 10,
			Status: "retrieving",
		}
		response, err := client.ListPinByCidJobs(options)

		require.NoError(t, err)
		require.NotNil(t, response)
		require.Equal(t, 1, response.Count)
		require.Len(t, response.Rows, 1)
		require.Equal(t, "job3", response.Rows[0].ID)
		require.Equal(t, "Qm789", response.Rows[0].IPFSPinHash)
		require.Equal(t, "retrieving", response.Rows[0].Status)
		require.Equal(t, "2023-05-12T12:00:00Z", response.Rows[0].DateQueued)
	})

	t.Run("empty response", func(t *testing.T) {
		auth := &auth{jwt: "valid_jwt_token"}
		client := New(auth)

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"count":0,"rows":[]}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		response, err := client.ListPinByCidJobs(nil)

		require.NoError(t, err)
		require.NotNil(t, response)
		require.Equal(t, 0, response.Count)
		require.Empty(t, response.Rows)
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

		response, err := client.ListPinByCidJobs(nil)

		require.Error(t, err)
		require.Nil(t, response)
		require.Contains(t, err.Error(), "Internal server error")
	})
}

func TestUpdateFileMetadata(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		auth := &auth{jwt: "valid_jwt_token"}
		client := New(auth)

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/pinning/hashMetadata", r.URL.Path)
			require.Equal(t, http.MethodPut, r.Method)
			require.Equal(t, "Bearer valid_jwt_token", r.Header.Get("Authorization"))
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))

			var payload map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&payload)
			require.NoError(t, err)
			require.Equal(t, "QmTestHash123", payload["ipfsPinHash"])
			require.Equal(t, "Updated File", payload["name"])
			require.Equal(t, map[string]interface{}{"key1": "value1", "key2": "value2"}, payload["keyvalues"])

			w.WriteHeader(http.StatusOK)
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		options := &PinMetadataUpdateOptions{
			Name: "Updated File",
			KeyValues: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
		}
		err := client.UpdateFileMetadata("QmTestHash123", options)

		require.NoError(t, err)
	})

	t.Run("empty file hash", func(t *testing.T) {
		auth := &auth{jwt: "valid_jwt_token"}
		client := New(auth)

		options := &PinMetadataUpdateOptions{
			Name: "Updated File",
		}
		err := client.UpdateFileMetadata("", options)

		require.Error(t, err)
		require.Contains(t, err.Error(), "fileHash and options are required")
	})

	t.Run("nil options", func(t *testing.T) {
		auth := &auth{jwt: "valid_jwt_token"}
		client := New(auth)

		err := client.UpdateFileMetadata("QmTestHash123", nil)

		require.Error(t, err)
		require.Contains(t, err.Error(), "fileHash and options are required")
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

		options := &PinMetadataUpdateOptions{
			Name: "Updated File",
		}
		err := client.UpdateFileMetadata("QmTestHash123", options)

		require.Error(t, err)
		require.Contains(t, err.Error(), "Internal server error")
	})
}

func TestDeleteFile(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		auth := &auth{jwt: "valid_jwt_token"}
		client := New(auth)

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/pinning/unpin/QmTestCID123", r.URL.Path)
			require.Equal(t, http.MethodDelete, r.Method)
			require.Equal(t, "Bearer valid_jwt_token", r.Header.Get("Authorization"))

			w.WriteHeader(http.StatusOK)
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		err := client.DeleteFile("QmTestCID123")

		require.NoError(t, err)
	})

	t.Run("empty CID", func(t *testing.T) {
		auth := &auth{jwt: "valid_jwt_token"}
		client := New(auth)

		err := client.DeleteFile("")

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

		err := client.DeleteFile("QmTestCID456")

		require.Error(t, err)
		require.Contains(t, err.Error(), "Internal server error")
	})

	t.Run("not found error", func(t *testing.T) {
		auth := &auth{jwt: "valid_jwt_token"}
		client := New(auth)

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":"File not found"}`))
		}))
		defer mockServer.Close()
		client.baseURL = mockServer.URL

		err := client.DeleteFile("QmNonExistentCID")

		require.Error(t, err)
		require.Contains(t, err.Error(), "File not found")
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

		err := client.DeleteFile("QmTestCID789")

		require.Error(t, err)
		require.Contains(t, err.Error(), "Unauthorized")
	})
}
