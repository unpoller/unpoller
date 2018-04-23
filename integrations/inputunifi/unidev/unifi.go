package unidev

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Debug ....
var Debug = false

const (
	// ClientPath is Unifi Clients API Path
	ClientPath = "/api/s/default/stat/sta"
	// DevicePath is where we get data about Unifi devices.
	DevicePath = "/api/s/default/stat/device"
	// NetworkPath contains network-configuration data. Not really graphable.
	NetworkPath = "/api/s/default/rest/networkconf"
	// UserGroupPath contains usergroup configurations.
	UserGroupPath = "/api/s/default/rest/usergroup"
)

// GetUnifiClients returns a response full of clients' data from the Unifi Controller.
func (c *AuthedReq) GetUnifiClients() ([]Asset, error) {
	var response struct {
		Clients []UCL `json:"data"`
		Meta    struct {
			Rc string `json:"rc"`
		} `json:"meta"`
	}
	if req, err := c.UniReq(ClientPath, ""); err != nil {
		return nil, err
	} else if resp, err := c.Do(req); err != nil {
		return nil, err
	} else if body, err := ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
	} else if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	} else if err = resp.Body.Close(); err != nil {
		log.Println("resp.Body.Close():", err) // Not fatal? Just log it.
	}
	clients := []Asset{}
	for _, r := range response.Clients {
		clients = append(clients, r)
	}
	return clients, nil
}

// GetUnifiDevices returns a response full of devices' data from the Unifi Controller.
func (c *AuthedReq) GetUnifiDevices() ([]Asset, error) {
	var parsed struct {
		Data []json.RawMessage `json:"data"`
		Meta struct {
			Rc string `json:"rc"`
		} `json:"meta"`
	}
	assets := []Asset{}
	if req, err := c.UniReq(DevicePath, ""); err != nil {
		return nil, err
	} else if resp, err := c.Do(req); err != nil {
		return nil, err
	} else if body, err := ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
	} else if err = json.Unmarshal(body, &parsed); err != nil {
		return nil, err
	} else if err = resp.Body.Close(); err != nil {
		log.Println("resp.Body.Close():", err) // Not fatal? Just log it.
	}
	for _, r := range parsed.Data {
		// Unamrshal into a map and check "type"
		var obj map[string]interface{}
		if err := json.Unmarshal(r, &obj); err != nil {
			return nil, err
		}
		assetType := ""
		if t, ok := obj["type"].(string); ok {
			assetType = t
		}
		// Unmarshal again into the correct type..
		var asset Asset
		switch assetType {
		case "uap":
			asset = &UAP{}
		case "ugw":
			asset = &USG{}
		case "usw":
			asset = &USW{}
		default:
			log.Println("unknown asset type -", assetType, "- skipping")
			continue
		}
		if Debug {
			log.Println("Unmarshalling", assetType)
		}
		if err := json.Unmarshal(r, asset); err != nil {
			return nil, err
		}
		assets = append(assets, asset)
	}
	return assets, nil
}
