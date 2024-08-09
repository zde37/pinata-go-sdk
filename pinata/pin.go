package pinata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type SortOrder string

const (
	SortOrderASC  SortOrder = "ASC"
	SortOrderDESC SortOrder = "DESC"
)

type PinStatus string

const (
	PinStatusPrechecking   PinStatus = "prechecking"
	PinStatusRetrieving    PinStatus = "retrieving"
	PinStatusExpired       PinStatus = "expired"
	PinStatusOverFreeLimit PinStatus = "over_free_limit"
	PinStatusOverMaxSize   PinStatus = "over_max_size"
	PinStatusInvalidObject PinStatus = "invalid_object"
	PinStatusBadHostNode   PinStatus = "bad_host_node"
)

type PinOptions struct {
	PinataMetadata PinataMetadata `json:"pinataMetadata,omitempty"`
	PinataOptions  struct {
		CidVersion int `json:"cidVersion,omitempty"`
	} `json:"pinataOptions,omitempty"`
}

type PinByCidOptions struct { 
	PinataOptions struct {
		GroupId   string   `json:"groupId,omitempty"`
		HostNodes []string `json:"hostNodes,omitempty"`
	} `json:"pinataOptions,omitempty"`
	PinataMetadata PinataMetadata `json:"pinataMetadata,omitempty"`
}

type PinByCidResponse struct {
	Id       string `json:"id,omitempty"`
	IpfsHash string `json:"ipfsHash,omitempty"`
	Status   string `json:"status,omitempty"`
	Name     string `json:"name,omitempty"`
}

type PinataMetadata struct {
	Name      string                 `json:"name,omitempty"`
	KeyValues map[string]interface{} `json:"keyvalues,omitempty"`
}

type PinResponse struct {
	IpfsHash    string `json:"IpfsHash,omitempty"`
	PinSize     int    `json:"PinSize,omitempty"`
	Timestamp   string `json:"Timestamp,omitempty"`
	IsDuplicate bool   `json:"IsDuplicate,omitempty"`
}

type PinMetadataUpdateOptions struct {
	Name      string                 `json:"name,omitempty"`
	KeyValues map[string]interface{} `json:"keyvalues,omitempty"`
}

type ListFilesOptions struct {
	Cid          string                 `json:"cid,omitempty"`
	GroupID      string                 `json:"groupId,omitempty"`
	Status       string                 `json:"status,omitempty"`
	PageLimit    int                    `json:"pageLimit,omitempty"`
	PageOffset   int                    `json:"pageOffset,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	PinSizeMin   int64                  `json:"pinSizeMin,omitempty"`
	PinSizeMax   int64                  `json:"pinSizeMax,omitempty"`
	PinStart     *time.Time             `json:"pinStart,omitempty"`
	PinEnd       *time.Time             `json:"pinEnd,omitempty"`
	UnpinStart   *time.Time             `json:"unpinStart,omitempty"`
	UnpinEnd     *time.Time             `json:"unpinEnd,omitempty"`
	IncludeCount bool                   `json:"includeCount,omitempty"`
}

type ListFilesResponse struct {
	Count int   `json:"count,omitempty"`
	Rows  []Pin `json:"rows,omitempty"`
}

type Pin struct {
	ID            string                 `json:"id,omitempty"`
	IPFSPinHash   string                 `json:"ipfs_pin_hash,omitempty"`
	Size          int                    `json:"size,omitempty"`
	UserID        string                 `json:"user_id,omitempty"`
	DatePinned    string                 `json:"date_pinned,omitempty"`
	DateUnpinned  string                 `json:"date_unpinned,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	Regions       []Region               `json:"regions,omitempty"`
	MimeType      string                 `json:"mime_type,omitempty"`
	NumberOfFiles int                    `json:"number_of_files,omitempty"`
}

type Region struct {
	RegionID                string `json:"regionId,omitempty"`
	CurrentReplicationCount int    `json:"currentReplicationCount,omitempty"`
	DesiredReplicationCount int    `json:"desiredReplicationCount,omitempty"`
}

type ListPinByCidOptions struct {
	Sort        SortOrder `json:"sort,omitempty"`
	Status      PinStatus `json:"status,omitempty"`
	IPFSPinHash string    `json:"ipfs_pin_hash,omitempty"`
	Limit       int       `json:"limit,omitempty"`
	Offset      int       `json:"offset,omitempty"`
}

type ListPinByCidResponse struct {
	Count int        `json:"count"`
	Rows  []PinEntry `json:"rows"`
}

type PinEntry struct {
	ID          string      `json:"id"`
	IPFSPinHash string      `json:"ipfs_pin_hash"`
	DateQueued  string      `json:"date_queued"`
	Name        string      `json:"name"`
	Status      string      `json:"status"`
	KeyValues   interface{} `json:"keyvalues"`
	HostNodes   []string    `json:"host_nodes"`
	PinPolicy   PinPolicy   `json:"pin_policy"`
}

type PinPolicy struct {
	Regions []struct {
		ID                      string `json:"id"`
		DesiredReplicationCount int    `json:"desiredReplicationCount"`
	} `json:"regions"`
	Version int `json:"version"`
}

func (c *Client) PinFileToIPFS(path string, options *PinOptions) (*PinResponse, error) {
	if path == "" {
		return nil, fmt.Errorf("ERR: filepath is required")
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("ERR: failed to open file: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		return nil, fmt.Errorf("ERR: failed to create form file: %w", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("ERR: failed to copy file content: %w", err)
	}

	if options != nil {
		optionsJSON, err := json.Marshal(options)
		if err != nil {
			return nil, fmt.Errorf("ERR: failed to marshal options: %w", err)
		}
		err = writer.WriteField("pinataOptions", string(optionsJSON))
		if err != nil {
			return nil, fmt.Errorf("ERR: failed to write pinataOptions field: %w", err)
		}
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("ERR: failed to close multipart writer: %w", err)
	}

	var response PinResponse
	err = c.NewRequest("POST", "/pinning/pinFileToIPFS").
		SetBody(body, writer.FormDataContentType()).
		Send(&response)

	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) PinJSONToIPFS(data interface{}, options *PinOptions) (*PinResponse, error) {
	if data == nil {
		return nil, fmt.Errorf("ERR: jsonData is required")
	}
	payload := map[string]interface{}{
		"pinataContent": data,
	}

	if options != nil {
		payload["pinataOptions"] = options.PinataOptions
		payload["pinataMetadata"] = options.PinataMetadata
	}

	req, err := c.NewRequest("POST", "/pinning/pinJSONToIPFS").SetJSONBody(payload)
	if err != nil {
		return nil, fmt.Errorf("ERR: failed to set JSON body: %w", err)
	}

	var response PinResponse
	err = req.Send(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) PinByCid(hashToPin string, options *PinByCidOptions) (*PinByCidResponse, error) {
	if hashToPin == "" {
		return nil, fmt.Errorf("ERR: hashToPin is required")
	}
	payload := map[string]interface{}{
		"hashToPin": hashToPin,
	}

	if options != nil {
		payload["pinataOptions"] = options.PinataOptions
		payload["pinataMetadata"] = options.PinataMetadata
	}

	req, err := c.NewRequest("POST", "/pinning/pinByHash").SetJSONBody(payload)
	if err != nil {
		return nil, fmt.Errorf("ERR: failed to set JSON body: %w", err)
	}

	var response PinByCidResponse
	err = req.Send(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) ListFiles(options *ListFilesOptions) (*ListFilesResponse, error) {
	req := c.NewRequest("GET", "/data/pinList")
	if options != nil {
		req.addListPinsQueryParams(options)
	}

	var response ListFilesResponse
	err := req.Send(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) ListPinByCidJobs(options *ListPinByCidOptions) (*ListPinByCidResponse, error) {
	req := c.NewRequest("GET", "/pinning/pinJobs")
	if options != nil {
		req.addListPinsByCidQueryParams(options)
	}

	var response ListPinByCidResponse
	err := req.Send(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) UpdateFileMetadata(fileHash string, options *PinMetadataUpdateOptions) error {
	if fileHash == "" || options == nil {
		return fmt.Errorf("ERR: fileHash and options are required")
	}

	payload := make(map[string]interface{})
	payload["ipfsPinHash"] = fileHash // "ipfsPinHash" wasn't shown in the docs site, inform pinata team
	payload["name"] = options.Name
	payload["keyvalues"] = options.KeyValues

	req, err := c.NewRequest("PUT", "/pinning/hashMetadata").SetJSONBody(payload)
	if err != nil {
		return fmt.Errorf("ERR: failed to set JSON body: %w", err)
	}

	err = req.Send(nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteFile(cid string) error {
	if cid == "" {
		return fmt.Errorf("cid is required")
	}

	err := c.NewRequest("DELETE", "/pinning/unpin/{cid}").
		AddPathParam("cid", cid).
		Send(nil)

	if err != nil {
		return err
	}
	return nil
}
