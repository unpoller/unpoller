package unifi

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

// GetUnifiClients returns a response full of clients' data from the Unifi Controller.
func (u *Unifi) GetUnifiClients() ([]UCL, error) {
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
	return response.Clients, nil
}

// GetUnifiDevices returns a response full of devices' data from the Unifi Controller.
func (u *Unifi) GetUnifiDevices() (*Devices, error) {
	var parsed struct {
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
	} else if err = json.Unmarshal(body, &parsed); err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal([]json.RawMessage)")
	}
	return u.parseUnifiDevices(parsed.Data), nil
}

// parseUnifiDevices parses the raw JSON from the Unifi Controller into device structures.
func (u *Unifi) parseUnifiDevices(data []json.RawMessage) *Devices {
	devices := new(Devices)
	// Loop each item in the raw JSON message, detect its type and unmarshal it.
	for i, r := range data {
		var usg USG
		var usw USW
		var uap UAP
		// Unamrshal into a map and check "type"
		var obj map[string]interface{}
		if err := json.Unmarshal(r, &obj); err != nil {
			u.eLogf("%d: json.Unmarshal(interfce{}): %v", i, err)
			continue
		}
		assetType := "<missing>"
		if t, ok := obj["type"].(string); ok {
			assetType = t
		}
		u.dLogf("Unmarshalling Device Type:", assetType)
		// Unmarshal again into the correct type..
		switch assetType {
		case "uap":
			if err := json.Unmarshal(r, &uap); err != nil {
				u.eLogf("%d: json.Unmarshal([]UAP): %v", i, err)
				continue
			}
			devices.UAPs = append(devices.UAPs, uap)
		case "ugw", "usg": // in case they ever fix the name in the api.
			if err := json.Unmarshal(r, &usg); err != nil {
				u.eLogf("%d: json.Unmarshal([]USG): %v", i, err)
				continue
			}
			devices.USGs = append(devices.USGs, usg)
		case "usw":
			if err := json.Unmarshal(r, &usw); err != nil {
				u.eLogf("%d: json.Unmarshal([]USW): %v", i, err)
				continue
			}
			devices.USWs = append(devices.USWs, usw)
		default:
			u.dLogf("unknown asset type -", assetType, "- skipping")
			continue
		}
	}
	return devices
}
