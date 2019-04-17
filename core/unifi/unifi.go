package unifi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
)

// GetClients returns a response full of clients' data from the Unifi Controller.
func (u *Unifi) GetClients(sites []Site) (*Clients, error) {
	var data []UCL
	for _, site := range sites {
		var response struct {
			Data []UCL `json:"data"`
		}
		u.dLogf("Polling Site '%s' (%s) Clients", site.Name, site.Desc)
		u.dLogf("Unmarshalling Device Type: ucl")
		clientPath := fmt.Sprintf(ClientPath, site.Name)
		if err := u.GetData(clientPath, &response); err != nil {
			return nil, err
		}
		data = append(data, response.Data...)
	}
	return &Clients{UCLs: data}, nil
}

// GetDevices returns a response full of devices' data from the Unifi Controller.
func (u *Unifi) GetDevices(sites []Site) (*Devices, error) {
	var data []json.RawMessage
	for _, site := range sites {
		u.dLogf("Polling Site '%s' (%s) Devices", site.Name, site.Desc)
		var response struct {
			Data []json.RawMessage `json:"data"`
		}
		devicePath := fmt.Sprintf(DevicePath, site.Name)
		if err := u.GetData(devicePath, &response); err != nil {
			return nil, err
		}
		data = append(data, response.Data...)
	}
	return u.parseDevices(data), nil
}

// Site represents a site's data. There are more pieces to this.
type Site struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

// GetSites returns a list of configured sites on the Unifi controller.
func (u *Unifi) GetSites() ([]Site, error) {
	var response struct {
		Data []Site `json:"data"`
	}
	if err := u.GetData(SiteList, &response); err != nil {
		return nil, err
	}
	var sites []string
	for i := range response.Data {
		sites = append(sites, response.Data[i].Name)
	}
	u.dLogf("Found %d sites: %s", len(sites), strings.Join(sites, ","))
	return response.Data, nil
}

// GetData makes a unifi request and unmarshal the response into a provided pointer.
func (u *Unifi) GetData(methodPath string, v interface{}) error {
	req, err := u.UniReq(methodPath, "")
	if err != nil {
		return errors.Wrapf(err, "c.UniReq(%s)", methodPath)
	}
	resp, err := u.Do(req)
	if err != nil {
		return errors.Wrapf(err, "c.Do(%s)", methodPath)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		return errors.Wrapf(err, "ioutil.ReadAll(%s)", methodPath)
	} else if err = json.Unmarshal(body, v); err != nil {
		return errors.Wrapf(err, "json.Unmarshal(%s)", methodPath)
	}
	return nil
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
