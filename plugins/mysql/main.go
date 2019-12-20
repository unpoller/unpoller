package main

import (
	"fmt"

	"github.com/davidnewhall/unifi-poller/pkg/poller"
	"golift.io/cnfg"
)

// mysqlConfig represents the data that is unmarshalled from the up.conf config file for this plugins.
type mysqlConfig struct {
	Interval cnfg.Duration `json:"interval" toml:"interval" xml:"interval" yaml:"interval"`
	Host     string        `json:"host" toml:"host" xml:"host" yaml:"host"`
	User     string        `json:"user" toml:"user" xml:"user" yaml:"user"`
	Pass     string        `json:"pass" toml:"pass" xml:"pass" yaml:"pass"`
	DB       string        `json:"db" toml:"db" xml:"db" yaml:"db"`
	Table    string        `json:"table" toml:"table" xml:"table" yaml:"table"`
	// Maps do not work with ENV VARIABLES yet, but may in the future.
	Fields []string `json:"fields" toml:"fields" xml:"field" yaml:"fields"`
}

// Pointers are ignored during ENV variable unmarshal, avoid pointers to your config.
// Only capital (exported) members are unmarshaled when passed into poller.NewOutput().
type application struct {
	Config mysqlConfig `json:"mysql" toml:"mysql" xml:"mysql" yaml:"mysql"`
}

func init() {
	u := &application{Config: mysqlConfig{}}

	poller.NewOutput(&poller.Output{
		Name:   "mysql",
		Config: u, // pass in the struct *above* your config (so it can see the struct tags).
		Method: u.Run,
	})
}

func main() {
	fmt.Println("this is a unifi-poller plugin; not an application")
}

func (a *application) Run(c poller.Collect) error {
	c.Logf("mysql plugin is not finished")
	return nil
}
