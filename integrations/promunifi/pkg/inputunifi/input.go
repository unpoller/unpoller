// Package inputunifi implements the poller.Input interface and bridges the gap between
// metrics from the unifi library, and the augments required to pump them into unifi-poller.
package inputunifi

import (
	"fmt"

	"sync"

	"github.com/davidnewhall/unifi-poller/pkg/poller"
	"golift.io/unifi"
)

// InputUnifi contains the running data.
type InputUnifi struct {
	Config Config `json:"unifi" toml:"unifi" xml:"unifi" yaml:"unifi"`
	poller.Logger
}

// Controller represents the configuration for a UniFi Controller.
// Each polled controller may have its own configuration.
type Controller struct {
	VerifySSL bool         `json:"verify_ssl" toml:"verify_ssl" xml:"verify_ssl" yaml:"verify_ssl"`
	SaveIDS   bool         `json:"save_ids" toml:"save_ids" xml:"save_ids" yaml:"save_ids"`
	SaveSites bool         `json:"save_sites" toml:"save_sites" xml:"save_sites" yaml:"save_sites"`
	Name      string       `json:"name" toml:"name" xml:"name,attr" yaml:"name"`
	User      string       `json:"user" toml:"user" xml:"user" yaml:"user"`
	Pass      string       `json:"pass" toml:"pass" xml:"pass" yaml:"pass"`
	URL       string       `json:"url" toml:"url" xml:"url" yaml:"url"`
	Sites     []string     `json:"sites,omitempty" toml:"sites,omitempty" xml:"sites" yaml:"sites"`
	Unifi     *unifi.Unifi `json:"-" toml:"-" xml:"-" yaml:"-"`
}

// Config contains our configuration data
type Config struct {
	sync.RWMutex              // locks the Unifi struct member when re-authing to unifi.
	Disable      bool         `json:"disable" toml:"disable" xml:"disable" yaml:"disable"`
	Controllers  []Controller `json:"controller" toml:"controller" xml:"controller" yaml:"controller"`
}

func init() {
	u := &InputUnifi{}

	poller.NewInput(&poller.InputPlugin{
		Input:  u, // this library implements poller.Input interface for Metrics().
		Config: u, // Defines our config data interface.
	})
}

// getUnifi (re-)authenticates to a unifi controller.
func (u *InputUnifi) getUnifi(c Controller) error {
	var err error

	u.Config.Lock()
	defer u.Config.Unlock()

	if c.Unifi != nil {
		c.Unifi.CloseIdleConnections()
	}

	// Create an authenticated session to the Unifi Controller.
	c.Unifi, err = unifi.NewUnifi(&unifi.Config{
		User:      c.User,
		Pass:      c.Pass,
		URL:       c.URL,
		VerifySSL: c.VerifySSL,
		ErrorLog:  u.LogErrorf, // Log all errors.
		DebugLog:  u.LogDebugf, // Log debug messages.
	})
	if err != nil {
		c.Unifi = nil
		return fmt.Errorf("unifi controller: %v", err)
	}

	u.LogDebugf("Authenticated with controller successfully, %s", c.URL)

	return nil
}
