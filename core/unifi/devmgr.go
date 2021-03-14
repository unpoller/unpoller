package unifi

import (
	"encoding/json"
	"fmt"
)

// Known commands that can be sent to device manager.
//nolint:lll
const (
	DevMgrPowerCycle      = "power-cycle"      // mac = switch mac (required), port_idx = PoE port to cycle (required)
	DevMgrAdopt           = "adopt"            // mac = device mac (required)
	DevMgrRestart         = "restart"          // mac = device mac (required)
	DevMgrForceProvision  = "force-provision"  // mac = device mac (required)
	DevMgrSpeedTest       = "speedtest"        // Start a speed test
	DevMgrSpeedTestStatus = "speedtest-status" // Get current state of the speed test
	DevMgrSetLocate       = "set-locate"       // mac = device mac (required): blink unit to locate
	DevMgrUnsetLocate     = "unset-locate"     // mac = device mac (required): led to normal state
	DevMgrUpgrade         = "upgrade"          // mac = device mac (required): upgrade firmware
	DevMgrUpgradeExternal = "upgrade-external" // mac = device mac (required), url = firmware URL (required)
	DevMgrMigrate         = "migrate"          // mac = device mac (required), inform_url = New Inform URL for device (required)
	DevMgrCancelMigrate   = "cancel-migrate"   // mac = device mac (required)
	DevMgrSpectrumScan    = "spectrum-scan"    // mac = AP mac     (required): trigger RF scan
)

// command is the type marshalled and sent to APIDevMgrPath.
type devMgrCommand struct {
	Command   string `json:"cmd"`
	Mac       string `json:"mac"`
	URL       string `json:"url,omitempty"`
	InformURL string `json:"inform_url,omitempty"`
	PortIndex int    `json:"port_idx,omitempty"`
}

// PowerCycle shuts off the PoE and turns it back on for a specific port.
// Get a USW from the device list to call this.
func (u *USW) PowerCycle(portIndex int) error {
	data, err := json.Marshal(&devMgrCommand{
		Command:   DevMgrPowerCycle,
		Mac:       u.Mac,
		PortIndex: portIndex,
	})
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	_, err = u.controller.GetJSON(fmt.Sprintf(APIDevMgrPath, u.SiteName), string(data))
	if err != nil {
		return fmt.Errorf("controller: %w", err)
	}

	return nil
}

// ScanRF begins a spectrum scan on an access point.
func (u *UAP) ScanRF() error {
	data, err := json.Marshal(&devMgrCommand{Command: DevMgrSpectrumScan, Mac: u.Mac})
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	_, err = u.controller.GetJSON(fmt.Sprintf(APIDevMgrPath, u.SiteName), string(data))
	if err != nil {
		return fmt.Errorf("controller: %w", err)
	}

	return nil
}

// Restart a device by MAC address on your site.
func (s *Site) Restart(mac string) error {
	data, err := json.Marshal(&devMgrCommand{Command: DevMgrRestart, Mac: mac})
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	_, err = s.controller.GetJSON(fmt.Sprintf(APIDevMgrPath, s.Name), string(data))
	if err != nil {
		return fmt.Errorf("controller: %w", err)
	}

	return nil
}

func (u *UAP) Restart() error {
	return u.site.Restart(u.Mac)
}

func (u *USW) Restart() error {
	return u.site.Restart(u.Mac)
}

func (u *USG) Restart() error {
	return u.site.Restart(u.Mac)
}

func (u *UDM) Restart() error {
	return u.site.Restart(u.Mac)
}

func (u *UXG) Restart() error {
	return u.site.Restart(u.Mac)
}

// Adopt a device by MAC address to your site.
func (s *Site) Adopt(mac string) error {
	data, err := json.Marshal(&devMgrCommand{Command: DevMgrAdopt, Mac: mac})
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	_, err = s.controller.GetJSON(fmt.Sprintf(APIDevMgrPath, s.Name), string(data))
	if err != nil {
		return fmt.Errorf("controller: %w", err)
	}

	return nil
}

// SpeedTest begins a speed test.
func (s *Site) SpeedTest() error {
	data, err := json.Marshal(&devMgrCommand{Command: DevMgrSpeedTest})
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	_, err = s.controller.GetJSON(fmt.Sprintf(APIDevMgrPath, s.Name), string(data))
	if err != nil {
		return fmt.Errorf("controller: %w", err)
	}

	return nil
}

// SpeedTestStatus returns the raw response for the status of a speed test.
func (s *Site) SpeedTestStatus() ([]byte, error) {
	data, err := json.Marshal(&devMgrCommand{Command: DevMgrSpeedTestStatus})
	if err != nil {
		return nil, fmt.Errorf("json marshal: %w", err)
	}

	b, err := s.controller.GetJSON(fmt.Sprintf(APIDevMgrPath, s.Name), string(data))
	if err != nil {
		return nil, fmt.Errorf("controller: %w", err)
	}

	return b, nil
}
