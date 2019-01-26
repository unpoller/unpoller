package unifi

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

// GetClients returns a response full of clients' data from the Unifi Controller.
func (u *Unifi) GetClients() (*Clients, error) {
	var response struct {
		Clients []UCL `json:"data"`
	}
	req, err := u.UniReq(ClientPath, "")
	if err != nil {
		return nil, errors.Wrap(err, "c.UniReq(ClientPath)")
	}
	resp, err := u.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "c.Do(req)")
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		return nil, errors.Wrap(err, "ioutil.ReadAll(resp.Body)")
	} else if err = json.Unmarshal(body, &response); err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal([]UCL)")
	}
	return &Clients{UCLs: response.Clients}, nil
}

// GetDevices returns a response full of devices' data from the Unifi Controller.
func (u *Unifi) GetDevices() (*Devices, error) {
	var response struct {
		Data []json.RawMessage `json:"data"`
	}
	req, err := u.UniReq(DevicePath, "")
	if err != nil {
		return nil, errors.Wrap(err, "c.UniReq(DevicePath)")
	}
	resp, err := u.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "c.Do(req)")
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		return nil, errors.Wrap(err, "ioutil.ReadAll(resp.Body)")
	} else if err = json.Unmarshal(body, &response); err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal([]json.RawMessage)")
	}
	return u.parseDevices(response.Data), nil
}

// parseDevices parses the raw JSON from the Unifi Controller into device structures.
func (u *Unifi) parseDevices(data []json.RawMessage) *Devices {
	devices := new(Devices)
	for i, r := range data {
		// Loop each item in the raw JSON message, detect its type and unmarshal it.
		var obj map[string]interface{}
		if err := json.Unmarshal(r, &obj); err != nil {
			u.eLogf("%d: json.Unmarshal(interfce{}): %v", i, err)
			continue
		}
		assetType := "<missing>"
		if t, ok := obj["type"].(string); ok {
			assetType = t
		}
		u.dLogf("Unmarshalling Device Type: %v", assetType)
		// Unmarshal again into the correct type..
		switch assetType {
		case "uap":
			var uap UAP
			if err := json.Unmarshal(r, &uap); err != nil {
				u.eLogf("%d: json.Unmarshal([]UAP): %v", i, err)
				continue
			}
			devices.UAPs = append(devices.UAPs, uap)
		case "ugw", "usg": // in case they ever fix the name in the api.
			var usg USG
			if err := json.Unmarshal(r, &usg); err != nil {
				u.eLogf("%d: json.Unmarshal([]USG): %v", i, err)
				continue
			}
			devices.USGs = append(devices.USGs, usg)
		case "usw":
			var usw USW
			if err := json.Unmarshal(r, &usw); err != nil {
				u.eLogf("%d: json.Unmarshal([]USW): %v", i, err)
				continue
			}
			devices.USWs = append(devices.USWs, usw)
		default:
			u.dLogf("unknown asset type - " + assetType + " - skipping")
			continue
		}
	}
	return devices
}
