package pinata

import (
	"fmt"
)

// PinataGroup represents a group in the Pinata platform.
// It contains information about the group, such as its ID, owner ID, name, creation time, and last update time.
type PinataGroup struct {
	ID        string `json:"id,omitempty"`
	OwnerID   string `json:"user_id,omitempty"`
	GroupName string `json:"name,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
}

// ListGroupsOptions represents the options for listing Pinata groups.
// The NameContains field filters the groups by name, the Limit field sets the maximum number of groups to return,
// and the Offset field sets the starting index for the returned groups.
type ListGroupsOptions struct {
	NameContains string `json:"nameContains,omitempty"`
	Limit        int    `json:"limit,omitempty"`
	Offset       int    `json:"offset,omitempty"`
}

// CreateGroup creates a new Pinata group with the specified name.
// It returns the newly created PinataGroup object, or an error if the creation failed.
// The group name is required and cannot be an empty string.
func (c *Client) CreateGroup(groupName string) (*PinataGroup, error) {
	if groupName == "" {
		return nil, fmt.Errorf("group name is required")
	}

	payload := make(map[string]string)
	payload["name"] = groupName

	req, err := c.NewRequest("POST", "/groups").SetJSONBody(payload)
	if err != nil {
		return nil, fmt.Errorf("ERR: failed to set JSON body: %w", err)
	}

	var response PinataGroup
	err = req.Send(&response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetGroup retrieves a Pinata group by its ID.
//
// If the provided groupID is empty, an error is returned.
// Otherwise, the function makes a GET request to the "/groups/{id}" endpoint
// and returns the corresponding PinataGroup struct, or an error if the request fails.
func (c *Client) GetGroup(groupID string) (*PinataGroup, error) {
	if groupID == "" {
		return nil, fmt.Errorf("group id is required")
	}

	var response PinataGroup
	err := c.NewRequest("GET", "/groups/{id}").
		AddPathParam("id", groupID).
		Send(&response)

	if err != nil {
		return nil, err
	}
	return &response, nil
}

// ListGroups retrieves a list of Pinata groups based on the provided options.
// If options is nil, the function will return all groups without any filtering or pagination.
// Otherwise, the function will apply the specified limit and offset to the list of groups.
func (c *Client) ListGroups(options *ListGroupsOptions) ([]PinataGroup, error) {
	req := c.NewRequest("GET", "/groups")
	if options != nil {
		req.addListGroupsQueryParams(options)
	}

	var response []PinataGroup
	err := req.Send(&response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// UpdateGroup updates the name of the Pinata group with the specified ID.
//
// If the provided groupID or newGroupName is empty, an error is returned.
// Otherwise, the function makes a PUT request to the "/groups/{id}" endpoint
// with the new group name in the request body, and returns the updated
// PinataGroup struct, or an error if the request fails.
func (c *Client) UpdateGroup(groupID, newGroupName string) (*PinataGroup, error) {
	if groupID == "" || newGroupName == "" {
		return nil, fmt.Errorf("group id and new group name are required")
	}

	payload := make(map[string]string)
	payload["name"] = newGroupName

	req, err := c.NewRequest("PUT", "/groups/{id}").
		AddPathParam("id", groupID).
		SetJSONBody(payload)
	if err != nil {
		return nil, fmt.Errorf("ERR: failed to set JSON body: %w", err)
	}

	var response PinataGroup
	err = req.Send(&response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// AddCidToGroup adds the specified CIDs to the group with the given ID.
// If the group ID or the list of CIDs is empty, an error is returned.
func (c *Client) AddCidToGroup(groupID string, cids []string) error {
	if groupID == "" || len(cids) == 0 {
		return fmt.Errorf("group id and at least one cid is required")
	}

	payload := make(map[string][]string)
	payload["cids"] = cids

	req, err := c.NewRequest("PUT", "/groups/{id}/cids").
		AddPathParam("id", groupID).
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

// RemoveCidFromGroup removes the specified CIDs from the group with the given ID.
// If the group ID or the list of CIDs is empty, an error is returned.
func (c *Client) RemoveCidFromGroup(groupID string, cids []string) error {
	if groupID == "" || len(cids) == 0 {
		return fmt.Errorf("group id and at least one cid is required")
	}

	payload := make(map[string][]string)
	payload["cids"] = cids

	req, err := c.NewRequest("DELETE", "/groups/{id}/cids").
		AddPathParam("id", groupID).
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

// RemoveGroup removes the group with the specified ID.
// If the group ID is empty, an error is returned.
func (c *Client) RemoveGroup(groupID string) error {
	if groupID == "" {
		return fmt.Errorf("group id is required")
	}

	err := c.NewRequest("DELETE", "/groups/{id}").
		AddPathParam("id", groupID).
		Send(nil)
		
	if err != nil {
		return err
	}
	return nil
}
