package main

import (
	"fmt"

	"github.com/unifi-poller/poller"
	"golift.io/cnfg"
)

// Config represents the data that is unmarshalled from the up.conf config file for this plugins.
// See up.conf.example.mysql for sample input data.
type Config struct {
	Disable  bool          `json:"disable" toml:"disable" xml:"disable" yaml:"disable"`
	Interval cnfg.Duration `json:"interval" toml:"interval" xml:"interval" yaml:"interval"`
	Host     string        `json:"host" toml:"host" xml:"host" yaml:"host"`
	User     string        `json:"user" toml:"user" xml:"user" yaml:"user"`
	Pass     string        `json:"pass" toml:"pass" xml:"pass" yaml:"pass"`
	DB       string        `json:"db" toml:"db" xml:"db" yaml:"db"`
	Devices  []Device      `json:"devices" toml:"devices" xml:"device" yaml:"devices"`
	Clients  *Clients      `json:"clients" toml:"clients" xml:"clients" yaml:"clients"`
}

// Device represents the configuration to save a devices' data.
// Type is one of uap, usw, ugw, udm.
// Table represents the mysql table name we save these fields to.
// Fields is a map of api response data key -> mysql column.
type Device struct {
	Type   string            `json:"type" toml:"type" xml:"type" yaml:"type"`
	Table  string            `json:"table" toml:"table" xml:"table" yaml:"table"`
	Fields map[string]string `json:"fields" toml:"fields" xml:"field" yaml:"fields"`
}

// Clients represents the configuration to save clients' data.
// Table represents the mysql table name we save these fields to.
// Fields is a map of api response data key -> mysql column.
type Clients struct {
	Table  string            `json:"table" toml:"table" xml:"table" yaml:"table"`
	Fields map[string]string `json:"fields" toml:"fields" xml:"field" yaml:"fields"`
}

// Pointers are ignored during ENV variable unmarshal, avoid pointers to your config.
// Only capital (exported) members are unmarshaled when passed into poller.NewOutput().
type plugin struct {
	*Config `json:"mysql" toml:"mysql" xml:"mysql" yaml:"mysql"`
	poller.Collect
}

func init() {
	u := &plugin{Config: &Config{}}

	poller.NewOutput(&poller.Output{
		Name:   "mysql",
		Config: u, // pass in the struct *above* your config (so it can see the struct tags).
		Method: u.Run,
	})
}

// main() is required, but it shouldn't do much as it's not used in plugin mode.
func main() {
	fmt.Println("this is a unifi-poller plugin; not an application")
}

// Run gets called by poller core code. Return when the plugin stops working or has an error.
// In other words, don't run your code in a go routine, it already is.
func (p *plugin) Run(c poller.Collect) error {
	if p.Collect = c; c == nil || p.Config == nil || p.Disable {
		return nil // no config or disabled, bail out.
	}

	return p.runCollector()
}
