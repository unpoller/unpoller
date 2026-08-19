// Package inputunas implements the poller.Input interface for UniFi UNAS Pro storage
// consoles. A UNAS Pro is a standalone UniFi OS console with no Network application: no
// sites, no /status endpoint, its own credentials and its own host. That is why it lives
// here rather than inside inputunifi, whose sites-then-devices flow has nothing to offer it.
//
// Endpoint discovery and the JSON shapes are derived from alexgreenbank's unaspoller
// (https://github.com/alexgreenbank/unaspoller, MIT).
package inputunas

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/unpoller/unifi/v5"
	"github.com/unpoller/unpoller/pkg/poller"
	"golift.io/cnfg"
)

// PluginName is the name of this input plugin.
const PluginName = "unas"

const (
	defaultUser    = "unpoller"
	defaultTimeout = 60 * time.Second
)

// InputUNAS contains the running data.
type InputUNAS struct {
	*Config `json:"unas" toml:"unas" xml:"unas" yaml:"unas"`
	Logger  poller.Logger
}

// Device represents one UNAS Pro console to poll. Each may have its own credentials.
type Device struct {
	URL       string            `json:"url"            toml:"url"            xml:"url"           yaml:"url"`
	User      string            `json:"user"           toml:"user"           xml:"user"          yaml:"user"`
	Pass      string            `json:"pass"           toml:"pass"           xml:"pass"          yaml:"pass"`
	VerifySSL *bool             `json:"verify_ssl"     toml:"verify_ssl"     xml:"verify_ssl"    yaml:"verify_ssl"`
	CertPaths []string          `json:"ssl_cert_paths" toml:"ssl_cert_paths" xml:"ssl_cert_path" yaml:"ssl_cert_paths"`
	Timeout   cnfg.Duration     `json:"timeout"        toml:"timeout"        xml:"timeout"       yaml:"timeout"`
	unas      *unifi.UNASClient `json:"-"              toml:"-"              xml:"-"             yaml:"-"`
}

// Config contains our configuration data.
type Config struct {
	sync.RWMutex           // locks a Device's client while it re-authenticates.
	Default      Device    `json:"defaults" toml:"defaults" xml:"default"      yaml:"defaults"`
	Disable      bool      `json:"disable"  toml:"disable"  xml:"disable,attr" yaml:"disable"`
	Devices      []*Device `json:"devices"  toml:"device"   xml:"device"       yaml:"devices"`
}

func init() { // nolint: gochecknoinits
	u := &InputUNAS{}

	poller.NewInput(&poller.InputPlugin{
		Name:   PluginName,
		Input:  u, // this package implements poller.Input for Metrics().
		Config: u, // defines our config data interface.
	})
}

// setDefaults fills a device's unset fields from the defaults block, then from package
// defaults. It never invents a URL: an empty URL is how "not configured" is expressed, and
// synthesizing one (as inputunifi does with localhost:8443) would make the plugin poll a
// host the operator never named. That is what keeps UNAS support opt-in.
func (u *InputUNAS) setDefaults(d *Device) *Device {
	if d.User == "" {
		d.User = u.Default.User
	}

	if d.User == "" {
		d.User = defaultUser
	}

	if d.Pass == "" {
		d.Pass = u.Default.Pass
	}

	if len(d.CertPaths) == 0 {
		d.CertPaths = u.Default.CertPaths
	}

	if d.Timeout.Duration == 0 {
		d.Timeout = u.Default.Timeout
	}

	if d.Timeout.Duration == 0 {
		d.Timeout.Duration = defaultTimeout
	}

	if d.VerifySSL == nil {
		if u.Default.VerifySSL != nil {
			d.VerifySSL = u.Default.VerifySSL
		} else {
			f := false
			d.VerifySSL = &f
		}
	}

	return d
}

// getCerts reads in cert files from disk and stores them as a slice of byte slices.
func (d *Device) getCerts() ([][]byte, error) {
	if len(d.CertPaths) == 0 {
		return nil, nil
	}

	b := make([][]byte, len(d.CertPaths))

	for i, f := range d.CertPaths {
		c, err := os.ReadFile(f)
		if err != nil {
			return nil, fmt.Errorf("reading SSL cert file: %w", err)
		}

		b[i] = c
	}

	return b, nil
}

// login (re-)authenticates to a UNAS console, replacing any existing client.
//
// NewUNASClient stops after Login: unlike NewUnifi it does not finish with GetServerData,
// because that requests /status, which path() rewrites to /proxy/network/status -- and a
// storage-only console has no Network application to answer it.
func (u *InputUNAS) login(d *Device) error {
	u.Lock()
	defer u.Unlock()

	if d.unas != nil && d.unas.Unifi != nil {
		d.unas.CloseIdleConnections()
	}

	certs, err := d.getCerts()
	if err != nil {
		return err
	}

	client, err := unifi.NewUNASClient(&unifi.Config{
		User:      d.User,
		Pass:      d.Pass,
		URL:       d.URL,
		SSLCert:   certs,
		VerifySSL: *d.VerifySSL,
		Timeout:   d.Timeout.Duration,
		ErrorLog:  u.LogErrorf,
		DebugLog:  u.LogDebugf,
	})
	if err != nil {
		d.unas = nil

		return fmt.Errorf("unas console %s: %w", d.URL, err)
	}

	d.unas = client
	u.LogDebugf("Authenticated with UNAS console successfully, %s", d.URL)

	return nil
}
