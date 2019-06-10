package main

import "time"

// Version is injected by the Makefile
var Version = "development"

const (
	// App defaults in case they're missing from the config.
	defaultConfFile = "/etc/unifi-poller/up.conf"
	defaultInterval = 30 * time.Second
	defaultInfxDb   = "unifi"
	defaultInfxUser = "unifi"
	defaultInfxPass = "unifi"
	defaultInfxURL  = "http://127.0.0.1:8086"
	defaultUnifUser = "influx"
	defaultUnifURL  = "https://127.0.0.1:8443"
)

// Config represents the data needed to poll a controller and report to influxdb.
type Config struct {
	Interval   Dur      `json:"interval,_omitempty" toml:"interval,_omitempty" xml:"interval" yaml:"interval"`
	Debug      bool     `json:"debug" toml:"debug" xml:"debug" yaml:"debug"`
	Quiet      bool     `json:"quiet" toml:"quiet" xml:"quiet" yaml:"quiet"`
	VerifySSL  bool     `json:"verify_ssl" toml:"verify_ssl" xml:"verify_ssl" yaml:"verify_ssl"`
	InfluxURL  string   `json:"influx_url,_omitempty" toml:"influx_url,_omitempty" xml:"influx_url" yaml:"influx_url"`
	InfluxUser string   `json:"influx_user,_omitempty" toml:"influx_user,_omitempty" xml:"influx_user" yaml:"influx_user"`
	InfluxPass string   `json:"influx_pass,_omitempty" toml:"influx_pass,_omitempty" xml:"influx_pass" yaml:"influx_pass"`
	InfluxDB   string   `json:"influx_db,_omitempty" toml:"influx_db,_omitempty" xml:"influx_db" yaml:"influx_db"`
	UnifiUser  string   `json:"unifi_user,_omitempty" toml:"unifi_user,_omitempty" xml:"unifi_user" yaml:"unifi_user"`
	UnifiPass  string   `json:"unifi_pass,_omitempty" toml:"unifi_pass,_omitempty" xml:"unifi_pass" yaml:"unifi_pass"`
	UnifiBase  string   `json:"unifi_url,_omitempty" toml:"unifi_url,_omitempty" xml:"unifi_url" yaml:"unifi_url"`
	Sites      []string `json:"sites,_omitempty" toml:"sites,_omitempty" xml:"sites" yaml:"sites"`
}

// Dur is used to UnmarshalTOML into a time.Duration value.
type Dur struct{ value time.Duration }

// UnmarshalTOML parses a duration type from a config file.
func (v *Dur) UnmarshalTOML(data []byte) error {
	unquoted := string(data[1 : len(data)-1])
	dur, err := time.ParseDuration(unquoted)
	if err == nil {
		v.value = dur
	}
	return err
}
