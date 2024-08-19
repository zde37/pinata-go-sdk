package pinata

import (
	"fmt"
	"net/http"
	"time"
)

// apiKeyResponse represents the response from an API key related request.
// It contains a slice of ApiKey structs and a count of the total number of keys.
type apiKeyResponse struct {
	Keys  []apiKey `json:"keys,omitempty"`
	Count int      `json:"count,omitempty"`
}

// apiKey represents an API key for the Pinata service.
type apiKey struct {
	ID        string    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Key       string    `json:"key,omitempty"`
	Secret    string    `json:"secret,omitempty"`
	MaxUses   int       `json:"max_uses,omitempty"`
	Uses      int       `json:"uses,omitempty"`
	UserID    string    `json:"user_id,omitempty"`
	Scopes    scope     `json:"scopes,omitempty"`
	Revoked   bool      `json:"revoked,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

// scope represents the permissions and access scopes for an API key.
// The Endpoints field contains the specific permissions for different API endpoints,
// while the Admin field indicates if the API key has full administrative access.
type scope struct {
	Endpoints struct {
		Data    scopeData    `json:"data"`
		Pinning scopePinning `json:"pinning"`
		Psa     scopePsa     `json:"psa"`
	} `json:"endpoints"`
	Admin bool `json:"admin,omitempty"`
}

// scopeData represents the data-related permissions for an API key.
// The PinList field indicates if the API key can list the user's pinned data.
// The UserPinnedDataTotal field indicates if the API key can retrieve the total amount of the user's pinned data.
type scopeData struct {
	PinList             bool `json:"pinList"`
	UserPinnedDataTotal bool `json:"userPinnedDataTotal"`
}

// scopePsa represents the permissions and access scopes for the PSA (Pinata Secure API) endpoints.
// The Pins field contains the specific permissions for different PSA pin-related operations.
type scopePsa struct {
	Pins struct {
		AddPinObject     bool `json:"addPinObject"`
		GetPinObject     bool `json:"getPinObject"`
		ListPinObjects   bool `json:"listPinObjects"`
		RemovePinObject  bool `json:"removePinObject"`
		ReplacePinObject bool `json:"replacePinObject"`
	} `json:"pins"`
}

// scopePinning represents the permissions and access scopes for the Pinning API endpoints.
// The HashMetadata field indicates if the API key can hash metadata.
// The HashPinPolicy field indicates if the API key can hash pin policies.
// The PinByHash field indicates if the API key can pin data by hash.
// The PinFileToIPFS field indicates if the API key can pin files to IPFS.
// The PinJSONToIPFS field indicates if the API key can pin JSON data to IPFS.
// The PinJobs field indicates if the API key can manage pin jobs.
// The Unpin field indicates if the API key can unpin data.
// The UserPinPolicy field indicates if the API key can manage user pin policies.
type scopePinning struct {
	HashMetadata  bool `json:"hashMetadata"`
	HashPinPolicy bool `json:"hashPinPolicy"`
	PinByHash     bool `json:"pinByHash"`
	PinFileToIPFS bool `json:"pinFileToIPFS"`
	PinJSONToIPFS bool `json:"pinJSONToIPFS"`
	PinJobs       bool `json:"pinJobs"`
	Unpin         bool `json:"unpin"`
	UserPinPolicy bool `json:"userPinPolicy"`
}

// GenerateApiKeyOptions represents the options for generating a new API key.
// KeyName is the name of the new API key.
// Permissions specifies the permissions and access scopes for the new API key.
// MaxUses specifies the maximum number of times the API key can be used.
type GenerateApiKeyOptions struct {
	KeyName     string      `json:"keyName,omitempty"`
	Permissions Permissions `json:"permissions,omitempty"`
	MaxUses     int         `json:"maxUses,omitempty"`
}

// secret represents the secret information returned when generating a new API key.
// The JWT field contains the JSON Web Token for the generated API key.
// The PinataApiKey field contains the API key itself.
// The PinataApiSecret field contains the API secret for the generated API key.
type secret struct {
	JWT             string `json:"JWT,omitempty"`
	PinataApiKey    string `json:"pinata_api_key,omitempty"`
	PinataApiSecret string `json:"pinata_api_secret,omitempty"`
}

// Permissions represents the permissions and access scopes for an API key.
// The Admin field indicates if the API key has administrative permissions.
// The Endpoints field specifies the permissions for different API endpoints.
type Permissions struct {
	Admin     bool      `json:"admin,omitempty"`
	Endpoints *EndPoint `json:"endpoints,omitempty"` // used a pointer for backward compatibility (generateApiKey endpoint)
}

// EndPoint represents the permissions and access scopes for different API endpoints.
// The Data field specifies the permissions for data-related operations.
// The Pinning field specifies the permissions for pinning-related operations.
type EndPoint struct {
	Data    Data    `json:"data,omitempty"`
	Pinning Pinning `json:"pinning,omitempty"`
}

// Data represents the permissions and access scopes for the "data" API endpoint.
// The PinList field indicates if the API key can list pinned data.
// The UserPinnedDataTotal field indicates if the API key can retrieve the total amount of user-pinned data.
type Data struct {
	PinList             bool `json:"pinList,omitempty"`
	UserPinnedDataTotal bool `json:"userPinnedDataTotal,omitempty"`
}

// Pinning represents the permissions and access scopes for pinning-related operations.
// HashMetadata indicates if the API key can access hash metadata.
// HashPinPolicy indicates if the API key can access the hash pin policy.
// PinByHash indicates if the API key can pin data by hash.
// PinFileToIPFS indicates if the API key can pin files to IPFS.
// PinJSONToIPFS indicates if the API key can pin JSON data to IPFS.
// PinJobs indicates if the API key can access pin jobs.
// UnPin indicates if the API key can unpin data.
// UserPinPolicy indicates if the API key can access the user's pin policy.
type Pinning struct {
	HashMetadata  bool `json:"hashMetadata,omitempty"`
	HashPinPolicy bool `json:"hashPinPolicy,omitempty"`
	PinByHash     bool `json:"pinByHash,omitempty"`
	PinFileToIPFS bool `json:"pinFileToIPFS,omitempty"`
	PinJSONToIPFS bool `json:"pinJSONToIPFS,omitempty"`
	PinJobs       bool `json:"pinJobs,omitempty"`
	UnPin         bool `json:"unpin,omitempty"`
	UserPinPolicy bool `json:"userPinPolicy,omitempty"`
}

// ListApiKeysOptions represents the options for listing API keys.
// Revoked indicates whether to include revoked API keys in the response.
// LimitedUse indicates whether to include API keys with limited use in the response.
// Exhausted indicates whether to include exhausted API keys in the response.
// Name is a filter to only include API keys with the specified name.
// Offset is the number of API keys to skip before returning the results.
type ListApiKeysOptions struct {
	Revoked    *bool  `json:"revoked,omitempty"`
	LimitedUse *bool  `json:"limitedUse,omitempty"`
	Exhausted  *bool  `json:"exhausted,omitempty"`
	Name       string `json:"name,omitempty"`
	Offset     int    `json:"offset,omitempty"`
}

// GenerateApiKey generates a new API key for the Pinata platform.
//
// The provided GenerateApiKeyOptions struct specifies the options for the new API key, such as the name, permissions, and expiration.
// If the options are nil, an error will be returned.
//
// The function returns a Secret struct containing the new API key and secret.
// If there is an error generating the API key, an error will be returned.
func (c *Client) GenerateApiKey(options *GenerateApiKeyOptions) (*secret, error) {
	if options == nil {
		return nil, fmt.Errorf("options cannot be nil")
	}

	req, err := c.NewRequest(http.MethodPost, "/users/generateApiKey").
		SetJSONBody(options)

	if err != nil {
		return nil, fmt.Errorf("ERR: failed to set JSON body: %w", err)
	}

	var response secret
	err = req.Send(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GenerateApiKeyV3 generates a new API key for the Pinata platform.
//
// The provided GenerateApiKeyOptions struct specifies the options for the new API key, such as the name, permissions, and expiration.
// If the options are nil, an error will be returned.
//
// The function returns a Secret struct containing the new API key and secret.
// If there is an error generating the API key, an error will be returned.
func (c *Client) GenerateApiKeyV3(options *GenerateApiKeyOptions) (*secret, error) {
	if options == nil {
		return nil, fmt.Errorf("options cannot be nil")
	}

	req, err := c.NewRequest(http.MethodPost, "/v3/pinata/keys").
		SetJSONBody(options)

	if err != nil {
		return nil, fmt.Errorf("ERR: failed to set JSON body: %w", err)
	}

	var response secret
	err = req.Send(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// ListApiKeys returns a list of API keys associated with the current user.
// The response includes information about each API key, such as whether it is revoked, limited use, or exhausted.
// The options parameter can be used to filter the results by various criteria.
func (c *Client) ListApiKeys() (*apiKeyResponse, error) {
	var response apiKeyResponse
	err := c.NewRequest(http.MethodGet, "/users/apiKeys").
		Send(&response)

	if err != nil {
		return nil, err
	}
	return &response, nil
}

// ListApiKeyV3 returns a list of API keys associated with the current user.
// The response includes information about each API key, such as whether it is revoked, limited use, or exhausted.
// The options parameter can be used to filter the results by various criteria.
func (c *Client) ListApiKeyV3(options *ListApiKeysOptions) (*apiKeyResponse, error) {
	req := c.NewRequest(http.MethodGet, "/v3/pinata/keys")
	if options != nil {
		req.setListApiKeysQueryParams(options)
	}

	var response apiKeyResponse
	err := req.Send(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// RevokeApiKey revokes the specified API key.
// If the apiKey parameter is empty, an error is returned.
func (c *Client) RevokeApiKey(apiKey string) error {
	if apiKey == "" {
		return fmt.Errorf("api key is required")
	}

	payload := make(map[string]string)
	payload["apiKey"] = apiKey

	req, err := c.NewRequest(http.MethodPut, "/users/revokeApiKey").
		SetJSONBody(payload)
	if err != nil {
		return fmt.Errorf("ERR: failed to set JSON body: %w", err)
	}

	err = req.Send(nil)
	if err != nil {
		return err
	}
	return nil
}

// RevokeApiKeyV3 revokes the specified API key.
// The key parameter is required and must be a valid API key.
// If the key is successfully revoked, this method returns nil. Otherwise, it returns an error.
func (c *Client) RevokeApiKeyV3(key string) error {
	if key == "" {
		return fmt.Errorf("key is required")
	}

	err := c.NewRequest(http.MethodPut, "/v3/pinata/keys/{key}").
		AddPathParam("key", key).
		Send(nil)

	if err != nil {
		return err
	}
	return nil
}
