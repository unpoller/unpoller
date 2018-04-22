package main

import (
	"net/http"
	"time"
)

// Version will be injected at build time.
var Version = "v0.1"

const (
	// LoginPath is Unifi Controller Login API Path
	LoginPath = "/api/login"
	// ClientPath is Unifi Clients API Path
	ClientPath = "/api/s/default/stat/sta"
	// DevicePath is where we get data about Unifi devices.
	DevicePath = "/api/s/default/stat/device"
	// NetworkPath contains network-configuration data. Not really graphable.
	NetworkPath = "/api/s/default/rest/networkconf"
	// UserGroupPath contains usergroup configurations.
	UserGroupPath = "/api/s/default/rest/usergroup"
	// App defaults in case they're missing from the config.
	defaultConfFile = "/usr/local/etc/unifi-poller/up.conf"
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
	Interval   Dur    `json:"interval" toml:"interval" xml:"interval" yaml:"interval"`
	InfluxURL  string `json:"influx_url" toml:"influx_url" xml:"influx_url" yaml:"influx_url"`
	InfluxUser string `json:"influx_user" toml:"influx_user" xml:"influx_user" yaml:"influx_user"`
	InfluxPass string `json:"influx_pass" toml:"influx_pass" xml:"influx_pass" yaml:"influx_pass"`
	InfluxDB   string `json:"influx_db" toml:"influx_db" xml:"influx_db" yaml:"influx_db"`
	UnifiUser  string `json:"unifi_user" toml:"unifi_user" xml:"unifi_user" yaml:"unifi_user"`
	UnifiPass  string `json:"unifi_pass" toml:"unifi_pass" xml:"unifi_pass" yaml:"unifi_pass"`
	UnifiBase  string `json:"unifi_url" toml:"unifi_url" xml:"unifi_url" yaml:"unifi_url"`
	uniClient  *http.Client
}

// Dur is used to UnmarshalTOML into a time.Duration value.
type Dur struct {
	value time.Duration
}

// UnmarshalTOML parses a duration type from a config file.
func (v *Dur) UnmarshalTOML(data []byte) error {
	unquoted := string(data[1 : len(data)-1])
	dur, err := time.ParseDuration(unquoted)
	if err == nil {
		v.value = dur
	}
	return err
}
