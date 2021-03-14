package unifi

import (
	"encoding/json"
	"fmt"
)

// Known commands that can be sent to device manager. All of these are implemented.
//nolint:lll // https://ubntwiki.com/products/software/unifi-controller/api#callable
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

// devMgrCmd is the type marshalled and sent to APIDevMgrPath.
type devMgrCmd struct {
	Cmd    string `json:"cmd"`                  // Required.
	Mac    string `json:"mac"`                  // Device MAC (required for most, but not all).
	URL    string `json:"url,omitempty"`        // External Upgrade only.
	Inform string `json:"inform_url,omitempty"` // Migration only.
	Port   int    `json:"port_idx,omitempty"`   // Power Cycle only.
}

// devMgrCommandReply is for commands with a return value.
func (s *Site) devMgrCommandReply(cmd *devMgrCmd) ([]byte, error) {
	data, err := json.Marshal(cmd)
	if err != nil {
		return nil, fmt.Errorf("json marshal: %w", err)
	}

	b, err := s.controller.GetJSON(fmt.Sprintf(APIDevMgrPath, s.Name), string(data))
	if err != nil {
		return nil, fmt.Errorf("controller: %w", err)
	}

	return b, nil
}

// devMgrCommandSimple is for commands with no return value.
func (s *Site) devMgrCommandSimple(cmd *devMgrCmd) error {
	_, err := s.devMgrCommandReply(cmd)
	return err
}

// PowerCycle shuts off the PoE and turns it back on for a specific port.
// Get a USW from the device list to call this.
func (u *USW) PowerCycle(portIndex int) error {
	return u.site.devMgrCommandSimple(&devMgrCmd{
		Cmd:  DevMgrPowerCycle,
		Mac:  u.Mac,
		Port: portIndex,
	})
}

// ScanRF begins a spectrum scan on an access point.
func (u *UAP) ScanRF() error {
	return u.site.devMgrCommandSimple(&devMgrCmd{Cmd: DevMgrSpectrumScan, Mac: u.Mac})
}

// Restart a device by MAC address on your site.
func (s *Site) Restart(mac string) error {
	return s.devMgrCommandSimple(&devMgrCmd{Cmd: DevMgrRestart, Mac: mac})
}

// Restart an access point.
func (u *UAP) Restart() error {
	return u.site.Restart(u.Mac)
}

// Restart a switch.
func (u *USW) Restart() error {
	return u.site.Restart(u.Mac)
}

// Restart a security gateway.
func (u *USG) Restart() error {
	return u.site.Restart(u.Mac)
}

// Restart a dream machine.
func (u *UDM) Restart() error {
	return u.site.Restart(u.Mac)
}

// Restart a 10Gb security gateway.
func (u *UXG) Restart() error {
	return u.site.Restart(u.Mac)
}

// Locate a device by MAC address on your site. This makes it blink.
func (s *Site) Locate(mac string) error {
	return s.devMgrCommandSimple(&devMgrCmd{Cmd: DevMgrSetLocate, Mac: mac})
}

// Locate an access point.
func (u *UAP) Locate() error {
	return u.site.Locate(u.Mac)
}

// Locate a switch.
func (u *USW) Locate() error {
	return u.site.Locate(u.Mac)
}

// Locate a security gateway.
func (u *USG) Locate() error {
	return u.site.Locate(u.Mac)
}

// Locate a dream machine.
func (u *UDM) Locate() error {
	return u.site.Locate(u.Mac)
}

// Locate a 10Gb security gateway.
func (u *UXG) Locate() error {
	return u.site.Locate(u.Mac)
}

// Unlocate a device by MAC address on your site. This makes it stop blinking.
func (s *Site) Unlocate(mac string) error {
	return s.devMgrCommandSimple(&devMgrCmd{Cmd: DevMgrUnsetLocate, Mac: mac})
}

// Unlocate an access point (stop blinking).
func (u *UAP) Unlocate() error {
	return u.site.Unlocate(u.Mac)
}

// Unlocate a switch (stop blinking).
func (u *USW) Unlocate() error {
	return u.site.Unlocate(u.Mac)
}

// Unlocate a security gateway (stop blinking).
func (u *USG) Unlocate() error {
	return u.site.Unlocate(u.Mac)
}

// Unlocate a dream machine (stop blinking).
func (u *UDM) Unlocate() error {
	return u.site.Unlocate(u.Mac)
}

// Unlocate a 10Gb security gateway (stop blinking).
func (u *UXG) Unlocate() error {
	return u.site.Unlocate(u.Mac)
}

// Provision force provisions a device by MAC address on your site.
func (s *Site) Provision(mac string) error {
	return s.devMgrCommandSimple(&devMgrCmd{Cmd: DevMgrForceProvision, Mac: mac})
}

// Provision an access point forcefully.
func (u *UAP) Provision() error {
	return u.site.Provision(u.Mac)
}

// Provision a switch forcefully.
func (u *USW) Provision() error {
	return u.site.Provision(u.Mac)
}

// Provision a security gateway forcefully.
func (u *USG) Provision() error {
	return u.site.Provision(u.Mac)
}

// Provision a dream machine forcefully.
func (u *UDM) Provision() error {
	return u.site.Provision(u.Mac)
}

// Provision a 10Gb security gateway forcefully.
func (u *UXG) Provision() error {
	return u.site.Provision(u.Mac)
}

// Upgrade starts a firmware upgrade on a device by MAC address on your site.
// URL is optional. If URL is not "" an external upgrade is performed.
func (s *Site) Upgrade(mac string, url string) error {
	if url == "" {
		return s.devMgrCommandSimple(&devMgrCmd{Cmd: DevMgrUpgrade, Mac: mac})
	}

	return s.devMgrCommandSimple(&devMgrCmd{Cmd: DevMgrUpgradeExternal, Mac: mac, URL: url})
}

// Upgrade firmware on an access point.
// URL is optional. If URL is not "" an external upgrade is performed.
func (u *UAP) Upgrade(url string) error {
	return u.site.Upgrade(u.Mac, url)
}

// Upgrade firmware on a switch.
// URL is optional. If URL is not "" an external upgrade is performed.
func (u *USW) Upgrade(url string) error {
	return u.site.Upgrade(u.Mac, url)
}

// Upgrade firmware on a security gateway.
// URL is optional. If URL is not "" an external upgrade is performed.
func (u *USG) Upgrade(url string) error {
	return u.site.Upgrade(u.Mac, url)
}

// Upgrade firmware on a dream machine.
// URL is optional. If URL is not "" an external upgrade is performed.
func (u *UDM) Upgrade(url string) error {
	return u.site.Upgrade(u.Mac, url)
}

// Upgrade formware on a 10Gb security gateway.
// URL is optional. If URL is not "" an external upgrade is performed.
func (u *UXG) Upgrade(url string) error {
	return u.site.Upgrade(u.Mac, url)
}

// Migrate sends a device to another controller's URL.
// Probably does not work on devices with built-in controllers like UDM & UXG.
func (s *Site) Migrate(mac string, url string) error {
	return s.devMgrCommandSimple(&devMgrCmd{Cmd: DevMgrMigrate, Mac: mac, Inform: url})
}

// Migrate sends an access point to another controller's URL.
func (u *UAP) Migrate(url string) error {
	return u.site.Migrate(u.Mac, url)
}

// Migrate sends a switch to another controller's URL.
func (u *USW) Migrate(url string) error {
	return u.site.Migrate(u.Mac, url)
}

// Migrate sends a security gateway to another controller's URL.
func (u *USG) Migrate(url string) error {
	return u.site.Migrate(u.Mac, url)
}

// Migrate sends a 10Gb gateway to another controller's URL.
func (u *UXG) Migrate(url string) error {
	return u.site.Migrate(u.Mac, url)
}

// CancelMigrate stops a migration in progress.
// Probably does not work on devices with built-in controllers like UDM & UXG.
func (s *Site) CancelMigrate(mac string) error {
	return s.devMgrCommandSimple(&devMgrCmd{Cmd: DevMgrCancelMigrate, Mac: mac})
}

// CancelMigrate stops an access point migration in progress.
func (u *UAP) CancelMigrate() error {
	return u.site.CancelMigrate(u.Mac)
}

// CancelMigrate stops a switch migration in progress.
func (u *USW) CancelMigrate() error {
	return u.site.CancelMigrate(u.Mac)
}

// CancelMigrate stops a security gateway migration in progress.
func (u *USG) CancelMigrate() error {
	return u.site.CancelMigrate(u.Mac)
}

// CancelMigrate stops 10Gb gateway a migration in progress.
func (u *UXG) CancelMigrate() error {
	return u.site.CancelMigrate(u.Mac)
}

// Adopt a device by MAC address to your site.
func (s *Site) Adopt(mac string) error {
	return s.devMgrCommandSimple(&devMgrCmd{Cmd: DevMgrAdopt, Mac: mac})
}

// SpeedTest begins a speed test on a site.
func (s *Site) SpeedTest() error {
	return s.devMgrCommandSimple(&devMgrCmd{Cmd: DevMgrSpeedTest})
}

// SpeedTestStatus returns the raw response for the status of a speed test.
// XXX: marshal the response into a data structure. This method will change!
func (s *Site) SpeedTestStatus() ([]byte, error) {
	body, err := s.devMgrCommandReply(&devMgrCmd{Cmd: DevMgrSpeedTestStatus})
	// marshal into struct here.
	return body, err
}
