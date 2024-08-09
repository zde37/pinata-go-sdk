package pinata

import "fmt"

// CidSignature represents the response from the Pinata API for a CID signature.
type CidSignature struct {
	Data struct {
		Cid       string `json:"cid,omitempty"`
		Signature string `json:"signature,omitempty"`
	} `json:"data,omitempty"`
}

// AddCidSignature adds a signature for the given CID. If either the CID or the
// signature is empty, an error is returned.
func (c *Client) AddCidSignature(cid, signature string) (*CidSignature, error) {
	if cid == "" || signature == "" {
		return nil, fmt.Errorf("cid and signature is required")
	}

	payload := make(map[string]string)
	payload["signature"] = signature

	req, err := c.NewRequest("POST", "/v3/ipfs/signature/{cid}").
		AddPathParam("cid", cid).
		SetJSONBody(payload)
	if err != nil {
		return nil, fmt.Errorf("ERR: failed to set JSON body: %w", err)
	}

	var response CidSignature
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
func (c *Client) GetCidSignature(cid string) (*CidSignature, error) {
	if cid == "" {
		return nil, fmt.Errorf("cid is required")
	}

	var response CidSignature
	err := c.NewRequest("GET", "/v3/ipfs/signature/{cid}").
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
func (c *Client) RemoveCidSignature(cid string) error {
	if cid == "" {
		return fmt.Errorf("cid is required")
	}

	err := c.NewRequest("DELETE", "/v3/ipfs/signature/{cid}").
		AddPathParam("cid", cid).
		Send(nil)

	if err != nil {
		return err
	}
	return nil
}
