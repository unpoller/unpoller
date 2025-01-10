// Package inputunifi implements the poller.Input interface and bridges the gap between
// metrics from the unifi library, and the augments required to pump them into unifi-poller.
package inputunifi

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/unpoller/unifi/v5"
	"github.com/unpoller/unpoller/pkg/poller"
)

// PluginName is the name of this input plugin.
const PluginName = "unifi"

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
	Logger     poller.Logger
}

// Controller represents the configuration for a UniFi Controller.
// Each polled controller may have its own configuration.
type Controller struct {
	VerifySSL  *bool        `json:"verify_ssl"     toml:"verify_ssl"     xml:"verify_ssl"     yaml:"verify_ssl"`
	SaveAnomal *bool        `json:"save_anomalies" toml:"save_anomalies" xml:"save_anomalies" yaml:"save_anomalies"`
	SaveAlarms *bool        `json:"save_alarms"    toml:"save_alarms"    xml:"save_alarms"    yaml:"save_alarms"`
	SaveEvents *bool        `json:"save_events"    toml:"save_events"    xml:"save_events"    yaml:"save_events"`
	SaveIDs    *bool        `json:"save_ids"       toml:"save_ids"       xml:"save_ids"       yaml:"save_ids"`
	SaveDPI    *bool        `json:"save_dpi"       toml:"save_dpi"       xml:"save_dpi"       yaml:"save_dpi"`
	SaveRogue  *bool        `json:"save_rogue"     toml:"save_rogue"     xml:"save_rogue"     yaml:"save_rogue"`
	HashPII    *bool        `json:"hash_pii"       toml:"hash_pii"       xml:"hash_pii"       yaml:"hash_pii"`
	DropPII    *bool        `json:"drop_pii"       toml:"drop_pii"       xml:"drop_pii"       yaml:"drop_pii"`
	SaveSites  *bool        `json:"save_sites"     toml:"save_sites"     xml:"save_sites"     yaml:"save_sites"`
	CertPaths  []string     `json:"ssl_cert_paths" toml:"ssl_cert_paths" xml:"ssl_cert_path"  yaml:"ssl_cert_paths"`
	User       string       `json:"user"           toml:"user"           xml:"user"           yaml:"user"`
	Pass       string       `json:"pass"           toml:"pass"           xml:"pass"           yaml:"pass"`
	APIKey     string       `json:"api_key"        toml:"api_key"        xml:"api_key"        yaml:"api_key"`
	URL        string       `json:"url"            toml:"url"            xml:"url"            yaml:"url"`
	Sites      []string     `json:"sites"          toml:"sites"          xml:"site"           yaml:"sites"`
	Unifi      *unifi.Unifi `json:"-"              toml:"-"              xml:"-"              yaml:"-"`
	ID         string       `json:"id,omitempty"` // this is an output, not an input.
}

// Config contains our configuration data.
type Config struct {
	sync.RWMutex               // locks the Unifi struct member when re-authing to unifi.
	Default      Controller    `json:"defaults"    toml:"defaults"   xml:"default"      yaml:"defaults"`
	Disable      bool          `json:"disable"     toml:"disable"    xml:"disable,attr" yaml:"disable"`
	Dynamic      bool          `json:"dynamic"     toml:"dynamic"    xml:"dynamic,attr" yaml:"dynamic"`
	Controllers  []*Controller `json:"controllers" toml:"controller" xml:"controller"   yaml:"controllers"`
}

// Metrics is simply a useful container for everything.
type Metrics struct {
	TS         time.Time
	Sites      []*unifi.Site
	Clients    []*unifi.Client
	SitesDPI   []*unifi.DPITable
	ClientsDPI []*unifi.DPITable
	RogueAPs   []*unifi.RogueAP
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

// getCerts reads in cert files from disk and stores them as a slice of of byte slices.
func (c *Controller) getCerts() ([][]byte, error) {
	if len(c.CertPaths) == 0 {
		return nil, nil
	}

	b := make([][]byte, len(c.CertPaths))

	for i, f := range c.CertPaths {
		d, err := os.ReadFile(f)
		if err != nil {
			return nil, fmt.Errorf("reading SSL cert file: %w", err)
		}

		b[i] = d
	}

	return b, nil
}

// getUnifi (re-)authenticates to a unifi controller.
// If certificate files are provided, they are re-read.
func (u *InputUnifi) getUnifi(c *Controller) error {
	u.Lock()
	defer u.Unlock()

	if c.Unifi != nil {
		c.Unifi.CloseIdleConnections()
	}

	certs, err := c.getCerts()
	if err != nil {
		return err
	}

	// Create an authenticated session to the Unifi Controller.
	c.Unifi, err = unifi.NewUnifi(&unifi.Config{
		User:      c.User,
		Pass:      c.Pass,
		APIKey:    c.APIKey,
		URL:       c.URL,
		SSLCert:   certs,
		VerifySSL: *c.VerifySSL,
		ErrorLog:  u.LogErrorf, // Log all errors.
		DebugLog:  u.LogDebugf, // Log debug messages.
	})
	if err != nil {
		c.Unifi = nil

		return fmt.Errorf("unifi controller: %w", err)
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
		return fmt.Errorf("controller: %w", err)
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
	b, err := os.ReadFile(filename)
	if err != nil {
		u.LogErrorf("Reading UniFi Password File: %v", err)
	}

	return strings.TrimSpace(string(b))
}

// setDefaults sets the default defaults.
func (u *InputUnifi) setDefaults(c *Controller) { //nolint:cyclop
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

	if c.DropPII == nil {
		c.DropPII = &f
	}

	if c.SaveDPI == nil {
		c.SaveDPI = &f
	}

	if c.SaveRogue == nil {
		c.SaveRogue = &f
	}

	if c.SaveIDs == nil {
		c.SaveIDs = &f
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

	if strings.HasPrefix(c.APIKey, "file://") {
		c.APIKey = u.getPassFromFile(strings.TrimPrefix(c.APIKey, "file://"))
	}

	if c.APIKey == "" {
		if c.Pass == "" {
			c.Pass = defaultPass
		}

		if c.User == "" {
			c.User = defaultUser
		}
	} else {
		// clear out user/pass combo, only use API-key
		c.User = ""
		c.Pass = ""
	}

	if len(c.Sites) == 0 {
		c.Sites = []string{defaultSite}
	}
}

// setControllerDefaults sets defaults for the for controllers.
// Any missing values come from defaults (above).
func (u *InputUnifi) setControllerDefaults(c *Controller) *Controller { //nolint:cyclop,funlen
	// Configured controller defaults.
	if c.SaveSites == nil {
		c.SaveSites = u.Default.SaveSites
	}

	if c.VerifySSL == nil {
		c.VerifySSL = u.Default.VerifySSL
	}

	if c.CertPaths == nil {
		c.CertPaths = u.Default.CertPaths
	}

	if c.HashPII == nil {
		c.HashPII = u.Default.HashPII
	}

	if c.DropPII == nil {
		c.DropPII = u.Default.DropPII
	}

	if c.SaveDPI == nil {
		c.SaveDPI = u.Default.SaveDPI
	}

	if c.SaveIDs == nil {
		c.SaveIDs = u.Default.SaveIDs
	}

	if c.SaveRogue == nil {
		c.SaveRogue = u.Default.SaveRogue
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

	if strings.HasPrefix(c.APIKey, "file://") {
		c.APIKey = u.getPassFromFile(strings.TrimPrefix(c.APIKey, "file://"))
	}

	if c.APIKey == "" {
		if c.Pass == "" {
			c.Pass = defaultPass
		}

		if c.User == "" {
			c.User = defaultUser
		}
	} else {
		// clear out user/pass combo, only use API-key
		c.User = ""
		c.Pass = ""
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
