package unifi

import (
	"encoding/json"
	"fmt"
)

// Known commands that can be sent to device manager.
const (
	DevMgrPowerCycle = "power-cycle"
)

// command is the type marshalled and sent to APIDevMgrPath.
type command struct {
	Command   string `json:"cmd"`
	Mac       string `json:"mac"`
	PortIndex int    `json:"port_idx"`
}

// PowerCycle shuts off the PoE and turns it back on for a specific port.
// Get a USW from the device list to call this.
func (u *USW) PowerCycle(portIndex int) error {
	data, err := json.Marshal(&command{
		Command:   DevMgrPowerCycle,
		Mac:       u.Mac,
		PortIndex: portIndex,
	})
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	_, err = u.controller.GetJSON(APIDevMgrPath, string(data))
	if err != nil {
		return fmt.Errorf("controller: %w", err)
	}

	return nil
}
