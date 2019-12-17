// Package inputunifi implements the poller.Input interface and bridges the gap between
// metrics from the unifi library, and the augments required to pump them into unifi-poller.
package inputunifi

import (
	"fmt"
	"os"
	"strings"

	"sync"

	"github.com/davidnewhall/unifi-poller/pkg/poller"
	"golift.io/unifi"
)

// InputUnifi contains the running data.
type InputUnifi struct {
	Config *Config `json:"unifi" toml:"unifi" xml:"unifi" yaml:"unifi"`
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
	Sites     []string     `json:"sites,omitempty" toml:"sites,omitempty" xml:"site" yaml:"sites"`
	Unifi     *unifi.Unifi `json:"-" toml:"-" xml:"-" yaml:"-"`
}

// Config contains our configuration data
type Config struct {
	sync.RWMutex               // locks the Unifi struct member when re-authing to unifi.
	Disable      bool          `json:"disable" toml:"disable" xml:"disable" yaml:"disable"`
	Controllers  []*Controller `json:"controllers" toml:"controller" xml:"controller" yaml:"controllers"`
}

func init() {
	u := &InputUnifi{}

	poller.NewInput(&poller.InputPlugin{
		Name:   "unifi",
		Input:  u, // this library implements poller.Input interface for Metrics().
		Config: u, // Defines our config data interface.
	})
}

// getUnifi (re-)authenticates to a unifi controller.
func (u *InputUnifi) getUnifi(c *Controller) error {
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

// checkSites makes sure the list of provided sites exists on the controller.
// This only runs once during initialization.
func (u *InputUnifi) checkSites(c *Controller) error {
	u.Config.RLock()
	defer u.Config.RUnlock()

	if len(c.Sites) < 1 || c.Sites[0] == "" {
		c.Sites = []string{"all"}
	}

	u.LogDebugf("Checking Controller Sites List")

	sites, err := c.Unifi.GetSites()
	if err != nil {
		return err
	}

	msg := []string{}
	for _, site := range sites {
		msg = append(msg, site.Name+" ("+site.Desc+")")
	}

	u.Logf("Found %d site(s) on controller %s: %v", len(msg), c.Name, strings.Join(msg, ", "))

	if StringInSlice("all", c.Sites) {
		c.Sites = []string{"all"}
		return nil
	}

	keep := []string{}

FIRST:
	for _, s := range c.Sites {
		for _, site := range sites {
			if s == site.Name {
				keep = append(keep, s)
				continue FIRST
			}
		}
		u.LogErrorf("Configured site not found on controller %s: %v", c.Name, s)
	}

	if c.Sites = keep; len(keep) < 1 {
		c.Sites = []string{"all"}
	}

	return nil
}

func (u *InputUnifi) dumpSitesJSON(c *Controller, path, name string, sites unifi.Sites) ([]byte, error) {
	allJSON := []byte{}

	for _, s := range sites {
		apiPath := fmt.Sprintf(path, s.Name)
		_, _ = fmt.Fprintf(os.Stderr, "[INFO] Dumping %s: '%s' JSON for site: %s (%s):\n", name, apiPath, s.Desc, s.Name)

		body, err := c.Unifi.GetJSON(apiPath)
		if err != nil {
			return allJSON, err
		}

		allJSON = append(allJSON, body...)
	}

	return allJSON, nil
}

// StringInSlice returns true if a string is in a slice.
func StringInSlice(str string, slice []string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, str) {
			return true
		}
	}

	return false
}
