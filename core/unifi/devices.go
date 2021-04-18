package unifi

import (
	"encoding/json"
	"fmt"
	"strings"
)

// GetDevices returns a response full of devices' data from the UniFi Controller.
func (u *Unifi) GetDevices(sites []*Site) (*Devices, error) {
	devices := new(Devices)

	for _, site := range sites {
		var response struct {
			Data []json.RawMessage `json:"data"`
		}

		devicePath := fmt.Sprintf(APIDevicePath, site.Name)
		if err := u.GetData(devicePath, &response); err != nil {
			return nil, err
		}

		loopDevices := u.parseDevices(response.Data, site)
		devices.UAPs = append(devices.UAPs, loopDevices.UAPs...)
		devices.USGs = append(devices.USGs, loopDevices.USGs...)
		devices.USWs = append(devices.USWs, loopDevices.USWs...)
		devices.UDMs = append(devices.UDMs, loopDevices.UDMs...)
		devices.UXGs = append(devices.UXGs, loopDevices.UXGs...)
	}

	return devices, nil
}

// GetUSWs returns all switches, an error, or nil if there are no switches.
func (u *Unifi) GetUSWs(site *Site) ([]*USW, error) {
	var response struct {
		Data []json.RawMessage `json:"data"`
	}

	err := u.GetData(fmt.Sprintf(APIDevicePath, site.Name), &response)
	if err != nil {
		return nil, err
	}

	return u.parseDevices(response.Data, site).USWs, nil
}

// GetUAPs returns all access points, an error, or nil if there are no APs.
func (u *Unifi) GetUAPs(site *Site) ([]*UAP, error) {
	var response struct {
		Data []json.RawMessage `json:"data"`
	}

	err := u.GetData(fmt.Sprintf(APIDevicePath, site.Name), &response)
	if err != nil {
		return nil, err
	}

	return u.parseDevices(response.Data, site).UAPs, nil
}

// GetUDMs returns all dream machines, an error, or nil if there are no UDMs.
func (u *Unifi) GetUDMs(site *Site) ([]*UDM, error) {
	var response struct {
		Data []json.RawMessage `json:"data"`
	}

	err := u.GetData(fmt.Sprintf(APIDevicePath, site.Name), &response)
	if err != nil {
		return nil, err
	}

	return u.parseDevices(response.Data, site).UDMs, nil
}

// GetUXGs returns all 10Gb gateways, an error, or nil if there are no UXGs.
func (u *Unifi) GetUXGs(site *Site) ([]*UXG, error) {
	var response struct {
		Data []json.RawMessage `json:"data"`
	}

	err := u.GetData(fmt.Sprintf(APIDevicePath, site.Name), &response)
	if err != nil {
		return nil, err
	}

	return u.parseDevices(response.Data, site).UXGs, nil
}

// GetUSGs returns all 1Gb gateways, an error, or nil if there are no USGs.
func (u *Unifi) GetUSGs(site *Site) ([]*USG, error) {
	var response struct {
		Data []json.RawMessage `json:"data"`
	}

	err := u.GetData(fmt.Sprintf(APIDevicePath, site.Name), &response)
	if err != nil {
		return nil, err
	}

	return u.parseDevices(response.Data, site).USGs, nil
}

// parseDevices parses the raw JSON from the Unifi Controller into device structures.
func (u *Unifi) parseDevices(data []json.RawMessage, site *Site) *Devices {
	devices := new(Devices)

	for _, r := range data {
		// Loop each item in the raw JSON message, detect its type and unmarshal it.
		o := make(map[string]interface{})
		if u.unmarshalDevice("map", r, &o) != nil {
			u.ErrorLog("unknown asset type - cannot find asset type in payload - skipping")
			continue
		}

		assetType, _ := o["type"].(string)
		u.DebugLog("Unmarshalling Device Type: %v, site %s ", assetType, site.SiteName)
		// Choose which type to unmarshal into based on the "type" json key.

		switch assetType { // Unmarshal again into the correct type..
		case "uap":
			u.unmarshallUAP(site, r, devices)
		case "ugw", "usg": // in case they ever fix the name in the api.
			u.unmarshallUSG(site, r, devices)
		case "usw":
			u.unmarshallUSW(site, r, devices)
		case "udm":
			u.unmarshallUDM(site, r, devices)
		case "uxg":
			u.unmarshallUXG(site, r, devices)
		default:
			u.ErrorLog("unknown asset type - %v - skipping", assetType)
		}
	}

	return devices
}

func (u *Unifi) unmarshallUAP(site *Site, payload json.RawMessage, devices *Devices) {
	dev := &UAP{SiteName: site.SiteName, SourceName: u.URL}
	if u.unmarshalDevice("uap", payload, dev) == nil {
		dev.Name = strings.TrimSpace(pick(dev.Name, dev.Mac))
		dev.site = site
		devices.UAPs = append(devices.UAPs, dev)
	}
}

func (u *Unifi) unmarshallUSG(site *Site, payload json.RawMessage, devices *Devices) {
	dev := &USG{SiteName: site.SiteName, SourceName: u.URL}
	if u.unmarshalDevice("ugw", payload, dev) == nil {
		dev.Name = strings.TrimSpace(pick(dev.Name, dev.Mac))
		dev.site = site
		devices.USGs = append(devices.USGs, dev)
	}
}

func (u *Unifi) unmarshallUSW(site *Site, payload json.RawMessage, devices *Devices) {
	dev := &USW{SiteName: site.SiteName, SourceName: u.URL}
	if u.unmarshalDevice("usw", payload, dev) == nil {
		dev.Name = strings.TrimSpace(pick(dev.Name, dev.Mac))
		dev.site = site
		devices.USWs = append(devices.USWs, dev)
	}
}

func (u *Unifi) unmarshallUXG(site *Site, payload json.RawMessage, devices *Devices) {
	dev := &UXG{SiteName: site.SiteName, SourceName: u.URL}
	if u.unmarshalDevice("uxg", payload, dev) == nil {
		dev.Name = strings.TrimSpace(pick(dev.Name, dev.Mac))
		dev.site = site
		devices.UXGs = append(devices.UXGs, dev)
	}
}

func (u *Unifi) unmarshallUDM(site *Site, payload json.RawMessage, devices *Devices) {
	dev := &UDM{SiteName: site.SiteName, SourceName: u.URL}
	if u.unmarshalDevice("udm", payload, dev) == nil {
		dev.Name = strings.TrimSpace(pick(dev.Name, dev.Mac))
		dev.site = site
		devices.UDMs = append(devices.UDMs, dev)
	}
}

// unmarshalDevice handles logging for the unmarshal operations in parseDevices().
func (u *Unifi) unmarshalDevice(dev string, data json.RawMessage, v interface{}) (err error) {
	if err = json.Unmarshal(data, v); err != nil {
		u.ErrorLog("json.Unmarshal(%v): %v", dev, err)
		u.ErrorLog("Enable Debug Logging to output the failed payload.")

		json, err := data.MarshalJSON()
		u.DebugLog("Failed Payload: %s (marshal err: %v)", json, err)
		u.DebugLog("The above payload can prove useful during torubleshooting when you open an Issue:")
		u.DebugLog("==- https://github.com/unifi-poller/unifi/issues/new -==")
	}

	if err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}

	return nil
}

// pick returns the first non empty string in a list.
// used in a few places around this library.
func pick(strings ...string) string {
	for _, s := range strings {
		if s != "" {
			return s
		}
	}

	return ""
}
