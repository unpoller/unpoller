// Package poller provides the CLI interface to setup unifi-poller.
package poller

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/prometheus/common/version"
	"github.com/spf13/pflag"
)

// New returns a new poller struct.
func New() *UnifiPoller {
	return &UnifiPoller{Config: &Config{Poller: &Poller{}}, Flags: &Flags{}}
}

// Start begins the application from a CLI.
// Parses cli flags, parses config file, parses env vars, sets up logging, then:
// - dumps a json payload OR - executes Run().
func (u *UnifiPoller) Start() error {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags)
	u.Flags.Parse(os.Args[1:])

	if u.Flags.ShowVer {
		fmt.Printf("%s v%s\n", AppName, version.Version)
		return nil // don't run anything else w/ version request.
	}

	cfile, err := getFirstFile(strings.Split(u.Flags.ConfigFile, ","))
	if err != nil {
		return err
	}

	u.Flags.ConfigFile = cfile
	if u.Flags.DumpJSON == "" { // do not print this when dumping JSON.
		u.Logf("Loading Configuration File: %s", u.Flags.ConfigFile)
	}

	// Parse config file and ENV variables.
	if err := u.ParseConfigs(); err != nil {
		return err
	}

	return u.Run()
}

// Parse turns CLI arguments into data structures. Called by Start() on startup.
func (f *Flags) Parse(args []string) {
	f.FlagSet = pflag.NewFlagSet(AppName, pflag.ExitOnError)
	f.Usage = func() {
		fmt.Printf("Usage: %s [--config=/path/to/up.conf] [--version]", AppName)
		f.PrintDefaults()
	}

	f.StringVarP(&f.DumpJSON, "dumpjson", "j", "",
		"This debug option prints a json payload and exits. See man page for more info.")
	f.StringVarP(&f.ConfigFile, "config", "c", DefaultConfFile,
		"Poller config file path. Separating multiple paths with a comma will load the first config file found.")
	f.BoolVarP(&f.ShowVer, "version", "v", false, "Print the version and exit.")
	_ = f.FlagSet.Parse(args) // pflag.ExitOnError means this will never return error.
}

// Run picks a mode and executes the associated functions. This will do one of three things:
// 1. Start the collector routine that polls unifi and reports to influx on an interval. (default)
// 2. Run the collector one time and report the metrics to influxdb. (lambda)
// 3. Start a web server and wait for Prometheus to poll the application for metrics.
func (u *UnifiPoller) Run() error {
	if u.Flags.DumpJSON != "" {
		u.Config.Quiet = true
		if err := u.InitializeInputs(); err != nil {
			return err
		}

		return u.PrintRawMetrics()
	}

	if u.Debug {
		log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ldate)
		u.LogDebugf("Debug Logging Enabled")
	}

	log.Printf("[INFO] UniFi Poller v%v Starting Up! PID: %d", version.Version, os.Getpid())

	if err := u.InitializeInputs(); err != nil {
		return err
	}

	return u.InitializeOutputs()
}
