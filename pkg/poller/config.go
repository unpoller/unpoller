package poller

/*
	I consider this file the pinacle example of how to allow a Go application to be configured from a file.
	You can put your configuration into any file format: XML, YAML, JSON, TOML, and you can override any
	struct member using an environment variable. The Duration type is also supported. All of the Config{}
	and Duration{} types and methods are reusable in other projects. Just adjust the data in the struct to
	meet your app's needs. See the New() procedure and Start() method in start.go for example usage.
*/

import (
	"sync"
	"time"

	"github.com/spf13/pflag"
	"golift.io/config"
	"golift.io/unifi"
)

// App defaults in case they're missing from the config.
const (
	// AppName is the name of the application.
	AppName          = "unifi-poller"
	defaultUnifiUser = "influx"
	defaultUnifiURL  = "https://127.0.0.1:8443"
)

// ENVConfigPrefix is the prefix appended to an env variable tag
// name before retrieving the value from the OS.
const ENVConfigPrefix = "UP"

// UnifiPoller contains the application startup data, and auth info for UniFi & Influx.
type UnifiPoller struct {
	Flags      *Flags
	Config     *Config
	sync.Mutex // locks the Unifi struct member when re-authing to unifi.
}

// Flags represents the CLI args available and their settings.
type Flags struct {
	ConfigFile string
	DumpJSON   string
	ShowVer    bool
	*pflag.FlagSet
}

// Metrics is a type shared by the exporting and reporting packages.
type Metrics struct {
	TS time.Time
	unifi.Sites
	unifi.IDSList
	unifi.Clients
	*unifi.Devices
}

// Controller represents the configuration for a UniFi Controller.
// Each polled controller may have its own configuration.
type Controller struct {
	VerifySSL bool         `json:"verify_ssl" toml:"verify_ssl" xml:"verify_ssl" yaml:"verify_ssl"`
	SaveIDS   bool         `json:"save_ids" toml:"save_ids" xml:"save_ids" yaml:"save_ids"`
	SaveSites bool         `json:"save_sites,omitempty" toml:"save_sites,omitempty" xml:"save_sites" yaml:"save_sites"`
	Name      string       `json:"name" toml:"name" xml:"name,attr" yaml:"name"`
	User      string       `json:"user,omitempty" toml:"user,omitempty" xml:"user" yaml:"user"`
	Pass      string       `json:"pass,omitempty" toml:"pass,omitempty" xml:"pass" yaml:"pass"`
	URL       string       `json:"url,omitempty" toml:"url,omitempty" xml:"url" yaml:"url"`
	Sites     []string     `json:"sites,omitempty" toml:"sites,omitempty" xml:"sites" yaml:"sites"`
	Unifi     *unifi.Unifi `json:"-" toml:"-" xml:"-" yaml:"-"`
}

// Config represents the data needed to poll a controller and report to influxdb.
// This is all of the data stored in the config file.
// Any with explicit defaults have omitempty on json and toml tags.
type Config struct {
	Poller      `json:"poller" toml:"poller" xml:"poller" yaml:"poller"`
	Controllers []Controller `json:"controller,omitempty" toml:"controller,omitempty" xml:"controller" yaml:"controller"`
}

// Poller is the global config values.
type Poller struct {
	Debug bool `json:"debug" toml:"debug" xml:"debug,attr" yaml:"debug"`
	Quiet bool `json:"quiet,omitempty" toml:"quiet,omitempty" xml:"quiet,attr" yaml:"quiet"`
}

// ParseConfigs parses the poller config and the config for each registered output plugin.
func (u *UnifiPoller) ParseConfigs() error {
	// Parse config file.
	if err := config.ParseFile(u.Config, u.Flags.ConfigFile); err != nil {
		u.Flags.Usage()
		return err
	}

	// Update Config with ENV variable overrides.
	if _, err := config.ParseENV(u.Config, ENVConfigPrefix); err != nil {
		return err
	}

	outputSync.Lock()
	defer outputSync.Unlock()

	for _, o := range outputs {
		// Parse config file for each output plugin.
		if err := config.ParseFile(o.Config, u.Flags.ConfigFile); err != nil {
			return err
		}

		// Update Config for each output with ENV variable overrides.
		if _, err := config.ParseENV(o.Config, ENVConfigPrefix); err != nil {
			return err
		}
	}

	return nil
}
