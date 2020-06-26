// Package inputunifi implements the poller.Input interface and bridges the gap between
// metrics from the unifi library, and the augments required to pump them into unifi-poller.
package inputunifi

import (
	"io/ioutil"
	"strings"
	"time"

	"sync"

	"github.com/pkg/errors"
	"github.com/unifi-poller/poller"
	"github.com/unifi-poller/unifi"
)

const (
	PluginName  = "unifi" // PluginName is the name of this input plugin.
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
	VerifySSL  *bool        `json:"verify_ssl" toml:"verify_ssl" xml:"verify_ssl" yaml:"verify_ssl"`
	SaveAnomal *bool        `json:"save_anomalies" toml:"save_anomalies" xml:"save_anomalies" yaml:"save_anomalies"`
	SaveAlarms *bool        `json:"save_alarms" toml:"save_alarms" xml:"save_alarms" yaml:"save_alarms"`
	SaveEvents *bool        `json:"save_events" toml:"save_events" xml:"save_events" yaml:"save_events"`
	SaveIDS    *bool        `json:"save_ids" toml:"save_ids" xml:"save_ids" yaml:"save_ids"`
	SaveDPI    *bool        `json:"save_dpi" toml:"save_dpi" xml:"save_dpi" yaml:"save_dpi"`
	HashPII    *bool        `json:"hash_pii" toml:"hash_pii" xml:"hash_pii" yaml:"hash_pii"`
	SaveSites  *bool        `json:"save_sites" toml:"save_sites" xml:"save_sites" yaml:"save_sites"`
	User       string       `json:"user" toml:"user" xml:"user" yaml:"user"`
	Pass       string       `json:"pass" toml:"pass" xml:"pass" yaml:"pass"`
	URL        string       `json:"url" toml:"url" xml:"url" yaml:"url"`
	Sites      []string     `json:"sites,omitempty" toml:"sites,omitempty" xml:"site" yaml:"sites"`
	Unifi      *unifi.Unifi `json:"-" toml:"-" xml:"-" yaml:"-"`
}

// Config contains our configuration data.
type Config struct {
	sync.RWMutex               // locks the Unifi struct member when re-authing to unifi.
	Default      Controller    `json:"defaults" toml:"defaults" xml:"default" yaml:"defaults"`
	Disable      bool          `json:"disable" toml:"disable" xml:"disable,attr" yaml:"disable"`
	Dynamic      bool          `json:"dynamic" toml:"dynamic" xml:"dynamic,attr" yaml:"dynamic"`
	Controllers  []*Controller `json:"controllers" toml:"controller" xml:"controller" yaml:"controllers"`
}

type Metrics struct {
	TS         time.Time
	Sites      []*unifi.Site
	Clients    []*unifi.Client
	SitesDPI   []*unifi.DPITable
	ClientsDPI []*unifi.DPITable
	Devices    *unifi.Devices
}

func init() { // nolint: gochecknoinits
	u := &InputUnifi{
		dynamic: make(map[string]*Controller),
	}

	poller.NewInput(&poller.InputPlugin{
		Name:   PluginName,
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
		VerifySSL: *c.VerifySSL,
		ErrorLog:  u.LogErrorf, // Log all errors.
		DebugLog:  u.LogDebugf, // Log debug messages.
	})
	if err != nil {
		c.Unifi = nil
		return errors.Wrap(err, "unifi controller")
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

	u.Logf("Found %d site(s) on controller %s: %v", len(msg), c.URL, strings.Join(msg, ", "))

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
		u.LogErrorf("Configured site not found on controller %s: %v", c.URL, s)
	}

	if c.Sites = keep; len(keep) == 0 {
		c.Sites = []string{"all"}
	}

	return nil
}

func (u *InputUnifi) getPassFromFile(filename string) string {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		u.LogErrorf("Reading UniFi Password File: %v", err)
	}

	return strings.TrimSpace(string(b))
}

// setDefaults sets the default defaults.
func (u *InputUnifi) setDefaults(c *Controller) {
	t := true
	f := false

	// Default defaults.
	if c.SaveSites == nil {
		c.SaveSites = &t
	}

	if c.VerifySSL == nil {
		c.VerifySSL = &f
	}

	if c.HashPII == nil {
		c.HashPII = &f
	}

	if c.SaveDPI == nil {
		c.SaveDPI = &f
	}

	if c.SaveIDS == nil {
		c.SaveIDS = &f
	}

	if c.SaveEvents == nil {
		c.SaveEvents = &f
	}

	if c.SaveAlarms == nil {
		c.SaveAlarms = &f
	}

	if c.SaveAnomal == nil {
		c.SaveAnomal = &f
	}

	if c.URL == "" {
		c.URL = defaultURL
	}

	if strings.HasPrefix(c.Pass, "file://") {
		c.Pass = u.getPassFromFile(strings.TrimPrefix(c.Pass, "file://"))
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
}

// setControllerDefaults sets defaults for the for controllers.
// Any missing values come from defaults (above).
func (u *InputUnifi) setControllerDefaults(c *Controller) *Controller {
	// Configured controller defaults.
	if c.SaveSites == nil {
		c.SaveSites = u.Default.SaveSites
	}

	if c.VerifySSL == nil {
		c.VerifySSL = u.Default.VerifySSL
	}

	if c.HashPII == nil {
		c.HashPII = u.Default.HashPII
	}

	if c.SaveDPI == nil {
		c.SaveDPI = u.Default.SaveDPI
	}

	if c.SaveIDS == nil {
		c.SaveIDS = u.Default.SaveIDS
	}

	if c.SaveEvents == nil {
		c.SaveEvents = u.Default.SaveEvents
	}

	if c.SaveAlarms == nil {
		c.SaveAlarms = u.Default.SaveAlarms
	}

	if c.SaveAnomal == nil {
		c.SaveAnomal = u.Default.SaveAnomal
	}

	if c.URL == "" {
		c.URL = u.Default.URL
	}

	if strings.HasPrefix(c.Pass, "file://") {
		c.Pass = u.getPassFromFile(strings.TrimPrefix(c.Pass, "file://"))
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

	return c
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
