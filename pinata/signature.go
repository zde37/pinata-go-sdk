package pinata

import (
	"fmt"
	"net/http"
)

// cidSignature represents the response from the Pinata API for a CID signature.
type cidSignature struct {
	Data struct {
		Cid       string `json:"cid,omitempty"`
		Signature string `json:"signature,omitempty"`
	} `json:"data,omitempty"`
}

// AddCidSignature adds a signature for the given CID. If either the CID or the
// signature is empty, an error is returned.
func (c *client) AddCidSignature(cid, signature string) (*cidSignature, error) {
	if cid == "" || signature == "" {
		return nil, fmt.Errorf("cid and signature is required")
	}

	payload := make(map[string]string)
	payload["signature"] = signature

	req, err := c.NewRequest(http.MethodPost, "/v3/ipfs/signature/{cid}").
		AddPathParam("cid", cid).
		SetJSONBody(payload)
	if err != nil {
		return nil, fmt.Errorf("ERR: failed to set JSON body: %w", err)
	}

	var response cidSignature
	err = req.Send(&response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetCidSignature retrieves the signature for the given CID from the Pinata API.
// If the CID is empty, an error is returned.
// The CidSignature struct is returned, which contains the CID and its signature.
// If an error occurs during the API request, the error is returned.
func (c *client) GetCidSignature(cid string) (*cidSignature, error) {
	if cid == "" {
		return nil, fmt.Errorf("cid is required")
	}

	var response cidSignature
	err := c.NewRequest(http.MethodGet, "/v3/ipfs/signature/{cid}").
		AddPathParam("cid", cid).
		Send(&response)

	if err != nil {
		return nil, err
	}
	return &response, nil
}

// RemoveCidSignature removes the signature for the given CID from the Pinata API.
// If the CID is empty, an error is returned.
// If an error occurs during the API request, the error is returned.
func (c *client) RemoveCidSignature(cid string) error {
	if cid == "" {
		return fmt.Errorf("cid is required")
	}

	err := c.NewRequest(http.MethodDelete, "/v3/ipfs/signature/{cid}").
		AddPathParam("cid", cid).
		Send(nil)

	if err != nil {
		return err
	}
	return nil
}
