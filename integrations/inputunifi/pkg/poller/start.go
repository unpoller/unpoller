// Package poller provides the CLI interface to setup unifi-poller.
package poller

import (
	"fmt"
	"log"
	"os"

	"github.com/github/hub/version"
	"github.com/spf13/pflag"
)

// New returns a new poller struct preloaded with default values.
// No need to call this if you call Start.c
func New() *UnifiPoller {
	return &UnifiPoller{
		Config: &Config{},
		Flag: &Flag{
			ConfigFile: DefaultConfFile,
		},
	}
}

// Start begins the application from a CLI.
// Parses cli flags, parses config file, parses env vars, sets up logging, then:
// - dumps a json payload OR - executes Run().
func (u *UnifiPoller) Start() error {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags)
	u.Flag.Parse(os.Args[1:])

	if u.Flag.ShowVer {
		fmt.Printf("%s v%s\n", AppName, version.Version)
		return nil // don't run anything else w/ version request.
	}

	if u.Flag.DumpJSON == "" { // do not print this when dumping JSON.
		u.Logf("Loading Configuration File: %s", u.Flag.ConfigFile)
	}

	// Parse config file and ENV variables.
	if err := u.ParseConfigs(); err != nil {
		return err
	}

	if len(u.Config.Controllers) < 1 {
		u.Config.Controllers = []Controller{{
			Sites:     []string{"all"},
			User:      defaultUnifiUser,
			Pass:      "",
			URL:       defaultUnifiURL,
			SaveSites: true,
		}}
	}

	if u.Flag.DumpJSON != "" {
		return u.DumpJSONPayload()
	}

	if u.Config.Debug {
		log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ldate)
		u.LogDebugf("Debug Logging Enabled")
	}

	return u.Run()
}

// Parse turns CLI arguments into data structures. Called by Start() on startup.
func (f *Flag) Parse(args []string) {
	f.FlagSet = pflag.NewFlagSet(AppName, pflag.ExitOnError)
	f.Usage = func() {
		fmt.Printf("Usage: %s [--config=/path/to/up.conf] [--version]", AppName)
		f.PrintDefaults()
	}

	f.StringVarP(&f.DumpJSON, "dumpjson", "j", "",
		"This debug option prints a json payload and exits. See man page for more info.")
	f.StringVarP(&f.ConfigFile, "config", "c", DefaultConfFile, "Poller config file path.")
	f.BoolVarP(&f.ShowVer, "version", "v", false, "Print the version and exit.")
	_ = f.FlagSet.Parse(args) // pflag.ExitOnError means this will never return error.
}

// Run picks a mode and executes the associated functions. This will do one of three things:
// 1. Start the collector routine that polls unifi and reports to influx on an interval. (default)
// 2. Run the collector one time and report the metrics to influxdb. (lambda)
// 3. Start a web server and wait for Prometheus to poll the application for metrics.
func (u *UnifiPoller) Run() error {
	log.Printf("[INFO] UniFi Poller v%v Starting Up! PID: %d", version.Version, os.Getpid())

	for i, c := range u.Config.Controllers {
		if c.Name == "" {
			u.Config.Controllers[i].Name = c.URL
		}

		switch err := u.GetUnifi(c); err {
		case nil:
			u.Logf("Polling UniFi Controller at %s v%s as user %s. Sites: %v",
				c.URL, c.Unifi.ServerVersion, c.User, c.Sites)
		default:
			u.LogErrorf("Controller Auth or Connection failed, but continuing to retry! %s: %v", c.URL, err)
		}
	}

	return u.InitializeOutputs()
}
