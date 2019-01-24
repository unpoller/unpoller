package unidev

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/pkg/errors"
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
func (c *AuthedReq) GetUnifiClients() ([]UCL, error) {
	var response struct {
		Clients []UCL `json:"data"`
		Meta    struct {
			Rc string `json:"rc"`
		} `json:"meta"`
	}
	req, err := c.UniReq(ClientPath, "")
	if err != nil {
		return nil, errors.Wrap(err, "c.UniReq(ClientPath)")
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "c.Do(req)")
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("resp.Body.Close():", err) // Not fatal? Just log it.
		}
	}()
	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		return nil, errors.Wrap(err, "ioutil.ReadAll(resp.Body)")
	} else if err = json.Unmarshal(body, &response); err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal([]UCL)")
	}
	return response.Clients, nil
}

// GetUnifiClientAssets provides an interface to return common asset types.
func (c *AuthedReq) GetUnifiClientAssets() ([]Asset, error) {
	clients, err := c.GetUnifiClients()
	assets := []Asset{}
	if err == nil {
		for _, r := range clients {
			assets = append(assets, r)
		}
	}
	return assets, err
}

// GetUnifiDevices returns a response full of devices' data from the Unifi Controller.
func (c *AuthedReq) GetUnifiDevices() ([]USG, []USW, []UAP, error) {
	var parsed struct {
		Data []json.RawMessage `json:"data"`
		Meta struct {
			Rc string `json:"rc"`
		} `json:"meta"`
	}
	req, err := c.UniReq(DevicePath, "")
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "c.UniReq(DevicePath)")
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "c.Do(req)")
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("resp.Body.Close():", err) // Not fatal? Just log it.
		}
	}()
	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		return nil, nil, nil, errors.Wrap(err, "ioutil.ReadAll(resp.Body)")
	} else if err = json.Unmarshal(body, &parsed); err != nil {
		return nil, nil, nil, errors.Wrap(err, "json.Unmarshal([]json.RawMessage)")
	}

	var usgs []USG
	var usws []USW
	var uaps []UAP
	// Loop each item in the raw JSON message, detect its type and unmarshal it.
	for i, r := range parsed.Data {
		var usg USG
		var usw USW
		var uap UAP
		// Unamrshal into a map and check "type"
		var obj map[string]interface{}
		if err := json.Unmarshal(r, &obj); err != nil {
			return nil, nil, nil, errors.Wrapf(err, "[%d] json.Unmarshal(interfce{})", i)
		}
		assetType := "<missing>"
		if t, ok := obj["type"].(string); ok {
			assetType = t
		}
		if Debug {
			log.Println("Unmarshalling Device Type:", assetType)
		}
		// Unmarshal again into the correct type..
		switch assetType {
		case "uap":
			if err := json.Unmarshal(r, &uap); err != nil {
				return nil, nil, nil, errors.Wrapf(err, "[%d] json.Unmarshal([]UAP)", i)
			}
			uaps = append(uaps, uap)
		case "ugw", "usg": // in case they ever fix the name in the api.
			if err := json.Unmarshal(r, &usg); err != nil {
				return nil, nil, nil, errors.Wrapf(err, "[%d] json.Unmarshal([]USG)", i)
			}
			usgs = append(usgs, usg)
		case "usw":
			if err := json.Unmarshal(r, &usw); err != nil {
				return nil, nil, nil, errors.Wrapf(err, "[%d] json.Unmarshal([]USW)", i)
			}
			usws = append(usws, usw)
		default:
			log.Println("unknown asset type -", assetType, "- skipping")
			continue
		}
	}
	return usgs, usws, uaps, nil
}

// GetUnifiDeviceAssets provides an interface to return common asset types.
func (c *AuthedReq) GetUnifiDeviceAssets() ([]Asset, error) {
	usgs, usws, uaps, err := c.GetUnifiDevices()
	assets := []Asset{}
	if err == nil {
		for _, r := range usgs {
			assets = append(assets, r)
		}
		for _, r := range usws {
			assets = append(assets, r)
		}
		for _, r := range uaps {
			assets = append(assets, r)
		}
	}
	return assets, err
}
