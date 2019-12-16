package poller

/*
	I consider this file the pinacle example of how to allow a Go application to be configured from a file.
	You can put your configuration into any file format: XML, YAML, JSON, TOML, and you can override any
	struct member using an environment variable. The Duration type is also supported. All of the Config{}
	and Duration{} types and methods are reusable in other projects. Just adjust the data in the struct to
	meet your app's needs. See the New() procedure and Start() method in start.go for example usage.
*/

import (
	"time"

	"github.com/spf13/pflag"
	"golift.io/config"
	"golift.io/unifi"
)

const (
	// AppName is the name of the application.
	AppName = "unifi-poller"
	// ENVConfigPrefix is the prefix appended to an env variable tag name.
	ENVConfigPrefix = "UP"
)

// UnifiPoller contains the application startup data, and auth info for UniFi & Influx.
type UnifiPoller struct {
	Flags *Flags
	*Config
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

// Config represents the core library input data.
type Config struct {
	Poller `json:"poller" toml:"poller" xml:"poller" yaml:"poller"`
}

// Poller is the global config values.
type Poller struct {
	Debug bool `json:"debug" toml:"debug" xml:"debug,attr" yaml:"debug"`
	Quiet bool `json:"quiet,omitempty" toml:"quiet,omitempty" xml:"quiet,attr" yaml:"quiet"`
}

// ParseConfigs parses the poller config and the config for each registered output plugin.
func (u *UnifiPoller) ParseConfigs() error {
	// Parse core config.
	if err := u.ParseInterface(u.Config); err != nil {
		return err
	}

	// Parse output plugin configs.
	outputSync.Lock()
	defer outputSync.Unlock()

	for _, o := range outputs {
		if err := u.ParseInterface(o.Config); err != nil {
			return err
		}
	}

	// Parse input plugin configs.
	inputSync.Lock()
	defer inputSync.Unlock()

	for _, i := range inputs {
		if err := u.ParseInterface(i.Config); err != nil {
			return err
		}
	}

	return nil
}

// ParseInterface parses the config file and environment variables into the provided interface.
func (u *UnifiPoller) ParseInterface(i interface{}) error {
	// Parse config file into provided interface.
	if err := config.ParseFile(i, u.Flags.ConfigFile); err != nil {
		return err
	}

	// Parse environment variables into provided interface.
	_, err := config.ParseENV(i, ENVConfigPrefix)

	return err
}
