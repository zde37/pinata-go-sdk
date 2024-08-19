package pinata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
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

// PinOptions represents the options for pinning a file or directory to Pinata.
// PinataMetadata contains metadata about the file or directory being pinned.
// PinataOptions contains options specific to the Pinata platform, such as the CID version.
type PinOptions struct {
	PinataMetadata PinataMetadata `json:"pinataMetadata,omitempty"`
	PinataOptions  struct {
		CidVersion int `json:"cidVersion,omitempty"`
	} `json:"pinataOptions,omitempty"`
}

// PinByCidOptions represents the options for pinning a file or directory to Pinata by its CID.
// PinataOptions contains options specific to the Pinata platform, such as the group ID and host nodes.
// PinataMetadata contains metadata about the file or directory being pinned.
type PinByCidOptions struct {
	PinataOptions struct {
		GroupId   string   `json:"groupId,omitempty"`
		HostNodes []string `json:"hostNodes,omitempty"`
	} `json:"pinataOptions,omitempty"`
	PinataMetadata PinataMetadata `json:"pinataMetadata,omitempty"`
}

// pinByCidResponse represents the response from pinning a file or directory to Pinata by its CID.
// ID is the unique identifier for the pin.
// IpfsHash is the IPFS hash of the pinned content.
// Status is the status of the pin operation.
// Name is the name of the pinned content.
type pinByCidResponse struct {
	ID       string `json:"id,omitempty"`
	IpfsHash string `json:"ipfsHash,omitempty"`
	Status   string `json:"status,omitempty"`
	Name     string `json:"name,omitempty"`
}

// PinataMetadata represents metadata associated with a file or directory pinned to Pinata.
// Name is the name of the pinned content.
// KeyValues is a map of key-value pairs containing additional metadata about the pinned content.
type PinataMetadata struct {
	Name      string                 `json:"name,omitempty"`
	KeyValues map[string]interface{} `json:"keyvalues,omitempty"`
}

// pinResponse represents the response from pinning a file or directory to Pinata.
// IpfsHash is the IPFS hash of the pinned content.
// PinSize is the size of the pinned content in bytes.
// Timestamp is the timestamp of when the content was pinned.
// IsDuplicate indicates whether the pinned content is a duplicate of an existing pin.
type pinResponse struct {
	IpfsHash    string `json:"IpfsHash,omitempty"`
	PinSize     int    `json:"PinSize,omitempty"`
	Timestamp   string `json:"Timestamp,omitempty"`
	IsDuplicate bool   `json:"IsDuplicate,omitempty"`
}

// PinMetadataUpdateOptions represents the options for updating the metadata of a file or directory pinned to Pinata.
// Name is the new name for the pinned content.
// KeyValues is a map of new key-value pairs containing additional metadata about the pinned content.
type PinMetadataUpdateOptions struct {
	Name      string                 `json:"name,omitempty"`
	KeyValues map[string]interface{} `json:"keyvalues,omitempty"`
}

// ListFilesOptions represents the options for listing files pinned to Pinata.
// Cid is the IPFS content identifier to filter pins by.
// GroupID is the ID of the group to filter pins by.
// Status is the status to filter pins by.
// PageLimit is the maximum number of pins to return per page.
// PageOffset is the number of pins to skip before returning results.
// Metadata is a map of key-value pairs to filter pins by.
// PinSizeMin is the minimum size in bytes of pins to return.
// PinSizeMax is the maximum size in bytes of pins to return.
// PinStart is the earliest date that pins were created.
// PinEnd is the latest date that pins were created.
// UnpinStart is the earliest date that pins were unpinned.
// UnpinEnd is the latest date that pins were unpinned.
// IncludeCount indicates whether to include the total count of matching pins.
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

// listFilesResponse represents the response from listing files pinned to Pinata.
// Count is the total number of pinned files.
// Rows is a slice of Pin structs representing the pinned files.
type listFilesResponse struct {
	Count int   `json:"count,omitempty"`
	Rows  []pin `json:"rows,omitempty"`
}

// pin represents a file or directory that has been pinned to Pinata.
// ID is the unique identifier for the pinned content.
// IPFSPinHash is the IPFS content identifier for the pinned content.
// Size is the size of the pinned content in bytes.
// UserID is the ID of the user who pinned the content.
// DatePinned is the date the content was pinned.
// DateUnpinned is the date the content was unpinned, if applicable.
// Metadata is a map of key-value pairs containing additional metadata about the pinned content.
// Regions is a slice of Region structs representing the regions where the pinned content is replicated.
// MimeType is the MIME type of the pinned content.
// NumberOfFiles is the number of files in the pinned content.
type pin struct {
	ID            string                 `json:"id,omitempty"`
	IPFSPinHash   string                 `json:"ipfs_pin_hash,omitempty"`
	Size          int                    `json:"size,omitempty"`
	UserID        string                 `json:"user_id,omitempty"`
	DatePinned    string                 `json:"date_pinned,omitempty"`
	DateUnpinned  string                 `json:"date_unpinned,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	Regions       []region               `json:"regions,omitempty"`
	MimeType      string                 `json:"mime_type,omitempty"`
	NumberOfFiles int                    `json:"number_of_files,omitempty"`
}

// region represents a geographic region where a file is pinned.
// RegionID is the unique identifier for the region.
// CurrentReplicationCount is the current number of replicas of the file in the region.
// DesiredReplicationCount is the desired number of replicas of the file in the region.
type region struct {
	RegionID                string `json:"regionId,omitempty"`
	CurrentReplicationCount int    `json:"currentReplicationCount,omitempty"`
	DesiredReplicationCount int    `json:"desiredReplicationCount,omitempty"`
}

// ListPinByCidOptions represents the options for listing pins by IPFS content identifier (CID).
// Sort specifies the sort order for the results.
// Status specifies the status of the pins to include in the results.
// IPFSPinHash specifies the IPFS content identifier to filter the results by.
// Limit specifies the maximum number of results to return.
// Offset specifies the number of results to skip before returning results.
type ListPinByCidOptions struct {
	Sort        SortOrder `json:"sort,omitempty"`
	Status      PinStatus `json:"status,omitempty"`
	IPFSPinHash string    `json:"ipfs_pin_hash,omitempty"`
	Limit       int       `json:"limit,omitempty"`
	Offset      int       `json:"offset,omitempty"`
}

// listPinByCidResponse represents the response from a request to list pins by IPFS content identifier (CID).
// Count is the total number of pins returned.
// Rows is a slice of PinEntry structs representing the pins that match the request.
type listPinByCidResponse struct {
	Count int        `json:"count,omitempty"`
	Rows  []pinEntry `json:"rows,omitempty"`
}

// pinEntry represents a single entry in the list of pinned content.
// ID is the unique identifier for the pinned content.
// IPFSPinHash is the IPFS content identifier (CID) for the pinned content.
// DateQueued is the date the content was queued for pinning.
// Name is the name of the pinned content.
// Status is the current status of the pinned content (e.g. queued, pinned, failed).
// KeyValues is a map of key-value pairs containing additional metadata about the pinned content.
// HostNodes is a list of node IDs where the pinned content is currently hosted.
// PinPolicy is the policy that governs how the pinned content is replicated across regions.
type pinEntry struct {
	ID          string      `json:"id,omitempty"`
	IPFSPinHash string      `json:"ipfs_pin_hash,omitempty"`
	DateQueued  string      `json:"date_queued,omitempty"`
	Name        string      `json:"name,omitempty"`
	Status      string      `json:"status,omitempty"`
	KeyValues   interface{} `json:"keyvalues,omitempty"`
	HostNodes   []string    `json:"host_nodes,omitempty"`
	PinPolicy   pinPolicy   `json:"pin_policy,omitempty"`
}

// pinPolicy represents the policy for pinning a file to IPFS.
// Regions specifies the geographic regions where the file should be pinned, and the desired replication count for each region.
// Version specifies the version of the pin policy.
type pinPolicy struct {
	Regions []struct {
		ID                      string `json:"id,omitempty"`
		DesiredReplicationCount int    `json:"desiredReplicationCount,omitempty"`
	} `json:"regions,omitempty"`
	Version int `json:"version,omitempty"`
}

// PinFileToIPFS uploads a file to IPFS and pins it to the Pinata network.
//
// path specifies the local file path of the file to be uploaded and pinned.
// options is an optional PinOptions struct that can be used to specify additional
// metadata and options for the pin operation.
//
// Returns a PinResponse struct containing the IPFS hash and other details of the
// pinned file, or an error if the operation fails.
func (c *client) PinFileToIPFS(path string, options *PinOptions) (*pinResponse, error) {
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

	var response pinResponse
	err = c.NewRequest(http.MethodPost, "/pinning/pinFileToIPFS").
		SetBody(body, writer.FormDataContentType()).
		Send(&response)

	if err != nil {
		return nil, err
	}

	return &response, nil
}

// PinJSONToIPFS pins the provided JSON data to IPFS using the Pinata API.
//
// The data parameter should be a JSON-serializable Go value. The options parameter
// can be used to provide additional metadata and options for the pin operation.
//
// This function returns a PinResponse containing the IPFS hash and other details
// of the pinned data, or an error if the operation fails.
func (c *client) PinJSONToIPFS(data interface{}, options *PinOptions) (*pinResponse, error) {
	if data == nil {
		return nil, fmt.Errorf("ERR: jsonData is required")
	}
	payload := make(map[string]interface{})
	payload["pinataContent"] = data

	if options != nil {
		payload["pinataOptions"] = options.PinataOptions
		payload["pinataMetadata"] = options.PinataMetadata
	}

	req, err := c.NewRequest(http.MethodPost, "/pinning/pinJSONToIPFS").SetJSONBody(payload)
	if err != nil {
		return nil, fmt.Errorf("ERR: failed to set JSON body: %w", err)
	}

	var response pinResponse
	err = req.Send(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// PinByCid pins the content identified by the provided hashToPin to IPFS using the Pinata API.
// The optional PinByCidOptions can be used to provide additional metadata and options for the pin operation.
// Returns a PinByCidResponse containing information about the pinned content.
func (c *client) PinByCid(hashToPin string, options *PinByCidOptions) (*pinByCidResponse, error) {
	if hashToPin == "" {
		return nil, fmt.Errorf("ERR: hashToPin is required")
	}
	payload := make(map[string]interface{})
	payload["hashToPin"] = hashToPin

	if options != nil {
		payload["pinataOptions"] = options.PinataOptions
		payload["pinataMetadata"] = options.PinataMetadata
	}

	req, err := c.NewRequest(http.MethodPost, "/pinning/pinByHash").SetJSONBody(payload)
	if err != nil {
		return nil, fmt.Errorf("ERR: failed to set JSON body: %w", err)
	}

	var response pinByCidResponse
	err = req.Send(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// ListFiles returns a list of files that have been pinned to Pinata.
// The options parameter can be used to filter the list of files.
func (c *client) ListFiles(options *ListFilesOptions) (*listFilesResponse, error) {
	req := c.NewRequest(http.MethodGet, "/data/pinList")
	if options != nil {
		req.setListPinsQueryParams(options)
	}

	var response listFilesResponse
	err := req.Send(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// ListPinByCidJobs returns a list of pin jobs for the provided ListPinByCidOptions.
// The ListPinByCidOptions can be used to filter the list of pin jobs.
// Returns a listPinByCidResponse containing information about the pin jobs.
func (c *client) ListPinByCidJobs(options *ListPinByCidOptions) (*listPinByCidResponse, error) {
	req := c.NewRequest(http.MethodGet, "/pinning/pinJobs")
	if options != nil {
		req.setListPinsByCidQueryParams(options)
	}

	var response listPinByCidResponse
	err := req.Send(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// UpdateFileMetadata updates the metadata for a file that has been pinned to Pinata.
// The fileHash parameter specifies the hash of the file to update.
// The options parameter specifies the new metadata to apply, including the name and key-value pairs.
// Returns an error if the fileHash or options are not provided, or if there is an error updating the metadata.
func (c *client) UpdateFileMetadata(fileHash string, options *PinMetadataUpdateOptions) error {
	if fileHash == "" || options == nil {
		return fmt.Errorf("ERR: fileHash and options are required")
	}

	payload := make(map[string]interface{})
	payload["ipfsPinHash"] = fileHash // "ipfsPinHash" wasn't shown as a query param in the docs. Inform the pinata team
	payload["name"] = options.Name
	payload["keyvalues"] = options.KeyValues

	req, err := c.NewRequest(http.MethodPut, "/pinning/hashMetadata").SetJSONBody(payload)
	if err != nil {
		return fmt.Errorf("ERR: failed to set JSON body: %w", err)
	}

	err = req.Send(nil)
	if err != nil {
		return err
	}
	return nil
}

// DeleteFile deletes the file with the given CID (content identifier) from the Pinata service.
// If the cid parameter is an empty string, an error is returned.
// Returns an error if the file could not be deleted.
func (c *client) DeleteFile(cid string) error {
	if cid == "" {
		return fmt.Errorf("cid is required")
	}

	err := c.NewRequest(http.MethodDelete, "/pinning/unpin/{cid}").
		AddPathParam("cid", cid).
		Send(nil)

	if err != nil {
		return err
	}
	return nil
}
