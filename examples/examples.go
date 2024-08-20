package main

import (
	"log"

	"github.com/zde37/pinata-go-sdk/pinata"
)

type PinataClient struct {
	client *pinata.Client
}

func NewPinataClient(auth *pinata.Auth) *PinataClient {
	client := pinata.New(auth)
	return &PinataClient{
		client: client,
	}
}

func (p *PinataClient) TestAuthentication() error {
	res, err := p.client.TestAuthentication()
	if err != nil {
		return err
	}

	log.Println(res.Message)
	return nil
}

func (p *PinataClient) PinFile(file string) error {
	options := &pinata.PinOptions{
		PinataMetadata: pinata.PinataMetadata{
			Name: "hi.txt",
			KeyValues: map[string]interface{}{
				"category": "example",
				"version":  2,
			},
		},
		PinataOptions: pinata.Options{
			CidVersion: 1,
		},
	}
	response, err := p.client.PinFileToIPFS(file, options)
	if err != nil {
		return err
	}

	log.Printf("file pinned successfully. Details: %+v\n", response)
	return nil
}

func (p *PinataClient) PinJSON(jsonData map[string]interface{}) error {
	options := &pinata.PinOptions{
		PinataMetadata: pinata.PinataMetadata{
			Name: "important-docs.json",
			KeyValues: map[string]interface{}{
				"category": "example",
				"version":  2,
			},
		},
		PinataOptions: pinata.Options{
			CidVersion: 1,
		},
	}

	response, err := p.client.PinJSONToIPFS(jsonData, options)
	if err != nil {
		return err
	}

	log.Printf("json pinned successfully. Details: %+v\n", response)
	return nil
}

func (p *PinataClient) GetFile(cid string) error {
	options := &pinata.ListFilesOptions{
		Cid:          cid,
		IncludeCount: true,
	}
	response, err := p.client.ListFiles(options)
	if err != nil {
		return err
	}

	log.Printf("pin info: %+v", response.Rows[0])
	return nil
}

func (p *PinataClient) ListFiles() error {
	options := &pinata.ListFilesOptions{
		IncludeCount: true,
		Status:       "pinned",
	}
	response, err := p.client.ListFiles(options)
	if err != nil {
		return err
	}

	log.Printf("total pins: %d\n", response.Count)
	for _, pin := range response.Rows {
		log.Printf("pin info: %+v\n\n", pin)
	}
	return nil
}

func (p *PinataClient) UpdateFileMetadata(cid string) error {
	options := &pinata.PinMetadataUpdateOptions{
		Name: "1",
		KeyValues: map[string]interface{}{
			"category": "1",
			"version":  1,
		},
	}

	err := p.client.UpdateFileMetadata(cid, options)
	if err != nil {
		return err
	}

	log.Println("file metadata updated successfully")
	return nil
}

func (p *PinataClient) DeleteFile(cid string) error {
	err := p.client.DeleteFile(cid)
	if err != nil {
		return err
	}

	log.Println("file deleted successfully")
	return nil
}

func (p *PinataClient) PinByCid(cid string) error {
	options := &pinata.PinByCidOptions{
		PinataMetadata: pinata.PinataMetadata{
			Name: "hi.txt",
			KeyValues: map[string]interface{}{
				"category": "example",
				"version":  2,
			},
		},
		PinataOptions: pinata.PinOpts{
			HostNodes: []string{
				"/ip4/172.22.33.3/tcp/4001/p2p/12D3KooWKyePX78pS5dtxkEubRDd7iyB3ihkUHsdLXLxJRAAAZu8",
				"/ip4/172.22.33.3/udp/4001/quic-v1/p2p/12D3KooWKyePX78pS5dtxkEubRDd7iyB3ihkUHsdLXLxJRAAAZu8",
			},
		},
	}
	response, err := p.client.PinByCid(cid, options)
	if err != nil {
		return err
	}

	log.Printf("pin info: %+v\n", response)
	return nil
}

func (p *PinataClient) ListPinByCidJobs(cid string) error {
	options := &pinata.ListPinByCidOptions{
		Sort:        pinata.SortOrderASC,
		Status:      pinata.PinStatusPrechecking,
		IPFSPinHash: cid,
	}

	response, err := p.client.ListPinByCidJobs(options)
	if err != nil {
		return err
	}

	log.Printf("total pins: %d\n", response.Count)
	for _, pin := range response.Rows {
		log.Printf("pin info: %+v\n\n", pin)
	}
	return nil
}

func (p *PinataClient) CreateGroup(name string) error {
	response, err := p.client.CreateGroup(name)
	if err != nil {
		return err
	}

	log.Printf("group created successfully. Details: %+v\n", response)
	return nil
}

func (p *PinataClient) GetGroup(groupName string) error {
	response, err := p.client.GetGroup(groupName)
	if err != nil {
		return err
	}

	log.Printf("group info: %+v\n", response)
	return nil
}

func (p *PinataClient) ListGroups() error {
	options := &pinata.ListGroupsOptions{}

	response, err := p.client.ListGroups(options)
	if err != nil {
		return err
	}

	for _, group := range response {
		log.Printf("group info: %+v\n\n", group)
	}
	return nil
}

func (p *PinataClient) UpdateGroupName(groupID, newName string) error {
	response, err := p.client.UpdateGroup(groupID, newName)
	if err != nil {
		return err
	}

	log.Printf("group info: %+v\n", response)
	return nil
}

func (p *PinataClient) AddCidToGroup(groupID string, cids []string) error {
	err := p.client.AddCidToGroup(groupID, cids)
	if err != nil {
		return err
	}

	log.Println("cids added to group successfully")
	return nil
}

func (p *PinataClient) RemoveCidFromGroup(groupID string, cids []string) error {
	err := p.client.RemoveCidFromGroup(groupID, cids)
	if err != nil {
		return err
	}

	log.Println("cids removed from group successfully")
	return nil
}

func (p *PinataClient) RemoveGroup(groupID string) error {
	err := p.client.RemoveGroup(groupID)
	if err != nil {
		return err
	}

	log.Println("group removed successfully")
	return nil
}

func (p *PinataClient) AddCidSignature(cid, signature string) error {
	response, err := p.client.AddCidSignature(cid, signature)
	if err != nil {
		return err
	}

	log.Printf("cid signature added successfully. Details: %+v\n", response)
	return nil
}

func (p *PinataClient) GetCidSignature(cid string) error {
	response, err := p.client.GetCidSignature(cid)
	if err != nil {
		return err
	}

	log.Printf("cid signature info: %+v\n", response)
	return nil
}

func (p *PinataClient) RemoveCidSignature(cid string) error {
	err := p.client.RemoveCidSignature(cid)
	if err != nil {
		return err
	}

	log.Println("cid signature removed successfully")
	return nil
}

func (p *PinataClient) CreateAPIKey(name string) error {
	options := &pinata.GenerateApiKeyOptions{
		KeyName: name,
		Permissions: pinata.Permissions{
			Admin: true,
		},
		MaxUses: 100,
	}

	response, err := p.client.GenerateApiKey(options)
	if err != nil {
		return err
	}

	log.Printf("api key created successfully. Details: %+v\n", response)
	return nil
}

func (p *PinataClient) CreateAPIKeyV3(name string) error {
	options := &pinata.GenerateApiKeyOptions{
		KeyName: name,
		Permissions: pinata.Permissions{
			Admin: true,
		},
		MaxUses: 100,
	}

	res, err := p.client.GenerateApiKeyV3(options)
	if err != nil {
		return err
	}

	log.Printf("api key created successfully. Details: %+v\n", res)
	return nil
}

func (p *PinataClient) ListAPIKeys() error {
	response, err := p.client.ListApiKeys()
	if err != nil {
		return err
	}

	log.Printf("api keys: %+v\n", response)
	return nil
}

func (p *PinataClient) ListApiKeyV3() error {
	revoked := false
	options := &pinata.ListApiKeysOptions{
		Revoked: &revoked,
	}

	res, err := p.client.ListApiKeyV3(options)
	if err != nil {
		return err
	}

	log.Printf("api keys: %+v\n", res)
	return nil
}

func (p *PinataClient) RevokeApiKey(apiKey string) error {
	err := p.client.RevokeApiKey(apiKey)
	if err != nil {
		return err
	}

	log.Println("api key revoked successfully")
	return nil
}

func (p *PinataClient) RevokeApiKeyV3(apiKey string) error {
	err := p.client.RevokeApiKeyV3(apiKey)
	if err != nil {
		return err
	}

	log.Println("api key revoked successfully")
	return nil
}
