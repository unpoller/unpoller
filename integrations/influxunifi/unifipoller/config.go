package unifipoller

import (
	"time"

	influx "github.com/influxdata/influxdb1-client/v2"
	"github.com/spf13/pflag"
	"golift.io/unifi"
)

// Version is injected by the Makefile
var Version = "development"

const (
	// App defaults in case they're missing from the config.
	defaultInterval = 30 * time.Second
	defaultInfxDb   = "unifi"
	defaultInfxUser = "unifi"
	defaultInfxPass = "unifi"
	defaultInfxURL  = "http://127.0.0.1:8086"
	defaultUnifUser = "influx"
	defaultUnifURL  = "https://127.0.0.1:8443"
)

// These are environment variables that can be used to override configuration.
// Useful for Docker users.
const (
	ENVConfigMode       = "UP_POLLING_MODE"
	ENVConfigInfluxDB   = "UP_INFLUX_DB"
	ENVConfigInfluxUser = "UP_INFLUX_USER"
	ENVConfigInfluxPass = "UP_INFLUX_PASS"
	ENVConfigInfluxURL  = "UP_INFLUX_URL"
	ENVConfigUnifiUser  = "UP_UNIFI_USER"
	ENVConfigUnifiPass  = "UP_UNIFI_PASS"
	ENVConfigUnifiBase  = "UP_UNIFI_URL"
	ENVConfigReAuth     = "UP_REAUTHENTICATE"
	ENVConfigVerifySSL  = "UP_VERIFY_SSL"
	ENVConfigCollectIDS = "UP_COLLECT_IDS"
	ENVConfigQuiet      = "UP_QUIET_MODE"
	ENVConfigDebug      = "UP_DEBUG_MODE"
	ENVConfigInterval   = "UP_POLLING_INTERVAL"
	ENVConfigMaxErrors  = "UP_MAX_ERRORS"
	ENVConfigSites      = "UP_POLL_SITES"
)

// UnifiPoller contains the application startup data, and auth info for UniFi & Influx.
type UnifiPoller struct {
	Influx     influx.Client
	Unifi      *unifi.Unifi
	Flag       *Flag
	Config     *Config
	errorCount int
	LastCheck  time.Time
}

// Flag represents the CLI args available and their settings.
type Flag struct {
	ConfigFile string
	DumpJSON   string
	ShowVer    bool
	*pflag.FlagSet
}

// Metrics contains all the data from the controller and an influx endpoint to send it to.
type Metrics struct {
	TS time.Time
	unifi.Sites
	unifi.IDSList
	unifi.Clients
	*unifi.Devices
	influx.BatchPoints
}

// Config represents the data needed to poll a controller and report to influxdb.
// This is all of the data stored in the config file.
type Config struct {
	MaxErrors  int      `json:"max_errors,_omitempty" toml:"max_errors,_omitempty" xml:"max_errors" yaml:"max_errors"`
	Interval   Duration `json:"interval,_omitempty" toml:"interval,_omitempty" xml:"interval" yaml:"interval"`
	Debug      bool     `json:"debug" toml:"debug" xml:"debug" yaml:"debug"`
	Quiet      bool     `json:"quiet,_omitempty" toml:"quiet,_omitempty" xml:"quiet" yaml:"quiet"`
	VerifySSL  bool     `json:"verify_ssl" toml:"verify_ssl" xml:"verify_ssl" yaml:"verify_ssl"`
	CollectIDS bool     `json:"collect_ids" toml:"collect_ids" xml:"collect_ids" yaml:"collect_ids"`
	ReAuth     bool     `json:"reauthenticate" toml:"reauthenticate" xml:"reauthenticate" yaml:"reauthenticate"`
	Mode       string   `json:"mode" toml:"mode" xml:"mode" yaml:"mode"`
	InfluxURL  string   `json:"influx_url,_omitempty" toml:"influx_url,_omitempty" xml:"influx_url" yaml:"influx_url"`
	InfluxUser string   `json:"influx_user,_omitempty" toml:"influx_user,_omitempty" xml:"influx_user" yaml:"influx_user"`
	InfluxPass string   `json:"influx_pass,_omitempty" toml:"influx_pass,_omitempty" xml:"influx_pass" yaml:"influx_pass"`
	InfluxDB   string   `json:"influx_db,_omitempty" toml:"influx_db,_omitempty" xml:"influx_db" yaml:"influx_db"`
	UnifiUser  string   `json:"unifi_user,_omitempty" toml:"unifi_user,_omitempty" xml:"unifi_user" yaml:"unifi_user"`
	UnifiPass  string   `json:"unifi_pass,_omitempty" toml:"unifi_pass,_omitempty" xml:"unifi_pass" yaml:"unifi_pass"`
	UnifiBase  string   `json:"unifi_url,_omitempty" toml:"unifi_url,_omitempty" xml:"unifi_url" yaml:"unifi_url"`
	Sites      []string `json:"sites,_omitempty" toml:"sites,_omitempty" xml:"sites" yaml:"sites"`
}

// Duration is used to UnmarshalTOML into a time.Duration value.
type Duration struct{ time.Duration }

// UnmarshalText parses a duration type from a config file.
func (d *Duration) UnmarshalText(data []byte) (err error) {
	d.Duration, err = time.ParseDuration(string(data))
	return
}
