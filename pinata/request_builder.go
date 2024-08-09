package pinata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type RequestBuilder interface {
}

type requestBuilder struct {
	client      *Client
	method      string
	path        string
	pathParams  map[string]string
	queryParams map[string]string
	headers     map[string]string
	body        io.Reader
	contentType string
}

type AsyncResult struct {
	Response interface{}
	Error    error
}

func (rb *requestBuilder) AddPathParam(key, value string) *requestBuilder {
	if rb.pathParams == nil {
		rb.pathParams = make(map[string]string)
	}
	rb.pathParams[key] = value
	return rb
}

func (rb *requestBuilder) AddQueryParam(key string, value interface{}) *requestBuilder {
	if rb.queryParams == nil {
		rb.queryParams = make(map[string]string)
	}
	rb.queryParams[key] = fmt.Sprintf("%v", value)
	return rb
}

func (rb *requestBuilder) AddHeaders(key, value string) *requestBuilder {
	if rb.headers == nil {
		rb.headers = make(map[string]string)
	}
	rb.headers[key] = value
	return rb
}

func (rb *requestBuilder) SetBody(body io.Reader, contentType string) *requestBuilder {
	rb.body = body
	rb.contentType = contentType
	return rb
}

// Helper method for setting JSON body
func (rb *requestBuilder) SetJSONBody(body interface{}) (*requestBuilder, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return rb, err
	}
	return rb.SetBody(bytes.NewReader(jsonBody), "application/json"), nil
}

func (rb *requestBuilder) addListPinsQueryParams(options *ListFilesOptions) *requestBuilder {
	if options.Cid != "" {
		rb.AddQueryParam("cid", options.Cid)
	}
	if options.GroupID != "" {
		rb.AddQueryParam("groupId", options.GroupID)
	}
	if options.Status != "" {
		rb.AddQueryParam("status", options.Status)
	}
	if options.PageLimit > 0 {
		rb.AddQueryParam("pageLimit", options.PageLimit)
	}
	if options.PageOffset > 0 {
		rb.AddQueryParam("pageOffset", options.PageOffset)
	}
	if options.PinSizeMin > 0 {
		rb.AddQueryParam("pinSizeMin", options.PinSizeMin)
	}
	if options.PinSizeMax > 0 {
		rb.AddQueryParam("pinSizeMax", options.PinSizeMax)
	}
	if options.PinStart != nil {
		rb.AddQueryParam("pinStart", options.PinStart.Format(time.RFC3339))
	}
	if options.PinEnd != nil {
		rb.AddQueryParam("pinEnd", options.PinEnd.Format(time.RFC3339))
	}
	if options.UnpinStart != nil {
		rb.AddQueryParam("unpinStart", options.UnpinStart.Format(time.RFC3339))
	}
	if options.UnpinEnd != nil {
		rb.AddQueryParam("unpinEnd", options.UnpinEnd.Format(time.RFC3339))
	}
	rb.AddQueryParam("includeCount", options.IncludeCount)

	if options.Metadata != nil {
		metadataJSON, err := json.Marshal(options.Metadata)
		if err == nil {
			rb.AddQueryParam("metadata", string(metadataJSON))
		}
	}

	return rb
}

func (rb *requestBuilder) addListPinsByCidQueryParams(options *ListPinByCidOptions) *requestBuilder {
	if options.Sort != "" {
		rb.AddQueryParam("sort", string(options.Sort))
	}
	if options.Status != "" {
		rb.AddQueryParam("status", string(options.Status))
	}
	if options.IPFSPinHash != "" {
		rb.AddQueryParam("ipfs_pin_hash", options.IPFSPinHash)
	}
	if options.Limit > 0 {
		rb.AddQueryParam("limit", options.Limit)
	}
	if options.Offset > 0 {
		rb.AddQueryParam("offset", options.Offset)
	}
	return rb
}

func (rb *requestBuilder) buildURL() (string, error) {
	path := rb.path
	for key, value := range rb.pathParams {
		placeholder := "{" + key + "}"
		if !strings.Contains(path, placeholder) {
			return "", fmt.Errorf("path parameter %s not found in path", key)
		}
		path = strings.Replace(path, placeholder, url.PathEscape(value), -1)
	}

	reqURL, err := url.Parse(rb.client.BaseURL + path)
	if err != nil {
		return "", err
	}

	// Add query parameters
	q := reqURL.Query()
	for k, v := range rb.queryParams {
		q.Add(k, v)
	}
	reqURL.RawQuery = q.Encode()

	return reqURL.String(), nil
}

func (rb *requestBuilder) Send(v interface{}) error {
	reqURL, err := rb.buildURL()
	if err != nil {
		return err
	}

	req, err := http.NewRequest(rb.method, reqURL, rb.body)
	if err != nil {
		return err
	}

	// Set headers
	for k, v := range rb.headers {
		req.Header.Set(k, v)
	}

	// Set auth header
	rb.client.Auth.setAuthHeader(req)

	// Set content type if body is present
	if rb.body != nil {
		req.Header.Set("Content-Type", rb.contentType)
	}

	resp, err := rb.client.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errorMsg interface{} // TODO: use a concrete type here
		if err := json.NewDecoder(resp.Body).Decode(&errorMsg); err != nil {
			return err
		}
		return fmt.Errorf("%v", errorMsg)
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return err
		}
	}

	return nil
}
