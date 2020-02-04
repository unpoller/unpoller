// Package inputunifi implements the poller.Input interface and bridges the gap between
// metrics from the unifi library, and the augments required to pump them into unifi-poller.
package inputunifi

import (
	"fmt"
	"os"
	"strings"

	"sync"

	"github.com/unifi-poller/poller"
	"github.com/unifi-poller/unifi"
)

const (
	defaultURL  = "https://127.0.0.1:8443"
	defaultUser = "unifipoller"
	defaultPass = "unifipoller"
	defaultSite = "all"
)

// InputUnifi contains the running data.
type InputUnifi struct {
	*Config    `json:"unifi" toml:"unifi" xml:"unifi" yaml:"unifi"`
	dynamic    map[string]*Controller
	sync.Mutex // to lock the map above.
	poller.Logger
}

// Controller represents the configuration for a UniFi Controller.
// Each polled controller may have its own configuration.
type Controller struct {
	VerifySSL bool         `json:"verify_ssl" toml:"verify_ssl" xml:"verify_ssl" yaml:"verify_ssl"`
	SaveIDS   bool         `json:"save_ids" toml:"save_ids" xml:"save_ids" yaml:"save_ids"`
	SaveDPI   bool         `json:"save_dpi" toml:"save_dpi" xml:"save_dpi" yaml:"save_dpi"`
	SaveSites *bool        `json:"save_sites" toml:"save_sites" xml:"save_sites" yaml:"save_sites"`
	Role      string       `json:"role" toml:"role" xml:"role,attr" yaml:"role"`
	User      string       `json:"user" toml:"user" xml:"user" yaml:"user"`
	Pass      string       `json:"pass" toml:"pass" xml:"pass" yaml:"pass"`
	URL       string       `json:"url" toml:"url" xml:"url" yaml:"url"`
	New       bool         `json:"new" toml:"new" xml:"new" yaml:"new"`
	Sites     []string     `json:"sites,omitempty" toml:"sites,omitempty" xml:"site" yaml:"sites"`
	Unifi     *unifi.Unifi `json:"-" toml:"-" xml:"-" yaml:"-"`
}

// Config contains our configuration data
type Config struct {
	sync.RWMutex               // locks the Unifi struct member when re-authing to unifi.
	Default      Controller    `json:"defaults" toml:"defaults" xml:"default" yaml:"defaults"`
	Disable      bool          `json:"disable" toml:"disable" xml:"disable,attr" yaml:"disable"`
	Dynamic      bool          `json:"dynamic" toml:"dynamic" xml:"dynamic,attr" yaml:"dynamic"`
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

	u.Lock()
	defer u.Unlock()

	if c.Unifi != nil {
		c.Unifi.CloseIdleConnections()
	}

	// Create an authenticated session to the Unifi Controller.
	c.Unifi, err = unifi.NewUnifi(&unifi.Config{
		User:      c.User,
		Pass:      c.Pass,
		URL:       c.URL,
		New:       c.New,
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
	u.RLock()
	defer u.RUnlock()

	if len(c.Sites) == 0 || c.Sites[0] == "" {
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

	u.Logf("Found %d site(s) on controller %s: %v", len(msg), c.Role, strings.Join(msg, ", "))

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
		u.LogErrorf("Configured site not found on controller %s: %v", c.Role, s)
	}

	if c.Sites = keep; len(keep) == 0 {
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

// setDefaults sets defaults for the defaults and for the controllers.
// which one depends on the useDefaults boolean.
func (u *InputUnifi) setDefaults(c *Controller, useDefaults bool) {
	// Default defaults.
	if useDefaults {
		if c.SaveSites == nil {
			t := true
			c.SaveSites = &t
		}

		if c.URL == "" {
			c.URL = defaultURL
		}

		if c.Role == "" {
			c.Role = c.URL
		}

		if c.Pass == "" {
			c.Pass = defaultPass
		}

		if c.User == "" {
			c.User = defaultUser
		}

		if len(c.Sites) == 0 {
			c.Sites = []string{defaultSite}
		}

		return
	}

	// Configured controller defaults.
	if c.SaveSites == nil {
		c.SaveSites = u.Default.SaveSites
	}

	if c.URL == "" {
		c.URL = u.Default.URL
	}

	if c.Role == "" && u.Default.Role != u.Default.URL {
		c.Role = u.Default.Role
	} else if c.Role == "" {
		c.Role = c.URL
	}

	if c.Pass == "" {
		c.Pass = u.Default.Pass
	}

	if c.User == "" {
		c.User = u.Default.User
	}

	if len(c.Sites) == 0 {
		c.Sites = u.Default.Sites
	}
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
