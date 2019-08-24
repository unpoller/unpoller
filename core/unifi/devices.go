package unifi

import "encoding/json"

// parseDevices parses the raw JSON from the Unifi Controller into device structures.
func (u *Unifi) parseDevices(data []json.RawMessage, siteName string) *Devices {
	devices := new(Devices)
	for _, r := range data {
		// Loop each item in the raw JSON message, detect its type and unmarshal it.
		assetType := "<type key missing>"
		if o := make(map[string]interface{}); u.unmarshalDevice("map", r, &o) != nil {
			continue
		} else if t, ok := o["type"].(string); ok {
			assetType = t
		}
		u.DebugLog("Unmarshalling Device Type: %v, site %s ", assetType, siteName)
		// Choose which type to unmarshal into based on the "type" json key.
		switch assetType { // Unmarshal again into the correct type..
		case "uap":
			dev := &UAP{SiteName: siteName}
			if u.unmarshalDevice(assetType, r, dev) == nil {
				dev.Name = pick(dev.Name, dev.Mac)
				devices.UAPs = append(devices.UAPs, dev)
			}
		case "ugw", "usg": // in case they ever fix the name in the api.
			dev := &USG{SiteName: siteName}
			if u.unmarshalDevice(assetType, r, dev) == nil {
				dev.Name = pick(dev.Name, dev.Mac)
				devices.USGs = append(devices.USGs, dev)
			}
		case "usw":
			dev := &USW{SiteName: siteName}
			if u.unmarshalDevice(assetType, r, dev) == nil {
				dev.Name = pick(dev.Name, dev.Mac)
				devices.USWs = append(devices.USWs, dev)
			}
		case "udm":
			dev := &UDM{SiteName: siteName}
			if u.unmarshalDevice(assetType, r, dev) == nil {
				dev.Name = pick(dev.Name, dev.Mac)
				devices.UDMs = append(devices.UDMs, dev)
			}
		default:
			u.ErrorLog("unknown asset type - %v - skipping", assetType)
		}
	}
	return devices
}

// unmarshalDevice handles logging for the unmarshal operations in parseDevices().
func (u *Unifi) unmarshalDevice(dev string, data json.RawMessage, v interface{}) (err error) {
	if err = json.Unmarshal(data, v); err != nil {
		u.ErrorLog("json.Unmarshal(%v): %v", dev, err)
		u.ErrorLog("Enable Debug Logging to output the failed payload.")
		json, err := data.MarshalJSON()
		u.DebugLog("Failed Payload: %s (marshal err: %v)", json, err)
		u.DebugLog("The above payload can prove useful during torubleshooting when you open an Issue:")
		u.DebugLog("==- https://github.com/golift/unifi/issues/new -==")
	}
	return err
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
