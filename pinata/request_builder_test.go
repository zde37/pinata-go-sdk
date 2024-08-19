package pinata

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAddPathParam(t *testing.T) {
	t.Run("add first path param", func(t *testing.T) {
		rb := &requestBuilder{}
		result := rb.AddPathParam("key1", "value1")

		require.Equal(t, rb, result)
		require.Len(t, rb.pathParams, 1)
		require.Equal(t, "value1", rb.pathParams["key1"])
	})

	t.Run("add multiple path params", func(t *testing.T) {
		rb := &requestBuilder{}
		rb.AddPathParam("key1", "value1")
		rb.AddPathParam("key2", "value2")
		rb.AddPathParam("key3", "value3")

		require.Len(t, rb.pathParams, 3)
		require.Equal(t, "value1", rb.pathParams["key1"])
		require.Equal(t, "value2", rb.pathParams["key2"])
		require.Equal(t, "value3", rb.pathParams["key3"])
	})

	t.Run("overwrite existing path param", func(t *testing.T) {
		rb := &requestBuilder{}
		rb.AddPathParam("key1", "value1")
		rb.AddPathParam("key1", "new_value")

		require.Len(t, rb.pathParams, 1)
		require.Equal(t, "new_value", rb.pathParams["key1"])
	})

	t.Run("add path param with empty key", func(t *testing.T) {
		rb := &requestBuilder{}
		rb.AddPathParam("", "value")

		require.Len(t, rb.pathParams, 1)
		require.Equal(t, "value", rb.pathParams[""])
	})

	t.Run("add path param with empty value", func(t *testing.T) {
		rb := &requestBuilder{}
		rb.AddPathParam("key", "")

		require.Len(t, rb.pathParams, 1)
		require.Equal(t, "", rb.pathParams["key"])
	})
}

func TestAddQueryParam(t *testing.T) {
	t.Run("add first query param", func(t *testing.T) {
		rb := &requestBuilder{}
		result := rb.AddQueryParam("key1", "value1")

		require.Equal(t, rb, result)
		require.Len(t, rb.queryParams, 1)
		require.Equal(t, "value1", rb.queryParams["key1"])
	})

	t.Run("add multiple query params", func(t *testing.T) {
		rb := &requestBuilder{}
		rb.AddQueryParam("key1", "value1")
		rb.AddQueryParam("key2", 42)
		rb.AddQueryParam("key3", true)

		require.Len(t, rb.queryParams, 3)
		require.Equal(t, "value1", rb.queryParams["key1"])
		require.Equal(t, "42", rb.queryParams["key2"])
		require.Equal(t, "true", rb.queryParams["key3"])
	})

	t.Run("overwrite existing query param", func(t *testing.T) {
		rb := &requestBuilder{}
		rb.AddQueryParam("key1", "value1")
		rb.AddQueryParam("key1", "new_value")

		require.Len(t, rb.queryParams, 1)
		require.Equal(t, "new_value", rb.queryParams["key1"])
	})

	t.Run("add query param with empty key", func(t *testing.T) {
		rb := &requestBuilder{}
		rb.AddQueryParam("", "value")

		require.Len(t, rb.queryParams, 1)
		require.Equal(t, "value", rb.queryParams[""])
	})

	t.Run("add query param with nil value", func(t *testing.T) {
		rb := &requestBuilder{}
		rb.AddQueryParam("key", nil)

		require.Len(t, rb.queryParams, 1)
		require.Equal(t, "<nil>", rb.queryParams["key"])
	})

	t.Run("add query param with complex value", func(t *testing.T) {
		rb := &requestBuilder{}
		complexValue := struct {
			Name  string
			Value int
		}{
			Name:  "test",
			Value: 123,
		}
		rb.AddQueryParam("complex", complexValue)

		require.Len(t, rb.queryParams, 1)
		require.Equal(t, "{test 123}", rb.queryParams["complex"])
	})
}

func TestAddHeaders(t *testing.T) {
	t.Run("add first header", func(t *testing.T) {
		rb := &requestBuilder{}
		result := rb.AddHeaders("Content-Type", "application/json")

		require.Equal(t, rb, result)
		require.Len(t, rb.headers, 1)
		require.Equal(t, "application/json", rb.headers["Content-Type"])
	})

	t.Run("add multiple headers", func(t *testing.T) {
		rb := &requestBuilder{}
		rb.AddHeaders("Content-Type", "application/json")
		rb.AddHeaders("Authorization", "Bearer token")
		rb.AddHeaders("User-Agent", "TestAgent")

		require.Len(t, rb.headers, 3)
		require.Equal(t, "application/json", rb.headers["Content-Type"])
		require.Equal(t, "Bearer token", rb.headers["Authorization"])
		require.Equal(t, "TestAgent", rb.headers["User-Agent"])
	})

	t.Run("overwrite existing header", func(t *testing.T) {
		rb := &requestBuilder{}
		rb.AddHeaders("Content-Type", "application/json")
		rb.AddHeaders("Content-Type", "text/plain")

		require.Len(t, rb.headers, 1)
		require.Equal(t, "text/plain", rb.headers["Content-Type"])
	})

	t.Run("add header with empty key", func(t *testing.T) {
		rb := &requestBuilder{}
		rb.AddHeaders("", "value")

		require.Len(t, rb.headers, 1)
		require.Equal(t, "value", rb.headers[""])
	})

	t.Run("add header with empty value", func(t *testing.T) {
		rb := &requestBuilder{}
		rb.AddHeaders("EmptyHeader", "")

		require.Len(t, rb.headers, 1)
		require.Equal(t, "", rb.headers["EmptyHeader"])
	})

	t.Run("add header to existing headers", func(t *testing.T) {
		rb := &requestBuilder{
			headers: map[string]string{
				"Existing": "Header",
			},
		}
		rb.AddHeaders("New", "Header")

		require.Len(t, rb.headers, 2)
		require.Equal(t, "Header", rb.headers["Existing"])
		require.Equal(t, "Header", rb.headers["New"])
	})
}

func TestSetBody(t *testing.T) {
	t.Run("set body with string reader", func(t *testing.T) {
		rb := &requestBuilder{}
		body := strings.NewReader("test body")
		result := rb.SetBody(body, "text/plain")

		require.Equal(t, rb, result)
		require.Equal(t, body, rb.body)
		require.Equal(t, "text/plain", rb.contentType)
	})

	t.Run("set body with bytes buffer", func(t *testing.T) {
		rb := &requestBuilder{}
		body := bytes.NewBuffer([]byte("test body"))
		result := rb.SetBody(body, "application/octet-stream")

		require.Equal(t, rb, result)
		require.Equal(t, body, rb.body)
		require.Equal(t, "application/octet-stream", rb.contentType)
	})

	t.Run("set body with nil reader", func(t *testing.T) {
		rb := &requestBuilder{}
		result := rb.SetBody(nil, "application/json")

		require.Equal(t, rb, result)
		require.Nil(t, rb.body)
		require.Equal(t, "application/json", rb.contentType)
	})

	t.Run("overwrite existing body", func(t *testing.T) {
		rb := &requestBuilder{
			body:        strings.NewReader("old body"),
			contentType: "text/plain",
		}
		newBody := strings.NewReader("new body")
		result := rb.SetBody(newBody, "application/json")

		require.Equal(t, rb, result)
		require.Equal(t, newBody, rb.body)
		require.Equal(t, "application/json", rb.contentType)
	})

	t.Run("set body with empty content type", func(t *testing.T) {
		rb := &requestBuilder{}
		body := strings.NewReader("test body")
		result := rb.SetBody(body, "")

		require.Equal(t, rb, result)
		require.Equal(t, body, rb.body)
		require.Empty(t, rb.contentType)
	})
}

func TestSetJSONBody(t *testing.T) {
	t.Run("set JSON body with struct", func(t *testing.T) {
		rb := &requestBuilder{}
		testStruct := struct {
			Name  string
			Value int
		}{
			Name:  "test",
			Value: 123,
		}
		result, err := rb.SetJSONBody(testStruct)

		require.NoError(t, err)
		require.Equal(t, rb, result)
		require.NotNil(t, rb.body)
		require.Equal(t, "application/json", rb.contentType)

		bodyBytes, _ := io.ReadAll(rb.body)
		expectedJSON := `{"Name":"test","Value":123}`
		require.JSONEq(t, expectedJSON, string(bodyBytes))
	})

	t.Run("set JSON body with map", func(t *testing.T) {
		rb := &requestBuilder{}
		testMap := map[string]interface{}{
			"key1": "value1",
			"key2": 42,
		}
		result, err := rb.SetJSONBody(testMap)

		require.NoError(t, err)
		require.Equal(t, rb, result)
		require.NotNil(t, rb.body)
		require.Equal(t, "application/json", rb.contentType)

		bodyBytes, _ := io.ReadAll(rb.body)
		expectedJSON := `{"key1":"value1","key2":42}`
		require.JSONEq(t, expectedJSON, string(bodyBytes))
	})

	t.Run("set JSON body with slice", func(t *testing.T) {
		rb := &requestBuilder{}
		testSlice := []string{"item1", "item2", "item3"}
		result, err := rb.SetJSONBody(testSlice)

		require.NoError(t, err)
		require.Equal(t, rb, result)
		require.NotNil(t, rb.body)
		require.Equal(t, "application/json", rb.contentType)

		bodyBytes, _ := io.ReadAll(rb.body)
		expectedJSON := `["item1","item2","item3"]`
		require.JSONEq(t, expectedJSON, string(bodyBytes))
	})

	t.Run("set JSON body with nil", func(t *testing.T) {
		rb := &requestBuilder{}
		result, err := rb.SetJSONBody(nil)

		require.NoError(t, err)
		require.Equal(t, rb, result)
		require.NotNil(t, rb.body)
		require.Equal(t, "application/json", rb.contentType)

		bodyBytes, _ := io.ReadAll(rb.body)
		expectedJSON := `null`
		require.JSONEq(t, expectedJSON, string(bodyBytes))
	})

	t.Run("set JSON body with unmarshallable type", func(t *testing.T) {
		rb := &requestBuilder{}
		unmarshallable := make(chan int)
		result, err := rb.SetJSONBody(unmarshallable)

		require.Error(t, err)
		require.Equal(t, rb, result)
		require.Nil(t, rb.body)
		require.Empty(t, rb.contentType)
	})
}

func TestSetListPinsQueryParams(t *testing.T) {
	t.Run("with all fields set", func(t *testing.T) {
		rb := &requestBuilder{}
		options := &ListFilesOptions{
			Cid:          "testCid",
			GroupID:      "testGroupId",
			Status:       "testStatus",
			PageLimit:    10,
			PageOffset:   5,
			PinSizeMin:   100,
			PinSizeMax:   1000,
			PinStart:     &time.Time{},
			PinEnd:       &time.Time{},
			UnpinStart:   &time.Time{},
			UnpinEnd:     &time.Time{},
			IncludeCount: true,
			Metadata:     map[string]interface{}{"key": "value"},
		}

		result := rb.setListPinsQueryParams(options)

		require.Equal(t, rb, result)
		require.Equal(t, "testCid", rb.queryParams["cid"])
		require.Equal(t, "testGroupId", rb.queryParams["groupId"])
		require.Equal(t, "testStatus", rb.queryParams["status"])
		require.Equal(t, "10", rb.queryParams["pageLimit"])
		require.Equal(t, "5", rb.queryParams["pageOffset"])
		require.Equal(t, "100", rb.queryParams["pinSizeMin"])
		require.Equal(t, "1000", rb.queryParams["pinSizeMax"])
		require.Contains(t, rb.queryParams["pinStart"], "0001-01-01")
		require.Contains(t, rb.queryParams["pinEnd"], "0001-01-01")
		require.Contains(t, rb.queryParams["unpinStart"], "0001-01-01")
		require.Contains(t, rb.queryParams["unpinEnd"], "0001-01-01")
		require.Equal(t, "true", rb.queryParams["includeCount"])
		require.Equal(t, `{"key":"value"}`, rb.queryParams["metadata"])
	})

	t.Run("with minimal fields set", func(t *testing.T) {
		rb := &requestBuilder{}
		options := &ListFilesOptions{
			Cid:          "testCid",
			IncludeCount: false,
		}

		result := rb.setListPinsQueryParams(options)

		require.Equal(t, rb, result)
		require.Equal(t, "testCid", rb.queryParams["cid"])
		require.Equal(t, "false", rb.queryParams["includeCount"])
		require.Len(t, rb.queryParams, 2)
	})

	t.Run("with zero values", func(t *testing.T) {
		rb := &requestBuilder{}
		options := &ListFilesOptions{
			PageLimit:    0,
			PageOffset:   0,
			PinSizeMin:   0,
			PinSizeMax:   0,
			IncludeCount: false,
		}

		result := rb.setListPinsQueryParams(options)

		require.Equal(t, rb, result)
		require.Equal(t, "false", rb.queryParams["includeCount"])
		require.Len(t, rb.queryParams, 1)
	})

	t.Run("with nil time pointers", func(t *testing.T) {
		rb := &requestBuilder{}
		options := &ListFilesOptions{
			PinStart:   nil,
			PinEnd:     nil,
			UnpinStart: nil,
			UnpinEnd:   nil,
		}

		result := rb.setListPinsQueryParams(options)

		require.Equal(t, rb, result)
		require.NotContains(t, rb.queryParams, "pinStart")
		require.NotContains(t, rb.queryParams, "pinEnd")
		require.NotContains(t, rb.queryParams, "unpinStart")
		require.NotContains(t, rb.queryParams, "unpinEnd")
	})

	t.Run("with invalid metadata", func(t *testing.T) {
		rb := &requestBuilder{}
		options := &ListFilesOptions{
			Metadata: nil,
		}

		result := rb.setListPinsQueryParams(options)

		require.Equal(t, rb, result)
		require.NotContains(t, rb.queryParams, "metadata")
	})
}

func TestSetListApiKeysQueryParams(t *testing.T) {
	t.Run("with all fields set", func(t *testing.T) {
		rb := &requestBuilder{}
		revoked := true
		limitedUse := false
		exhausted := true
		options := &ListApiKeysOptions{
			Name:       "testName",
			Offset:     10,
			Revoked:    &revoked,
			LimitedUse: &limitedUse,
			Exhausted:  &exhausted,
		}

		result := rb.setListApiKeysQueryParams(options)

		require.Equal(t, rb, result)
		require.Equal(t, "testName", rb.queryParams["name"])
		require.Equal(t, "10", rb.queryParams["offset"])
		require.Equal(t, "true", rb.queryParams["revoked"])
		require.Equal(t, "false", rb.queryParams["limitedUse"])
		require.Equal(t, "true", rb.queryParams["exhausted"])
	})

	t.Run("with only name set", func(t *testing.T) {
		rb := &requestBuilder{}
		options := &ListApiKeysOptions{
			Name: "testName",
		}

		result := rb.setListApiKeysQueryParams(options)

		require.Equal(t, rb, result)
		require.Equal(t, "testName", rb.queryParams["name"])
		require.Len(t, rb.queryParams, 1)
	})

	t.Run("with only offset set", func(t *testing.T) {
		rb := &requestBuilder{}
		options := &ListApiKeysOptions{
			Offset: 5,
		}

		result := rb.setListApiKeysQueryParams(options)

		require.Equal(t, rb, result)
		require.Equal(t, "5", rb.queryParams["offset"])
		require.Len(t, rb.queryParams, 1)
	})

	t.Run("with zero offset", func(t *testing.T) {
		rb := &requestBuilder{}
		options := &ListApiKeysOptions{
			Offset: 0,
		}

		result := rb.setListApiKeysQueryParams(options)

		require.Equal(t, rb, result)
		require.Len(t, rb.queryParams, 0)
	})

	t.Run("with only boolean fields set", func(t *testing.T) {
		rb := &requestBuilder{}
		revoked := false
		limitedUse := true
		exhausted := false
		options := &ListApiKeysOptions{
			Revoked:    &revoked,
			LimitedUse: &limitedUse,
			Exhausted:  &exhausted,
		}

		result := rb.setListApiKeysQueryParams(options)

		require.Equal(t, rb, result)
		require.Equal(t, "false", rb.queryParams["revoked"])
		require.Equal(t, "true", rb.queryParams["limitedUse"])
		require.Equal(t, "false", rb.queryParams["exhausted"])
		require.Len(t, rb.queryParams, 3)
	})

	t.Run("with nil boolean fields", func(t *testing.T) {
		rb := &requestBuilder{}
		options := &ListApiKeysOptions{
			Name:   "testName",
			Offset: 5,
		}

		result := rb.setListApiKeysQueryParams(options)

		require.Equal(t, rb, result)
		require.Equal(t, "testName", rb.queryParams["name"])
		require.Equal(t, "5", rb.queryParams["offset"])
		require.Len(t, rb.queryParams, 2)
		require.NotContains(t, rb.queryParams, "revoked")
		require.NotContains(t, rb.queryParams, "limitedUse")
		require.NotContains(t, rb.queryParams, "exhausted")
	})
}

func TestSetListGroupsQueryParams(t *testing.T) {
	t.Run("with all fields set", func(t *testing.T) {
		rb := &requestBuilder{}
		options := &ListGroupsOptions{
			NameContains: "test",
			Limit:        10,
			Offset:       5,
		}

		result := rb.setListGroupsQueryParams(options)

		require.Equal(t, rb, result)
		require.Equal(t, "test", rb.queryParams["nameContains"])
		require.Equal(t, "10", rb.queryParams["limit"])
		require.Equal(t, "5", rb.queryParams["offset"])
	})

	t.Run("with only nameContains set", func(t *testing.T) {
		rb := &requestBuilder{}
		options := &ListGroupsOptions{
			NameContains: "group",
		}

		result := rb.setListGroupsQueryParams(options)

		require.Equal(t, rb, result)
		require.Equal(t, "group", rb.queryParams["nameContains"])
		require.Len(t, rb.queryParams, 1)
	})

	t.Run("with only limit set", func(t *testing.T) {
		rb := &requestBuilder{}
		options := &ListGroupsOptions{
			Limit: 20,
		}

		result := rb.setListGroupsQueryParams(options)

		require.Equal(t, rb, result)
		require.Equal(t, "20", rb.queryParams["limit"])
		require.Len(t, rb.queryParams, 1)
	})

	t.Run("with only offset set", func(t *testing.T) {
		rb := &requestBuilder{}
		options := &ListGroupsOptions{
			Offset: 15,
		}

		result := rb.setListGroupsQueryParams(options)

		require.Equal(t, rb, result)
		require.Equal(t, "15", rb.queryParams["offset"])
		require.Len(t, rb.queryParams, 1)
	})

	t.Run("with zero values", func(t *testing.T) {
		rb := &requestBuilder{}
		options := &ListGroupsOptions{
			NameContains: "",
			Limit:        0,
			Offset:       0,
		}

		result := rb.setListGroupsQueryParams(options)

		require.Equal(t, rb, result)
		require.Len(t, rb.queryParams, 0)
	})

	t.Run("with negative limit and offset", func(t *testing.T) {
		rb := &requestBuilder{}
		options := &ListGroupsOptions{
			Limit:  -10,
			Offset: -5,
		}

		result := rb.setListGroupsQueryParams(options)

		require.Equal(t, rb, result)
		require.Len(t, rb.queryParams, 0)
	})
}

func TestSetListPinsByCidQueryParams(t *testing.T) {
	t.Run("with all fields set", func(t *testing.T) {
		rb := &requestBuilder{}
		options := &ListPinByCidOptions{
			Sort:        SortOrderASC,
			Status:      PinStatusRetrieving,
			IPFSPinHash: "QmTest123",
			Limit:       100,
			Offset:      10,
		}

		result := rb.setListPinsByCidQueryParams(options)

		require.Equal(t, rb, result)
		require.Equal(t, string(SortOrderASC), rb.queryParams["sort"])
		require.Equal(t, string(PinStatusRetrieving), rb.queryParams["status"])
		require.Equal(t, "QmTest123", rb.queryParams["ipfs_pin_hash"])
		require.Equal(t, "100", rb.queryParams["limit"])
		require.Equal(t, "10", rb.queryParams["offset"])
	})

	t.Run("with only sort and status", func(t *testing.T) {
		rb := &requestBuilder{}
		options := &ListPinByCidOptions{
			Sort:   SortOrderASC,
			Status: PinStatusRetrieving,
		}

		result := rb.setListPinsByCidQueryParams(options)

		require.Equal(t, rb, result)
		require.Equal(t, string(SortOrderASC), rb.queryParams["sort"])
		require.Equal(t, string(PinStatusRetrieving), rb.queryParams["status"])
		require.NotContains(t, rb.queryParams, "ipfs_pin_hash")
		require.NotContains(t, rb.queryParams, "limit")
		require.NotContains(t, rb.queryParams, "offset")
	})

	t.Run("with only IPFSPinHash", func(t *testing.T) {
		rb := &requestBuilder{}
		options := &ListPinByCidOptions{
			IPFSPinHash: "QmTest456",
		}

		result := rb.setListPinsByCidQueryParams(options)

		require.Equal(t, rb, result)
		require.Equal(t, "QmTest456", rb.queryParams["ipfs_pin_hash"])
		require.NotContains(t, rb.queryParams, "sort")
		require.NotContains(t, rb.queryParams, "status")
		require.NotContains(t, rb.queryParams, "limit")
		require.NotContains(t, rb.queryParams, "offset")
	})

	t.Run("with zero limit and offset", func(t *testing.T) {
		rb := &requestBuilder{}
		options := &ListPinByCidOptions{
			Limit:  0,
			Offset: 0,
		}

		result := rb.setListPinsByCidQueryParams(options)

		require.Equal(t, rb, result)
		require.NotContains(t, rb.queryParams, "limit")
		require.NotContains(t, rb.queryParams, "offset")
	})

	t.Run("with negative limit and offset", func(t *testing.T) {
		rb := &requestBuilder{}
		options := &ListPinByCidOptions{
			Limit:  -10,
			Offset: -5,
		}

		result := rb.setListPinsByCidQueryParams(options)

		require.Equal(t, rb, result)
		require.NotContains(t, rb.queryParams, "limit")
		require.NotContains(t, rb.queryParams, "offset")
	})

	t.Run("with empty options", func(t *testing.T) {
		rb := &requestBuilder{}
		options := &ListPinByCidOptions{}

		result := rb.setListPinsByCidQueryParams(options)

		require.Equal(t, rb, result)
		require.Empty(t, rb.queryParams)
	})
}

func TestBuildURL(t *testing.T) {
	t.Run("successful URL build with path params and query params", func(t *testing.T) {
		rb := &requestBuilder{
			client: &client{BaseURL: "https://api.pinata.cloud"},
			path:   "/v1/pinning/{pinType}/{hashToPin}",
			pathParams: map[string]string{
				"pinType":   "pinByHash",
				"hashToPin": "QmTest123",
			},
			queryParams: map[string]string{
				"pinataMetadata": `{"name":"TestFile"}`,
			},
		}

		url, err := rb.buildURL()

		require.NoError(t, err)
		require.Equal(t, "https://api.pinata.cloud/v1/pinning/pinByHash/QmTest123?pinataMetadata=%7B%22name%22%3A%22TestFile%22%7D", url)
	})

	t.Run("error when path parameter is wrong", func(t *testing.T) {
		rb := &requestBuilder{
			client: &client{BaseURL: "https://api.pinata.cloud"},
			path:   "/v1/pinning/{pinType}/{hashToPin1}",
			pathParams: map[string]string{
				"hashToPin": "hashToPin",
			},
		}

		_, err := rb.buildURL()

		require.Error(t, err)
		require.Contains(t, err.Error(), "path parameter hashToPin not found in path")
	})

	t.Run("URL encoding of path parameters", func(t *testing.T) {
		rb := &requestBuilder{
			client: &client{BaseURL: "https://api.pinata.cloud"},
			path:   "/v1/files/{fileName}",
			pathParams: map[string]string{
				"fileName": "test file with spaces.txt",
			},
		}

		url, err := rb.buildURL()

		require.NoError(t, err)
		require.Equal(t, "https://api.pinata.cloud/v1/files/test%20file%20with%20spaces.txt", url)
	})

	t.Run("multiple query parameters", func(t *testing.T) {
		rb := &requestBuilder{
			client: &client{BaseURL: "https://api.pinata.cloud"},
			path:   "/v1/data",
			queryParams: map[string]string{
				"status": "pinned",
				"limit":  "10",
				"offset": "0",
			},
		}

		url, err := rb.buildURL()

		require.NoError(t, err)
		require.Equal(t, "https://api.pinata.cloud/v1/data?limit=10&offset=0&status=pinned", url)
	})
}

func TestSend(t *testing.T) {
	t.Run("successful request with JSON response", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/test", r.URL.Path)
			require.Equal(t, "Bearer test_token", r.Header.Get("Authorization"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"key": "value"}`))
		}))
		defer mockServer.Close()

		client := &client{
			BaseURL:    mockServer.URL,
			HTTPClient: mockServer.Client(),
			Auth:       NewAuthWithJWT("test_token"),
		}

		rb := &requestBuilder{
			client: client,
			method: http.MethodGet,
			path:   "/test",
		}

		var result map[string]string
		err := rb.Send(&result)

		require.NoError(t, err)
		require.Equal(t, map[string]string{"key": "value"}, result)
	})

	t.Run("request with query parameters", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/test", r.URL.Path)
			require.Equal(t, "param1=value1&param2=value2", r.URL.RawQuery)

			w.WriteHeader(http.StatusOK)
		}))
		defer mockServer.Close()

		client := &client{
			BaseURL:    mockServer.URL,
			HTTPClient: mockServer.Client(),
			Auth:       NewAuthWithJWT("test_token"),
		}

		rb := &requestBuilder{
			client: client,
			method: http.MethodGet,
			path:   "/test",
			queryParams: map[string]string{
				"param1": "value1",
				"param2": "value2",
			},
		}

		err := rb.Send(nil)

		require.NoError(t, err)
	})

	t.Run("request with custom headers", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/test", r.URL.Path)
			require.Equal(t, "custom_value", r.Header.Get("Custom-Header"))

			w.WriteHeader(http.StatusOK)
		}))
		defer mockServer.Close()

		client := &client{
			BaseURL:    mockServer.URL,
			HTTPClient: mockServer.Client(),
			Auth:       NewAuthWithJWT("test_token"),
		}

		rb := &requestBuilder{
			client: client,
			method: http.MethodPost,
			path:   "/test",
			headers: map[string]string{
				"Custom-Header": "custom_value",
			},
		}

		err := rb.Send(nil)

		require.NoError(t, err)
	})

	t.Run("request with body", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/test", r.URL.Path)
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))

			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			require.Equal(t, `{"key":"value"}`, string(body))

			w.WriteHeader(http.StatusOK)
		}))
		defer mockServer.Close()

		client := &client{
			BaseURL:    mockServer.URL,
			HTTPClient: mockServer.Client(),
			Auth:       NewAuthWithJWT("test_token"),
		}

		rb := &requestBuilder{
			client:      client,
			method:      http.MethodPost,
			path:        "/test",
			contentType: "application/json",
			body:        strings.NewReader(`{"key":"value"}`),
		}

		err := rb.Send(nil)

		require.NoError(t, err)
	})

	t.Run("error response", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Bad Request"}`))
		}))
		defer mockServer.Close()

		client := &client{
			BaseURL:    mockServer.URL,
			HTTPClient: mockServer.Client(),
			Auth:       NewAuthWithJWT("test_token"),
		}

		rb := &requestBuilder{
			client: client,
			method: http.MethodGet,
			path:   "/test",
		}

		err := rb.Send(nil)

		require.Error(t, err)
		require.Contains(t, err.Error(), "Bad Request")
	})

	t.Run("network error", func(t *testing.T) {
		client := &client{
			BaseURL: "http://non-existent-url",
			HTTPClient: &http.Client{
				Timeout: time.Millisecond,
			},
			Auth: NewAuthWithJWT("test_token"),
		}

		rb := &requestBuilder{
			client: client,
			method: http.MethodGet,
			path:   "/test",
		}

		err := rb.Send(nil)

		require.Error(t, err)
	})
}
