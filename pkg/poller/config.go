package poller

/*
	I consider this file the pinacle example of how to allow a Go application to be configured from a file.
	You can put your configuration into any file format: XML, YAML, JSON, TOML, and you can override any
	struct member using an environment variable. The Duration type is also supported. All of the Config{}
	and Duration{} types and methods are reusable in other projects. Just adjust the data in the struct to
	meet your app's needs. See the New() procedure and Start() method in start.go for example usage.
*/

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/davidnewhall/unifi-poller/pkg/influxunifi"
	"github.com/spf13/pflag"
	"golift.io/unifi"
	yaml "gopkg.in/yaml.v2"
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
const ENVConfigPrefix = "UP_"

// UnifiPoller contains the application startup data, and auth info for UniFi & Influx.
type UnifiPoller struct {
	Influx     *influxunifi.InfluxUnifi
	Unifi      *unifi.Unifi
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

// Config represents the data needed to poll a controller and report to influxdb.
// This is all of the data stored in the config file.
// Any with explicit defaults have omitempty on json and toml tags.
type Config struct {
	Interval   Duration `json:"interval,omitempty" toml:"interval,omitempty" xml:"interval" yaml:"interval"`
	Debug      bool     `json:"debug" toml:"debug" xml:"debug" yaml:"debug"`
	Quiet      bool     `json:"quiet,omitempty" toml:"quiet,omitempty" xml:"quiet" yaml:"quiet"`
	VerifySSL  bool     `json:"verify_ssl" toml:"verify_ssl" xml:"verify_ssl" yaml:"verify_ssl"`
	SaveIDS    bool     `json:"save_ids" toml:"save_ids" xml:"save_ids" yaml:"save_ids"`
	ReAuth     bool     `json:"reauthenticate" toml:"reauthenticate" xml:"reauthenticate" yaml:"reauthenticate"`
	InfxBadSSL bool     `json:"influx_insecure_ssl" toml:"influx_insecure_ssl" xml:"influx_insecure_ssl" yaml:"influx_insecure_ssl"`
	SaveSites  bool     `json:"save_sites,omitempty" toml:"save_sites,omitempty" xml:"save_sites" yaml:"save_sites"`
	Mode       string   `json:"mode" toml:"mode" xml:"mode" yaml:"mode"`
	HTTPListen string   `json:"http_listen" toml:"http_listen" xml:"http_listen" yaml:"http_listen"`
	Namespace  string   `json:"namespace" toml:"namespace" xml:"namespace" yaml:"namespace"`
	InfluxURL  string   `json:"influx_url,omitempty" toml:"influx_url,omitempty" xml:"influx_url" yaml:"influx_url"`
	InfluxUser string   `json:"influx_user,omitempty" toml:"influx_user,omitempty" xml:"influx_user" yaml:"influx_user"`
	InfluxPass string   `json:"influx_pass,omitempty" toml:"influx_pass,omitempty" xml:"influx_pass" yaml:"influx_pass"`
	InfluxDB   string   `json:"influx_db,omitempty" toml:"influx_db,omitempty" xml:"influx_db" yaml:"influx_db"`
	UnifiUser  string   `json:"unifi_user,omitempty" toml:"unifi_user,omitempty" xml:"unifi_user" yaml:"unifi_user"`
	UnifiPass  string   `json:"unifi_pass,omitempty" toml:"unifi_pass,omitempty" xml:"unifi_pass" yaml:"unifi_pass"`
	UnifiBase  string   `json:"unifi_url,omitempty" toml:"unifi_url,omitempty" xml:"unifi_url" yaml:"unifi_url"`
	Sites      []string `json:"sites,omitempty" toml:"sites,omitempty" xml:"sites" yaml:"sites"`
}

// Duration is used to UnmarshalTOML into a time.Duration value.
type Duration struct{ time.Duration }

// UnmarshalText parses a duration type from a config file.
func (d *Duration) UnmarshalText(data []byte) (err error) {
	d.Duration, err = time.ParseDuration(string(data))
	return
}

// ParseFile parses and returns our configuration data.
func (c *Config) ParseFile(configFile string) error {
	switch buf, err := ioutil.ReadFile(configFile); {
	case err != nil:
		return err
	case strings.Contains(configFile, ".json"):
		return json.Unmarshal(buf, c)
	case strings.Contains(configFile, ".xml"):
		return xml.Unmarshal(buf, c)
	case strings.Contains(configFile, ".yaml"):
		return yaml.Unmarshal(buf, c)
	default:
		return toml.Unmarshal(buf, c)
	}
}

// ParseENV copies environment variables into configuration values.
// This is useful for Docker users that find it easier to pass ENV variables
// than a specific configuration file. Uses reflection to find struct tags.
// This method uses the json struct tag member to match environment variables.
// Use a custom tag name by changing "json" below, but that's overkill for this app.
func (c *Config) ParseENV() error {
	t := reflect.TypeOf(*c)             // Get "types" from the Config struct.
	for i := 0; i < t.NumField(); i++ { // Loop each Config struct member
		tag := t.Field(i).Tag.Get("json")                 // Get the ENV variable name from "json" struct tag
		tag = strings.Split(strings.ToUpper(tag), ",")[0] // Capitalize and remove ,omitempty suffix
		env := os.Getenv(ENVConfigPrefix + tag)           // Then pull value from OS.
		if tag == "" || env == "" {                       // Skip if either are empty.
			continue
		}

		// Reflect and update the u.Config struct member at position i.
		switch field := reflect.ValueOf(c).Elem().Field(i); field.Type().String() {
		// Handle each member type appropriately (differently).
		case "string":
			// This is a reflect package method to update a struct member by index.
			field.SetString(env)

		case "int":
			val, err := strconv.Atoi(env)
			if err != nil {
				return fmt.Errorf("%s: %v", tag, err)
			}
			field.Set(reflect.ValueOf(val))

		case "[]string":
			field.Set(reflect.ValueOf(strings.Split(env, ",")))

		case path.Base(t.PkgPath()) + ".Duration":
			val, err := time.ParseDuration(env)
			if err != nil {
				return fmt.Errorf("%s: %v", tag, err)
			}
			field.Set(reflect.ValueOf(Duration{val}))

		case "bool":
			val, err := strconv.ParseBool(env)
			if err != nil {
				return fmt.Errorf("%s: %v", tag, err)
			}
			field.SetBool(val)
		}
		// Add more types here if more types are added to the config struct.
	}

	return nil
}
