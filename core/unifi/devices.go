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
			if uap := (UAP{}); u.unmarshalDevice(assetType, r, &uap) == nil {
				uap.SiteName = siteName
				devices.UAPs = append(devices.UAPs, uap)
			}
		case "ugw", "usg": // in case they ever fix the name in the api.
			if usg := (USG{}); u.unmarshalDevice(assetType, r, &usg) == nil {
				usg.SiteName = siteName
				devices.USGs = append(devices.USGs, usg)
			}
		case "usw":
			if usw := (USW{}); u.unmarshalDevice(assetType, r, &usw) == nil {
				usw.SiteName = siteName
				devices.USWs = append(devices.USWs, usw)
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
