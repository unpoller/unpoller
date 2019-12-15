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

	"github.com/davidnewhall/unifi-poller/pkg/influxunifi"
	"github.com/spf13/pflag"
	"golift.io/config"
	"golift.io/unifi"
)

// Version is injected by the Makefile
var Version = "development"

const (
	// App defaults in case they're missing from the config.
	appName           = "unifi-poller"
	defaultInterval   = 30 * time.Second
	defaultInfluxDB   = "unifi"
	defaultInfluxUser = "unifi"
	defaultInfluxPass = "unifi"
	defaultInfluxURL  = "http://127.0.0.1:8086"
	defaultUnifiUser  = "influx"
	defaultUnifiURL   = "https://127.0.0.1:8443"
	defaultHTTPListen = "0.0.0.0:9130"
)

// ENVConfigPrefix is the prefix appended to an env variable tag
// name before retrieving the value from the OS.
const ENVConfigPrefix = "UP"

// UnifiPoller contains the application startup data, and auth info for UniFi & Influx.
type UnifiPoller struct {
	Influx     *influxunifi.InfluxUnifi
	Flag       *Flag
	Config     *Config
	LastCheck  time.Time
	sync.Mutex // locks the Unifi struct member when re-authing to unifi.
}

// Flag represents the CLI args available and their settings.
type Flag struct {
	ConfigFile string
	DumpJSON   string
	ShowVer    bool
	*pflag.FlagSet
}

// Controller represents the configuration for a UniFi Controller.
// Each polled controller may have its own configuration.
type Controller struct {
	VerifySSL bool         `json:"verify_ssl" toml:"verify_ssl" xml:"verify_ssl" yaml:"verify_ssl"`
	SaveIDS   bool         `json:"save_ids" toml:"save_ids" xml:"save_ids" yaml:"save_ids"`
	SaveSites bool         `json:"save_sites,omitempty" toml:"save_sites,omitempty" xml:"save_sites" yaml:"save_sites"`
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
	Interval    config.Duration `json:"interval,omitempty" toml:"interval,omitempty" xml:"interval" yaml:"interval"`
	Debug       bool            `json:"debug" toml:"debug" xml:"debug" yaml:"debug"`
	Quiet       bool            `json:"quiet,omitempty" toml:"quiet,omitempty" xml:"quiet" yaml:"quiet"`
	InfxBadSSL  bool            `json:"influx_insecure_ssl" toml:"influx_insecure_ssl" xml:"influx_insecure_ssl" yaml:"influx_insecure_ssl"`
	Mode        string          `json:"mode" toml:"mode" xml:"mode" yaml:"mode"`
	HTTPListen  string          `json:"http_listen" toml:"http_listen" xml:"http_listen" yaml:"http_listen"`
	Namespace   string          `json:"namespace" toml:"namespace" xml:"namespace" yaml:"namespace"`
	InfluxURL   string          `json:"influx_url,omitempty" toml:"influx_url,omitempty" xml:"influx_url" yaml:"influx_url"`
	InfluxUser  string          `json:"influx_user,omitempty" toml:"influx_user,omitempty" xml:"influx_user" yaml:"influx_user"`
	InfluxPass  string          `json:"influx_pass,omitempty" toml:"influx_pass,omitempty" xml:"influx_pass" yaml:"influx_pass"`
	InfluxDB    string          `json:"influx_db,omitempty" toml:"influx_db,omitempty" xml:"influx_db" yaml:"influx_db"`
	Controllers []Controller    `json:"controller,omitempty" toml:"controller,omitempty" xml:"controller" yaml:"controller"`
}
