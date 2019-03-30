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
	for _, r := range data {
		// Loop each item in the raw JSON message, detect its type and unmarshal it.
		var obj map[string]interface{}
		var uap UAP
		var usg USG
		var usw USW

		if u.unmarshalDevice("interface{}", &obj, r) != nil {
			continue
		}
		assetType := "<missing>"
		if t, ok := obj["type"].(string); ok {
			assetType = t
		}
		u.dLogf("Unmarshalling Device Type: %v", assetType)
		switch assetType { // Unmarshal again into the correct type..
		case "uap":
			if u.unmarshalDevice(assetType, &uap, r) == nil {
				devices.UAPs = append(devices.UAPs, uap)
			}
		case "ugw", "usg": // in case they ever fix the name in the api.
			if u.unmarshalDevice(assetType, &usg, r) == nil {
				devices.USGs = append(devices.USGs, usg)
			}
		case "usw":
			if u.unmarshalDevice(assetType, &usw, r) == nil {
				devices.USWs = append(devices.USWs, usw)
			}
		default:
			u.eLogf("unknown asset type - %v - skipping", assetType)
			continue
		}
	}
	return devices
}

// unmarshalDevice handles logging for the unmarshal operations in parseDevices().
func (u *Unifi) unmarshalDevice(device string, ptr interface{}, payload json.RawMessage) error {
	err := json.Unmarshal(payload, ptr)
	if err != nil {
		u.eLogf("json.Unmarshal(%v): %v", device, err)
		u.eLogf("Enable Debug Logging to output the failed payload.")
		json, err := payload.MarshalJSON()
		u.dLogf("Failed Payload: %s (marshal err: %v)", json, err)
		u.dLogf("The above payload can prove useful during torubleshooting when you open an Issue:")
		u.dLogf("==- https://github.com/golift/unifi/issues/new -==")
	}
	return err
}
