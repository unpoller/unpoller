package unifi

import "encoding/json"

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
