package pinata

import (
	"fmt"
	"net/http"
	"time"
)

// swapData represents the data for a single swap, including the mapped CID and the creation timestamp.
type swapData struct {
	MappedCid string    `json:"mappedCid"`
	CreatedAt time.Time `json:"createdAt"`
}

type addSwapResponse struct {
	Data swapData `json:"data"`
}

type deleteSwapResponse struct {
	Data interface{} `json:"data"`
}

type getSwapResponse struct {
	Data []swapData `json:"data"`
}

// AddSwap adds a new swap for the given CID. The swapCid parameter represents the CID
// that will be mapped to the original CID. If either the cid or swapCid is empty,
// an error is returned.
func (c *Client) AddSwap(cid, swapCid string) (*addSwapResponse, error) {
	if cid == "" || swapCid == "" {
		return nil, fmt.Errorf("cid and swapcid are required")
	}

	payload := make(map[string]string)
	payload["swapCid"] = swapCid

	req, err := c.NewRequest(http.MethodPut, "/v3/ipfs/swap/{cid}").
		AddPathParam("cid", cid).
		SetJSONBody(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to set JSON body: %w", err)
	}

	var response addSwapResponse
	err = req.Send(&response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetSwapHistory retrieves the swap history for the given CID and domain.
// The CID and domain parameters are required.
// The function returns a getSwapResponse containing the swap history data, or an error if the request fails.
func (c *Client) GetSwapHistory(cid, domain string) (*getSwapResponse, error) {
	if cid == "" || domain == "" {
		return nil, fmt.Errorf("cid and domain are required")
	}

	var response getSwapResponse
	err := c.NewRequest(http.MethodDelete, "/v3/ipfs/swap/{cid}").
		AddPathParam("cid", cid).
		AddQueryParam("domain", domain).
		Send(&response)

	if err != nil {
		return nil, err
	}
	return &response, nil
}

// RemoveSwap removes the swap for the given CID. If the cid is empty, an error is returned.
func (c *Client) RemoveSwap(cid string) (*deleteSwapResponse, error) {
	if cid == "" {
		return nil, fmt.Errorf("cid is required")
	}

	var response deleteSwapResponse
	err := c.NewRequest(http.MethodDelete, "/v3/ipfs/swap/{cid}").
		AddPathParam("cid", cid).
		Send(&response)

	if err != nil {
		return nil, err
	}
	return &response, nil
}
